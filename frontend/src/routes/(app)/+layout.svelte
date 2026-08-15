<script lang="ts">
	import Sidebar from '$lib/components/Sidebar.svelte';
	import TopHeader from '$lib/components/TopHeader.svelte';

	let { children } = $props();
</script>

<div class="min-h-screen w-full flex bg-[#FAFAFA] dark:bg-[#090D16] text-slate-900 dark:text-slate-100 font-sans transition-colors duration-200">
	<!-- Collapsible Left Sidebar -->
	<Sidebar />

	<!-- Main Workspace Area -->
	<div class="flex-1 flex flex-col min-w-0 h-screen overflow-y-auto">
		<!-- Top Header -->
		<TopHeader />

		<!-- Page Content Container -->
		<main class="flex-1 p-6 md:p-8 max-w-7xl w-full mx-auto">
			{@render children?.()}
		</main>
	</div>
</div>
