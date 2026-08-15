<script lang="ts">
	import { goto } from '$app/navigation';
	import { appState } from '$lib/state.svelte';
	import { Camera, Lock, User, Eye, EyeOff, ArrowRight, Loader2, Image as ImageIcon, Cloud } from 'lucide-svelte';

	// Pre-fill from env for development by default
	const devUsername = import.meta.env.VITE_DEV_USERNAME || 'admin@deepphotos.local';
	const devPassword = import.meta.env.VITE_DEV_PASSWORD || 'deepphotos2026';

	// Svelte 5 state runes
	let username = $state(devUsername);
	let password = $state(devPassword);
	let showPassword = $state(false);
	let isLoading = $state(false);

	// Minimal clean photo cards for left side display
	const photoCards = [
		{ id: 1, title: 'Cloudy Mountain Peak', date: 'August 2026', tag: 'Landscape', bg: 'from-sky-100 to-blue-50' },
		{ id: 2, title: 'Coastal Waves & Sky', date: 'July 2026', tag: 'Travel', bg: 'from-cyan-100 to-sky-50' },
		{ id: 3, title: 'Morning Horizon', date: 'June 2026', tag: 'Nature', bg: 'from-sky-200/50 to-cyan-50' }
	];

	function handleLogin(e: Event) {
		e.preventDefault();
		isLoading = true;

		setTimeout(() => {
			isLoading = false;
			appState.login();
			goto('/gallery');
		}, 600);
	}
</script>

<div class="min-h-screen w-full flex flex-col lg:flex-row bg-[#FFFFFF] text-slate-900 font-sans overflow-hidden relative">
	
	<!-- Soft Cloud & Sky Ambient Glow Background Elements -->
	<div class="absolute -top-32 -left-32 w-[500px] h-[500px] bg-sky-200/40 rounded-full blur-[120px] pointer-events-none animate-cloud-float"></div>
	<div class="absolute -bottom-32 -right-32 w-[500px] h-[500px] bg-sky-100/60 rounded-full blur-[140px] pointer-events-none animate-cloud-float" style="animation-delay: 4s;"></div>
	<div class="absolute top-1/2 left-1/3 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-cyan-100/30 rounded-full blur-[160px] pointer-events-none"></div>

	<!-- LEFT SIDE: Minimal Cloud Branding, Photo Cards & Copyright -->
	<div class="hidden lg:flex w-1/2 min-h-screen flex-col justify-between p-12 lg:p-16 border-r border-sky-100 relative bg-gradient-to-br from-white via-sky-50/40 to-sky-100/30">
		
		<!-- Top Header Branding -->
		<div class="relative z-10 flex items-center gap-3">
			<div class="w-10 h-10 rounded-xl bg-sky-400 text-white flex items-center justify-center shadow-md shadow-sky-300/40">
				<Camera class="w-5 h-5 text-white" />
			</div>
			<div>
				<h1 class="text-xl font-bold tracking-tight text-slate-900">DeepPhotos</h1>
				<p class="text-xs text-slate-500">Minimal Sky & Cloud Photo Vault</p>
			</div>
		</div>

		<!-- Center: Photo Cards Display -->
		<div class="relative z-10 my-auto py-6">
			<div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-sky-100/80 border border-sky-200 text-sky-800 text-xs font-medium mb-4">
				<Cloud class="w-3.5 h-3.5 text-sky-500" />
				<span>Clean & Minimal Experience</span>
			</div>

			<h2 class="text-3xl font-extrabold tracking-tight text-slate-900 mb-2">
				Store your memories under <br/>
				<span class="text-sky-500">clear blue skies.</span>
			</h2>
			<p class="text-sm text-slate-600 max-w-sm mb-8">
				A minimal, fast photo management gallery built for effortless clarity.
			</p>

			<!-- Stacked Cloud Photo Cards -->
			<div class="space-y-3.5 max-w-sm">
				{#each photoCards as card}
					<div class="p-3.5 rounded-2xl cloud-card flex items-center justify-between transition-all duration-300 hover:shadow-lg hover:-translate-y-0.5">
						<div class="flex items-center gap-3">
							<div class={`w-11 h-11 rounded-xl bg-gradient-to-br ${card.bg} border border-sky-200/70 flex items-center justify-center text-sky-600`}>
								<ImageIcon class="w-5 h-5" />
							</div>
							<div>
								<p class="text-xs font-bold text-slate-800">{card.title}</p>
								<p class="text-[11px] text-slate-500">{card.date}</p>
							</div>
						</div>
						<span class="text-[11px] text-sky-700 px-2.5 py-0.5 rounded-full bg-sky-100 border border-sky-200 font-medium">
							{card.tag}
						</span>
					</div>
				{/each}
			</div>
		</div>

		<!-- Left Bottom Copyright Text -->
		<div class="relative z-10">
			<p class="text-xs text-slate-500">
				© {new Date().getFullYear()} DeepPhotos. All rights reserved.
			</p>
		</div>

	</div>

	<!-- RIGHT SIDE: Minimal Cloud Login Form -->
	<div class="w-full lg:w-1/2 min-h-screen flex items-center justify-center p-6 md:p-12 relative z-10">
		
		<div class="w-full max-w-sm animate-fade-in">
			
			<!-- Mobile Header Branding (Visible on small screens) -->
			<div class="flex lg:hidden flex-col items-center mb-8 text-center">
				<div class="w-10 h-10 rounded-xl bg-sky-400 text-white flex items-center justify-center shadow-md shadow-sky-300/40 mb-2">
					<Camera class="w-5 h-5 text-white" />
				</div>
				<h1 class="text-xl font-bold text-slate-900">DeepPhotos</h1>
				<p class="text-xs text-slate-500">Minimal Photo Vault</p>
			</div>

			<!-- Minimal White Cloud Card Container -->
			<div class="p-8 rounded-3xl cloud-card">
				
				<div class="mb-6">
					<h3 class="text-xl font-bold text-slate-900">Sign In</h3>
					<p class="text-xs text-slate-500 mt-1">Enter your details to access your photos</p>
				</div>

				<form onsubmit={handleLogin} class="space-y-4">
					
					<!-- Username / Name Field -->
					<div class="space-y-1.5">
						<label for="username" class="text-xs font-semibold text-slate-700 flex items-center gap-1.5">
							<User class="w-3.5 h-3.5 text-sky-500" />
							<span>Username</span>
						</label>
						<input
							id="username"
							type="text"
							bind:value={username}
							placeholder="Username"
							class="w-full h-11 px-3.5 text-xs rounded-xl cloud-input"
							required
						/>
					</div>

					<!-- Password Field -->
					<div class="space-y-1.5">
						<label for="password" class="text-xs font-semibold text-slate-700 flex items-center gap-1.5">
							<Lock class="w-3.5 h-3.5 text-sky-500" />
							<span>Password</span>
						</label>
						<div class="relative">
							<input
								id="password"
								type={showPassword ? 'text' : 'password'}
								bind:value={password}
								placeholder="Password"
								class="w-full h-11 px-3.5 pr-10 text-xs rounded-xl cloud-input"
								required
							/>
							<button
								type="button"
								onclick={() => showPassword = !showPassword}
								class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors p-1"
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
						class="w-full h-11 mt-2 rounded-xl bg-sky-400 hover:bg-sky-500 text-white font-bold text-xs transition-all duration-200 flex items-center justify-center gap-2 cursor-pointer shadow-md shadow-sky-300/50 disabled:opacity-50"
					>
						{#if isLoading}
							<Loader2 class="w-4 h-4 animate-spin text-white" />
							<span>Signing In...</span>
						{:else}
							<span>Sign In</span>
							<ArrowRight class="w-4 h-4" />
						{/if}
					</button>

				</form>

			</div>

			<!-- Mobile Copyright (Visible on mobile only) -->
			<p class="block lg:hidden text-center text-xs text-slate-500 mt-6">
				© {new Date().getFullYear()} DeepPhotos. All rights reserved.
			</p>

		</div>

	</div>

</div>
