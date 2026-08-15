<script lang="ts">
	import { Image as ImageIcon, Heart, Calendar, Filter, Download, Trash2, Eye, ChevronLeft, ChevronRight, X } from 'lucide-svelte';

	const photos = [
		{ id: 1, title: 'Alpine Peak Sunrise', date: 'August 14, 2026', size: '4.2 MB', url: 'https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=1200&auto=format&fit=crop&q=80', favorite: true },
		{ id: 2, title: 'Coastal Waves', date: 'August 12, 2026', size: '3.8 MB', url: 'https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=1200&auto=format&fit=crop&q=80', favorite: false },
		{ id: 3, title: 'Forest Trail Path', date: 'August 10, 2026', size: '5.1 MB', url: 'https://images.unsplash.com/photo-1448375240586-882707db888b?w=1200&auto=format&fit=crop&q=80', favorite: true },
		{ id: 4, title: 'City Lights at Night', date: 'August 08, 2026', size: '2.9 MB', url: 'https://images.unsplash.com/photo-1519501025264-65ba15a82390?w=1200&auto=format&fit=crop&q=80', favorite: false },
		{ id: 5, title: 'Autumn Lake Reflection', date: 'August 04, 2026', size: '6.4 MB', url: 'https://images.unsplash.com/photo-1470071459604-3b5ec3a7fe05?w=1200&auto=format&fit=crop&q=80', favorite: true },
		{ id: 6, title: 'Desert Sand Dunes', date: 'July 28, 2026', size: '3.4 MB', url: 'https://images.unsplash.com/photo-1509316975850-ff9c5deb0cd9?w=1200&auto=format&fit=crop&q=80', favorite: false }
	];

	let selectedIndex = $state<number | null>(null);
	let activeFilter = $state<'all' | 'favorites'>('all');

	const filteredPhotos = $derived(
		activeFilter === 'favorites' ? photos.filter(p => p.favorite) : photos
	);

	const selectedPhoto = $derived(
		selectedIndex !== null && selectedIndex >= 0 && selectedIndex < filteredPhotos.length
			? filteredPhotos[selectedIndex]
			: null
	);

	function openPhoto(index: number) {
		selectedIndex = index;
	}

	function closePhoto() {
		selectedIndex = null;
	}

	function prevPhoto() {
		if (selectedIndex !== null) {
			selectedIndex = (selectedIndex - 1 + filteredPhotos.length) % filteredPhotos.length;
		}
	}

	function nextPhoto() {
		if (selectedIndex !== null) {
			selectedIndex = (selectedIndex + 1) % filteredPhotos.length;
		}
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
			<h1 class="text-2xl font-bold text-slate-900 dark:text-white">Photo Gallery</h1>
			<p class="text-xs text-slate-500 dark:text-slate-400">Timeline view of your imported photos</p>
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
				All Photos ({photos.length})
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
				Favorites ({photos.filter(p => p.favorite).length})
			</button>
		</div>
	</div>

	<!-- Photo Grid -->
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
					<img
						src={photo.url}
						alt={photo.title}
						class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
					/>
					<div class="absolute inset-0 bg-gradient-to-t from-slate-900/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity flex items-end p-4">
						<div class="text-white">
							<p class="text-xs font-bold">{photo.title}</p>
							<p class="text-[10px] text-slate-300">{photo.date} • {photo.size}</p>
						</div>
					</div>

					{#if photo.favorite}
						<div class="absolute top-3 right-3 p-1.5 rounded-full bg-white/80 dark:bg-slate-900/80 backdrop-blur-md text-rose-500 shadow-sm">
							<Heart class="w-3.5 h-3.5 fill-rose-500" />
						</div>
					{/if}
				</div>
			</div>
		{/each}
	</div>

	{#if selectedPhoto && selectedIndex !== null}
		<!-- Fullscreen Lightbox Modal with Next & Prev Buttons -->
		<div class="fixed inset-0 z-50 bg-slate-950/90 backdrop-blur-md flex items-center justify-center p-4">
			
			<!-- Close Button -->
			<button
				type="button"
				onclick={closePhoto}
				title="Close (Esc)"
				class="absolute top-5 right-5 z-20 w-10 h-10 rounded-full bg-white/10 text-white flex items-center justify-center hover:bg-white/20 transition-all cursor-pointer"
			>
				<X class="w-5 h-5" />
			</button>

			<!-- PREVIOUS PHOTO BUTTON -->
			<button
				type="button"
				onclick={prevPhoto}
				title="Previous Photo (Left Arrow)"
				class="absolute left-4 md:left-8 z-20 w-12 h-12 rounded-full bg-white/10 hover:bg-sky-400 text-white flex items-center justify-center transition-all cursor-pointer shadow-lg backdrop-blur-md group"
			>
				<ChevronLeft class="w-7 h-7 transition-transform group-hover:-translate-x-0.5" />
			</button>

			<!-- NEXT PHOTO BUTTON -->
			<button
				type="button"
				onclick={nextPhoto}
				title="Next Photo (Right Arrow)"
				class="absolute right-4 md:right-8 z-20 w-12 h-12 rounded-full bg-white/10 hover:bg-sky-400 text-white flex items-center justify-center transition-all cursor-pointer shadow-lg backdrop-blur-md group"
			>
				<ChevronRight class="w-7 h-7 transition-transform group-hover:translate-x-0.5" />
			</button>

			<!-- Lightbox Card Content -->
			<div class="bg-white dark:bg-slate-900 rounded-3xl overflow-hidden max-w-4xl w-full border border-slate-200/50 dark:border-slate-800 shadow-2xl animate-fade-in relative">
				<div class="relative bg-slate-950 flex items-center justify-center min-h-[350px] max-h-[65vh] overflow-hidden">
					<img
						src={selectedPhoto.url}
						alt={selectedPhoto.title}
						class="max-h-[65vh] w-auto max-w-full object-contain mx-auto"
					/>
				</div>

				<div class="p-5 md:p-6 flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900">
					<div>
						<div class="flex items-center gap-2">
							<h3 class="text-base font-bold text-slate-900 dark:text-white">{selectedPhoto.title}</h3>
							<span class="text-xs font-semibold px-2 py-0.5 rounded-full bg-sky-100 dark:bg-sky-950 text-sky-700 dark:text-sky-300">
								{selectedIndex + 1} of {filteredPhotos.length}
							</span>
						</div>
						<p class="text-xs text-slate-500 dark:text-slate-400 mt-1">{selectedPhoto.date} • {selectedPhoto.size}</p>
					</div>

					<div class="flex items-center gap-2">
						<button
							type="button"
							class="px-3.5 py-2 rounded-xl border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-200 text-xs font-semibold flex items-center gap-1.5 transition-colors cursor-pointer"
						>
							<Download class="w-4 h-4" /> Download Original
						</button>
					</div>
				</div>
			</div>

		</div>
	{/if}

</div>
