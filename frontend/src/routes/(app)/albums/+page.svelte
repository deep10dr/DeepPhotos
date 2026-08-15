<script lang="ts">
	import { appState } from '$lib/state.svelte';
	import { FolderClosed, Plus, Image as ImageIcon, X, Trash2, ArrowLeft, CheckCircle2, Check, PlusCircle } from 'lucide-svelte';

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

	async function fetchAlbums() {
		isLoading = true;
		try {
			const res = await fetch(`${appState.apiBaseUrl}/api/albums`);
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
			const res = await fetch(`${appState.apiBaseUrl}/api/albums/${albumId}`);
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
			const res = await fetch(`${appState.apiBaseUrl}/api/photos?deleted=false`);
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

	$effect(() => {
		fetchAlbums();
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
			const res = await fetch(`${appState.apiBaseUrl}/api/albums`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ name: albumName, description: albumDesc })
			});

			if (res.ok) {
				albumName = '';
				albumDesc = '';
				showCreateModal = false;
				await fetchAlbums();
			}
		} catch (err) {
			console.error('Error creating album:', err);
		}
	}

	async function handleDeleteAlbum(id: string, name: string, e: Event) {
		e.stopPropagation();
		if (!confirm(`Delete album "${name}"?`)) return;

		try {
			const res = await fetch(`${appState.apiBaseUrl}/api/albums/${id}`, { method: 'DELETE' });
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
			const res = await fetch(`${appState.apiBaseUrl}/api/albums/${selectedAlbum.id}/photos`, {
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

<div class="space-y-6 animate-fade-in">
	
	{#if selectedAlbum}
		<!-- ALBUM DETAIL VIEW -->
		<div class="space-y-6">
			<div class="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
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
						<div class="aspect-4/3 rounded-2xl overflow-hidden bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-800 shadow-sm relative group">
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
		<div class="flex items-center justify-between">
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
</div>
