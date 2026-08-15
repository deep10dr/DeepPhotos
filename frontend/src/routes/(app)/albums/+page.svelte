<script lang="ts">
	import { FolderClosed, Plus, Image as ImageIcon } from 'lucide-svelte';

	const albums = [
		{ name: 'Vacation & Travel', count: 184, cover: 'https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=600&auto=format&fit=crop&q=80' },
		{ name: 'Family & Friends', count: 320, cover: 'https://images.unsplash.com/photo-1511895426328-dc8714191300?w=600&auto=format&fit=crop&q=80' },
		{ name: 'Nature & Landscapes', count: 96, cover: 'https://images.unsplash.com/photo-1448375240586-882707db888b?w=600&auto=format&fit=crop&q=80' }
	];
</script>

<div class="space-y-6 animate-fade-in">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
				<FolderClosed class="w-6 h-6 text-sky-500" />
				Albums
			</h1>
			<p class="text-xs text-slate-500 dark:text-slate-400">Organized photo collections</p>
		</div>
		<button type="button" class="px-4 py-2 rounded-xl bg-sky-400 text-white text-xs font-bold flex items-center gap-1.5 shadow-sm shadow-sky-300/50 hover:bg-sky-500 transition-all cursor-pointer">
			<Plus class="w-4 h-4" />
			New Album
		</button>
	</div>

	<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
		{#each albums as album}
			<div class="p-4 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:shadow-md transition-all group cursor-pointer">
				<div class="aspect-4/3 rounded-xl overflow-hidden bg-slate-100 dark:bg-slate-800 mb-3">
					<img src={album.cover} alt={album.name} class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
				</div>
				<h3 class="text-sm font-bold text-slate-900 dark:text-white">{album.name}</h3>
				<p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">{album.count} photos</p>
			</div>
		{/each}
	</div>
</div>
