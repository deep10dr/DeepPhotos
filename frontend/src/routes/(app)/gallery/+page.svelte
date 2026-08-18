<script lang="ts">
	import { onMount } from 'svelte';
	import { appState } from '$lib/state.svelte';
	import { apiFetch, getMediaUrl } from '$lib/api';
	import {
		Image as ImageIcon,
		Heart,
		Calendar,
		Download,
		Trash2,
		ChevronLeft,
		ChevronRight,
		X,
		Info,
		Video,
		FileText,
		Lock,
		Play,
		ZoomIn,
		ZoomOut,
		RotateCcw,
		ShieldCheck,
		AlertCircle,
		ArrowLeft,
		Sliders
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

	async function fetchPhotos() {
		isLoading = true;
		try {
			const res = await apiFetch('/api/photos?deleted=false');
			if (res.ok) {
				const data: PhotoItem[] = await res.json();
				photos = data.filter(p => !p.locked_folder_id).map(p => ({
					...p,
					url: getMediaUrl(`/api/photos/${p.id}/file`),
					thumbnail_url: getMediaUrl(`/api/photos/${p.id}/thumbnail`)
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
			await apiFetch(`/api/photos/${photo.id}`, {
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
		if (!confirm(`Move "${photo.title}" to bin?`)) return;

		try {
			await apiFetch(`/api/photos/${photo.id}`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ is_deleted: true })
			});
			if (selectedIndex !== null) closePhoto();
			fetchPhotos();
		} catch (err) {
			console.error('Error deleting photo:', err);
		}
	}

	function openMoveToVaultDialog() {
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

			const updateRes = await apiFetch(`/api/photos/${selectedPhoto.id}`, {
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

	// Drag & Drop — tracks nested enter/leave to prevent flicker
	let isDraggingExternal = $state(false);
	let isUploadingExternal = $state(false);
	let dragCounter = 0;

	function handleDragEnter(e: DragEvent) {
		e.preventDefault();
		dragCounter++;
		if (dragCounter === 1) {
			isDraggingExternal = true;
		}
	}

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		if (e.dataTransfer) {
			e.dataTransfer.dropEffect = 'copy';
		}
	}

	function handleDragLeave(e: DragEvent) {
		e.preventDefault();
		dragCounter--;
		if (dragCounter <= 0) {
			dragCounter = 0;
			isDraggingExternal = false;
		}
	}

	async function handleDrop(e: DragEvent) {
		e.preventDefault();
		dragCounter = 0;
		isDraggingExternal = false;

		const dt = e.dataTransfer;
		if (!dt) return;

		// 1. LOCAL FILES — user dragged files from their file system
		if (dt.files && dt.files.length > 0) {
			isUploadingExternal = true;
			try {
				const formData = new FormData();
				for (const file of Array.from(dt.files)) {
					formData.append('files', file);
				}
				const res = await apiFetch('/api/photos/upload', {
					method: 'POST',
					body: formData
				});
				if (res.ok) {
					await fetchPhotos();
				} else {
					console.error('Upload failed:', await res.text());
				}
			} catch (err) {
				console.error('Drop upload error:', err);
			} finally {
				isUploadingExternal = false;
			}
			return;
		}

		// 2. EXTERNAL URL — user dragged an image from another browser tab
		let urlToUpload = '';

		const html = dt.getData('text/html');
		if (html) {
			const match = html.match(/src="([^"]+)"/);
			if (match && match[1] && match[1].startsWith('http')) {
				urlToUpload = match[1];
			}
		}

		if (!urlToUpload) {
			const uriList = dt.getData('text/uri-list');
			if (uriList && uriList.startsWith('http')) {
				urlToUpload = uriList.split('\n')[0].trim();
			}
		}

		if (!urlToUpload) {
			const plain = dt.getData('text/plain');
			if (plain && plain.startsWith('http')) {
				urlToUpload = plain.trim();
			}
		}

		if (urlToUpload) {
			isUploadingExternal = true;
			try {
				const res = await apiFetch('/api/photos/upload-url', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ url: urlToUpload })
				});
				if (res.ok) {
					await fetchPhotos();
				} else {
					alert('Failed to import image from URL.');
				}
			} catch (err) {
				console.error(err);
				alert('Network error during import.');
			} finally {
				isUploadingExternal = false;
			}
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div role="region" aria-label="Gallery Dropzone" class="h-full flex flex-col gap-6 animate-fade-in" ondragenter={handleDragEnter} ondragover={handleDragOver} ondragleave={handleDragLeave} ondrop={handleDrop}>

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

							<div class="absolute inset-0 bg-gradient-to-t from-slate-900/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity flex items-end justify-between p-4">
								<div class="text-white max-w-[65%]">
									<p class="text-xs font-bold truncate">{photo.title}</p>
									<p class="text-[10px] text-slate-300 truncate">{photo.taken_at || 'Recently'} • {formatBytes(photo.size)}</p>
								</div>

								<div class="flex items-center gap-1.5">
									<button
										type="button"
										onclick={(e) => toggleFavorite(photo, e)}
										class="p-1.5 rounded-full bg-white/20 hover:bg-white/40 text-white backdrop-blur-md transition-colors"
									>
										<Heart class={`w-4 h-4 ${photo.is_favorite ? 'fill-rose-500 text-rose-500' : 'text-white'}`} />
									</button>
									<button
										type="button"
										onclick={(e) => deletePhoto(photo, e)}
										class="p-1.5 rounded-full bg-white/20 hover:bg-rose-500/80 text-white backdrop-blur-md transition-colors"
									>
										<Trash2 class="w-4 h-4" />
									</button>
								</div>
							</div>

							{#if photo.is_favorite}
								<div class="absolute top-3 right-3 p-1.5 rounded-full bg-white/80 dark:bg-slate-900/80 backdrop-blur-md text-rose-500 shadow-sm">
									<Heart class="w-3.5 h-3.5 fill-rose-500" />
								</div>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>

	{#if selectedPhoto && selectedIndex !== null}
		<!-- GOOGLE PHOTOS FULL-SCREEN THEATER MODE (Z-INDEX 100 PREVENTS ANY OVERFLOW OR CLASHING) -->
		<div class="fixed inset-0 z-[100] w-screen h-screen bg-black/95 backdrop-blur-xl flex flex-col select-none overflow-hidden animate-fade-in">

			<!-- TOP GOOGLE PHOTOS NAVIGATION OVERLAY TOOLBAR (Z-INDEX 110) -->
			<div class="absolute top-0 left-0 right-0 z-[110] p-4 flex items-center justify-between pointer-events-none">

				<!-- Left: Back Arrow, Title & Counter -->
				<div class="flex items-center gap-3 pointer-events-auto bg-black/40 backdrop-blur-md pl-2 pr-4 py-2 rounded-full border border-white/10">
					<button
						type="button"
						onclick={closePhoto}
						title="Back to Gallery (Esc)"
						class="p-2 rounded-full text-slate-200 hover:text-white hover:bg-white/10 transition-colors cursor-pointer"
					>
						<ArrowLeft class="w-5 h-5" />
					</button>
					<div>
						<h3 class="text-sm font-bold text-white flex items-center gap-2">
							<span>{selectedPhoto.title}</span>
							<span class="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-white/15 text-slate-200">
								{selectedIndex + 1} of {filteredPhotos.length}
							</span>
						</h3>
						<p class="text-[11px] text-slate-400">{selectedPhoto.taken_at || 'Recently'} • {formatBytes(selectedPhoto.size)}</p>
					</div>
				</div>

				<!-- Right Action Buttons: Zoom, Favorite, Move to Vault, Details Info, Download, Trash, Close -->
				<div class="flex items-center gap-1.5 sm:gap-2 pointer-events-auto bg-black/40 backdrop-blur-md px-2 py-1.5 rounded-full border border-white/10">

					{#if selectedPhoto.file_type !== 'video' && !selectedPhoto.mime_type.startsWith('video/')}
						<button
							type="button"
							onclick={() => showZoomControls = !showZoomControls}
							class={`p-2 rounded-full transition-all cursor-pointer ${
								showZoomControls
									? 'bg-sky-400 text-white'
									: 'text-slate-200 hover:bg-white/10'
							}`}
							title="Zoom Slider"
						>
							<Sliders class="w-5 h-5" />
						</button>
					{/if}

					<button
						type="button"
						onclick={() => toggleFavorite(selectedPhoto)}
						class={`p-2 rounded-full transition-all cursor-pointer ${
							selectedPhoto.is_favorite
								? 'text-rose-500 bg-white/15'
								: 'text-slate-200 hover:bg-white/10'
						}`}
						title="Favorite"
					>
						<Heart class={`w-5 h-5 ${selectedPhoto.is_favorite ? 'fill-rose-500 text-rose-500' : ''}`} />
					</button>

					<button
						type="button"
						onclick={openMoveToVaultDialog}
						class="p-2 rounded-full text-slate-200 hover:text-sky-300 hover:bg-white/10 transition-colors cursor-pointer"
						title="Move to Locked Vault"
					>
						<Lock class="w-5 h-5" />
					</button>

					<button
						type="button"
						onclick={() => showDetailsModal = !showDetailsModal}
						class={`p-2 rounded-full transition-all cursor-pointer ${
							showDetailsModal
								? 'bg-sky-400 text-white'
								: 'text-slate-200 hover:bg-white/10'
						}`}
						title="Info & Technical Details"
					>
						<Info class="w-5 h-5" />
					</button>

					<a
						href={selectedPhoto.url}
						download={selectedPhoto.filename}
						target="_blank"
						class="p-2 rounded-full text-slate-200 hover:text-white hover:bg-white/10 transition-colors cursor-pointer"
						title="Download Original"
					>
						<Download class="w-5 h-5" />
					</a>

					<button
						type="button"
						onclick={() => deletePhoto(selectedPhoto)}
						class="p-2 rounded-full text-slate-200 hover:text-rose-400 hover:bg-white/10 transition-colors cursor-pointer"
						title="Move to Trash"
					>
						<Trash2 class="w-5 h-5" />
					</button>

					<button
						type="button"
						onclick={closePhoto}
						title="Close (Esc)"
						class="p-2 rounded-full text-slate-200 hover:text-white hover:bg-white/10 transition-colors cursor-pointer ml-1"
					>
						<X class="w-6 h-6" />
					</button>

				</div>

			</div>

			<!-- FLOATING PREVIOUS ARROW (Z-INDEX 110) -->
			<button
				type="button"
				onclick={prevPhoto}
				title="Previous (Left Arrow)"
				class="absolute left-4 sm:left-8 top-1/2 -translate-y-1/2 z-[110] w-12 h-12 rounded-full bg-slate-900/60 hover:bg-white/20 text-white flex items-center justify-center transition-all cursor-pointer backdrop-blur-md group"
			>
				<ChevronLeft class="w-7 h-7 transition-transform group-hover:-translate-x-0.5" />
			</button>

			<!-- FLOATING NEXT ARROW (Z-INDEX 110) -->
			<button
				type="button"
				onclick={nextPhoto}
				title="Next (Right Arrow)"
				class="absolute right-4 sm:right-8 top-1/2 -translate-y-1/2 z-[110] w-12 h-12 rounded-full bg-slate-900/60 hover:bg-white/20 text-white flex items-center justify-center transition-all cursor-pointer backdrop-blur-md group"
			>
				<ChevronRight class="w-7 h-7 transition-transform group-hover:translate-x-0.5" />
			</button>

			<!-- GOOGLE PHOTOS CENTER THEATER DISPLAY CANVAS -->
			<div class="flex-1 w-full h-full relative flex items-center justify-center p-0 overflow-hidden">
				{#if selectedPhoto.file_type === 'video' || selectedPhoto.mime_type.startsWith('video/')}
					<video src={selectedPhoto.url} controls autoplay class="w-full h-full object-contain"></video>
				{:else}
					<img
						src={selectedPhoto.url}
						alt={selectedPhoto.title}
						style={`transform: scale(${zoomScale / 100}); transition: transform 0.2s cubic-bezier(0.2, 0, 0, 1);`}
						class="w-full h-full object-contain select-none cursor-grab active:cursor-grabbing"
					/>
				{/if}
			</div>

			<!-- FLOATING BOTTOM ZOOM BAR (Z-INDEX 110) -->
			{#if showZoomControls && selectedPhoto.file_type !== 'video' && !selectedPhoto.mime_type.startsWith('video/')}
				<div class="absolute bottom-6 left-1/2 -translate-x-1/2 z-[110] p-3 px-6 bg-slate-900/90 backdrop-blur-xl border border-slate-800 rounded-2xl flex items-center justify-center gap-4 shadow-2xl animate-fade-in">
					<div class="flex items-center gap-3">
						<button
							type="button"
							onclick={zoomOut}
							title="Zoom Out (-25%)"
							class="p-1.5 rounded-lg bg-white/10 hover:bg-white/20 text-white transition-colors cursor-pointer"
						>
							<ZoomOut class="w-4 h-4" />
						</button>

						<input
							type="range"
							min="40"
							max="400"
							step="5"
							bind:value={zoomScale}
							class="w-40 md:w-64 h-2 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-sky-400"
							title={`Current Zoom: ${zoomScale}%`}
						/>

						<button
							type="button"
							onclick={zoomIn}
							title="Zoom In (+25%)"
							class="p-1.5 rounded-lg bg-white/10 hover:bg-white/20 text-white transition-colors cursor-pointer"
						>
							<ZoomIn class="w-4 h-4" />
						</button>

						<span class="text-xs font-mono font-bold text-sky-400 w-12 text-center">{zoomScale}%</span>

						<button
							type="button"
							onclick={resetZoom}
							title="Reset Zoom (100%)"
							class="px-3 py-1 rounded-lg bg-white/10 hover:bg-white/20 text-xs text-slate-200 font-semibold flex items-center gap-1 transition-colors cursor-pointer"
						>
							<RotateCcw class="w-3.5 h-3.5" />
							<span>Reset</span>
						</button>
					</div>
				</div>
			{/if}

			<!-- GOOGLE PHOTOS RIGHT SLIDE-OVER INFO DRAWER (Z-INDEX 120) -->
			{#if showDetailsModal}
				<div class="fixed top-0 right-0 bottom-0 w-80 sm:w-96 bg-slate-900/95 backdrop-blur-2xl z-[120] p-6 border-l border-slate-800 shadow-2xl overflow-y-auto space-y-6 animate-fade-in text-white">
					<div class="flex items-center justify-between border-b border-slate-800 pb-4">
						<h4 class="font-bold text-base flex items-center gap-2">
							<Info class="w-5 h-5 text-sky-400" />
							Info & Details
						</h4>
						<button
							type="button"
							onclick={() => showDetailsModal = false}
							class="p-1 text-slate-400 hover:text-white rounded-full hover:bg-white/10 transition-colors"
						>
							<X class="w-5 h-5" />
						</button>
					</div>

					<div class="space-y-4 text-xs">
						<div class="space-y-1">
							<p class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">File Name</p>
							<p class="font-bold text-slate-100 text-sm break-all">{selectedPhoto.filename}</p>
						</div>

						<div class="space-y-1">
							<p class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">MIME & Category</p>
							<p class="font-semibold text-slate-200">{selectedPhoto.mime_type} ({selectedPhoto.file_type})</p>
						</div>

						<div class="space-y-1">
							<p class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Resolution & Dimensions</p>
							<p class="font-semibold text-slate-200">{selectedPhoto.width > 0 ? `${selectedPhoto.width} × ${selectedPhoto.height} px` : 'Native Resolution'}</p>
						</div>

						<div class="space-y-1">
							<p class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">File Size</p>
							<p class="font-semibold text-slate-200">{formatBytes(selectedPhoto.size)}</p>
						</div>

						<div class="space-y-1">
							<p class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Device & Ingest Info</p>
							<p class="font-semibold text-slate-200">{selectedPhoto.exif_model || 'DeepPhotos Ingest'}</p>
						</div>
					</div>
				</div>
			{/if}
		</div>
	{/if}

	<!-- MOVE TO LOCKED FOLDER DIALOG (Z-INDEX 130) -->
	{#if showMoveToVaultModal && selectedPhoto}
		<div class="fixed inset-0 z-[130] bg-slate-950/70 backdrop-blur-sm flex items-center justify-center p-4">
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

	<!-- DRAG & DROP EXTERNAL OVERLAYS -->
	{#if isDraggingExternal}
		<div class="fixed inset-0 z-[200] bg-sky-500/10 backdrop-blur-sm border-[6px] border-dashed border-sky-400 flex flex-col items-center justify-center animate-fade-in pointer-events-none transition-all">
			<div class="bg-white dark:bg-slate-900 p-8 rounded-3xl shadow-2xl flex flex-col items-center justify-center gap-4">
				<div class="w-20 h-20 rounded-full bg-sky-100 dark:bg-sky-950 text-sky-500 flex items-center justify-center shadow-inner">
					<Download class="w-10 h-10 animate-bounce" />
				</div>
				<div class="text-center">
					<h2 class="text-xl font-bold text-slate-900 dark:text-white">Drop to import</h2>
					<p class="text-xs text-slate-500 dark:text-slate-400 mt-1">We'll download this image directly to your vault</p>
				</div>
			</div>
		</div>
	{/if}

	{#if isUploadingExternal}
		<div class="fixed inset-0 z-[200] bg-slate-950/80 backdrop-blur-md flex flex-col items-center justify-center animate-fade-in pointer-events-none">
			<div class="bg-white dark:bg-slate-900 p-8 rounded-3xl shadow-2xl flex flex-col items-center justify-center gap-4 w-64">
				<div class="w-16 h-16 rounded-full border-4 border-sky-100 border-t-sky-500 animate-spin"></div>
				<div class="text-center">
					<h2 class="text-base font-bold text-slate-900 dark:text-white">Importing Media...</h2>
					<p class="text-xs text-slate-500 dark:text-slate-400 mt-1">Downloading file from the web</p>
				</div>
			</div>
		</div>
	{/if}
</div>
