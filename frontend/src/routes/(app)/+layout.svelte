<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { appState } from '$lib/state.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import TopHeader from '$lib/components/TopHeader.svelte';

	let { children } = $props();
	let ready = $state(false);

	onMount(() => {
		// Restore session from localStorage on every page load
		const token = localStorage.getItem('deepphotos_token');
		if (!token) {
			goto('/');
			return;
		}
		// Parse role from JWT payload (no library needed — just base64 decode)
		try {
			const payload = JSON.parse(atob(token.split('.')[1]));
			if (payload.exp && payload.exp * 1000 < Date.now()) {
				// Token expired
				localStorage.removeItem('deepphotos_token');
				goto('/');
				return;
			}
			// Sync role into appState from token
			if (payload.role) appState.user.role = payload.role;
			if (payload.email) appState.user.email = payload.email;
		} catch {
			// Malformed token
			localStorage.removeItem('deepphotos_token');
			goto('/');
			return;
		}
		appState.isAuthenticated = true;
		ready = true;
	});
</script>

{#if ready}
<div class="min-h-screen w-full flex bg-[#FAFAFA] dark:bg-[#090D16] text-slate-900 dark:text-slate-100 font-sans transition-colors duration-200">
	<!-- Collapsible Left Sidebar -->
	<Sidebar />

	<!-- Main Workspace Area -->
	<div class="flex-1 flex flex-col min-w-0 h-screen overflow-y-auto">
		<!-- Top Header -->
		<TopHeader />

		<!-- Page Content Container -->
		<main class="flex-1 overflow-hidden p-6 md:p-8 max-w-7xl w-full mx-auto flex flex-col">
			{@render children?.()}
		</main>
	</div>
</div>
{/if}
