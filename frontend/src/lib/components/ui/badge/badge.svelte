<script lang="ts">
	import { cn } from "$lib/utils";

	let {
		class: className,
		variant = "default",
		children,
		...restProps
	}: {
		class?: string;
		variant?: "default" | "secondary" | "outline" | "glow";
		children?: any;
		[key: string]: any;
	} = $props();

	const variantStyles = {
		default: "bg-blue-500/10 text-blue-400 border-blue-500/20",
		secondary: "bg-slate-800 text-slate-300 border-slate-700",
		outline: "border-slate-700 text-slate-400",
		glow: "bg-gradient-to-r from-blue-500/15 to-indigo-500/15 text-blue-300 border-blue-500/30 shadow-sm shadow-blue-500/20"
	};
</script>

<div
	class={cn(
		"inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
		variantStyles[variant],
		className
	)}
	{...restProps}
>
	{@render children?.()}
</div>
