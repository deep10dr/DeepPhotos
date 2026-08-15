<script lang="ts">
	import type { HTMLButtonAttributes } from "svelte/elements";
	import { cn } from "$lib/utils";

	type Variant = "default" | "destructive" | "outline" | "secondary" | "ghost" | "link" | "glow" | "sky";
	type Size = "default" | "sm" | "lg" | "icon";

	let {
		class: className,
		variant = "sky",
		size = "default",
		type = "button",
		children,
		disabled = false,
		...restProps
	}: {
		class?: string;
		variant?: Variant;
		size?: Size;
		type?: "button" | "submit" | "reset";
		children?: any;
		disabled?: boolean;
		[key: string]: any;
	} = $props();

	const variantStyles: Record<Variant, string> = {
		default: "bg-sky-500 text-slate-950 font-semibold hover:bg-sky-400 shadow-lg shadow-sky-500/25 active:scale-[0.98]",
		sky: "bg-gradient-to-r from-sky-400 via-sky-500 to-cyan-500 text-slate-950 font-bold shadow-lg shadow-sky-500/30 hover:shadow-sky-400/50 hover:brightness-110 active:scale-[0.98]",
		glow: "bg-gradient-to-r from-sky-500 via-blue-600 to-cyan-400 text-white font-semibold shadow-lg shadow-sky-500/35 hover:shadow-sky-400/50 hover:brightness-110 active:scale-[0.98]",
		destructive: "bg-red-600 text-white hover:bg-red-500 shadow-sm active:scale-[0.98]",
		outline: "border border-sky-500/30 bg-sky-950/40 text-sky-200 hover:bg-sky-900/50 hover:border-sky-400/60 hover:text-white backdrop-blur-md active:scale-[0.98]",
		secondary: "bg-sky-950 text-sky-100 border border-sky-800/60 hover:bg-sky-900/80 shadow-sm active:scale-[0.98]",
		ghost: "text-sky-300 hover:bg-sky-900/50 hover:text-white active:scale-[0.98]",
		link: "text-sky-400 underline-offset-4 hover:underline"
	};

	const sizeStyles: Record<Size, string> = {
		default: "h-11 px-5 py-2.5 text-sm",
		sm: "h-9 rounded-md px-3 text-xs",
		lg: "h-12 rounded-lg px-8 text-base",
		icon: "h-10 w-10 p-0"
	};
</script>

<button
	{type}
	{disabled}
	class={cn(
		"inline-flex items-center justify-center rounded-xl font-medium transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-400/80 focus-visible:ring-offset-2 focus-visible:ring-offset-sky-950 disabled:pointer-events-none disabled:opacity-50 select-none cursor-pointer",
		variantStyles[variant],
		sizeStyles[size],
		className
	)}
	{...restProps}
>
	{@render children?.()}
</button>
