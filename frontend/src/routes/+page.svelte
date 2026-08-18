<script lang="ts">
	import { goto } from '$app/navigation';
	import { appState } from '$lib/state.svelte';
	import { Lock, User, Eye, EyeOff, ArrowRight, Loader2, Image as ImageIcon, Cloud, Sparkles, FolderLock } from 'lucide-svelte';

	// Pre-fill from env for development by default
	const devUsername = import.meta.env.VITE_DEV_USERNAME || 'admin@deepphotos.local';
	const devPassword = import.meta.env.VITE_DEV_PASSWORD || 'deepphotos2026';

	// Svelte 5 state runes
	let username = $state(devUsername);
	let password = $state(devPassword);
	let showPassword = $state(false);
	let isLoading = $state(false);

	// Minimal clean photo cards for left side display
	const features = [
		{ id: 1, title: 'Smart Gallery', desc: 'Timeline-based photo browsing', icon: ImageIcon, color: 'text-sky-500', bg: 'bg-sky-100 dark:bg-sky-950/50' },
		{ id: 2, title: 'Secure Vault', desc: 'Passcode-protected hidden folders', icon: FolderLock, color: 'text-sky-600', bg: 'bg-sky-200/50 dark:bg-sky-900/40' },
		{ id: 3, title: 'Original Quality', desc: 'Raw, uncompressed MinIO storage', icon: Sparkles, color: 'text-cyan-500', bg: 'bg-cyan-100 dark:bg-cyan-950/50' }
	];

	async function handleLogin(e: Event) {
		e.preventDefault();
		isLoading = true;

		await appState.login(username, password);

		isLoading = false;
		goto('/gallery');
	}
</script>

<div class="min-h-screen w-full flex flex-col lg:flex-row bg-[#FFFFFF] dark:bg-[#090D16] text-slate-900 dark:text-slate-100 font-sans overflow-hidden relative">

	<!-- Soft Ambient Glow Background Elements -->
	<div class="absolute -top-32 -left-32 w-[500px] h-[500px] bg-sky-200/40 dark:bg-sky-950/30 rounded-full blur-[120px] pointer-events-none animate-pulse" style="animation-duration: 8s;"></div>
	<div class="absolute -bottom-32 -right-32 w-[500px] h-[500px] bg-sky-100/60 dark:bg-sky-900/20 rounded-full blur-[140px] pointer-events-none animate-pulse" style="animation-duration: 12s; animation-delay: 2s;"></div>
	<div class="absolute top-1/2 left-1/3 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-cyan-100/30 dark:bg-cyan-950/20 rounded-full blur-[160px] pointer-events-none"></div>

	<!-- LEFT SIDE: Premium Branding & Features -->
	<div class="hidden lg:flex w-1/2 min-h-screen flex-col justify-between p-12 lg:p-16 border-r border-sky-100 dark:border-slate-800 relative z-10 backdrop-blur-[2px]">

	<div class="flex items-center">
  <Cloud class="w-16 h-16 text-sky-500 float-animation drop-shadow-md" fill="currentColor" />
  <div>
      <h1 class="text-2xl font-extrabold tracking-tight text-slate-900 dark:text-white flex items-baseline">
        Dee<span class="text-4xl bg-clip-text text-transparent bg-gradient-to-r from-sky-500 to-cyan-300 dark:from-sky-300 dark:to-white -mx-0.5">P</span><span class="bg-clip-text text-transparent bg-gradient-to-r from-sky-500 via-blue-400 to-cyan-300 dark:from-sky-300 dark:via-blue-200 dark:to-white">hotos</span>
      </h1>
  </div>
	</div>

		<!-- Center: Hero Text & Features -->
		<div class="my-auto py-12 max-w-lg">
			<h2 class="text-4xl md:text-5xl font-extrabold tracking-tight text-slate-900 dark:text-white mb-6 leading-tight">
				Your memories, <br/>
				<span class="text-sky-500">beautifully secured.</span>
			</h2>
			<p class="text-base text-slate-600 dark:text-slate-400 mb-12 leading-relaxed">
				Experience a lightning-fast, minimalist photo gallery built for privacy and performance. Keep your media in your own hands.
			</p>

			<!-- Feature Cards Stack -->
			<div class="space-y-4">
				{#each features as feature, i}
					<div
						class="p-4 rounded-2xl bg-white/60 dark:bg-slate-900/60 backdrop-blur-xl border border-white/80 dark:border-slate-800/80 shadow-sm flex items-center gap-4 transition-all duration-500 hover:-translate-y-1 hover:shadow-xl hover:shadow-sky-900/5"
						style={`animation: fade-in-up 0.6s ease-out ${i * 0.15}s both;`}
					>
						<div class={`w-12 h-12 rounded-xl ${feature.bg} flex items-center justify-center ${feature.color}`}>
							<feature.icon class="w-6 h-6" />
						</div>
						<div>
							<h3 class="text-sm font-bold text-slate-900 dark:text-white">{feature.title}</h3>
							<p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">{feature.desc}</p>
						</div>
					</div>
				{/each}
			</div>
		</div>

		<!-- Left Bottom Copyright Text -->
		<div>
			<p class="text-xs font-medium text-slate-400 dark:text-slate-500">
				© {new Date().getFullYear()} DeepPhotos. Open-Source Software.
			</p>
		</div>
	</div>

	<!-- RIGHT SIDE: Premium Login Form -->
	<div class="w-full lg:w-1/2 min-h-screen flex items-center justify-center p-6 md:p-12 relative z-20">

		<div class="w-full max-w-sm" style="animation: fade-in-scale 0.7s ease-out both;">

			<!-- Mobile Header Branding (Visible on small screens) -->
			<div class="flex lg:hidden flex-col items-center mb-10 text-center">
				<Cloud class="w-14 h-14 text-sky-500 mb-4 float-animation drop-shadow-md" fill="currentColor" />
				<h1 class="text-2xl font-extrabold tracking-tight text-slate-900 dark:text-white">DeepPhotos</h1>
				<p class="text-sm font-medium text-sky-600 dark:text-sky-400">Self-Hosted Vault</p>
			</div>

			<!-- Premium Glassmorphism Card Container -->
			<div class="p-8 md:p-10 rounded-[2rem] bg-white/80 dark:bg-slate-900/80 backdrop-blur-2xl border border-white dark:border-slate-800 shadow-2xl shadow-sky-900/10 dark:shadow-black/50">

				<div class="mb-8">
					<h3 class="text-2xl font-bold text-slate-900 dark:text-white">Welcome back</h3>
					<p class="text-sm text-slate-500 dark:text-slate-400 mt-2">Enter your credentials to access your vault</p>
				</div>

				<form onsubmit={handleLogin} class="space-y-5">

					<!-- Username Field -->
					<div class="space-y-2">
						<label for="username" class="text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider flex items-center gap-2">
							<User class="w-3.5 h-3.5 text-sky-500" />
							<span>Username</span>
						</label>
						<input
							id="username"
							type="text"
							bind:value={username}
							placeholder="admin@deepphotos.local"
							class="w-full h-12 px-4 text-sm rounded-xl bg-slate-50/50 dark:bg-slate-950/50 backdrop-blur-sm border border-slate-200 dark:border-slate-800 text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-slate-600 focus:outline-none focus:border-sky-500 focus:ring-4 focus:ring-sky-500/10 transition-all"
							required
						/>
					</div>

					<!-- Password Field -->
					<div class="space-y-2">
						<label for="password" class="text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider flex items-center gap-2">
							<Lock class="w-3.5 h-3.5 text-sky-500" />
							<span>Password</span>
						</label>
						<div class="relative group">
							<input
								id="password"
								type={showPassword ? 'text' : 'password'}
								bind:value={password}
								placeholder="••••••••"
								class="w-full h-12 px-4 pr-12 text-sm rounded-xl bg-slate-50/50 dark:bg-slate-950/50 backdrop-blur-sm border border-slate-200 dark:border-slate-800 text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-slate-600 focus:outline-none focus:border-sky-500 focus:ring-4 focus:ring-sky-500/10 transition-all"
								required
							/>
							<button
								type="button"
								onclick={() => showPassword = !showPassword}
								class="absolute right-2 top-1/2 -translate-y-1/2 w-8 h-8 flex items-center justify-center rounded-lg text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 hover:bg-slate-200/50 dark:hover:bg-slate-800/50 transition-all"
							>
								{#if showPassword}
									<EyeOff class="w-4 h-4" />
								{:else}
									<Eye class="w-4 h-4" />
								{/if}
							</button>
						</div>
					</div>

					<!-- Sign In Button -->
					<button
						type="submit"
						disabled={isLoading}
						class="w-full h-11 mt-4 rounded-xl bg-sky-400 hover:bg-sky-500 text-white font-bold text-sm transition-all duration-300 flex items-center justify-center gap-2 shadow-lg shadow-sky-300/50 hover:shadow-sky-400/50 hover:-translate-y-0.5 disabled:opacity-70 disabled:hover:translate-y-0 disabled:hover:shadow-sky-300/50 cursor-pointer relative overflow-hidden group"
					>
						<!-- Shine effect -->
						<div class="absolute inset-0 -translate-x-full group-hover:animate-shine bg-gradient-to-r from-transparent via-white/20 to-transparent skew-x-12"></div>


						{#if isLoading}
							<Loader2 class="w-5 h-5 animate-spin text-white" />
							<span>Authenticating...</span>
						{:else}
							<span>Sign In</span>
							<ArrowRight class="w-4 h-4 transition-transform group-hover:translate-x-1" />
						{/if}
					</button>

				</form>
			</div>
		</div>

	</div>
</div>

<style>
	/* Global animations for this page */
	:global(.float-animation) {
		animation: float-drift 6s ease-in-out infinite;
	}

	@keyframes float-drift {
		0%, 100% { transform: translateY(0px) scale(1); }
		50% { transform: translateY(-6px) scale(1.02); }
	}

	@keyframes fade-in-up {
		0% { opacity: 0; transform: translateY(20px); }
		100% { opacity: 1; transform: translateY(0); }
	}

	@keyframes fade-in-scale {
		0% { opacity: 0; transform: scale(0.95); }
		100% { opacity: 1; transform: scale(1); }
	}

	@keyframes -global-shine {
		100% { transform: translateX(100%); }
	}
	:global(.animate-shine) {
		animation: -global-shine 1.5s infinite;
	}
</style>
