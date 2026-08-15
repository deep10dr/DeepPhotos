<script lang="ts">
	import { goto } from '$app/navigation';
	import { appState } from '$lib/state.svelte';
	import { Search, Upload, Bell, Sun, Moon } from 'lucide-svelte';

	let searchQuery = $state('');

	function goToProfile() {
		goto('/profile');
	}
</script>

<header class="h-16 border-b border-sky-100 dark:border-slate-800 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md px-6 flex items-center justify-between sticky top-0 z-20 transition-colors">
	
	<!-- Search Bar -->
	<div class="flex items-center gap-3 w-72 md:w-96">
		<div class="relative w-full">
			<Search class="w-4 h-4 text-slate-400 dark:text-slate-500 absolute left-3.5 top-1/2 -translate-y-1/2" />
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Search photos, albums, dates..."
				class="w-full h-9 pl-9 pr-4 text-xs rounded-xl bg-slate-100/80 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 placeholder:text-slate-400 focus:outline-none focus:border-sky-400 focus:bg-white dark:focus:bg-slate-900 focus:ring-2 focus:ring-sky-100 dark:focus:ring-sky-900 transition-all"
			/>
		</div>
	</div>

	<!-- Top Right Action Items & Theme Toggle -->
	<div class="flex items-center gap-3">
		
		<!-- Upload Button -->
		<button
			type="button"
			class="hidden sm:flex items-center gap-2 px-3.5 py-1.5 rounded-xl bg-sky-400 hover:bg-sky-500 text-white text-xs font-semibold shadow-sm shadow-sky-300/50 transition-all cursor-pointer"
		>
			<Upload class="w-3.5 h-3.5" />
			<span>Upload</span>
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
