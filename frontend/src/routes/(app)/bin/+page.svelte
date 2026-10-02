<script lang="ts">
	import { onMount } from 'svelte';
	import { appState } from '$lib/state.svelte';
	import { apiFetch, getMediaUrl } from '$lib/api';
	import { notify, confirmDialog } from '$lib/notify.svelte';
	import { Trash2, ArchiveRestore, Image as ImageIcon, Play, Video, FileText } from 'lucide-svelte';

	interface DeletedPhoto {
		id: string;
		title: string;
		filename: string;
		file_type?: string;
		mime_type?: string;
		size: number;
		updated_at: string;
		url?: string;
		thumbnail_url?: string;
	}

	let binnedItems = $state<DeletedPhoto[]>([]);
	let isLoading = $state(true);

	async function fetchBinItems() {
		isLoading = true;
		try {
			const res = await apiFetch('/api/media?deleted=true');
			if (res.ok) {
				const data: DeletedPhoto[] = await res.json();
				binnedItems = data.map(item => ({
					...item,
					url: getMediaUrl(`/api/media/${item.id}/file`),
					thumbnail_url: getMediaUrl(`/api/media/${item.id}/thumbnail`)
				}));
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
		const target = binnedItems.find(i => i.id === id);
		const name = target ? (target.title || target.filename) : 'Item';

		const confirmed = await confirmDialog.ask({
			title: 'Restore Item to Gallery',
			message: `Restore "${name}" back to your photo gallery?`,
			confirmText: 'Yes, Restore',
			cancelText: 'Cancel',
			type: 'info'
		});
		if (!confirmed) return;

		try {
			const res = await apiFetch('/api/media/batch-restore', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ ids: [id] })
			});

			if (res.ok) {
				notify.success(`"${name}" restored to gallery!`);
				appState.refreshPhotos();
				fetchBinItems();
			} else {
				notify.error('Failed to restore item from bin.');
			}
		} catch (err) {
			console.error('Error restoring photo:', err);
			notify.error('Network error restoring photo.');
		}
	}

	async function restoreAll() {
		if (binnedItems.length === 0) return;
		const count = binnedItems.length;
		const confirmed = await confirmDialog.ask({
			title: 'Restore All Items',
			message: `Restore all ${count} items back to your photo gallery?`,
			confirmText: 'Yes, Restore All',
			cancelText: 'Cancel',
			type: 'info'
		});
		if (!confirmed) return;

		const ids = binnedItems.map((item) => item.id);
		try {
			const res = await apiFetch('/api/media/batch-restore', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ ids })
			});

			if (res.ok) {
				notify.success(`All ${count} items restored to gallery!`);
				appState.refreshPhotos();
				fetchBinItems();
			} else {
				notify.error('Failed to restore items.');
			}
		} catch (err) {
			console.error('Error restoring all items:', err);
			notify.error('Network error restoring items.');
		}
	}

	async function emptyBin() {
		if (binnedItems.length === 0) return;
		const count = binnedItems.length;
		const confirmed = await confirmDialog.ask({
			title: 'Empty Bin & Purge All Media',
			message: `Permanently purge all ${count} items from server storage? This action cannot be undone.`,
			confirmText: 'Yes, Purge All',
			cancelText: 'Cancel',
			type: 'danger'
		});
		if (!confirmed) return;

		const ids = binnedItems.map((item) => item.id);
		try {
			const res = await apiFetch('/api/media/batch-delete', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ ids })
			});

			if (res.ok) {
				notify.success(`Bin emptied. ${count} items permanently purged.`);
				appState.refreshPhotos();
				fetchBinItems();
			} else {
				notify.error('Failed to empty bin.');
			}
		} catch (err) {
			console.error('Error emptying bin:', err);
			notify.error('Network error purging bin.');
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
			<p class="text-xs text-slate-500 dark:text-slate-400">
				Items soft-deleted from gallery storage ({binnedItems.length} items)
			</p>
		</div>

		<div class="flex items-center gap-2.5">
			<button
				type="button"
				onclick={restoreAll}
				disabled={binnedItems.length === 0}
				class="px-3.5 py-2 rounded-xl bg-sky-50 dark:bg-sky-950/60 text-sky-600 dark:text-sky-400 hover:bg-sky-100 dark:hover:bg-sky-900 text-xs font-bold transition-all cursor-pointer border border-sky-200/60 dark:border-sky-800 disabled:opacity-50 flex items-center gap-1.5"
			>
				<ArchiveRestore class="w-3.5 h-3.5" />
				<span>Restore All</span>
			</button>

			<button
				type="button"
				onclick={emptyBin}
				disabled={binnedItems.length === 0}
				class="px-4 py-2 rounded-xl bg-rose-50 dark:bg-rose-950/60 text-rose-600 dark:text-rose-400 hover:bg-rose-100 dark:hover:bg-rose-900 text-xs font-bold transition-all cursor-pointer border border-rose-200/60 dark:border-rose-800 disabled:opacity-50 flex items-center gap-1.5"
			>
				<Trash2 class="w-3.5 h-3.5" />
				<span>Empty Bin</span>
			</button>
		</div>
	</div>

	<div class="flex-1 overflow-y-auto min-h-0">
		{#if !isLoading && binnedItems.length === 0}
			<div
				class="p-12 text-center rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm max-w-md mx-auto my-12 space-y-3"
			>
				<div
					class="w-12 h-12 rounded-2xl bg-rose-100 dark:bg-rose-950 text-rose-500 flex items-center justify-center mx-auto"
				>
					<Trash2 class="w-6 h-6" />
				</div>
				<h3 class="text-base font-bold text-slate-900 dark:text-white">
					Bin is empty
				</h3>
				<p class="text-xs text-slate-500 dark:text-slate-400">
					Deleted photos will appear here before permanent purging.
				</p>
			</div>
		{:else}
			<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
				{#each binnedItems as item}
					<div
						class="group relative rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:shadow-md transition-all p-3 flex flex-col justify-between gap-3"
					>
						<!-- Media Thumbnail Container -->
						<div class="aspect-4/3 w-full rounded-xl overflow-hidden bg-slate-100 dark:bg-slate-800 relative">
							{#if item.file_type === 'video' || (item.mime_type && item.mime_type.startsWith('video/'))}
								<div class="w-full h-full bg-slate-950 flex items-center justify-center relative">
									<video src={item.url} class="w-full h-full object-cover opacity-80" preload="metadata"></video>
									<div class="w-8 h-8 rounded-full bg-white/20 backdrop-blur-md text-white flex items-center justify-center absolute">
										<Play class="w-4 h-4 fill-white ml-0.5" />
									</div>
								</div>
							{:else}
								<img
									src={item.thumbnail_url}
									alt={item.title || item.filename}
									class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
									onerror={(e) => {
										const img = e.currentTarget as HTMLImageElement;
										img.onerror = null;
										if (item.url && img.src !== item.url) {
											img.src = item.url;
										} else {
											img.src = '/empty-folder.png';
										}
									}}
								/>
							{/if}
						</div>

						<!-- Details & Restore Action -->
						<div class="space-y-2">
							<div class="space-y-0.5">
								<p class="text-xs font-bold text-slate-800 dark:text-white truncate">
									{item.title || item.filename}
								</p>
								<p class="text-[11px] text-slate-400">
									Deleted {item.updated_at || 'Recently'} • {formatBytes(item.size)}
								</p>
							</div>

							<button
								type="button"
								onclick={() => restoreItem(item.id)}
								class="w-full py-1.5 px-3 rounded-xl bg-sky-50 dark:bg-sky-950/60 text-sky-600 dark:text-sky-400 hover:bg-sky-100 dark:hover:bg-sky-900 border border-sky-200 dark:border-sky-800 text-xs font-bold transition-all flex items-center justify-center gap-1.5 cursor-pointer"
								title="Restore to Gallery"
							>
								<ArchiveRestore class="w-3.5 h-3.5" />
								<span>Restore</span>
							</button>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>
