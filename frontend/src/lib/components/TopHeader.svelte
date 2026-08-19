<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiFetch, invalidateCache } from '$lib/api';
	import { appState } from '$lib/state.svelte';
	import { Search, CloudUpload, Bell, Sun, Moon, Loader2, } from 'lucide-svelte';

	let searchQuery = $state('');
	let fileInput = $state<HTMLInputElement | null>(null);
	let isCloudUploading = $state(false);
	let isDraggingOver = $state(false);

	function goToProfile() {
		goto('/profile');
	}

	function triggerCloudUpload() {
		fileInput?.click();
	}

	async function CloudUploadFiles(filesList: FileList | File[]) {
		if (!filesList || filesList.length === 0) return;

		isCloudUploading = true;
		const formData = new FormData();
		for (let i = 0; i < filesList.length; i++) {
			formData.append('files', filesList[i]);
		}

		try {
			const res = await apiFetch('/api/photos', {
				method: 'POST',
				body: formData
			});
			if (res.ok) {
				invalidateCache('/api/photos');
				appState.refreshPhotos();
			} else {
				console.error('Failed to CloudUpload media:', await res.text());
			}
		} catch (err) {
			console.error('Network error during media CloudUpload:', err);
		} finally {
			isCloudUploading = false;
			if (fileInput) fileInput.value = '';
		}
	}

	async function handleExternalUrlDrop(url: string) {
		if (!url || !url.startsWith('http')) return;
		isCloudUploading = true;

		try {
			const res = await apiFetch('/api/photos/CloudUpload-url', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ url })
			});
			if (res.ok) {
				invalidateCache('/api/photos');
				appState.refreshPhotos();
			}
		} catch (err) {
			console.error('Failed to CloudUpload external image URL:', err);
		} finally {
			isCloudUploading = false;
		}
	}

	function handleFileSelect(e: Event) {
		const target = e.target as HTMLInputElement;
		if (target.files) {
			CloudUploadFiles(target.files);
		}
	}

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		isDraggingOver = true;
	}

	function handleDragLeave(e: DragEvent) {
		e.preventDefault();
		isDraggingOver = false;
	}

	async function handleDrop(e: DragEvent) {
		e.preventDefault();
		isDraggingOver = false;

		// 1. Local files drag & drop
		if (e.dataTransfer?.files && e.dataTransfer.files.length > 0) {
			CloudUploadFiles(e.dataTransfer.files);
			return;
		}

		// 2. Drag & drop image from another website or browser tab
		const htmlData = e.dataTransfer?.getData('text/html');
		let imageUrl = e.dataTransfer?.getData('URL') || e.dataTransfer?.getData('text/uri-list');

		if (!imageUrl && htmlData) {
			const match = htmlData.match(/src=["'](https?:\/\/[^"']+)["']/i);
			if (match) imageUrl = match[1];
		}

		if (imageUrl) {
			await handleExternalUrlDrop(imageUrl);
		}
	}
</script>

<svelte:window
	ondragover={handleDragOver}
	ondragleave={handleDragLeave}
	ondrop={handleDrop}
/>

<!-- Drag & Drop Overlay Indicator -->
{#if isDraggingOver}
	<div class="fixed inset-0 z-50 bg-sky-500/30 dark:bg-sky-950/60 backdrop-blur-md border-4 border-dashed border-sky-400 flex flex-col items-center justify-center text-white space-y-4 animate-fade-in pointer-events-none">
		<div class="w-20 h-20 rounded-3xl bg-sky-400/80 flex items-center justify-center shadow-xl shadow-sky-500/50 animate-bounce">
			<CloudUpload class="w-10 h-10 text-white" />
		</div>
		<h2 class="text-2xl font-extrabold tracking-tight">Drop files or web images anywhere to CloudUpload</h2>
		<p class="text-xs text-sky-100 font-medium">Automatic classification & MinIO storage</p>
	</div>
{/if}

<header class="h-16 border-b border-sky-100 dark:border-slate-800 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md px-6 flex items-center justify-between sticky top-0 z-20 transition-colors">

	<!-- Hidden File Input for Multi-Media Ingestion -->
	<input
		type="file"
		bind:this={fileInput}
		onchange={handleFileSelect}
		multiple
		accept="image/*,video/*,application/pdf,text/*,.doc,.docx,.zip"
		class="hidden"
	/>

	<!-- Search Bar -->
	<div class="flex items-center gap-3 w-72 md:w-96">
		<div class="relative w-full">
			<Search class="w-4 h-4 text-slate-400 dark:text-slate-500 absolute left-3.5 top-1/2 -translate-y-1/2" />
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Search photos, videos, documents..."
				class="w-full h-9 pl-9 pr-4 text-xs rounded-xl bg-slate-100/80 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 placeholder:text-slate-400 focus:outline-none focus:border-sky-400 focus:bg-white dark:focus:bg-slate-900 focus:ring-2 focus:ring-sky-100 dark:focus:ring-sky-900 transition-all"
			/>
		</div>
	</div>

	<!-- Top Right Action Items & Theme Toggle -->
	<div class="flex items-center gap-3">

		<!-- CloudUpload Button (Triggers Real File CloudUpload to Go API) -->
		<button
			type="button"
			onclick={triggerCloudUpload}
			disabled={isCloudUploading}
			class="hidden sm:flex items-center gap-2 px-3.5 py-1.5 rounded-xl text-xs font-semibold shadow-sm shadow-sky-300/50 transition-all cursor-pointer disabled:opacity-50"
		>
			{#if isCloudUploading}
				<Loader2 class="w-3.5 h-3.5 animate-spin  text-sky-400 hover:text-sky-500 " />
			{:else}
				<CloudUpload class="w-3.5 h-3.5  text-sky-400 hover:text-sky-500 " />
			{/if}
		</button>

		<!-- DARK / LIGHT MODE TOGGLE BUTTON -->
		<button
			type="button"
			onclick={() => appState.toggleTheme()}
			title={appState.theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'}
			class="p-2 rounded-xl text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition-all cursor-pointer"
		>
			{#if appState.theme === 'dark'}
				<Sun class="w-4 h-4 text-amber-400 animate-fade-in" />
			{:else}
				<Moon class="w-4 h-4 text-slate-600 animate-fade-in" />
			{/if}
		</button>

		<!-- Notifications -->
		<button
			type="button"
			class="p-2 rounded-xl text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors relative cursor-pointer"
			title="Notifications"
		>
			<Bell class="w-4 h-4" />
			<span class="w-2 h-2 rounded-full bg-sky-400 absolute top-2 right-2 ring-2 ring-white dark:ring-slate-900"></span>
		</button>

		<!-- TOP RIGHT CORNER: User Profile Avatar -->
		<button
			type="button"
			onclick={goToProfile}
			title={`User Profile: ${appState.user.name} (${appState.user.role})`}
			class="relative p-0.5 rounded-full hover:ring-4 hover:ring-sky-200 dark:hover:ring-sky-900 transition-all cursor-pointer group"
		>
			<img
				src={appState.user.avatar}
				alt={appState.user.name}
				class="w-9 h-9 rounded-full object-cover ring-2 ring-sky-400 group-hover:scale-105 transition-transform"
			/>
			<span class="w-2.5 h-2.5 rounded-full bg-emerald-400 absolute bottom-0 right-0 ring-2 ring-white dark:ring-slate-900"></span>
		</button>

	</div>

</header>
