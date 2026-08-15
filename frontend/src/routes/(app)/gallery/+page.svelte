<script lang="ts">
	import { appState } from '$lib/state.svelte';
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
		Maximize2,
		Minimize2,
		ZoomIn,
		ZoomOut,
		RotateCcw,
		ShieldCheck,
		AlertCircle
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

	// Lightbox Zoom & Maximize Controls
	let zoomScale = $state(100);
	let isMaximized = $state(true); // Default to full screen view per user request

	// Move to Locked Folder Dialog State
	let showMoveToVaultModal = $state(false);
	let selectedVaultFolderId = $state('');
	let vaultPasscode = $state('');
	let vaultError = $state('');

	async function fetchPhotos() {
		isLoading = true;
		try {
			const res = await fetch(`${appState.apiBaseUrl}/api/photos?deleted=false`);
			if (res.ok) {
				const data: PhotoItem[] = await res.json();
				// Filter out photos that are moved to a locked folder
				photos = data.filter(p => !p.locked_folder_id).map(p => ({
					...p,
					url: `${appState.apiBaseUrl}/api/photos/${p.id}/file`,
					thumbnail_url: `${appState.apiBaseUrl}/api/photos/${p.id}/thumbnail`
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
			const res = await fetch(`${appState.apiBaseUrl}/api/locked-folders`);
			if (res.ok) {
				lockedFolders = await res.json();
			}
		} catch (e) {
			console.warn('API error fetching locked folders:', e);
		}
	}

	$effect(() => {
		const v = appState.uploadVersion;
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
		isMaximized = true;
		showDetailsModal = false;
		showMoveToVaultModal = false;
	}

	function closePhoto() {
		selectedIndex = null;
		zoomScale = 100;
		showDetailsModal = false;
		showMoveToVaultModal = false;
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
		if (zoomScale < 300) zoomScale += 20;
	}

	function zoomOut() {
		if (zoomScale > 50) zoomScale -= 20;
	}

	function resetZoom() {
		zoomScale = 100;
	}

	async function toggleFavorite(photo: PhotoItem, e?: Event) {
		if (e) e.stopPropagation();
		const updatedState = !photo.is_favorite;
		photo.is_favorite = updatedState;
		try {
			await fetch(`${appState.apiBaseUrl}/api/photos/${photo.id}`, {
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
			await fetch(`${appState.apiBaseUrl}/api/photos/${photo.id}`, {
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

		// 1. Verify 4-digit passcode of target locked folder
		try {
			const verifyRes = await fetch(`${appState.apiBaseUrl}/api/locked-folders/${selectedVaultFolderId}/verify`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ passcode: vaultPasscode })
			});

			if (!verifyRes.ok) {
				vaultError = 'Incorrect passcode for locked folder!';
				return;
			}

			// 2. Move photo to locked folder
			const updateRes = await fetch(`${appState.apiBaseUrl}/api/photos/${selectedPhoto.id}`, {
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

<div class="space-y-6 animate-fade-in">
	
	<!-- Header Bar -->
	<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
		<div>
			<h1 class="text-2xl font-bold text-slate-900 dark:text-white">Media Vault & Gallery</h1>
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

	<!-- Empty State when no photos exist -->
	{#if !isLoading && filteredPhotos.length === 0}
		<div class="p-12 text-center rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm max-w-md mx-auto my-12 space-y-3">
			<div class="w-12 h-12 rounded-2xl bg-sky-100 dark:bg-sky-950 text-sky-500 flex items-center justify-center mx-auto">
				<ImageIcon class="w-6 h-6" />
			</div>
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
									img.src = photo.url || '';
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

	{#if selectedPhoto && selectedIndex !== null}
		<!-- Fullscreen / Resizable Lightbox Modal -->
		<div class="fixed inset-0 z-50 bg-slate-950/95 backdrop-blur-md flex items-center justify-center p-2 sm:p-6 transition-all duration-300">
			
			<!-- PROMINENT TOP-RIGHT CANCEL / CLOSE BUTTON -->
			<button
				type="button"
				onclick={closePhoto}
				title="Close Lightbox (Esc)"
				class="absolute top-4 right-4 z-40 px-3.5 py-2 rounded-2xl bg-white/20 hover:bg-rose-500 text-white text-xs font-bold flex items-center gap-1.5 backdrop-blur-md shadow-xl transition-all cursor-pointer border border-white/20"
			>
				<X class="w-4 h-4" />
				<span>Close</span>
			</button>

			<!-- PREVIOUS MEDIA BUTTON -->
			<button
				type="button"
				onclick={prevPhoto}
				title="Previous (Left Arrow)"
				class="absolute left-4 md:left-8 z-30 w-12 h-12 rounded-full bg-white/10 hover:bg-sky-400 text-white flex items-center justify-center transition-all cursor-pointer shadow-lg backdrop-blur-md group"
			>
				<ChevronLeft class="w-7 h-7 transition-transform group-hover:-translate-x-0.5" />
			</button>

			<!-- NEXT MEDIA BUTTON -->
			<button
				type="button"
				onclick={nextPhoto}
				title="Next (Right Arrow)"
				class="absolute right-4 md:right-8 z-30 w-12 h-12 rounded-full bg-white/10 hover:bg-sky-400 text-white flex items-center justify-center transition-all cursor-pointer shadow-lg backdrop-blur-md group"
			>
				<ChevronRight class="w-7 h-7 transition-transform group-hover:translate-x-0.5" />
			</button>

			<!-- Lightbox Main Container Card (Defaults to Full Screen View) -->
			<div class={`bg-white dark:bg-slate-900 border border-slate-200/50 dark:border-slate-800 shadow-2xl transition-all duration-300 flex flex-col relative overflow-hidden ${
				isMaximized
					? 'w-full h-full rounded-none max-w-none'
					: 'max-w-5xl w-full rounded-3xl max-h-[90vh]'
			}`}>
				
				<!-- Media Display Area with Dynamic Zooming -->
				<div class="relative bg-slate-950 flex-1 flex items-center justify-center overflow-auto min-h-[380px] p-4">
					{#if selectedPhoto.file_type === 'video' || selectedPhoto.mime_type.startsWith('video/')}
						<video src={selectedPhoto.url} controls autoplay class="max-h-[78vh] w-auto max-w-full mx-auto shadow-xl rounded-xl"></video>
					{:else}
						<div class="overflow-auto max-w-full max-h-full flex items-center justify-center">
							<img
								src={selectedPhoto.url}
								alt={selectedPhoto.title}
								style={`transform: scale(${zoomScale / 100}); transition: transform 0.2s ease-out;`}
								class="max-h-[75vh] w-auto max-w-full object-contain mx-auto shadow-2xl"
							/>
						</div>
					{/if}
				</div>

				<!-- BOTTOM INTERACTIVE CONTROLS BAR: Zoom Slider & Actions -->
				<div class="p-4 md:px-6 bg-slate-100 dark:bg-slate-950 border-t border-slate-200 dark:border-slate-800 flex flex-wrap items-center justify-between gap-4">
					
					<!-- Left: Media Title & Counter -->
					<div class="flex items-center gap-3">
						<div>
							<h3 class="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
								<span>{selectedPhoto.title}</span>
								<span class="text-[11px] font-semibold px-2 py-0.5 rounded-full bg-sky-100 dark:bg-sky-950 text-sky-700 dark:text-sky-300">
									{selectedIndex + 1} of {filteredPhotos.length}
								</span>
							</h3>
							<p class="text-[11px] text-slate-500 dark:text-slate-400 mt-0.5">{selectedPhoto.taken_at || 'Recently'} • {formatBytes(selectedPhoto.size)}</p>
						</div>
					</div>

					<!-- Center: ZOOM SLIDER CONTROL (Min 50% to Max 300%) -->
					{#if selectedPhoto.file_type !== 'video' && !selectedPhoto.mime_type.startsWith('video/')}
						<div class="flex items-center gap-2.5 px-3.5 py-1.5 rounded-xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm">
							<button
								type="button"
								onclick={zoomOut}
								title="Zoom Out (-20%)"
								class="p-1 text-slate-500 hover:text-slate-800 dark:hover:text-slate-100 transition-colors cursor-pointer"
							>
								<ZoomOut class="w-4 h-4" />
							</button>

							<input
								type="range"
								min="50"
								max="300"
								step="5"
								bind:value={zoomScale}
								class="w-28 md:w-36 h-1.5 bg-slate-200 dark:bg-slate-700 rounded-lg appearance-none cursor-pointer accent-sky-400"
								title={`Current Zoom: ${zoomScale}%`}
							/>

							<button
								type="button"
								onclick={zoomIn}
								title="Zoom In (+20%)"
								class="p-1 text-slate-500 hover:text-slate-800 dark:hover:text-slate-100 transition-colors cursor-pointer"
							>
								<ZoomIn class="w-4 h-4" />
							</button>

							<span class="text-xs font-mono font-bold text-sky-600 dark:text-sky-400 w-10 text-center">{zoomScale}%</span>

							<button
								type="button"
								onclick={resetZoom}
								title="Reset Zoom (100%)"
								class="p-1 text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-colors cursor-pointer ml-1"
							>
								<RotateCcw class="w-3.5 h-3.5" />
							</button>
						</div>
					{/if}

					<!-- Right: Actions & Move to Locked Folder -->
					<div class="flex items-center gap-2">
						
						<!-- MOVE TO LOCKED FOLDER BUTTON -->
						<button
							type="button"
							onclick={openMoveToVaultDialog}
							class="px-3 py-1.5 rounded-xl border border-sky-200 dark:border-sky-800 bg-sky-50 dark:bg-sky-950/60 text-sky-700 dark:text-sky-300 hover:bg-sky-400 hover:text-white text-xs font-semibold flex items-center gap-1.5 transition-all cursor-pointer"
						>
							<Lock class="w-3.5 h-3.5" /> Move to Vault
						</button>

						<button
							type="button"
							onclick={() => isMaximized = !isMaximized}
							title={isMaximized ? 'Restore Window' : 'Maximize Full Screen'}
							class={`px-3 py-1.5 rounded-xl border text-xs font-semibold flex items-center gap-1.5 transition-colors cursor-pointer ${
								isMaximized
									? 'bg-sky-400 border-sky-400 text-white shadow-sm'
									: 'border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-200 hover:bg-white dark:hover:bg-slate-900'
							}`}
						>
							{#if isMaximized}
								<Minimize2 class="w-3.5 h-3.5" /> Restore
							{:else}
								<Maximize2 class="w-3.5 h-3.5" /> Full Screen
							{/if}
						</button>

						<button
							type="button"
							onclick={() => showDetailsModal = !showDetailsModal}
							class={`px-3 py-1.5 rounded-xl border text-xs font-semibold flex items-center gap-1.5 transition-colors cursor-pointer ${
								showDetailsModal
									? 'bg-sky-400 border-sky-400 text-white shadow-sm'
									: 'border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-200 hover:bg-white dark:hover:bg-slate-900'
							}`}
						>
							<Info class="w-3.5 h-3.5" /> Details
						</button>

						<a
							href={selectedPhoto.url}
							download={selectedPhoto.filename}
							target="_blank"
							class="px-3 py-1.5 rounded-xl border border-slate-200 dark:border-slate-700 hover:bg-white dark:hover:bg-slate-900 text-slate-700 dark:text-slate-200 text-xs font-semibold flex items-center gap-1.5 transition-colors cursor-pointer"
						>
							<Download class="w-3.5 h-3.5" /> Download
						</a>
					</div>
				</div>

				<!-- TECHNICAL METADATA INSPECTOR PANEL -->
				{#if showDetailsModal}
					<div class="p-5 bg-slate-50 dark:bg-slate-950 border-t border-slate-200 dark:border-slate-800 text-xs space-y-3 animate-fade-in">
						<h4 class="font-bold text-slate-900 dark:text-white flex items-center gap-1.5">
							<Info class="w-4 h-4 text-sky-500" />
							Ingest Metadata & Technical Details
						</h4>

						<div class="grid grid-cols-1 md:grid-cols-3 gap-3 text-slate-600 dark:text-slate-300">
							<div class="p-3 rounded-xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 space-y-1">
								<p class="text-[11px] font-semibold text-slate-400">Media Category & MIME</p>
								<p class="font-bold text-slate-800 dark:text-slate-100">{selectedPhoto.mime_type} ({selectedPhoto.file_type})</p>
							</div>

							<div class="p-3 rounded-xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 space-y-1">
								<p class="text-[11px] font-semibold text-slate-400">Resolution / Dimensions</p>
								<p class="font-bold text-slate-800 dark:text-slate-100">{selectedPhoto.width > 0 ? `${selectedPhoto.width} × ${selectedPhoto.height} px` : 'Native Resolution'}</p>
							</div>

							<div class="p-3 rounded-xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 space-y-1">
								<p class="text-[11px] font-semibold text-slate-400">Camera / Device Model</p>
								<p class="font-bold text-slate-800 dark:text-slate-100">{selectedPhoto.exif_model || 'DeepPhotos Ingest'}</p>
							</div>
						</div>
					</div>
				{/if}
			</div>

		</div>
	{/if}

	<!-- MOVE TO LOCKED FOLDER DIALOG -->
	{#if showMoveToVaultModal && selectedPhoto}
		<div class="fixed inset-0 z-50 bg-slate-950/70 backdrop-blur-sm flex items-center justify-center p-4">
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

</div>
