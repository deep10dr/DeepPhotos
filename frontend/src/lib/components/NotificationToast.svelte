<script lang="ts">
	import { notify, confirmDialog } from '$lib/notify.svelte';
	import { CheckCircle2, AlertCircle, AlertTriangle, Info, X, HelpCircle } from 'lucide-svelte';
</script>

<!-- Global Toast Container (Fixed Top-Right) -->
<div class="fixed top-5 right-5 z-999999 flex flex-col gap-3 max-w-sm w-full pointer-events-none px-4 sm:px-0">
	{#each notify.toasts as toast (toast.id)}
		<div
			class={`pointer-events-auto rounded-2xl p-4 shadow-xl border backdrop-blur-2xl transition-all duration-300 animate-fade-in flex items-start justify-between gap-3 ${
				toast.type === 'success'
					? 'bg-emerald-50 dark:bg-emerald-950/90 text-emerald-900 dark:text-emerald-100 border-emerald-200 dark:border-emerald-500/40 shadow-emerald-500/10'
					: toast.type === 'error'
					? 'bg-rose-50 dark:bg-rose-950/90 text-rose-900 dark:text-rose-100 border-rose-200 dark:border-rose-500/40 shadow-rose-500/10'
					: toast.type === 'warning'
					? 'bg-amber-50 dark:bg-amber-950/90 text-amber-900 dark:text-amber-100 border-amber-200 dark:border-amber-500/40 shadow-amber-500/10'
					: 'bg-white dark:bg-slate-900/90 text-slate-900 dark:text-slate-100 border-sky-200 dark:border-sky-500/40 shadow-sky-500/10'
			}`}
		>
			<div class="flex items-start gap-3 min-w-0">
				<!-- Icon Badge -->
				<div class="mt-0.5 shrink-0">
					{#if toast.type === 'success'}
						<CheckCircle2 class="w-5 h-5 text-emerald-400" />
					{:else if toast.type === 'error'}
						<AlertCircle class="w-5 h-5 text-rose-400" />
					{:else if toast.type === 'warning'}
						<AlertTriangle class="w-5 h-5 text-amber-400" />
					{:else}
						<Info class="w-5 h-5 text-sky-400" />
					{/if}
				</div>

				<!-- Content -->
				<div class="min-w-0 space-y-0.5">
					{#if toast.title}
						<h4 class="text-xs font-bold tracking-tight">{toast.title}</h4>
					{/if}
					<p class="text-xs font-medium text-slate-200 leading-relaxed break-words">{toast.message}</p>
				</div>
			</div>

			<!-- Dismiss Button -->
			<button
				type="button"
				onclick={() => notify.remove(toast.id)}
				class="p-1 rounded-lg text-slate-400 hover:text-white hover:bg-white/10 transition-colors shrink-0 cursor-pointer"
				title="Dismiss"
			>
				<X class="w-4 h-4" />
			</button>
		</div>
	{/each}
</div>

<!-- Global Interactive Yes / No Confirmation Dialog Modal -->
{#if confirmDialog.active}
	<div class="fixed inset-0 z-[999999] bg-slate-950/70 backdrop-blur-sm flex items-center justify-center p-4 animate-fade-in">
		<div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-2xl max-w-sm w-full space-y-4 animate-fade-in">
			<div class="flex items-start gap-3.5">
				<div class={`p-3 rounded-2xl shrink-0 ${
					confirmDialog.options.type === 'danger'
						? 'bg-rose-100 dark:bg-rose-950/80 text-rose-500'
						: confirmDialog.options.type === 'warning'
						? 'bg-amber-100 dark:bg-amber-950/80 text-amber-500'
						: 'bg-sky-100 dark:bg-sky-950/80 text-sky-500'
				}`}>
					{#if confirmDialog.options.type === 'danger'}
						<AlertCircle class="w-6 h-6" />
					{:else if confirmDialog.options.type === 'warning'}
						<AlertTriangle class="w-6 h-6" />
					{:else}
						<HelpCircle class="w-6 h-6" />
					{/if}
				</div>

				<div class="space-y-1 min-w-0">
					<h3 class="text-base font-bold text-slate-900 dark:text-white">
						{confirmDialog.options.title || 'Confirm Action'}
					</h3>
					<p class="text-xs text-slate-600 dark:text-slate-300 leading-relaxed">
						{confirmDialog.options.message}
					</p>
				</div>
			</div>

			<!-- Yes / No Action Buttons -->
			<div class="flex items-center justify-end gap-2.5 pt-2 border-t border-slate-100 dark:border-slate-800">
				<button
					type="button"
					onclick={() => confirmDialog.cancel()}
					class="px-4 py-2 rounded-xl text-xs font-semibold text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
				>
					{confirmDialog.options.cancelText || 'No'}
				</button>
				<button
					type="button"
					onclick={() => confirmDialog.confirm()}
					class={`px-5 py-2 rounded-xl text-xs font-bold text-white shadow-sm transition-all cursor-pointer ${
						confirmDialog.options.type === 'danger'
							? 'bg-rose-500 hover:bg-rose-600 shadow-rose-300/50 dark:shadow-none'
							: confirmDialog.options.type === 'warning'
							? 'bg-amber-500 hover:bg-amber-600 shadow-amber-300/50 dark:shadow-none'
							: 'bg-sky-400 hover:bg-sky-500 shadow-sky-300/50 dark:shadow-none'
					}`}
				>
					{confirmDialog.options.confirmText || 'Yes'}
				</button>
			</div>
		</div>
	</div>
{/if}
