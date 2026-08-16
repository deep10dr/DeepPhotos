<script lang="ts">
	import { appState } from '$lib/state.svelte';
	import { FileText, Download, Trash2, ExternalLink } from 'lucide-svelte';

	interface DocumentItem {
		id: string;
		title: string;
		filename: string;
		object_key: string;
		mime_type: string;
		size: number;
		taken_at: string;
		url?: string;
	}

	let docs = $state<DocumentItem[]>([]);
	let isLoading = $state(true);

	async function fetchDocuments() {
		isLoading = true;
		try {
			const res = await fetch(`${appState.apiBaseUrl}/api/photos?type=document&deleted=false`);
			if (res.ok) {
				const data: DocumentItem[] = await res.json();
				docs = data.map(d => ({
					...d,
					url: `${appState.apiBaseUrl}/api/photos/${d.id}/file`
				}));
			}
		} catch (e) {
			console.warn('API error fetching documents:', e);
		} finally {
			isLoading = false;
		}
	}

	$effect(() => {
		const v = appState.uploadVersion;
		fetchDocuments();
	});

	function formatBytes(bytes: number): string {
		if (!bytes) return '0 B';
		const k = 1024;
		const sizes = ['B', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
	}
</script>

<div class="space-y-6 animate-fade-in">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
				<FileText class="w-6 h-6 text-sky-500" />
				Documents Vault
			</h1>
			<p class="text-xs text-slate-500 dark:text-slate-400">PDFs, text files, and scanned documents (`document/uuid/...`)</p>
		</div>
	</div>

	<!-- EMPTY STATE WITH EMOJI -->
	{#if !isLoading && docs.length === 0}
		<div class="p-12 text-center rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm max-w-md mx-auto my-12 space-y-3">
			<div class="text-5xl mb-2">📄✨</div>
			<h3 class="text-base font-bold text-slate-900 dark:text-white">This documents vault is empty</h3>
			<p class="text-xs text-slate-500 dark:text-slate-400">Drag & drop PDFs or text files anywhere into the app to save them in MinIO!</p>
		</div>
	{:else}
		<div class="p-6 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm space-y-3">
			{#each docs as doc}
				<div class="p-4 rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700 flex items-center justify-between gap-4">
					<div class="flex items-center gap-3">
						<div class="w-10 h-10 rounded-lg bg-sky-100 dark:bg-sky-950 text-sky-600 dark:text-sky-400 flex items-center justify-center">
							<FileText class="w-5 h-5" />
						</div>
						<div>
							<p class="text-xs font-bold text-slate-800 dark:text-white">{doc.title || doc.filename}</p>
							<p class="text-[11px] text-slate-500 dark:text-slate-400">{doc.mime_type} • {formatBytes(doc.size)}</p>
						</div>
					</div>

					<div class="flex items-center gap-2">
						<a
							href={doc.url}
							target="_blank"
							download={doc.filename}
							class="px-3.5 py-1.5 rounded-lg border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-200 hover:bg-white dark:hover:bg-slate-900 text-xs font-semibold flex items-center gap-1.5 transition-colors cursor-pointer"
						>
							<Download class="w-3.5 h-3.5" /> Download
						</a>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
