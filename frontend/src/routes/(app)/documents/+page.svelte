<script lang="ts">
	import { FileText, Download, File, HardDrive } from 'lucide-svelte';

	const docs = [
		{ name: 'Passport Scan.pdf', size: '1.4 MB', date: 'August 01, 2026' },
		{ name: 'Insurance Policy.pdf', size: '3.2 MB', date: 'July 24, 2026' },
		{ name: 'Property Receipt.pdf', size: '890 KB', date: 'June 18, 2026' }
	];
</script>

<div class="space-y-6 animate-fade-in">
	<div>
		<h1 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
			<FileText class="w-6 h-6 text-sky-500" />
			Documents & Scans
		</h1>
		<p class="text-xs text-slate-500 dark:text-slate-400">Stored document scans and PDF files</p>
	</div>

	<div class="p-6 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm space-y-3">
		{#each docs as doc}
			<div class="p-3.5 rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700 flex items-center justify-between hover:bg-sky-50/50 dark:hover:bg-slate-700/50 transition-colors">
				<div class="flex items-center gap-3">
					<div class="w-10 h-10 rounded-lg bg-sky-100 dark:bg-sky-950 text-sky-600 dark:text-sky-400 flex items-center justify-center">
						<File class="w-5 h-5" />
					</div>
					<div>
						<p class="text-xs font-bold text-slate-800 dark:text-white">{doc.name}</p>
						<p class="text-[11px] text-slate-400">{doc.date} • {doc.size}</p>
					</div>
				</div>
				<button type="button" class="p-2 rounded-lg border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 hover:text-sky-600 dark:hover:text-sky-400 hover:border-sky-300 dark:hover:border-sky-700 text-xs flex items-center gap-1 cursor-pointer">
					<Download class="w-4 h-4" /> Download
				</button>
			</div>
		{/each}
	</div>
</div>
