<script lang="ts">
	import { cn } from "$lib/utils";

	let { class: className, children, ...restProps }: { class?: string; children?: any; [key: string]: any } = $props();
</script>

<div class={cn("flex flex-col space-y-1.5 p-6 md:p-8", className)} {...restProps}>
	{@render children?.()}
</div>
