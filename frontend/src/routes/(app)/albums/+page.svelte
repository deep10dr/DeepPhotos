<script lang="ts">
	import { onMount } from 'svelte';
	import { appState } from '$lib/state.svelte';
	import { apiFetch } from '$lib/api';
	import { FolderClosed, Plus, Image as ImageIcon, X, Trash2, ArrowLeft, CheckCircle2, Check, PlusCircle, Lock, ChevronLeft, ChevronRight, ShieldCheck, AlertCircle } from 'lucide-svelte';

	interface Album {
		id: string;
		name: string;
		description: string;
		cover_url?: string;
		photos_count: number;
	}

	interface PhotoItem {
		id: string;
		title: string;
		thumbnail_url?: string;
		url?: string;
	}

	let albums = $state<Album[]>([]);
	let isLoading = $state(true);
	let showCreateModal = $state(false);
	let albumName = $state('');
	let albumDesc = $state('');

	// Selected Album Detail View State
	let selectedAlbum = $state<Album | null>(null);
	let albumPhotos = $state<PhotoItem[]>([]);
	let showAddPhotosModal = $state(false);
	let availablePhotos = $state<PhotoItem[]>([]);
	let selectedPhotoIds = $state<string[]>([]);

	// Lightbox & Move to Vault State
	let lightboxPhoto = $state<PhotoItem | null>(null);
	let lightboxIndex = $state<number | null>(null);
	let showMoveToVaultModal = $state(false);
	let lockedFolders = $state<any[]>([]);
	let selectedVaultFolderId = $state('');
	let vaultPasscode = $state('');
	let vaultError = $state('');

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

	function openLightbox(index: number) {
		lightboxIndex = index;
		lightboxPhoto = {
			...albumPhotos[index],
			url: albumPhotos[index].url || `${appState.apiBaseUrl}/api/photos/${albumPhotos[index].id}/file`
		};
	}

	function closeLightbox() {
		lightboxPhoto = null;
		lightboxIndex = null;
	}

	function prevLightboxPhoto() {
		if (lightboxIndex !== null && lightboxIndex > 0) {
			openLightbox(lightboxIndex - 1);
		}
	}

	function nextLightboxPhoto() {
		if (lightboxIndex !== null && lightboxIndex < albumPhotos.length - 1) {
			openLightbox(lightboxIndex + 1);
		}
	}

	function openMoveToVaultDialog() {
		if (lockedFolders.length === 0) {
			alert('No locked folders available. Create one first in Locked Vault.');
			return;
		}
		selectedVaultFolderId = lockedFolders[0].id;
		vaultPasscode = '';
		vaultError = '';
		showMoveToVaultModal = true;
	}

	async function handleMoveToLockedFolder(e: Event) {
		e.preventDefault();
		if (!lightboxPhoto || !selectedVaultFolderId || !vaultPasscode) return;
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

			const updateRes = await apiFetch(`/api/photos/${lightboxPhoto.id}`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ locked_folder_id: selectedVaultFolderId })
			});

			if (updateRes.ok) {
				showMoveToVaultModal = false;
				closeLightbox();
				if (selectedAlbum) fetchAlbumDetail(selectedAlbum.id);
			} else {
				vaultError = 'Failed to move photo to locked folder.';
			}
		} catch (err) {
			console.error('Error moving to locked folder:', err);
			vaultError = 'Network error.';
		}
	}

	async function fetchAlbums() {
		isLoading = true;
		try {
			const res = await apiFetch('/api/albums');
			if (res.ok) {
				albums = await res.json();
			}
		} catch (e) {
			console.warn('API error fetching albums:', e);
		} finally {
			isLoading = false;
		}
	}

	async function fetchAlbumDetail(albumId: string) {
		try {
			const res = await apiFetch(`/api/albums/${albumId}`);
			if (res.ok) {
				const albumData = await res.json();
				selectedAlbum = albumData;
				// In a real API call, album photos are embedded or fetched
				albumPhotos = albumData.photos || [];
			}
		} catch (e) {
			console.warn('API error fetching album detail:', e);
		}
	}

	async function fetchAvailablePhotos() {
		try {
			const res = await apiFetch('/api/photos?deleted=false');
			if (res.ok) {
				const data = await res.json();
				availablePhotos = data.map((p: any) => ({
					id: p.id,
					title: p.title,
					thumbnail_url: `${appState.apiBaseUrl}/api/photos/${p.id}/thumbnail`
				}));
			}
		} catch (e) {
			console.warn('Error fetching available photos:', e);
		}
	}

	onMount(() => {
		fetchAlbums();
		fetchLockedFolders();
	});

	function openAlbum(album: Album) {
		selectedAlbum = album;
		fetchAlbumDetail(album.id);
	}

	function closeAlbum() {
		selectedAlbum = null;
		albumPhotos = [];
		fetchAlbums();
	}

	async function handleCreateAlbum(e: Event) {
		e.preventDefault();
		if (!albumName) return;

		try {
			const res = await apiFetch('/api/albums', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ name: albumName, description: albumDesc })
			});

			if (res.ok) {
				const newAlbum = await res.json();
				albums = [...albums, newAlbum];
				albumName = '';
				albumDesc = '';
				showCreateModal = false;
			}
		} catch (err) {
			console.error('Error creating album:', err);
		}
	}

	async function handleDeleteAlbum(id: string, name: string, e: Event) {
		e.stopPropagation();
		if (!confirm(`Delete album "${name}"?`)) return;

		try {
			const res = await apiFetch(`/api/albums/${id}`, { method: 'DELETE' });
			if (res.ok) {
				if (selectedAlbum?.id === id) closeAlbum();
				await fetchAlbums();
			}
		} catch (err) {
			console.error('Error deleting album:', err);
		}
	}

	function openAddPhotosModal() {
		selectedPhotoIds = [];
		fetchAvailablePhotos();
		showAddPhotosModal = true;
	}

	function togglePhotoSelection(id: string) {
		if (selectedPhotoIds.includes(id)) {
			selectedPhotoIds = selectedPhotoIds.filter(i => i !== id);
		} else {
			selectedPhotoIds = [...selectedPhotoIds, id];
		}
	}

	async function handleAddPhotosToAlbum() {
		if (!selectedAlbum || selectedPhotoIds.length === 0) return;

		try {
			const res = await apiFetch(`/api/albums/${selectedAlbum.id}/photos`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ photo_ids: selectedPhotoIds })
			});

			if (res.ok) {
				showAddPhotosModal = false;
				await fetchAlbumDetail(selectedAlbum.id);
				await fetchAlbums();
			}
		} catch (err) {
			console.error('Error adding photos to album:', err);
		}
	}
</script>

<div class="h-full flex flex-col gap-6 animate-fade-in">
	
	{#if selectedAlbum}
		<!-- ALBUM DETAIL VIEW -->
		<div class="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4 shrink-0">
			<div class="flex items-center gap-3">
				<button
					type="button"
					onclick={closeAlbum}
					class="p-2 rounded-xl border border-slate-200 dark:border-slate-800 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300 transition-colors"
				>
					<ArrowLeft class="w-5 h-5" />
				</button>
				<div>
					<h1 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
						<FolderClosed class="w-6 h-6 text-sky-500" />
						{selectedAlbum.name}
					</h1>
					<p class="text-xs text-slate-500 dark:text-slate-400">{selectedAlbum.description || 'Custom Collection'}</p>
				</div>
			</div>

			<button
				type="button"
				onclick={openAddPhotosModal}
				class="px-4 py-2 rounded-xl bg-sky-400 text-white text-xs font-bold flex items-center gap-1.5 shadow-sm shadow-sky-300/50 hover:bg-sky-500 transition-all cursor-pointer"
			>
				<Plus class="w-4 h-4" />
				Add Photos
			</button>
		</div>

		<div class="flex-1 overflow-y-auto min-h-0">
			<!-- EMPTY ALBUM VIEW ENHANCEMENT -->
			{#if albumPhotos.length === 0}
				<div class="p-12 text-center rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm max-w-lg mx-auto my-12 space-y-4">
					<div class="w-16 h-16 rounded-3xl bg-sky-100 dark:bg-sky-950 text-sky-500 flex items-center justify-center mx-auto">
						<FolderClosed class="w-8 h-8" />
					</div>
					<div class="space-y-1">
						<h3 class="text-lg font-bold text-slate-900 dark:text-white">This album is empty</h3>
						<p class="text-xs text-slate-500 dark:text-slate-400">Please add photos to this album to build your custom collection.</p>
					</div>

					<button
						type="button"
						onclick={openAddPhotosModal}
						class="px-6 py-2.5 rounded-xl bg-sky-400 text-white text-xs font-bold inline-flex items-center gap-2 shadow-md shadow-sky-300/50 hover:bg-sky-500 transition-all cursor-pointer"
					>
						<PlusCircle class="w-4 h-4" />
						Add Photos Now
					</button>
				</div>
			{:else}
				<!-- Album Grid -->
				<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
					{#each albumPhotos as photo}
						<div class="aspect-4/3 rounded-2xl overflow-hidden bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-800 shadow-sm relative group cursor-pointer" onclick={() => openLightbox(albumPhotos.indexOf(photo))}>
							<img src={photo.thumbnail_url} alt={photo.title} class="w-full h-full object-cover group-hover:scale-105 transition-transform" />
							<div class="absolute inset-0 bg-gradient-to-t from-slate-900/60 via-transparent to-transparent p-3 flex items-end">
								<p class="text-xs font-bold text-white truncate">{photo.title}</p>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>

	{:else}
		<!-- ALBUMS OVERVIEW LIST -->
		<div class="flex items-center justify-between shrink-0">
			<div>
				<h1 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
					<FolderClosed class="w-6 h-6 text-sky-500" />
					Albums
				</h1>
				<p class="text-xs text-slate-500 dark:text-slate-400">Organized photo collections in your vault</p>
			</div>
			<button
				type="button"
				onclick={() => showCreateModal = true}
				class="px-4 py-2 rounded-xl bg-sky-400 text-white text-xs font-bold flex items-center gap-1.5 shadow-sm shadow-sky-300/50 hover:bg-sky-500 transition-all cursor-pointer"
			>
				<Plus class="w-4 h-4" />
				New Album
			</button>
		</div>

		<div class="flex-1 overflow-y-auto min-h-0">
			{#if !isLoading && albums.length === 0}
				<div class="p-12 text-center rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm max-w-md mx-auto my-12 space-y-3">
					<div class="w-12 h-12 rounded-2xl bg-sky-100 dark:bg-sky-950 text-sky-500 flex items-center justify-center mx-auto">
						<FolderClosed class="w-6 h-6" />
					</div>
					<h3 class="text-base font-bold text-slate-900 dark:text-white">No albums created</h3>
					<p class="text-xs text-slate-500 dark:text-slate-400">Click <strong class="text-sky-500">New Album</strong> to organize your photos into custom collections.</p>
				</div>
			{:else}
				<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
					{#each albums as album}
						<div
							role="button"
							tabindex="0"
							onclick={() => openAlbum(album)}
							onkeydown={(e) => e.key === 'Enter' && openAlbum(album)}
							class="p-4 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:shadow-md transition-all group cursor-pointer relative"
						>
							<div class="aspect-4/3 rounded-xl overflow-hidden bg-slate-100 dark:bg-slate-800 mb-3 relative flex items-center justify-center text-slate-400">
								{#if album.cover_url}
									<img src={album.cover_url} alt={album.name} class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
								{:else}
									<FolderClosed class="w-10 h-10 text-sky-400/60" />
								{/if}

								<button
									type="button"
									onclick={(e) => handleDeleteAlbum(album.id, album.name, e)}
									class="absolute top-3 right-3 p-1.5 rounded-full bg-slate-900/60 hover:bg-rose-500 text-white backdrop-blur-md transition-colors opacity-0 group-hover:opacity-100"
									title="Delete Album"
								>
									<Trash2 class="w-3.5 h-3.5" />
								</button>
							</div>
							<h3 class="text-sm font-bold text-slate-900 dark:text-white">{album.name}</h3>
							<p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">{album.photos_count || 0} photos • {album.description || 'Collection'}</p>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	{/if}

	{#if showCreateModal}
		<!-- Create Album Modal -->
		<div class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-sm flex items-center justify-center p-4">
			<div class="bg-white dark:bg-slate-900 rounded-3xl p-6 md:p-8 max-w-md w-full border border-slate-200 dark:border-slate-800 shadow-2xl animate-fade-in relative">
				<button
					type="button"
					onclick={() => showCreateModal = false}
					class="absolute top-5 right-5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1"
				>
					<X class="w-5 h-5" />
				</button>

				<h3 class="text-lg font-bold text-slate-900 dark:text-white mb-1 flex items-center gap-2">
					<FolderClosed class="w-5 h-5 text-sky-500" />
					Create New Album
				</h3>
				<p class="text-xs text-slate-500 dark:text-slate-400 mb-6">Group photos into a custom collection</p>

				<form onsubmit={handleCreateAlbum} class="space-y-4">
					<div class="space-y-1.5">
						<label for="album-name" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Album Name</label>
						<input
							id="album-name"
							type="text"
							bind:value={albumName}
							placeholder="e.g. Summer Vacation 2026"
							class="w-full h-10 px-3.5 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400"
							required
						/>
					</div>

					<div class="space-y-1.5">
						<label for="album-desc" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Description</label>
						<input
							id="album-desc"
							type="text"
							bind:value={albumDesc}
							placeholder="e.g. California trip highlights"
							class="w-full h-10 px-3.5 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400"
						/>
					</div>

					<div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100 dark:border-slate-800">
						<button
							type="button"
							onclick={() => showCreateModal = false}
							class="px-4 py-2 rounded-xl text-xs font-semibold text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-5 py-2 rounded-xl bg-sky-400 hover:bg-sky-500 text-white text-xs font-bold shadow-sm shadow-sky-300/50"
						>
							Create Album
						</button>
					</div>
				</form>
			</div>
		</div>
	{/if}

	{#if showAddPhotosModal && selectedAlbum}
		<!-- Add Photos Modal -->
		<div class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-sm flex items-center justify-center p-4">
			<div class="bg-white dark:bg-slate-900 rounded-3xl p-6 max-w-2xl w-full border border-slate-200 dark:border-slate-800 shadow-2xl animate-fade-in relative flex flex-col max-h-[80vh]">
				<button
					type="button"
					onclick={() => showAddPhotosModal = false}
					class="absolute top-5 right-5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1"
				>
					<X class="w-5 h-5" />
				</button>

				<h3 class="text-lg font-bold text-slate-900 dark:text-white mb-1 flex items-center gap-2">
					<ImageIcon class="w-5 h-5 text-sky-500" />
					Add Photos to "{selectedAlbum.name}"
				</h3>
				<p class="text-xs text-slate-500 dark:text-slate-400 mb-4">Select photos from your timeline gallery</p>

				<div class="flex-1 overflow-y-auto grid grid-cols-3 gap-3 p-1">
					{#each availablePhotos as photo}
						<div
							role="button"
							tabindex="0"
							onclick={() => togglePhotoSelection(photo.id)}
							onkeydown={(e) => e.key === 'Enter' && togglePhotoSelection(photo.id)}
							class={`aspect-square rounded-xl overflow-hidden relative cursor-pointer border-2 transition-all ${
								selectedPhotoIds.includes(photo.id)
									? 'border-sky-400 ring-2 ring-sky-200 dark:ring-sky-900'
									: 'border-transparent'
							}`}
						>
							<img src={photo.thumbnail_url} alt={photo.title} class="w-full h-full object-cover" />
							{#if selectedPhotoIds.includes(photo.id)}
								<div class="absolute top-2 right-2 p-1 rounded-full bg-sky-400 text-white shadow-md">
									<Check class="w-3.5 h-3.5 stroke-[3]" />
								</div>
							{/if}
						</div>
					{/each}
				</div>

				<div class="flex items-center justify-between pt-4 border-t border-slate-100 dark:border-slate-800 mt-4">
					<span class="text-xs text-slate-500 font-semibold">{selectedPhotoIds.length} photos selected</span>
					<div class="flex items-center gap-3">
						<button
							type="button"
							onclick={() => showAddPhotosModal = false}
							class="px-4 py-2 rounded-xl text-xs font-semibold text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800"
						>
							Cancel
						</button>
						<button
							type="button"
							onclick={handleAddPhotosToAlbum}
							disabled={selectedPhotoIds.length === 0}
							class="px-5 py-2 rounded-xl bg-sky-400 hover:bg-sky-500 text-white text-xs font-bold shadow-sm shadow-sky-300/50 disabled:opacity-50"
						>
							Add Selected Photos
						</button>
					</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- ALBUM LIGHTBOX -->
	{#if lightboxPhoto && lightboxIndex !== null}
		<div class="fixed inset-0 z-[100] w-screen h-screen bg-black/95 backdrop-blur-xl flex flex-col select-none overflow-hidden animate-fade-in">
			<!-- TOP BAR -->
			<div class="absolute top-0 left-0 right-0 z-[110] p-4 flex items-center justify-between pointer-events-none">
				<div class="flex items-center gap-3 pointer-events-auto bg-black/40 backdrop-blur-md pl-2 pr-4 py-2 rounded-full border border-white/10">
					<button type="button" onclick={closeLightbox} class="p-2 rounded-full text-slate-200 hover:text-white hover:bg-white/10 transition-colors cursor-pointer">
						<ArrowLeft class="w-5 h-5" />
					</button>
					<div>
						<h3 class="text-sm font-bold text-white flex items-center gap-2">
							<span>{lightboxPhoto.title}</span>
							<span class="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-white/15 text-slate-200">
								{lightboxIndex + 1} of {albumPhotos.length}
							</span>
						</h3>
					</div>
				</div>

				<div class="flex items-center gap-1.5 sm:gap-2 pointer-events-auto bg-black/40 backdrop-blur-md px-2 py-1.5 rounded-full border border-white/10">
					<button type="button" onclick={openMoveToVaultDialog} class="p-2 rounded-full text-slate-200 hover:text-sky-300 hover:bg-white/10 transition-colors cursor-pointer" title="Move to Locked Vault">
						<Lock class="w-5 h-5" />
					</button>
					<button type="button" onclick={closeLightbox} class="p-2 rounded-full text-slate-200 hover:text-white hover:bg-white/10 transition-colors cursor-pointer ml-1" title="Close (Esc)">
						<X class="w-6 h-6" />
					</button>
				</div>
			</div>

			<!-- ARROWS -->
			{#if lightboxIndex > 0}
				<button type="button" onclick={prevLightboxPhoto} class="absolute left-4 sm:left-8 top-1/2 -translate-y-1/2 z-[110] w-12 h-12 rounded-full bg-slate-900/60 hover:bg-white/20 text-white flex items-center justify-center transition-all cursor-pointer backdrop-blur-md">
					<ChevronLeft class="w-7 h-7" />
				</button>
			{/if}
			{#if lightboxIndex < albumPhotos.length - 1}
				<button type="button" onclick={nextLightboxPhoto} class="absolute right-4 sm:right-8 top-1/2 -translate-y-1/2 z-[110] w-12 h-12 rounded-full bg-slate-900/60 hover:bg-white/20 text-white flex items-center justify-center transition-all cursor-pointer backdrop-blur-md">
					<ChevronRight class="w-7 h-7" />
				</button>
			{/if}

			<!-- IMAGE CANVAS -->
			<div class="flex-1 w-full h-full relative flex items-center justify-center p-0 overflow-hidden">
				<img src={lightboxPhoto.url} alt={lightboxPhoto.title} class="w-full h-full object-contain select-none cursor-grab active:cursor-grabbing" />
			</div>
		</div>
	{/if}

	<!-- MOVE TO LOCKED FOLDER DIALOG (Z-INDEX 130) -->
	{#if showMoveToVaultModal && lightboxPhoto}
		<div class="fixed inset-0 z-[130] bg-slate-950/70 backdrop-blur-sm flex items-center justify-center p-4">
			<div class="bg-white dark:bg-slate-900 rounded-3xl p-6 md:p-8 max-w-md w-full border border-slate-200 dark:border-slate-800 shadow-2xl animate-fade-in relative">
				<button type="button" onclick={() => showMoveToVaultModal = false} class="absolute top-5 right-5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1">
					<X class="w-5 h-5" />
				</button>

				<div class="flex items-center gap-2 mb-1">
					<Lock class="w-5 h-5 text-sky-500" />
					<h3 class="text-lg font-bold text-slate-900 dark:text-white">Move to Locked Folder</h3>
				</div>
				<p class="text-xs text-slate-500 dark:text-slate-400 mb-6">
					Move "<strong class="text-slate-800 dark:text-slate-200">{lightboxPhoto.title}</strong>" into a passcode-protected folder
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
						<select id="target-folder" bind:value={selectedVaultFolderId} class="w-full h-10 px-3 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400">
							{#each lockedFolders as folder}
								<option value={folder.id}>{folder.name} ({folder.description || 'Locked'})</option>
							{/each}
						</select>
					</div>

					<div class="space-y-1.5">
						<label for="vault-passcode" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Folder 4-Digit Passcode</label>
						<input id="vault-passcode" type="password" maxLength={4} bind:value={vaultPasscode} placeholder="••••" class="w-full h-11 text-center font-mono text-lg tracking-widest rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400" required />
					</div>

					<div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100 dark:border-slate-800">
						<button type="button" onclick={() => showMoveToVaultModal = false} class="px-4 py-2 rounded-xl text-xs font-semibold text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800">Cancel</button>
						<button type="submit" class="px-5 py-2 rounded-xl bg-sky-400 hover:bg-sky-500 text-white text-xs font-bold shadow-sm shadow-sky-300/50 flex items-center gap-1.5"><ShieldCheck class="w-4 h-4" /> Move to Vault</button>
					</div>
				</form>
			</div>
		</div>
	{/if}
</div>
