<script lang="ts">
	import { Sparkles, Calendar, Heart, MapPin } from 'lucide-svelte';

	const memoryHighlights = [
		{ title: 'Summer Trip 2026', subtitle: '1 year ago today', location: 'California, USA', itemsCount: 42, cover: 'https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=600&auto=format&fit=crop&q=80' },
		{ title: 'Mountain Hiking Weekend', subtitle: '2 years ago today', location: 'Alps Region', itemsCount: 28, cover: 'https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=600&auto=format&fit=crop&q=80' }
	];
</script>

<div class="h-full flex flex-col gap-6 animate-fade-in">
	<div class="shrink-0">
		<h1 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
			<Sparkles class="w-6 h-6 text-sky-500" />
			Memories & Highlights
		</h1>
		<p class="text-xs text-slate-500 dark:text-slate-400">Rediscover special moments from past years</p>
	</div>

	<div class="flex-1 overflow-y-auto min-h-0">
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			{#each memoryHighlights as memory}
				<div class="rounded-3xl overflow-hidden bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:shadow-md transition-all group cursor-pointer">
					<div class="h-48 overflow-hidden relative">
						<img src={memory.cover} alt={memory.title} class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
						<div class="absolute inset-0 bg-gradient-to-t from-slate-950/70 via-slate-950/20 to-transparent p-6 flex flex-col justify-end text-white">
							<span class="text-[11px] font-semibold text-sky-300 flex items-center gap-1">
								<Calendar class="w-3.5 h-3.5" /> {memory.subtitle}
							</span>
							<h3 class="text-lg font-bold mt-1">{memory.title}</h3>
							<p class="text-xs text-slate-300 flex items-center gap-1 mt-0.5">
								<MapPin class="w-3 h-3 text-sky-400" /> {memory.location} • {memory.itemsCount} photos
							</p>
						</div>
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>
