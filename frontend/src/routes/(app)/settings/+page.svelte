<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	onMount(() => {
		goto('/profile');
	});
</script>

<div class="p-8 text-center text-slate-500">Redirecting to Settings & Profile...</div>
