<script lang="ts">
	import type { Snippet } from 'svelte';
	import { CloudUpload, Loader2 } from 'lucide-svelte';
	import { apiFetch, invalidateCache } from '$lib/api';
	import { appState } from '$lib/state.svelte';
	import { notify } from '$lib/notify.svelte';

	interface Props {
		children?: Snippet;
		class?: string;
		fullScreenOverlay?: boolean;
		onUploadStart?: () => void;
		onUploadComplete?: () => void;
		onUploadError?: (error: string) => void;
	}

	let {
		children,
		class: className = '',
		fullScreenOverlay = true,
		onUploadStart,
		onUploadComplete,
		onUploadError
	}: Props = $props();

	let dragCounter = $state(0);
	let isDraggingOver = $state(false);
	let isUploading = $state(false);

	function handleDragEnter(e: DragEvent) {
		e.preventDefault();
		dragCounter++;
		if (dragCounter === 1) {
			isDraggingOver = true;
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
			isDraggingOver = false;
		}
	}

	async function handleDrop(e: DragEvent) {
		e.preventDefault();
		dragCounter = 0;
		isDraggingOver = false;

		const dt = e.dataTransfer;
		console.log(dt)
		
		if (!dt) return;

		// 1. LOCAL FILES — Dragged from filesystem
		if (dt.files && dt.files.length > 0) {
			isUploading = true;
			onUploadStart?.();
			try {
				const formData = new FormData();
				for (const file of Array.from(dt.files)) {
					formData.append('files', file);
				}

				const res = await apiFetch('/api/media/upload', {
					method: 'POST',
					body: formData
				});

				if (res.ok) {
					invalidateCache('/api/media');
					invalidateCache('/api/gallery');
					appState.refreshPhotos();
					notify.success('Media uploaded successfully!');
					onUploadComplete?.();
				} else {
					const errText = await res.text();
					notify.error(errText || 'Failed to upload drop files.');
					onUploadError?.(errText || 'Failed to upload drop files.');
				}
			} catch (err: any) {
				console.error('Drop zone upload error:', err);
				notify.error(err?.message || 'Network error during file drop upload.');
				onUploadError?.(err?.message || 'Network error during file drop upload.');
			} finally {
				isUploading = false;
			}
			return;
		}

		// 2. EXTERNAL WEB URL — Dragged from another browser tab or site
		let urlToUpload = '';

		const html = dt.getData('text/html');
		if (html) {
			const match = html.match(/src=["'](https?:\/\/[^"']+)["']/i);
			if (match && match[1]) {
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
			isUploading = true;
			onUploadStart?.();
			try {
				const res = await apiFetch('/api/media/upload-url', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ url: urlToUpload })
				});

				if (res.ok) {
					invalidateCache('/api/media');
					invalidateCache('/api/gallery');
					appState.refreshPhotos();
					notify.success('Media imported from URL!');
					onUploadComplete?.();
				} else {
					notify.error('Failed to import media from URL.');
					onUploadError?.('Failed to import media from URL.');
				}
			} catch (err: any) {
				console.error('Drop zone URL import error:', err);
				notify.error(err?.message || 'Network error during URL import.');
				onUploadError?.(err?.message || 'Network error during URL import.');
			} finally {
				isUploading = false;
			}
		}
	}
</script>

<div
	role="region"
	aria-label="Dropzone"
	class={`relative ${className}`}
	ondragenter={handleDragEnter}
	ondragover={handleDragOver}
	ondragleave={handleDragLeave}
	ondrop={handleDrop}
>
	{@render children?.()}

	<!-- Drag & Drop Visual Overlay -->
	{#if isDraggingOver || isUploading}
		<div
			class={`${
				fullScreenOverlay ? 'fixed inset-0 z-10000' : 'absolute inset-0 z-50'
			} bg-slate-950/80 backdrop-blur-xl border-4 border-dashed border-sky-400 flex flex-col items-center justify-center text-white space-y-4 animate-fade-in pointer-events-none p-6 text-center shadow-2xl rounded-2xl`}
		>
			{#if isUploading}
				<div class="w-20 h-20 rounded-3xl bg-sky-500/80 flex items-center justify-center shadow-xl shadow-sky-500/40 animate-spin">
					<Loader2 class="w-10 h-10 text-white" />
				</div>
				<h2 class="text-2xl font-extrabold tracking-tight">Uploading & Processing...</h2>
				<p class="text-xs text-sky-200 font-medium">Classifying media & saving to MinIO storage</p>
			{:else}
				<div class="w-20 h-20 rounded-3xl bg-sky-500/80 flex items-center justify-center shadow-xl shadow-sky-500/40 animate-bounce">
					<CloudUpload class="w-10 h-10 text-white" />
				</div>
				<h2 class="text-2xl font-extrabold tracking-tight">Drop files or web images anywhere to upload</h2>
				<p class="text-xs text-sky-200 font-medium">Automatic classification & uncompressed storage</p>
			{/if}
		</div>
	{/if}
</div>
