/**
 * DeepPhotos API Utility — Central fetch wrapper with:
 *  - JWT token injection (reads from localStorage automatically)
 *  - In-memory response cache with configurable TTL
 *  - Request deduplication (same in-flight request won't fire twice)
 *  - Selective cache invalidation by key prefix
 */

import { appState } from '$lib/state.svelte';

// ─── Cache Store ─────────────────────────────────────────────────────────────

interface CacheEntry {
	data: unknown;
	expiresAt: number;
}

const cache = new Map<string, CacheEntry>();
const inFlight = new Map<string, Promise<unknown>>();

/** Default TTL per endpoint (in milliseconds) */
const TTL_MAP: Record<string, number> = {
	'/api/photos':         30_000,  // 30 sec  — gallery data
	'/api/albums':         60_000,  // 60 sec  — album list
	'/api/locked-folders': 60_000,  // 60 sec  — locked folders
	'/api/audit-logs':     120_000, // 2 min   — audit logs (rarely change)
	'/api/users':          60_000,  // 60 sec  — user list
};

function getTTL(path: string): number {
	for (const prefix of Object.keys(TTL_MAP)) {
		if (path.startsWith(prefix)) return TTL_MAP[prefix];
	}
	return 0; // No cache for unknown endpoints
}

function getCacheKey(path: string, init?: RequestInit): string {
	const method = init?.method?.toUpperCase() || 'GET';
	return `${method}:${path}`;
}

// ─── Cache Invalidation ───────────────────────────────────────────────────────

/**
 * Invalidate all cached entries that start with the given prefix.
 * Call this after mutations (POST/PUT/DELETE) to force fresh data.
 *
 * Example: invalidateCache('/api/photos') clears all photo-related caches.
 */
export function invalidateCache(prefix: string) {
	for (const key of cache.keys()) {
		if (key.includes(prefix)) {
			cache.delete(key);
		}
	}
}

/** Clears the entire cache (e.g. on logout) */
export function clearAllCache() {
	cache.clear();
	inFlight.clear();
}

// ─── Core API Fetch ───────────────────────────────────────────────────────────

/**
 * Central API fetch utility. Automatically:
 *  1. Prepends apiBaseUrl
 *  2. Injects JWT Authorization header
 *  3. Caches GET responses with TTL (returns stale data while revalidating)
 *  4. Deduplicates concurrent identical GET requests
 *
 * @param path   - API path starting with /api/...
 * @param init   - Standard RequestInit options
 * @param noCache - Set true to bypass cache (for forced refreshes)
 */
export async function apiFetch(
	path: string,
	init?: RequestInit,
	noCache = false
): Promise<Response> {
	const token =
		typeof localStorage !== 'undefined'
			? localStorage.getItem('deepphotos_token') || ''
			: '';

	const headers = new Headers(init?.headers);
	headers.set('Accept', 'application/json');
	if (token && !headers.has('Authorization')) {
		headers.set('Authorization', `Bearer ${token}`);
	}

	const method = init?.method?.toUpperCase() || 'GET';
	const url = `${appState.apiBaseUrl}${path}`;

	// ── Non-GET requests: bypass cache, invalidate related entries ──
	if (method !== 'GET') {
		const res = await fetch(url, { ...init, headers });
		// Invalidate cache entries related to this resource on success
		if (res.ok) {
			const basePath = path.split('?')[0].split('/').slice(0, 3).join('/');
			invalidateCache(basePath);
		}
		return res;
	}

	// ── GET requests: use cache ──
	if (!noCache) {
		const cacheKey = getCacheKey(path, init);
		const ttl = getTTL(path);

		// Return valid cached entry immediately
		const cached = cache.get(cacheKey);
		if (cached && cached.expiresAt > Date.now()) {
			// Return a synthetic Response from cached data
			return new Response(JSON.stringify(cached.data), {
				status: 200,
				headers: { 'Content-Type': 'application/json', 'X-Cache': 'HIT' }
			});
		}

		// Deduplicate: if same request is already in-flight, await it
		if (inFlight.has(cacheKey)) {
			const data = await inFlight.get(cacheKey)!;
			return new Response(JSON.stringify(data), {
				status: 200,
				headers: { 'Content-Type': 'application/json', 'X-Cache': 'DEDUP' }
			});
		}

		if (ttl > 0) {
			// Start fetch and register as in-flight
			const fetchPromise = fetch(url, { ...init, headers })
				.then(async (res) => {
					if (res.ok) {
						const data = await res.json();
						cache.set(cacheKey, { data, expiresAt: Date.now() + ttl });
						inFlight.delete(cacheKey);
						return data;
					}
					inFlight.delete(cacheKey);
					return null;
				})
				.catch((e) => {
					inFlight.delete(cacheKey);
					throw e;
				});

			inFlight.set(cacheKey, fetchPromise);
			const data = await fetchPromise;

			if (data === null) {
				// Backend returned error — return raw response for caller to handle
				return fetch(url, { ...init, headers });
			}

			return new Response(JSON.stringify(data), {
				status: 200,
				headers: { 'Content-Type': 'application/json', 'X-Cache': 'MISS' }
			});
		}
	}

	// Fallthrough: no cache configured for this path
	return fetch(url, { ...init, headers });
}
