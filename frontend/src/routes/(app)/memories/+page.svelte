<script lang="ts">
	import { onMount } from 'svelte';
	import { apiFetch, getMediaUrl } from '$lib/api';
	import { notify, confirmDialog } from '$lib/notify.svelte';
	import { Sparkles, Calendar, Heart, MapPin, Image as ImageIcon, Loader2, Plus, Trash2, X, ChevronLeft, ChevronRight, Play } from 'lucide-svelte';

	interface PhotoItem {
		id: string;
		title: string;
		filename: string;
		file_type?: string;
		mime_type?: string;
		created_at?: string;
		taken_at?: string;
		is_favorite?: boolean;
		url?: string;
		thumbnail_url?: string;
	}

	interface MemoryHighlight {
		id: string;
		title: string;
		subtitle: string;
		itemsCount: number;
		cover: string;
		photos: PhotoItem[];
		isCustom?: boolean;
	}

	let memoryHighlights = $state<MemoryHighlight[]>([]);
	let isLoading = $state(true);

	// Create Memory Modal State
	let showCreateModal = $state(false);
	let newMemoryTitle = $state('');
	let newMemoryDesc = $state('');
	let isSubmitting = $state(false);

	// Selected Memory Detail Modal State
	let activeMemory = $state<MemoryHighlight | null>(null);
	let selectedPhotoIndex = $state<number | null>(null);

	async function fetchMemories() {
		isLoading = true;
		try {
			const highlights: MemoryHighlight[] = [];

			// 1. Fetch Custom User Memories from /api/memories
			const memRes = await apiFetch('/api/memories');
			if (memRes.ok) {
				const userMemories = await memRes.json();
				for (const m of userMemories) {
					// Get photos for this memory
					const detailRes = await apiFetch(`/api/memories/${m.id}`);
					let photos: PhotoItem[] = [];
					if (detailRes.ok) {
						const detail = await detailRes.json();
						photos = (detail.photos || []).map((p: any) => ({
							...p,
							url: getMediaUrl(`/api/photos/${p.id}/file`),
							thumbnail_url: getMediaUrl(`/api/photos/${p.id}/thumbnail`)
						}));
					}

					const cover = photos.length > 0
						? (photos[0].thumbnail_url || photos[0].url || '')
						: '/empty-folder.png';

					highlights.push({
						id: m.id,
						title: m.title,
						subtitle: m.description || `${photos.length} photos`,
						itemsCount: photos.length,
						cover: cover,
						photos: photos,
						isCustom: true
					});
				}
			}

			// 2. Fetch Automatic Highlights from /api/media?type=gallery&deleted=false
			const photoRes = await apiFetch('/api/media?type=gallery&deleted=false');
			if (photoRes.ok) {
				const rawPhotos: PhotoItem[] = await photoRes.json();
				const mappedPhotos = rawPhotos.map(p => ({
					...p,
					url: getMediaUrl(`/api/media/${p.id}/file`),
					thumbnail_url: getMediaUrl(`/api/media/${p.id}/thumbnail`)
				}));

				// Favorites Highlight
				const favs = mappedPhotos.filter(p => p.is_favorite);
				if (favs.length > 0) {
					highlights.push({
						id: 'favs',
						title: 'Favorites Showcase',
						subtitle: 'Starred Moments',
						itemsCount: favs.length,
						cover: favs[0].thumbnail_url || favs[0].url || '',
						photos: favs,
						isCustom: false
					});
				}

				// Group Photos by Year
				const yearGroups = new Map<string, PhotoItem[]>();
				for (const p of mappedPhotos) {
					const dateStr = p.taken_at || p.created_at || '';
					const year = dateStr ? new Date(dateStr).getFullYear().toString() : 'Recent';
					if (!yearGroups.has(year)) {
						yearGroups.set(year, []);
					}
					yearGroups.get(year)!.push(p);
				}

				for (const [year, items] of yearGroups.entries()) {
					highlights.push({
						id: `year-${year}`,
						title: year === 'Recent' ? 'Recent Highlights' : `Memories from ${year}`,
						subtitle: `${items.length} moments saved`,
						itemsCount: items.length,
						cover: items[0].thumbnail_url || items[0].url || '',
						photos: items,
						isCustom: false
					});
				}
			}

			memoryHighlights = highlights;
		} catch (err) {
			console.warn('Error fetching memories:', err);
		} finally {
			isLoading = false;
		}
	}

	onMount(() => {
		fetchMemories();
	});

	async function handleCreateMemory(e: Event) {
		e.preventDefault();
		if (!newMemoryTitle.trim()) {
			notify.error('Please enter a memory title.');
			return;
		}

		isSubmitting = true;
		try {
			const res = await apiFetch('/api/memories', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					title: newMemoryTitle.trim(),
					description: newMemoryDesc.trim()
				})
			});

			if (res.ok) {
				notify.success(`Created Memory "${newMemoryTitle}"! You can now add photos to it from Gallery.`);
				showCreateModal = false;
				newMemoryTitle = '';
				newMemoryDesc = '';
				await fetchMemories();
			} else {
				notify.error('Failed to create memory.');
			}
		} catch (err) {
			console.error('Error creating memory:', err);
			notify.error('Network error creating memory.');
		} finally {
			isSubmitting = false;
		}
	}

	async function handleDeleteMemory(memory: MemoryHighlight, e: Event) {
		e.stopPropagation();
		const confirmed = await confirmDialog.ask({
			title: 'Delete Memory Collection',
			message: `Delete memory collection "${memory.title}"?`,
			confirmText: 'Yes, Delete Memory',
			cancelText: 'Cancel',
			type: 'danger'
		});
		if (!confirmed) return;

		try {
			const res = await apiFetch(`/api/memories/${memory.id}`, { method: 'DELETE' });
			if (res.ok) {
				notify.success(`Memory "${memory.title}" deleted.`);
				if (activeMemory?.id === memory.id) activeMemory = null;
				await fetchMemories();
			} else {
				notify.error('Failed to delete memory.');
			}
		} catch (err) {
			console.error('Error deleting memory:', err);
			notify.error('Network error deleting memory.');
		}
	}
</script>

<div class="h-full flex flex-col gap-6 animate-fade-in">
	<div class="flex items-center justify-between shrink-0">
		<div>
			<h1 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
				<Sparkles class="w-6 h-6 text-sky-500" />
				Memories & Highlights
			</h1>
			<p class="text-xs text-slate-500 dark:text-slate-400">Rediscover special moments from past years & custom collections</p>
		</div>

		<button
			type="button"
			onclick={() => showCreateModal = true}
			class="px-4 py-2 rounded-xl bg-sky-400 hover:bg-sky-500 text-white text-xs font-bold shadow-sm shadow-sky-300/50 flex items-center gap-1.5 cursor-pointer transition-all"
		>
			<Plus class="w-4 h-4" />
			<span>Create Memory</span>
		</button>
	</div>

	<div class="flex-1 overflow-y-auto min-h-0">
		{#if isLoading}
			<div class="p-12 text-center flex flex-col items-center justify-center gap-3">
				<Loader2 class="w-8 h-8 animate-spin text-sky-500" />
				<p class="text-xs text-slate-500 dark:text-slate-400">Loading your memories...</p>
			</div>
		{:else if memoryHighlights.length === 0}
			<div class="p-12 text-center rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm max-w-md mx-auto my-12 space-y-3">
				<div class="w-12 h-12 rounded-2xl bg-sky-100 dark:bg-sky-950 text-sky-500 flex items-center justify-center mx-auto">
					<Sparkles class="w-6 h-6" />
				</div>
				<h3 class="text-base font-bold text-slate-900 dark:text-white">No memories yet</h3>
				<p class="text-xs text-slate-500 dark:text-slate-400">Click <strong>Create Memory</strong> above or select photos in Gallery to build memory collections.</p>
			</div>
		{:else}
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				{#each memoryHighlights as memory (memory.id)}
					<div
						role="button"
						tabindex="0"
						onclick={() => activeMemory = memory}
						onkeydown={(e) => e.key === 'Enter' && (activeMemory = memory)}
						class="rounded-3xl overflow-hidden bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:shadow-md transition-all group cursor-pointer relative"
					>
						<div class="h-48 overflow-hidden relative">
							<img
								src={memory.cover}
								alt={memory.title}
								class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
								onerror={(e) => {
									const img = e.currentTarget as HTMLImageElement;
									img.onerror = null;
									img.src = '/empty-folder.png';
								}}
							/>
							<div class="absolute inset-0 bg-gradient-to-t from-slate-950/85 via-slate-950/30 to-transparent p-6 flex flex-col justify-end text-white">
								<div class="flex items-center justify-between">
									<span class="text-[11px] font-semibold text-sky-300 flex items-center gap-1">
										<Calendar class="w-3.5 h-3.5" /> {memory.subtitle}
									</span>

									{#if memory.isCustom}
										<button
											type="button"
											onclick={(e) => handleDeleteMemory(memory, e)}
											class="p-1.5 rounded-full bg-white/20 hover:bg-rose-500 text-white backdrop-blur-md transition-colors"
											title="Delete Memory"
										>
											<Trash2 class="w-3.5 h-3.5" />
										</button>
									{/if}
								</div>

								<h3 class="text-lg font-bold mt-1">{memory.title}</h3>
								<p class="text-xs text-slate-300 flex items-center gap-1 mt-0.5">
									<ImageIcon class="w-3.5 h-3.5 text-sky-400" /> {memory.itemsCount} photos
								</p>
							</div>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<!-- Create Memory Modal -->
{#if showCreateModal}
	<div class="fixed inset-0 z-[10001] bg-slate-950/80 backdrop-blur-md flex items-center justify-center p-4 animate-fade-in">
		<div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-2xl max-w-md w-full space-y-5 animate-fade-in">
			<div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-4">
				<div class="flex items-center gap-2.5">
					<div class="w-10 h-10 rounded-2xl bg-sky-100 dark:bg-sky-950 text-sky-500 flex items-center justify-center">
						<Sparkles class="w-5 h-5" />
					</div>
					<div>
						<h3 class="text-base font-bold text-slate-900 dark:text-white">Create New Memory</h3>
						<p class="text-xs text-slate-500 dark:text-slate-400">Organize moments into a Memory Collection</p>
					</div>
				</div>
				<button type="button" onclick={() => showCreateModal = false} class="p-1 rounded-xl text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 cursor-pointer">
					<X class="w-5 h-5" />
				</button>
			</div>

			<form onsubmit={handleCreateMemory} class="space-y-4">
				<div class="space-y-1.5">
					<label for="new-memory-title" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Memory Title *</label>
					<input
						id="new-memory-title"
						type="text"
						bind:value={newMemoryTitle}
						placeholder="e.g. Summer Trip 2026, Family Gathering"
						class="w-full h-10 px-3.5 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400"
						required
					/>
				</div>

				<div class="space-y-1.5">
					<label for="new-memory-desc" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Description (Optional)</label>
					<textarea
						id="new-memory-desc"
						bind:value={newMemoryDesc}
						rows={2}
						placeholder="Add notes about this memory..."
						class="w-full p-3 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400 resize-none"
					></textarea>
				</div>

				<div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-100 dark:border-slate-800">
					<button
						type="button"
						onclick={() => showCreateModal = false}
						class="px-4 py-2 rounded-xl text-xs font-semibold text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 cursor-pointer"
					>
						Cancel
					</button>
					<button
						type="submit"
						disabled={isSubmitting}
						class="px-5 py-2 rounded-xl bg-sky-400 hover:bg-sky-500 text-white text-xs font-bold shadow-sm shadow-sky-300/50 flex items-center gap-1.5 disabled:opacity-50 cursor-pointer"
					>
						<Sparkles class="w-4 h-4" />
						<span>Create Memory</span>
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Memory Detail Photos View Modal -->
{#if activeMemory}
	<div class="fixed inset-0 z-[10000] bg-slate-950/80 backdrop-blur-md flex flex-col p-4 md:p-8 animate-fade-in">
		<div class="flex items-center justify-between text-white shrink-0 mb-6 max-w-7xl w-full mx-auto">
			<div class="flex items-center gap-3">
				<button
					type="button"
					onclick={() => activeMemory = null}
					class="p-2 rounded-xl bg-white/10 hover:bg-white/20 text-white transition-colors cursor-pointer"
				>
					<ChevronLeft class="w-5 h-5" />
				</button>
				<div>
					<h2 class="text-xl font-bold flex items-center gap-2">
						<Sparkles class="w-5 h-5 text-sky-400" />
						{activeMemory.title}
					</h2>
					<p class="text-xs text-slate-300">{activeMemory.subtitle} ({activeMemory.photos.length} photos)</p>
				</div>
			</div>

			<button
				type="button"
				onclick={() => activeMemory = null}
				class="p-2 rounded-xl bg-white/10 hover:bg-white/20 text-white transition-colors cursor-pointer"
			>
				<X class="w-5 h-5" />
			</button>
		</div>

		<div class="flex-1 overflow-y-auto max-w-7xl w-full mx-auto min-h-0">
			{#if activeMemory.photos.length === 0}
				<div class="p-12 text-center text-white/70 bg-slate-900/60 rounded-3xl max-w-md mx-auto my-12 space-y-2 border border-slate-800">
					<ImageIcon class="w-10 h-10 mx-auto text-sky-400 opacity-80" />
					<h3 class="text-base font-bold text-white">No photos in this memory</h3>
					<p class="text-xs text-slate-400">Select photos in Gallery and click <strong>Add to Memory</strong> to include them here.</p>
				</div>
			{:else}
				<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
					{#each activeMemory.photos as photo, idx}
						<div
							role="button"
							tabindex="0"
							onclick={() => selectedPhotoIndex = idx}
							onkeydown={(e) => e.key === 'Enter' && (selectedPhotoIndex = idx)}
							class="aspect-4/3 rounded-2xl overflow-hidden bg-slate-900 relative group cursor-pointer border border-slate-800"
						>
							{#if photo.file_type === 'video' || (photo.mime_type && photo.mime_type.startsWith('video/'))}
								<div class="w-full h-full bg-slate-950 flex items-center justify-center relative">
									<video src={photo.url} class="w-full h-full object-cover opacity-80" preload="metadata"></video>
									<div class="w-8 h-8 rounded-full bg-white/20 backdrop-blur-md text-white flex items-center justify-center absolute">
										<Play class="w-4 h-4 fill-white ml-0.5" />
									</div>
								</div>
							{:else}
								<img
									src={photo.thumbnail_url}
									alt={photo.title}
									class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
									onerror={(e) => {
										const img = e.currentTarget as HTMLImageElement;
										img.onerror = null;
										img.src = '/empty-folder.png';
									}}
								/>
							{/if}
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
{/if}
