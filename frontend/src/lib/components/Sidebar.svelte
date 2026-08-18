<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { appState } from '$lib/state.svelte';
	import {
		Camera,
		Image as ImageIcon,
		Sparkles,
		FolderClosed,
		Lock,
		FileText,
		Trash2,
		Settings,
		ChevronLeft,
		ChevronRight,
		LogOut
	} from 'lucide-svelte';
    import { Cloud } from '@lucide/svelte';

	const navItems = [
		{ name: 'Gallery', path: '/gallery', icon: ImageIcon },
		{ name: 'Memories', path: '/memories', icon: Sparkles },
		{ name: 'Albums', path: '/albums', icon: FolderClosed },
		{ name: 'Vault', path: '/vault', icon: Lock },
		{ name: 'Documents', path: '/documents', icon: FileText },
		{ name: 'Bin', path: '/bin', icon: Trash2 }
	];

	let isSettingsActive = $derived(page.url.pathname === '/profile' || page.url.pathname === '/settings');

	function handleLogout() {
		appState.logout();
		goto('/');
	}
</script>

<aside
	class={`h-screen bg-white/90 dark:bg-slate-900/95 border-r border-sky-100 dark:border-slate-800 flex flex-col justify-between transition-all duration-300 relative z-30 shadow-sm ${
		appState.isSidebarCollapsed ? 'w-20' : 'w-64'
	}`}
>
	<!-- Top Section: App Branding & Sidebar Minimize/Maximize Toggle -->
	<div>
		<div class="h-16 px-4 flex items-center justify-between border-b border-sky-100/80 dark:border-slate-800">
			{#if !appState.isSidebarCollapsed}
				<a href="/gallery" class="flex items-center gap-1.5 group">
					<Cloud class="w-7 h-7 text-sky-500 float-animation drop-shadow-sm group-hover:scale-110 transition-transform" fill="currentColor" />
					<div>
						<h1 class="text-[17px] font-extrabold tracking-tight text-slate-900 dark:text-white flex items-baseline">
							Dee<span class="text-2xl bg-clip-text text-transparent bg-gradient-to-r from-sky-500 to-cyan-300 dark:from-sky-300 dark:to-white -mx-[1px]">P</span><span class="bg-clip-text text-transparent bg-gradient-to-r from-sky-500 via-blue-400 to-cyan-300 dark:from-sky-300 dark:via-blue-200 dark:to-white">hotos</span>
						</h1>
					</div>
				</a>
			{:else}
				<a href="/gallery" class="mx-auto group">
					<Cloud class="w-7 h-7 text-sky-500 float-animation drop-shadow-sm group-hover:scale-110 transition-transform" fill="currentColor" />
				</a>
			{/if}

			<!-- Minimize / Maximize Toggle Button -->
			<button
				type="button"
				onclick={() => appState.toggleSidebar()}
				title={appState.isSidebarCollapsed ? 'Maximize Sidebar' : 'Minimize Sidebar'}
				class={`p-1.5 rounded-lg text-slate-400 dark:text-slate-500 hover:text-slate-700 dark:hover:text-slate-200 hover:bg-sky-50 dark:hover:bg-slate-800 transition-colors cursor-pointer ${
					appState.isSidebarCollapsed ? 'mx-auto mt-2' : ''
				}`}
			>
				{#if appState.isSidebarCollapsed}
					<ChevronRight class="w-4 h-4 text-sky-500" />
				{:else}
					<ChevronLeft class="w-4 h-4" />
				{/if}
			</button>
		</div>

		<!-- Navigation Menu -->
		<nav class="p-3 space-y-1.5 mt-2">
			{#each navItems as item}
				{@const isActive = page.url.pathname === item.path || (item.path === '/gallery' && page.url.pathname === '/')}
				<a
					href={item.path}
					title={appState.isSidebarCollapsed ? item.name : undefined}
					class={`flex items-center gap-3 px-3.5 py-2.5 rounded-xl font-semibold text-xs transition-all duration-200 group ${
						isActive
							? 'bg-sky-400 text-white shadow-md shadow-sky-300/50 dark:shadow-sky-900/50'
							: 'text-slate-600 dark:text-slate-300 hover:bg-sky-50 dark:hover:bg-slate-800 hover:text-sky-600 dark:hover:text-sky-400'
					} ${appState.isSidebarCollapsed ? 'justify-center px-0' : ''}`}
				>
					<item.icon class={`w-4 h-4 shrink-0 ${isActive ? 'text-white' : 'text-slate-400 dark:text-slate-500 group-hover:text-sky-500'}`} />

					{#if !appState.isSidebarCollapsed}
						<span>{item.name}</span>
					{/if}
				</a>
			{/each}
		</nav>
	</div>

	<!-- Bottom Section: Settings & Logout -->
	<div class="p-3 border-t border-sky-100/80 dark:border-slate-800 space-y-1.5">
		<a
			href="/profile"
			title={appState.isSidebarCollapsed ? 'Settings' : undefined}
			class={`flex items-center gap-3 px-3.5 py-2.5 rounded-xl font-semibold text-xs transition-all duration-200 group ${
				isSettingsActive
					? 'bg-sky-400 text-white shadow-md shadow-sky-300/50 dark:shadow-sky-900/50'
					: 'text-slate-600 dark:text-slate-300 hover:bg-sky-50 dark:hover:bg-slate-800 hover:text-sky-600 dark:hover:text-sky-400'
			} ${appState.isSidebarCollapsed ? 'justify-center px-0' : ''}`}
		>
			<Settings class={`w-4 h-4 shrink-0 ${isSettingsActive ? 'text-white' : 'text-slate-400 dark:text-slate-500 group-hover:text-sky-500'}`} />
			{#if !appState.isSidebarCollapsed}
				<span>Settings</span>
			{/if}
		</a>

		<button
			type="button"
			onclick={handleLogout}
			title={appState.isSidebarCollapsed ? 'Log Out' : undefined}
			class={`w-full flex items-center gap-3 px-3.5 py-2.5 rounded-xl font-semibold text-xs text-rose-600 dark:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/40 transition-all duration-200 cursor-pointer ${
				appState.isSidebarCollapsed ? 'justify-center px-0' : ''
			}`}
		>
			<LogOut class="w-4 h-4 shrink-0" />
			{#if !appState.isSidebarCollapsed}
				<span>Log Out</span>
			{/if}
		</button>
	</div>
</aside>

<style>
	.float-animation {
		animation: float-drift 6s ease-in-out infinite;
	}
	@keyframes float-drift {
		0%, 100% { transform: translateY(0px) scale(1); }
		50% { transform: translateY(-3px) scale(1.02); }
	}
</style>
