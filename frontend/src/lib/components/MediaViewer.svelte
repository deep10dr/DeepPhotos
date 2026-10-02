<script lang="ts">
	import { getMediaUrl } from '$lib/api';
	import {
		X,
		ChevronLeft,
		ChevronRight,
		Download,
		Heart,
		Lock,
		Trash2,
		ZoomIn,
		ZoomOut,
		RotateCcw,
		Info,
		Maximize2,
		Scaling,
		Play
	} from 'lucide-svelte';

	interface PhotoItem {
		id: string;
		title: string;
		filename: string;
		object_key?: string;
		mime_type: string;
		file_type: string;
		size: number;
		width?: number;
		height?: number;
		exif_model?: string;
		taken_at?: string;
		is_favorite?: boolean;
		is_deleted?: boolean;
		url?: string;
		thumbnail_url?: string;
	}

	interface Props {
		photos: PhotoItem[];
		selectedIndex: number | null;
		onclose: () => void;
		onnavigate?: (index: number) => void;
		ontoggleFavorite?: (photo: any) => void;
		ondelete?: (photo: any) => void;
		onmoveToVault?: (photo: any) => void;
	}

	let {
		photos,
		selectedIndex,
		onclose,
		onnavigate,
		ontoggleFavorite,
		ondelete,
		onmoveToVault
	}: Props = $props();

	let currentPhoto = $derived(
		selectedIndex !== null && selectedIndex >= 0 && selectedIndex < photos.length
			? photos[selectedIndex]
			: null
	);

	let showDetailsModal = $state(false);
	let showZoomControls = $state(false);
	let zoomScale = $state(100);
	let videoRef = $state<HTMLVideoElement | null>(null);
	let isPlaying = $state(false);

	function prevPhoto() {
		if (selectedIndex !== null && selectedIndex > 0) {
			resetState();
			onnavigate ? onnavigate(selectedIndex - 1) : null;
		}
	}

	function nextPhoto() {
		if (selectedIndex !== null && selectedIndex < photos.length - 1) {
			resetState();
			onnavigate ? onnavigate(selectedIndex + 1) : null;
		}
	}

	function resetState() {
		zoomScale = 100;
		showZoomControls = false;
		isPlaying = false;
	}

	function zoomIn() {
		zoomScale = Math.min(zoomScale + 25, 400);
	}

	function zoomOut() {
		zoomScale = Math.max(zoomScale - 25, 50);
	}

	function resetZoom() {
		zoomScale = 100;
	}

	function togglePlay() {
		if (!videoRef) return;
		if (videoRef.paused) {
			videoRef.play();
		} else {
			videoRef.pause();
		}
	}

	function formatBytes(bytes: number) {
		if (!bytes || bytes === 0) return '0 B';
		const k = 1024;
		const sizes = ['B', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	}

	function checkIsVideo(item: PhotoItem | null): boolean {
		if (!item) return false;
		if (item.file_type === 'video') return true;
		if (item.mime_type && item.mime_type.toLowerCase().startsWith('video/')) return true;
		if (item.filename) {
			const ext = item.filename.split('.').pop()?.toLowerCase();
			if (['mp4', 'webm', 'mov', 'mkv', 'avi', 'm4v', 'ogv', '3gp'].includes(ext || '')) return true;
		}
		return false;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (selectedIndex === null) return;
		if (e.key === 'ArrowLeft') prevPhoto();
		if (e.key === 'ArrowRight') nextPhoto();
		if (e.key === 'Escape') onclose();
		if (e.key === ' ' && checkIsVideo(currentPhoto)) {
			e.preventDefault();
			togglePlay();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if currentPhoto && selectedIndex !== null}
	<div class="fixed inset-0 z-9999 w-screen h-screen bg-black/95 dark:bg-slate-950/95 backdrop-blur-xl flex flex-col select-none overflow-hidden animate-fade-in">
		<!-- TOP OVERLAY TOOLBAR -->
		<div class="absolute top-0 left-0 right-0 z-10000 p-4 flex flex-wrap items-center justify-between gap-3 pointer-events-none">
			<!-- Left: Back Arrow & Photo Info -->
			<div class="flex items-center gap-2 pointer-events-auto p-2 ">
				<button
					type="button"
					onclick={onclose}
					title="Back (Esc)"
					class="p-1.5 rounded-full text-slate-200 hover:text-white hover:bg-white/10 transition-colors cursor-pointer"
				>
					<ChevronLeft class="w-5 h-5" />
				</button>

			</div>

			<!-- Right: Actions -->
			<div class="flex items-center gap-2.5 pointer-events-auto">
				<!-- Resized Badge with Tooltip -->
				<div class="relative group hidden md:flex items-center gap-1.5 px-3 py-1.5 rounded-full text-[11px] font-semibold text-slate-300  cursor-pointer">
					<Scaling class="w-3.5 h-3.5 text-sky-400" />
					<span
						class="pointer-events-none absolute top-full right-0 mt-2 whitespace-nowrap rounded-xl bg-slate-900/95 border border-white/10 px-3 py-1.5 text-[11px] font-medium text-slate-200 opacity-0 group-hover:opacity-100 transition-opacity duration-150 shadow-2xl z-[10001]"
						role="tooltip"
					>
						Image is resized because of the screen size
					</span>
				</div>

				<!-- Toolbar buttons -->
				<div class="flex items-center gap-1.5 bg-slate-900/80 backdrop-blur-xl p-1.5 rounded-full border border-white/10 shadow-lg">
					{#if currentPhoto.file_type !== 'video' && !currentPhoto.mime_type?.startsWith('video/')}
						<button
							type="button"
							onclick={() => showZoomControls = !showZoomControls}
							title="Zoom Controls"
							class={`p-2 rounded-full transition-colors cursor-pointer ${showZoomControls ? 'bg-sky-500 text-white' : 'text-slate-300 hover:text-white hover:bg-white/10'}`}
						>
							<ZoomIn class="w-4 h-4" />
						</button>
					{/if}

					{#if ontoggleFavorite}
						<button
							type="button"
							onclick={() => currentPhoto && ontoggleFavorite(currentPhoto)}
							title="Toggle Favorite"
							class="p-2 rounded-full text-slate-300 hover:text-white hover:bg-white/10 transition-colors cursor-pointer"
						>
							<Heart class={`w-4 h-4 ${currentPhoto.is_favorite ? 'fill-rose-500 text-rose-500' : ''}`} />
						</button>
					{/if}

					{#if onmoveToVault}
						<button
							type="button"
							onclick={() => currentPhoto && onmoveToVault(currentPhoto)}
							title="Move to Vault"
							class="p-2 rounded-full text-slate-300 hover:text-sky-300 hover:bg-sky-500/80 transition-colors cursor-pointer"
						>
							<Lock class="w-4 h-4" />
						</button>
					{/if}

					<button
						type="button"
						onclick={() => showDetailsModal = !showDetailsModal}
						title="Info & Exif"
						class={`p-2 rounded-full transition-colors cursor-pointer ${showDetailsModal ? 'bg-sky-500 text-white' : 'text-slate-300 hover:text-white hover:bg-white/10'}`}
					>
						<Info class="w-4 h-4" />
					</button>

					<a
						href={getMediaUrl(currentPhoto.url || `/api/media/${currentPhoto.id}/file`)}
						download={currentPhoto.filename}
						target="_blank"
						title="Download Original"
						class="p-2 rounded-full text-slate-300 hover:text-white hover:bg-white/10 transition-colors cursor-pointer"
					>
						<Download class="w-4 h-4" />
					</a>

					{#if ondelete}
						<button
							type="button"
							onclick={() => currentPhoto && ondelete(currentPhoto)}
							title="Delete"
							class="p-2 rounded-full text-slate-300 hover:text-rose-400 hover:bg-rose-500/20 transition-colors cursor-pointer"
						>
							<Trash2 class="w-4 h-4" />
						</button>
					{/if}

					<div class="w-px h-4 bg-white/20 mx-0.5"></div>

					<button
						type="button"
						onclick={onclose}
						title="Close (Esc)"
						class="p-2 rounded-full text-slate-300 hover:text-white hover:bg-white/10 transition-colors cursor-pointer"
					>
						<X class="w-4 h-4" />
					</button>
				</div>
			</div>
		</div>

		<!-- PREVIOUS / NEXT ARROWS -->
		{#if selectedIndex > 0}
			<button
				type="button"
				onclick={prevPhoto}
				title="Previous (Left Arrow)"
				class="absolute left-4 sm:left-8 top-1/2 -translate-y-1/2 z-10000 w-12 h-12 rounded-full bg-slate-900/80 hover:bg-slate-800 text-white flex items-center justify-center transition-all cursor-pointer backdrop-blur-md border border-white/10 shadow-xl"
			>
				<ChevronLeft class="w-7 h-7" />
			</button>
		{/if}

		{#if selectedIndex < photos.length - 1}
			<button
				type="button"
				onclick={nextPhoto}
				title="Next (Right Arrow)"
				class="absolute right-4 sm:right-8 top-1/2 -translate-y-1/2 z-10000 w-12 h-12 rounded-full bg-slate-900/80 hover:bg-slate-800 text-white flex items-center justify-center transition-all cursor-pointer backdrop-blur-md border border-white/10 shadow-xl"
			>
				<ChevronRight class="w-7 h-7" />
			</button>
		{/if}

		<!-- CENTER MEDIA CANVAS -->
		<div class="flex-1 w-full h-full min-h-0 min-w-0 relative flex items-center justify-center pt-20 pb-20 px-4 md:px-16 overflow-hidden">
			{#if checkIsVideo(currentPhoto)}
				{#key currentPhoto.id}
					<div class="relative max-w-[calc(100vw-64px)] max-h-[calc(100vh-140px)] flex items-center justify-center group">
						<video
							bind:this={videoRef}
							src={getMediaUrl(currentPhoto.url || `/api/media/${currentPhoto.id}/file`)}
							controls
							autoplay
							playsinline
							preload="auto"
							onplay={() => isPlaying = true}
							onpause={() => isPlaying = false}
							onended={() => isPlaying = false}
							class="max-w-[calc(100vw-64px)] max-h-[calc(100vh-140px)] w-auto h-auto object-contain rounded-2xl shadow-2xl bg-black border border-white/10"
						>
							<track kind="captions" />
						</video>

						{#if !isPlaying}
							<button
								type="button"
								onclick={togglePlay}
								class="absolute w-20 h-20 rounded-full bg-sky-500 hover:bg-sky-400 text-white shadow-2xl flex items-center justify-center transition-transform hover:scale-110 cursor-pointer z-10 border-2 border-white/40"
								title="Play Video"
							>
								<Play class="w-10 h-10 fill-white ml-1 text-white" />
							</button>
						{/if}
					</div>
				{/key}
			{:else}
				<img
					src={getMediaUrl(currentPhoto.url || `/api/media/${currentPhoto.id}/file`)}
					alt={currentPhoto.title}
					style={`transform: scale(${zoomScale / 100}); transition: transform 0.2s cubic-bezier(0.2, 0, 0, 1);`}
					class="max-w-[calc(100vw-64px)] max-h-[calc(100vh-140px)] w-auto h-auto object-contain select-none rounded-2xl shadow-2xl"
				/>
			{/if}
		</div>

		<!-- BOTTOM ZOOM BAR -->
		{#if showZoomControls && !checkIsVideo(currentPhoto)}
			<div class="absolute bottom-6 left-1/2 -translate-x-1/2 z-10000 p-3 px-6 bg-slate-900/90 backdrop-blur-xl border border-slate-800 rounded-2xl flex items-center justify-center gap-4 shadow-2xl animate-fade-in">
				<div class="flex items-center gap-3">
					<button
						type="button"
						onclick={zoomOut}
						disabled={zoomScale <= 50}
						class="p-2 rounded-xl bg-slate-800 text-white hover:bg-slate-700 disabled:opacity-50 cursor-pointer"
					>
						<ZoomOut class="w-4 h-4" />
					</button>
					<span class="text-xs font-mono font-bold text-white w-12 text-center">{zoomScale}%</span>
					<button
						type="button"
						onclick={zoomIn}
						disabled={zoomScale >= 300}
						class="p-2 rounded-xl bg-slate-800 text-white hover:bg-slate-700 disabled:opacity-50 cursor-pointer"
					>
						<ZoomIn class="w-4 h-4" />
					</button>
					<button
						type="button"
						onclick={resetZoom}
						class="p-2 rounded-xl bg-slate-800 text-white hover:bg-slate-700 cursor-pointer ml-2"
						title="Reset Zoom"
					>
						<RotateCcw class="w-4 h-4" />
					</button>
				</div>
			</div>
		{/if}

		<!-- INFO DRAWER -->
		{#if showDetailsModal}
			<div class="absolute right-0 top-0 bottom-0 z-10000 w-80 bg-slate-900/95 border-l border-slate-800 p-6 text-white overflow-y-auto backdrop-blur-2xl shadow-2xl animate-fade-in flex flex-col gap-6">
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
						<p class="font-bold text-slate-100 text-sm break-all">{currentPhoto.filename}</p>
					</div>

					<div class="space-y-1">
						<p class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">MIME & Category</p>
						<p class="font-semibold text-slate-200">{currentPhoto.mime_type} ({currentPhoto.file_type})</p>
					</div>

					<div class="space-y-1">
						<p class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Resolution & Dimensions</p>
						<p class="font-semibold text-slate-200">{currentPhoto.width && currentPhoto.width > 0 ? `${currentPhoto.width} × ${currentPhoto.height} px` : 'Native Resolution'}</p>
					</div>

					<div class="space-y-1">
						<p class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">File Size</p>
						<p class="font-semibold text-slate-200">{formatBytes(currentPhoto.size)}</p>
					</div>

					<div class="space-y-1">
						<p class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Device & Ingest Info</p>
						<p class="font-semibold text-slate-200">{currentPhoto.exif_model || 'DeepPhotos Ingest'}</p>
					</div>
				</div>
			</div>
		{/if}
	</div>
{/if}
