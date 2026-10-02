<script lang="ts">
	import { CloudUpload, Loader2 } from 'lucide-svelte';
	import { apiFetch, invalidateCache } from '$lib/api';
	import { appState } from '$lib/state.svelte';
	import { notify } from '$lib/notify.svelte';

	interface Props {
		label?: string;
		accept?: string;
		multiple?: boolean;
		variant?: 'primary' | 'secondary' | 'iconOnly';
		class?: string;
		onUploadStart?: () => void;
		onUploadComplete?: () => void;
		onUploadError?: (error: string) => void;
	}

	let {
		label,
		accept = 'image/*,video/*,application/pdf,text/*,.doc,.docx,.zip',
		multiple = true,
		variant = 'primary',
		class: className = '',
		onUploadStart,
		onUploadComplete,
		onUploadError
	}: Props = $props();

	let fileInput = $state<HTMLInputElement | null>(null);
	let isUploading = $state(false);

	function triggerUpload() {
		fileInput?.click();
	}

	async function handleFileSelect(e: Event) {
		const target = e.target as HTMLInputElement;
		const files = target.files;
		if (!files || files.length === 0) return;

		isUploading = true;
		onUploadStart?.();

		const formData = new FormData();
		for (let i = 0; i < files.length; i++) {
			formData.append('files', files[i]);
		}

		try {
			const res = await apiFetch('/api/media/upload', {
				method: 'POST',
				body: formData
			});

			if (res.ok) {
				invalidateCache('/api/media');
				invalidateCache('/api/gallery');
				appState.refreshPhotos();
				notify.success('Files uploaded successfully!');
				onUploadComplete?.();
			} else {
				const errText = await res.text();
				notify.error(errText || 'Failed to upload selected files.');
				onUploadError?.(errText || 'Failed to upload selected files.');
			}
		} catch (err: any) {
			console.error('Upload error:', err);
			notify.error(err?.message || 'Network error during file upload.');
			onUploadError?.(err?.message || 'Network error during file upload.');
		} finally {
			isUploading = false;
			if (fileInput) fileInput.value = '';
		}
	}
</script>

<!-- Hidden Native File Input -->
<input
	type="file"
	bind:this={fileInput}
	onchange={handleFileSelect}
	{multiple}
	{accept}
	class="hidden"
/>

<!-- Upload Action Button -->
{#if variant === 'iconOnly'}
	<button
		type="button"
		onclick={triggerUpload}
		disabled={isUploading}
		title={label || 'Upload Media'}
		class={`p-2 rounded-xl text-slate-500 dark:text-slate-400 hover:text-sky-500 hover:bg-slate-100 dark:hover:bg-slate-800 transition-all cursor-pointer disabled:opacity-50 ${className}`}
	>
		{#if isUploading}
			<Loader2 class="w-4 h-4 animate-spin text-sky-500" />
		{:else}
			<CloudUpload class="w-4 h-4 text-sky-500" />
		{/if}
	</button>
{:else if variant === 'secondary'}
	<button
		type="button"
		onclick={triggerUpload}
		disabled={isUploading}
		class={`px-3.5 py-1.5 rounded-xl text-xs font-semibold bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-200 hover:bg-slate-200 dark:hover:bg-slate-700 transition-all flex items-center gap-1.5 cursor-pointer disabled:opacity-50 ${className}`}
	>
		{#if isUploading}
			<Loader2 class="w-3.5 h-3.5 animate-spin text-sky-500" />
			<span>Uploading...</span>
		{:else}
			<CloudUpload class="w-3.5 h-3.5 text-sky-500" />
			<span>{label || 'Upload'}</span>
		{/if}
	</button>
{:else}
	<!-- Primary Variant (Default) -->
	<button
		type="button"
		onclick={triggerUpload}
		disabled={isUploading}
		class={`px-4 py-2 rounded-xl bg-sky-400 text-white text-xs font-bold flex items-center gap-1.5 shadow-sm shadow-sky-300/50 hover:bg-sky-500 transition-all cursor-pointer disabled:opacity-50 ${className}`}
	>
		{#if isUploading}
			<Loader2 class="w-4 h-4 animate-spin" />
			<span>Uploading...</span>
		{:else}
			<CloudUpload class="w-4 h-4" />
			<span>{label || 'Upload'}</span>
		{/if}
	</button>
{/if}
