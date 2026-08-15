<script lang="ts">
	import { cn } from "$lib/utils";

	let {
		class: className,
		type = "text",
		value = $bindable(""),
		placeholder = "",
		disabled = false,
		id,
		...restProps
	}: {
		class?: string;
		type?: string;
		value?: string;
		placeholder?: string;
		disabled?: boolean;
		id?: string;
		[key: string]: any;
	} = $props();
</script>

<input
	{id}
	{type}
	bind:value
	{placeholder}
	{disabled}
	class={cn(
		"flex h-11 w-full rounded-xl glass-input px-4 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 transition-all duration-200",
		className
	)}
	{...restProps}
/>
