<script lang="ts">
	import { onMount } from 'svelte';
	import { appState } from '$lib/state.svelte';
	import { apiFetch, getMediaUrl } from '$lib/api';
	import { notify, confirmDialog } from '$lib/notify.svelte';
	import DropZone from '$lib/components/DropZone.svelte';
	import UploadButton from '$lib/components/UploadButton.svelte';
	import MediaViewer from '$lib/components/MediaViewer.svelte';
	import {
		Image as ImageIcon,
		Heart,
		Calendar,
		Download,
		X,
		ChevronLeft,
		ChevronRight,
		Trash2,
		Video,
		Lock,
		Play,
		ZoomIn,
		ZoomOut,
		RotateCcw,
		ShieldCheck,
		AlertCircle,
		ArrowLeft,
		Sliders,
		Scaling,
		CheckSquare,
		Square,
		Sparkles,
		Check,
		FolderPlus,
		Info
	} from 'lucide-svelte';

	interface PhotoItem {
		id: string;
		title: string;
		filename: string;
		object_key: string;
		mime_type: string;
		file_type: string;
		size: number;
		width: number;
		height: number;
		exif_model?: string;
		taken_at: string;
		is_favorite: boolean;
		is_deleted: boolean;
		locked_folder_id?: string;
		url?: string;
		thumbnail_url?: string;
	}

	interface LockedFolder {
		id: string;
		name: string;
		description: string;
	}

	interface Memory {
		id: string;
		title: string;
		description: string;
		items_count: number;
		cover_url?: string;
	}

	let photos = $state<PhotoItem[]>([]);
	let lockedFolders = $state<LockedFolder[]>([]);
	let isLoading = $state(true);
	let selectedIndex = $state<number | null>(null);
	let activeFilter = $state<'all' | 'favorites' | 'videos'>('all');
	let showDetailsModal = $state(false);
	let showZoomControls = $state(false);

	// Lightbox Zoom Controls
	let zoomScale = $state(100);

	// Move to Locked Folder Dialog State
	let showMoveToVaultModal = $state(false);
	let selectedVaultFolderId = $state('');
	let vaultPasscode = $state('');
	let vaultError = $state('');

	// Multi-Select State
	let isSelectionMode = $state(false);
	let selectedIds = $state<string[]>([]);

	// Add to Memory Modal State
	let showAddToMemoryModal = $state(false);
	let existingMemories = $state<Memory[]>([]);
	let selectedMemoryId = $state('new');
	let memoryTitle = $state('');
	let memoryDesc = $state('');
	let isSavingMemory = $state(false);

	function toggleSelectPhoto(id: string, e?: Event) {
		if (e) e.stopPropagation();
		if (selectedIds.includes(id)) {
			selectedIds = selectedIds.filter(i => i !== id);
		} else {
			selectedIds = [...selectedIds, id];
		}
	}

	function selectAllPhotos() {
		selectedIds = filteredPhotos.map(p => p.id);
	}

	function clearSelection() {
		selectedIds = [];
		isSelectionMode = false;
	}

	async function batchDeleteSelected() {
		if (selectedIds.length === 0) return;
		const count = selectedIds.length;
		const confirmed = await confirmDialog.ask({
			title: 'Move Selected Items to Bin',
			message: `Move ${count} selected items to the bin?`,
			confirmText: 'Yes, Move to Bin',
			cancelText: 'Cancel',
			type: 'warning'
		});
		if (!confirmed) return;

		try {
			for (const id of selectedIds) {
				await apiFetch(`/api/media/${id}`, {
					method: 'PUT',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ is_deleted: true })
				});
			}
			notify.success(`${count} items moved to bin.`);
			clearSelection();
			fetchPhotos();
		} catch (err) {
			console.error('Error batch deleting:', err);
			notify.error('Failed to move selected items to bin.');
		}
	}

	async function openAddToMemoryModal() {
		if (selectedIds.length === 0) return;
		try {
			const res = await apiFetch('/api/memories');
			if (res.ok) {
				existingMemories = await res.json();
			}
		} catch (err) {
			console.error('Error loading memories:', err);
		}
		selectedMemoryId = existingMemories.length > 0 ? existingMemories[0].id : 'new';
		memoryTitle = '';
		memoryDesc = '';
		showAddToMemoryModal = true;
	}

	async function saveToMemory() {
		if (selectedIds.length === 0) return;
		isSavingMemory = true;
		try {
			if (selectedMemoryId === 'new') {
				if (!memoryTitle.trim()) {
					notify.error('Please enter a Memory Title');
					isSavingMemory = false;
					return;
				}
				const res = await apiFetch('/api/memories', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({
						title: memoryTitle.trim(),
						description: memoryDesc.trim(),
						photo_ids: selectedIds
					})
				});
				if (res.ok) {
					notify.success(`Created Memory "${memoryTitle}" with ${selectedIds.length} photos!`);
					showAddToMemoryModal = false;
					clearSelection();
				} else {
					notify.error('Failed to create memory.');
				}
			} else {
				const res = await apiFetch(`/api/memories/${selectedMemoryId}/photos`, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ photo_ids: selectedIds })
				});
				if (res.ok) {
					const targetMem = existingMemories.find(m => m.id === selectedMemoryId);
					notify.success(`Added ${selectedIds.length} photos to Memory "${targetMem?.title || ''}"!`);
					showAddToMemoryModal = false;
					clearSelection();
				} else {
					notify.error('Failed to add photos to memory.');
				}
			}
		} catch (err) {
			console.error('Error saving memory:', err);
			notify.error('Network error saving memory.');
		} finally {
			isSavingMemory = false;
		}
	}

	async function fetchPhotos() {
		isLoading = true;
		try {
			const res = await apiFetch('/api/media?type=gallery&deleted=false');
			if (res.ok) {
				const data: PhotoItem[] = await res.json();
				photos = data.filter(p => !p.locked_folder_id).map(p => ({
					...p,
					url: getMediaUrl(`/api/media/${p.id}/file`),
					thumbnail_url: getMediaUrl(`/api/media/${p.id}/thumbnail`)
				}));
			}
		} catch (e) {
			console.warn('API error fetching photos:', e);
		} finally {
			isLoading = false;
		}
	}

	async function fetchLockedFolders() {
		try {
			const res = await apiFetch('/api/locked-folders');
			if (res.ok) {
				lockedFolders = await res.json();
			}
		} catch (e) {
			console.warn('API error fetching locked folders:', e);
		}
	}

	onMount(() => {
		fetchPhotos();
		fetchLockedFolders();
	});

	const filteredPhotos = $derived(
		activeFilter === 'favorites'
			? photos.filter(p => p.is_favorite)
			: activeFilter === 'videos'
			? photos.filter(p => p.file_type === 'video' || p.mime_type.startsWith('video/'))
			: photos
	);

	const selectedPhoto = $derived(
		selectedIndex !== null && selectedIndex >= 0 && selectedIndex < filteredPhotos.length
			? filteredPhotos[selectedIndex]
			: null
	);

	function openPhoto(index: number) {
		if (isSelectionMode) {
			const item = filteredPhotos[index];
			if (item) toggleSelectPhoto(item.id);
			return;
		}
		selectedIndex = index;
		zoomScale = 100;
		showDetailsModal = false;
		showMoveToVaultModal = false;
		showZoomControls = false;
	}

	function closePhoto() {
		selectedIndex = null;
		zoomScale = 100;
		showDetailsModal = false;
		showMoveToVaultModal = false;
		showZoomControls = false;
	}

	function prevPhoto() {
		if (selectedIndex !== null && filteredPhotos.length > 0) {
			selectedIndex = (selectedIndex - 1 + filteredPhotos.length) % filteredPhotos.length;
			zoomScale = 100;
			showDetailsModal = false;
			showMoveToVaultModal = false;
		}
	}

	function nextPhoto() {
		if (selectedIndex !== null && filteredPhotos.length > 0) {
			selectedIndex = (selectedIndex + 1) % filteredPhotos.length;
			zoomScale = 100;
			showDetailsModal = false;
			showMoveToVaultModal = false;
		}
	}

	function zoomIn() {
		if (zoomScale < 400) zoomScale += 25;
	}

	function zoomOut() {
		if (zoomScale > 40) zoomScale -= 25;
	}

	function resetZoom() {
		zoomScale = 100;
	}

	async function toggleFavorite(photo: PhotoItem, e?: Event) {
		if (e) e.stopPropagation();
		const updatedState = !photo.is_favorite;
		photo.is_favorite = updatedState;
		try {
			await apiFetch(`/api/media/${photo.id}`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ is_favorite: updatedState })
			});
		} catch (err) {
			console.error('Error toggling favorite:', err);
		}
	}

	async function deletePhoto(photo: PhotoItem, e?: Event) {
		if (e) e.stopPropagation();
		const confirmed = await confirmDialog.ask({
			title: 'Move Item to Bin',
			message: `Move "${photo.title}" to bin? You can restore it anytime from the Bin.`,
			confirmText: 'Yes, Move to Bin',
			cancelText: 'Cancel',
			type: 'warning'
		});
		if (!confirmed) return;

		try {
			const res = await apiFetch(`/api/media/${photo.id}`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ is_deleted: true })
			});
			if (res.ok) {
				notify.success(`"${photo.title}" moved to bin.`);
			}
			if (selectedIndex !== null) closePhoto();
			fetchPhotos();
		} catch (err) {
			console.error('Error deleting photo:', err);
			notify.error('Failed to move item to bin.');
		}
	}

	function openMoveToVaultDialog(photo?: PhotoItem, e?: Event) {
		if (e) e.stopPropagation();
		if (photo) {
			// Find photo index or set selectedPhoto
			const idx = filteredPhotos.findIndex(p => p.id === photo.id);
			if (idx !== -1) selectedIndex = idx;
		}
		if (lockedFolders.length === 0) {
			alert('No locked folders available. Please create a locked folder first in Locked Vault!');
			return;
		}
		selectedVaultFolderId = lockedFolders[0].id;
		vaultPasscode = '';
		vaultError = '';
		showMoveToVaultModal = true;
	}

	async function handleMoveToLockedFolder(e: Event) {
		e.preventDefault();
		if (!selectedPhoto || !selectedVaultFolderId || !vaultPasscode) return;
		vaultError = '';

		try {
			const verifyRes = await apiFetch(`/api/locked-folders/${selectedVaultFolderId}/verify`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ passcode: vaultPasscode })
			});

			if (!verifyRes.ok) {
				vaultError = 'Incorrect passcode for locked folder!';
				return;
			}

			const updateRes = await apiFetch(`/api/media/${selectedPhoto.id}`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ locked_folder_id: selectedVaultFolderId })
			});

			if (updateRes.ok) {
				showMoveToVaultModal = false;
				closePhoto();
				fetchPhotos();
			} else {
				vaultError = 'Failed to move photo to locked folder.';
			}
		} catch (err) {
			console.error('Error moving to locked folder:', err);
			vaultError = 'Network error.';
		}
	}

	function formatBytes(bytes: number): string {
		if (!bytes) return '0 B';
		const k = 1024;
		const sizes = ['B', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
	}

	function handleKeydown(e: KeyboardEvent) {
		if (selectedIndex === null) return;
		if (e.key === 'ArrowLeft') prevPhoto();
		if (e.key === 'ArrowRight') nextPhoto();
		if (e.key === 'Escape') closePhoto();
	}

</script>

<svelte:window onkeydown={handleKeydown} />

<DropZone class="h-full flex flex-col gap-6 animate-fade-in" onUploadComplete={fetchPhotos}>

	<!-- Header Bar -->
	<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 shrink-0">
		<div>
			<h1 class="text-2xl font-bold text-slate-900 dark:text-white">Gallery</h1>
			<p class="text-xs text-slate-500 dark:text-slate-400">Timeline view of your private photo & video storage</p>
		</div>

		<!-- Filter Tabs -->
		<div class="flex items-center gap-1.5 p-1 rounded-xl bg-slate-200/60 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-xs font-semibold">
			<button
				type="button"
				onclick={() => { activeFilter = 'all'; selectedIndex = null; }}
				class={`px-3 py-1.5 rounded-lg transition-all cursor-pointer ${
					activeFilter === 'all'
						? 'bg-sky-400 text-white shadow-sm'
						: 'text-slate-600 dark:text-slate-300 hover:text-slate-900 dark:hover:text-white'
				}`}
			>
				All Items ({photos.length})
			</button>

			<button
				type="button"
				onclick={() => { activeFilter = 'favorites'; selectedIndex = null; }}
				class={`px-3 py-1.5 rounded-lg transition-all cursor-pointer ${
					activeFilter === 'favorites'
						? 'bg-sky-400 text-white shadow-sm'
						: 'text-slate-600 dark:text-slate-300 hover:text-slate-900 dark:hover:text-white'
				}`}
			>
				<Heart class="w-3.5 h-3.5 inline mr-1" />
				Favorites ({photos.filter(p => p.is_favorite).length})
			</button>

			<button
				type="button"
				onclick={() => { activeFilter = 'videos'; selectedIndex = null; }}
				class={`px-3 py-1.5 rounded-lg transition-all cursor-pointer ${
					activeFilter === 'videos'
						? 'bg-sky-400 text-white shadow-sm'
						: 'text-slate-600 dark:text-slate-300 hover:text-slate-900 dark:hover:text-white'
				}`}
			>
				<Video class="w-3.5 h-3.5 inline mr-1" />
				Videos ({photos.filter(p => p.file_type === 'video' || p.mime_type.startsWith('video/')).length})
			</button>
		</div>

		<!-- Action Controls -->
		<div class="flex items-center gap-2">
			<!-- Select Toggle Button -->
			<button
				type="button"
				onclick={() => { isSelectionMode = !isSelectionMode; if (!isSelectionMode) selectedIds = []; }}
				class={`px-3.5 py-1.5 rounded-xl text-xs font-semibold border transition-all cursor-pointer flex items-center gap-1.5 ${
					isSelectionMode
						? 'bg-sky-500 text-white border-sky-400 shadow-sm'
						: 'bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-200 border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700'
				}`}
			>
				<CheckSquare class="w-3.5 h-3.5" />
				<span>{isSelectionMode ? 'Cancel Select' : 'Select'}</span>
			</button>

			<!-- Reusable Upload Button -->
			<UploadButton label="Upload" variant="primary" />
		</div>
	</div>

	<div class="flex-1 overflow-y-auto min-h-0">
		<!-- Empty State when no photos exist -->
		{#if !isLoading && filteredPhotos.length === 0}
			<div class="p-12 text-center rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm max-w-md mx-auto my-12 space-y-3">
				<img src="/empty-folder.png" alt="Empty" class="w-20 h-20 mx-auto mb-4 opacity-60 dark:opacity-40 select-none pointer-events-none drop-shadow-sm" />
				<h3 class="text-base font-bold text-slate-900 dark:text-white">No media in gallery</h3>
				<p class="text-xs text-slate-500 dark:text-slate-400">Drag & drop files anywhere or click <strong class="text-sky-500">Upload</strong> in header!</p>
			</div>
		{:else}
			<!-- Media Grid -->
			<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-3 gap-6">
				{#each filteredPhotos as photo, index}
					<div
						role="button"
						tabindex="0"
						onclick={() => openPhoto(index)}
						onkeydown={(e) => e.key === 'Enter' && openPhoto(index)}
						class="group relative rounded-2xl overflow-hidden bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:shadow-md transition-all duration-300 cursor-pointer"
					>
						<div class="aspect-4/3 overflow-hidden bg-slate-100 dark:bg-slate-800 relative">
							{#if photo.file_type === 'video' || photo.mime_type.startsWith('video/')}
								<div class="w-full h-full bg-slate-950 flex items-center justify-center relative">
									<video src={photo.url} class="w-full h-full object-cover opacity-80" preload="metadata"></video>
									<div class="w-10 h-10 rounded-full bg-white/20 backdrop-blur-md text-white flex items-center justify-center absolute">
										<Play class="w-5 h-5 fill-white ml-0.5" />
									</div>
								</div>
							{:else}
								<img
									src={photo.thumbnail_url}
									alt={photo.title}
									class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
									onerror={(e) => {
										const img = e.currentTarget as HTMLImageElement;
										img.onerror = null;
										if (photo.url && img.src !== photo.url) {
											img.src = photo.url;
										} else {
											img.src = '/empty-folder.png';
										}
									}}
								/>
							{/if}

							<div class="absolute inset-0 bg-linear-to-t from-slate-900/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity flex items-end justify-between p-4">
								<div class="text-white max-w-[65%]">
									<p class="text-xs font-bold truncate">{photo.title}</p>
									<p class="text-[10px] text-slate-300 truncate">{photo.taken_at || 'Recently'} • {formatBytes(photo.size)}</p>
								</div>

								<div class="flex items-center gap-1.5">
									<button
										type="button"
										onclick={(e) => toggleFavorite(photo, e)}
										class="p-1.5 rounded-full bg-white/20 hover:bg-white/40 text-white backdrop-blur-md transition-colors"
										title="Favorite"
									>
										<Heart class={`w-4 h-4 ${photo.is_favorite ? 'fill-rose-500 text-rose-500' : 'text-white'}`} />
									</button>
									<button
										type="button"
										onclick={(e) => openMoveToVaultDialog(photo, e)}
										class="p-1.5 rounded-full bg-white/20 hover:bg-sky-500/80 text-white backdrop-blur-md transition-colors"
										title="Move to Locked Vault"
									>
										<Lock class="w-4 h-4" />
									</button>
									<button
										type="button"
										onclick={(e) => deletePhoto(photo, e)}
										class="p-1.5 rounded-full bg-white/20 hover:bg-rose-500/80 text-white backdrop-blur-md transition-colors"
										title="Move to Bin"
									>
										<Trash2 class="w-4 h-4" />
									</button>
								</div>
							</div>

							<!-- Selection Overlay Checkbox -->
							{#if isSelectionMode}
								<div
									role="button"
									tabindex="0"
									onclick={(e) => toggleSelectPhoto(photo.id, e)}
									onkeydown={(e) => e.key === 'Enter' && toggleSelectPhoto(photo.id, e)}
									class={`absolute top-3 left-3 z-10 p-1.5 rounded-xl backdrop-blur-md transition-all cursor-pointer ${
										selectedIds.includes(photo.id)
											? 'bg-sky-500 text-white shadow-lg scale-105'
											: 'bg-slate-950/40 text-white/70 hover:bg-slate-950/70'
									}`}
								>
									{#if selectedIds.includes(photo.id)}
										<CheckSquare class="w-5 h-5 text-white" />
									{:else}
										<Square class="w-5 h-5" />
									{/if}
								</div>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</DropZone>

<MediaViewer
	photos={filteredPhotos}
	selectedIndex={selectedIndex}
	onclose={closePhoto}
	onnavigate={(idx) => selectedIndex = idx}
	ontoggleFavorite={(photo) => toggleFavorite(photo)}
	ondelete={(photo) => deletePhoto(photo)}
	onmoveToVault={(photo) => openMoveToVaultDialog(photo)}
/>
	<!-- MOVE TO LOCKED FOLDER DIALOG (Z-INDEX 10001) -->
	{#if showMoveToVaultModal && selectedPhoto}
		<div class="fixed inset-0 z-[10001] bg-slate-950/70 backdrop-blur-sm flex items-center justify-center p-4">
			<div class="bg-white dark:bg-slate-900 rounded-3xl p-6 md:p-8 max-w-md w-full border border-slate-200 dark:border-slate-800 shadow-2xl animate-fade-in relative">
				<button
					type="button"
					onclick={() => showMoveToVaultModal = false}
					class="absolute top-5 right-5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1"
				>
					<X class="w-5 h-5" />
				</button>

				<div class="flex items-center gap-2 mb-1">
					<Lock class="w-5 h-5 text-sky-500" />
					<h3 class="text-lg font-bold text-slate-900 dark:text-white">Move to Locked Folder</h3>
				</div>
				<p class="text-xs text-slate-500 dark:text-slate-400 mb-6">
					Move "<strong class="text-slate-800 dark:text-slate-200">{selectedPhoto.title}</strong>" into a passcode-protected folder
				</p>

				{#if vaultError}
					<div class="p-3 mb-4 rounded-xl bg-rose-50 dark:bg-rose-950/60 text-rose-600 dark:text-rose-400 text-xs font-semibold flex items-center gap-2 border border-rose-200 dark:border-rose-800">
						<AlertCircle class="w-4 h-4 shrink-0" />
						<span>{vaultError}</span>
					</div>
				{/if}

				<form onsubmit={handleMoveToLockedFolder} class="space-y-4">
					<div class="space-y-1.5">
						<label for="target-folder" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Select Target Locked Folder</label>
						<select
							id="target-folder"
							bind:value={selectedVaultFolderId}
							class="w-full h-10 px-3 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400"
						>
							{#each lockedFolders as folder}
								<option value={folder.id}>{folder.name} ({folder.description || 'Locked'})</option>
							{/each}
						</select>
					</div>

					<div class="space-y-1.5">
						<label for="vault-passcode" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Folder 4-Digit Passcode</label>
						<input
							id="vault-passcode"
							type="password"
							maxLength={4}
							bind:value={vaultPasscode}
							placeholder="••••"
							class="w-full h-11 text-center font-mono text-lg tracking-widest rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400"
							required
						/>
					</div>

					<div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100 dark:border-slate-800">
						<button
							type="button"
							onclick={() => showMoveToVaultModal = false}
							class="px-4 py-2 rounded-xl text-xs font-semibold text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-5 py-2 rounded-xl bg-sky-400 hover:bg-sky-500 text-white text-xs font-bold shadow-sm shadow-sky-300/50 flex items-center gap-1.5"
						>
							<ShieldCheck class="w-4 h-4" /> Move to Vault
						</button>
					</div>
				</form>
			</div>
		</div>
	{/if}

	<!-- Floating Multi-Selection Action Bar -->
	{#if isSelectionMode && selectedIds.length > 0}
		<div class="fixed bottom-6 left-1/2 -translate-x-1/2 z-[1000] bg-white/95 dark:bg-slate-900/95 text-slate-900 dark:text-white backdrop-blur-xl border border-slate-200 dark:border-slate-700/80 rounded-2xl px-5 py-3 shadow-2xl flex items-center gap-4 animate-fade-in max-w-max">
			<span class="text-xs font-bold text-sky-500 dark:text-sky-400">{selectedIds.length} Selected</span>

			<div class="h-4 w-px bg-slate-200 dark:bg-slate-700"></div>

			<!-- Add to Memory Button -->
			<button
				type="button"
				onclick={openAddToMemoryModal}
				class="px-3.5 py-1.5 rounded-xl bg-sky-500 hover:bg-sky-600 text-white text-xs font-bold transition-all flex items-center gap-1.5 cursor-pointer shadow-sm shadow-sky-500/30"
			>
				<Sparkles class="w-3.5 h-3.5" />
				<span>Add to Memory</span>
			</button>

			<!-- Move to Bin Button -->
			<button
				type="button"
				onclick={batchDeleteSelected}
				class="px-3.5 py-1.5 rounded-xl bg-rose-50 dark:bg-rose-500/20 hover:bg-rose-500 text-rose-600 dark:text-rose-300 hover:text-white border border-rose-200 dark:border-rose-500/30 text-xs font-semibold transition-all flex items-center gap-1.5 cursor-pointer"
			>
				<Trash2 class="w-3.5 h-3.5" />
				<span>Move to Bin</span>
			</button>

			<div class="h-4 w-px bg-slate-200 dark:bg-slate-700"></div>

			<button
				type="button"
				onclick={selectAllPhotos}
				class="text-xs text-slate-600 dark:text-slate-300 hover:text-slate-900 dark:hover:text-white font-medium cursor-pointer"
			>
				Select All
			</button>

			<button
				type="button"
				onclick={() => selectedIds = []}
				class="text-xs text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 cursor-pointer"
			>
				Deselect
			</button>
		</div>
	{/if}

	<!-- Add to Memory Dialog Modal -->
	{#if showAddToMemoryModal}
		<div class="fixed inset-0 z-[10000] bg-slate-950/80 backdrop-blur-md flex items-center justify-center p-4 animate-fade-in w-full h-full">
			<div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-2xl max-w-md w-full space-y-5 animate-fade-in">
				<div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-4">
					<div class="flex items-center gap-2.5">
						<div class="w-10 h-10 rounded-2xl bg-sky-100 dark:bg-sky-950 text-sky-500 flex items-center justify-center">
							<Sparkles class="w-5 h-5" />
						</div>
						<div>
							<h3 class="text-base font-bold text-slate-900 dark:text-white">Add {selectedIds.length} Photos to Memory</h3>
							<p class="text-xs text-slate-500 dark:text-slate-400">Save to an existing memory or create a new collection</p>
						</div>
					</div>
					<button type="button" onclick={() => showAddToMemoryModal = false} class="p-1 rounded-xl text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 cursor-pointer">
						<X class="w-5 h-5" />
					</button>
				</div>

				<div class="space-y-4">
					<div class="space-y-1.5">
						<label for="memory-target-select" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Target Memory</label>
						<select
							id="memory-target-select"
							bind:value={selectedMemoryId}
							class="w-full h-10 px-3 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400"
						>
							<option value="new">✨ + Create New Memory Collection</option>
							{#each existingMemories as mem}
								<option value={mem.id}>📖 {mem.title} ({mem.items_count} items)</option>
							{/each}
						</select>
					</div>

					{#if selectedMemoryId === 'new'}
						<div class="space-y-1.5">
							<label for="memory-title-input" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Memory Title *</label>
							<input
								id="memory-title-input"
								type="text"
								bind:value={memoryTitle}
								placeholder="e.g. Summer Vacation 2026, Birthday Trip"
								class="w-full h-10 px-3.5 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400"
							/>
						</div>

						<div class="space-y-1.5">
							<label for="memory-desc-input" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Description (Optional)</label>
							<textarea
								id="memory-desc-input"
								bind:value={memoryDesc}
								rows={2}
								placeholder="Add details or notes about this memory..."
								class="w-full p-3 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400 resize-none"
							></textarea>
						</div>
					{/if}
				</div>

				<div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-100 dark:border-slate-800">
					<button
						type="button"
						onclick={() => showAddToMemoryModal = false}
						class="px-4 py-2 rounded-xl text-xs font-semibold text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 cursor-pointer"
					>
						Cancel
					</button>
					<button
						type="button"
						onclick={saveToMemory}
						disabled={isSavingMemory}
						class="px-5 py-2 rounded-xl bg-sky-400 hover:bg-sky-500 text-white text-xs font-bold shadow-sm shadow-sky-300/50 flex items-center gap-1.5 disabled:opacity-50 cursor-pointer"
					>
						<Sparkles class="w-4 h-4" />
						<span>{selectedMemoryId === 'new' ? 'Create & Add Photos' : 'Add Photos'}</span>
					</button>
				</div>
			</div>
		</div>
	{/if}
