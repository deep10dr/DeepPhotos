<script lang="ts">
	import { Trash2, RotateCcw, Image as ImageIcon } from 'lucide-svelte';

	const binnedItems = [
		{ name: 'Blurry Shot 001.jpg', size: '2.1 MB', deletedDate: '2 days ago' },
		{ name: 'Duplicate Capture.jpg', size: '4.5 MB', deletedDate: '5 days ago' }
	];
</script>

<div class="space-y-6 animate-fade-in">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
				<Trash2 class="w-6 h-6 text-rose-500" />
				Bin / Trash
			</h1>
			<p class="text-xs text-slate-500 dark:text-slate-400">Deleted items are permanently removed after 30 days</p>
		</div>
		<button type="button" class="px-4 py-2 rounded-xl bg-rose-50 dark:bg-rose-950/60 text-rose-600 dark:text-rose-400 hover:bg-rose-100 dark:hover:bg-rose-900 text-xs font-bold transition-all cursor-pointer border border-transparent dark:border-rose-800">
			Empty Bin
		</button>
	</div>

	<div class="p-6 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm space-y-3">
		{#each binnedItems as item}
			<div class="p-3.5 rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700 flex items-center justify-between">
				<div class="flex items-center gap-3">
					<div class="w-10 h-10 rounded-lg bg-rose-100 dark:bg-rose-950 text-rose-600 dark:text-rose-400 flex items-center justify-center">
						<ImageIcon class="w-5 h-5" />
					</div>
					<div>
						<p class="text-xs font-bold text-slate-800 dark:text-white">{item.name}</p>
						<p class="text-[11px] text-slate-400">Deleted {item.deletedDate} • {item.size}</p>
					</div>
				</div>
				<button type="button" class="px-3 py-1.5 rounded-lg border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 hover:text-sky-600 dark:hover:text-sky-400 hover:border-sky-300 dark:hover:border-sky-700 text-xs font-semibold flex items-center gap-1 cursor-pointer">
					<RotateCcw class="w-3.5 h-3.5" /> Restore
				</button>
			</div>
		{/each}
	</div>
</div>
