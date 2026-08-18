<script lang="ts">
	import { onMount } from 'svelte';
	import { appState } from '$lib/state.svelte';
	import { apiFetch } from '$lib/api';
	import { Trash2, RotateCcw, Image as ImageIcon } from 'lucide-svelte';

	interface DeletedPhoto {
		id: string;
		title: string;
		filename: string;
		size: number;
		updated_at: string;
	}

	let binnedItems = $state<DeletedPhoto[]>([]);
	let isLoading = $state(true);

	async function fetchBinItems() {
		isLoading = true;
		try {
			const res = await apiFetch('/api/photos?deleted=true');
			if (res.ok) {
				binnedItems = await res.json();
			}
		} catch (e) {
			console.warn('API error fetching bin items:', e);
		} finally {
			isLoading = false;
		}
	}

	onMount(() => {
		fetchBinItems();
	});

	async function restoreItem(id: string) {
		try {
			await apiFetch('/api/photos/batch-restore', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ ids: [id] })
			});
			appState.refreshPhotos();
			fetchBinItems();
		} catch (err) {
			console.error('Error restoring photo:', err);
		}
	}

	async function emptyBin() {
		if (binnedItems.length === 0) return;
		if (!confirm(`Permanently purge ${binnedItems.length} items from server storage? This action cannot be undone.`)) return;

		const ids = binnedItems.map(item => item.id);
		try {
			await apiFetch('/api/photos/batch-delete', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ ids })
			});
			appState.refreshPhotos();
			fetchBinItems();
		} catch (err) {
			console.error('Error emptying bin:', err);
		}
	}

	function formatBytes(bytes: number): string {
		if (!bytes) return '0 B';
		const k = 1024;
		const sizes = ['B', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
	}
</script>

<div class="h-full flex flex-col gap-6 animate-fade-in">
	<div class="flex items-center justify-between shrink-0">
		<div>
			<h1 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
				<Trash2 class="w-6 h-6 text-rose-500" />
				Bin / Trash
			</h1>
			<p class="text-xs text-slate-500 dark:text-slate-400">Items soft-deleted from gallery storage</p>
		</div>
		<button
			type="button"
			onclick={emptyBin}
			disabled={binnedItems.length === 0}
			class="px-4 py-2 rounded-xl bg-rose-50 dark:bg-rose-950/60 text-rose-600 dark:text-rose-400 hover:bg-rose-100 dark:hover:bg-rose-900 text-xs font-bold transition-all cursor-pointer border border-transparent dark:border-rose-800 disabled:opacity-50"
		>
			Empty Bin
		</button>
	</div>

	<div class="flex-1 overflow-y-auto min-h-0">
		{#if !isLoading && binnedItems.length === 0}
			<div class="p-12 text-center rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm max-w-md mx-auto my-12 space-y-3">
				<div class="w-12 h-12 rounded-2xl bg-rose-100 dark:bg-rose-950 text-rose-500 flex items-center justify-center mx-auto">
					<Trash2 class="w-6 h-6" />
				</div>
				<h3 class="text-base font-bold text-slate-900 dark:text-white">Bin is empty</h3>
				<p class="text-xs text-slate-500 dark:text-slate-400">Deleted photos will appear here before permanent purging.</p>
			</div>
		{:else}
			<div class="p-6 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm space-y-3">
				{#each binnedItems as item}
					<div class="p-3.5 rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700 flex items-center justify-between">
						<div class="flex items-center gap-3">
							<div class="w-10 h-10 rounded-lg bg-rose-100 dark:bg-rose-950 text-rose-600 dark:text-rose-400 flex items-center justify-center">
								<ImageIcon class="w-5 h-5" />
							</div>
							<div>
								<p class="text-xs font-bold text-slate-800 dark:text-white">{item.title || item.filename}</p>
								<p class="text-[11px] text-slate-400">Deleted {item.updated_at || 'Recently'} • {formatBytes(item.size)}</p>
							</div>
						</div>
						<button
							type="button"
							onclick={() => restoreItem(item.id)}
							class="px-3 py-1.5 rounded-lg border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 hover:text-sky-600 dark:hover:text-sky-400 hover:border-sky-300 dark:hover:border-sky-700 text-xs font-semibold flex items-center gap-1 cursor-pointer"
						>
							<RotateCcw class="w-3.5 h-3.5" /> Restore
						</button>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>
