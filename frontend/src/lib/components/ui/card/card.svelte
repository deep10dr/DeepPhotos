<script lang="ts">
	import { cn } from "$lib/utils";

	let { class: className, children, ...restProps }: { class?: string; children?: any; [key: string]: any } = $props();
</script>

<div
	class={cn("rounded-2xl glass-card text-slate-100 shadow-2xl relative overflow-hidden", className)}
	{...restProps}
>
	{@render children?.()}
</div>
