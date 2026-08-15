<script lang="ts">
	import { cn } from "$lib/utils";

	let { class: className, children, ...restProps }: { class?: string; children?: any; [key: string]: any } = $props();
</script>

<h3 class={cn("text-2xl font-bold tracking-tight text-white", className)} {...restProps}>
	{@render children?.()}
</h3>
