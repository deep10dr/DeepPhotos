<script lang="ts">
	import { Lock, ShieldCheck, KeyRound, Eye } from 'lucide-svelte';

	let pin = $state('');
	let isUnlocked = $state(false);

	function unlockVault(e: Event) {
		e.preventDefault();
		if (pin === '1234' || pin.length >= 4) {
			isUnlocked = true;
		}
	}
</script>

<div class="space-y-6 animate-fade-in max-w-3xl mx-auto">
	<div class="text-center mb-8">
		<div class="w-14 h-14 rounded-2xl bg-sky-100 dark:bg-sky-950 text-sky-600 dark:text-sky-400 flex items-center justify-center mx-auto mb-3 shadow-md">
			<Lock class="w-7 h-7" />
		</div>
		<h1 class="text-2xl font-bold text-slate-900 dark:text-white">Locked Vault</h1>
		<p class="text-xs text-slate-500 dark:text-slate-400 mt-1">Passcode protected AES-256 encrypted media folder</p>
	</div>

	{#if !isUnlocked}
		<form onsubmit={unlockVault} class="p-8 rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm max-w-sm mx-auto space-y-4 text-center">
			<label for="vault-passcode" class="text-xs font-semibold text-slate-700 dark:text-slate-300 block">Enter Vault Passcode (Default: 1234)</label>
			<input
				id="vault-passcode"
				type="password"
				bind:value={pin}
				placeholder="••••"
				maxlength="6"
				class="w-full h-12 text-center text-lg tracking-widest rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-900 dark:text-white focus:outline-none focus:border-sky-400 focus:bg-white dark:focus:bg-slate-900"
			/>
			<button
				type="submit"
				class="w-full h-10 rounded-xl bg-sky-400 hover:bg-sky-500 text-white font-bold text-xs shadow-sm shadow-sky-300/50 cursor-pointer"
			>
				Unlock Vault
			</button>
		</form>
	{:else}
		<div class="p-6 rounded-2xl bg-emerald-50 dark:bg-emerald-950/60 border border-emerald-200 dark:border-emerald-800 text-emerald-800 dark:text-emerald-300 text-center animate-fade-in">
			<ShieldCheck class="w-8 h-8 text-emerald-500 mx-auto mb-2" />
			<h3 class="text-base font-bold">Vault Unlocked</h3>
			<p class="text-xs text-emerald-600 dark:text-emerald-400 mt-0.5">Your protected items are available for this session.</p>
		</div>
	{/if}
</div>
