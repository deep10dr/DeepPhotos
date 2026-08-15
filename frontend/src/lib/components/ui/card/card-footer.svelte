<script lang="ts">
	import { cn } from "$lib/utils";

	let { class: className, children, ...restProps }: { class?: string; children?: any; [key: string]: any } = $props();
</script>

<div class={cn("flex items-center p-6 pt-0 md:p-8 md:pt-0", className)} {...restProps}>
	{@render children?.()}
</div>
