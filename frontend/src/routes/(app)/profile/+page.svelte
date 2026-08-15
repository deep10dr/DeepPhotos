<script lang="ts">
	import { appState, type RegisteredUser } from '$lib/state.svelte';
	import {
		User,
		ShieldCheck,
		Key,
		HardDrive,
		Database,
		Lock,
		CheckCircle2,
		XCircle,
		Mail,
		Calendar,
		Sparkles,
		Save,
		UserPlus,
		Trash2,
		Clock,
		Globe,
		Laptop,
		AlertTriangle,
		X
	} from 'lucide-svelte';

	let activeSection = $state<'profile' | 'users' | 'history' | 'permissions' | 'server'>('profile');
	let isSaved = $state(false);

	let nameInput = $state(appState.user.name);
	let emailInput = $state(appState.user.email);

	// Add User Modal State
	let showAddUserModal = $state(false);
	let newUserName = $state('');
	let newUserEmail = $state('');
	let newUserRole = $state<'Administrator' | 'Editor' | 'Viewer'>('Editor');
	let newUserAvatar = $state('https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=150&auto=format&fit=crop&q=80');

	function saveProfile(e: Event) {
		e.preventDefault();
		appState.user.name = nameInput;
		appState.user.email = emailInput;
		isSaved = true;
		setTimeout(() => isSaved = false, 2500);
	}

	function handleAddUser(e: Event) {
		e.preventDefault();
		if (!newUserName || !newUserEmail) return;

		appState.addUser({
			name: newUserName,
			email: newUserEmail,
			role: newUserRole,
			avatar: newUserAvatar,
			status: 'Active'
		});

		newUserName = '';
		newUserEmail = '';
		showAddUserModal = false;
	}

	function handleDeleteUser(userId: string, userName: string) {
		if (confirm(`Are you sure you want to delete user "${userName}"?`)) {
			appState.deleteUser(userId);
		}
	}
</script>

<div class="space-y-6 max-w-5xl animate-fade-in">
	
	<!-- Header Banner with User Summary -->
	<div class="p-6 md:p-8 rounded-3xl bg-gradient-to-r from-sky-400 via-sky-500 to-blue-600 text-white shadow-lg shadow-sky-300/30 dark:shadow-sky-900/30 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
		<div class="flex items-center gap-4">
			<img
				src={appState.user.avatar}
				alt={appState.user.name}
				class="w-16 h-16 rounded-2xl object-cover ring-4 ring-white/30 shadow-md"
			/>
			<div>
				<div class="flex items-center gap-2">
					<h1 class="text-2xl font-bold">{appState.user.name}</h1>
					<span class="px-2.5 py-0.5 rounded-full bg-white/20 text-white text-xs font-semibold backdrop-blur-md">
						{appState.user.role}
					</span>
				</div>
				<p class="text-xs text-sky-100 mt-1 flex flex-wrap items-center gap-3">
					<span class="flex items-center gap-1"><Mail class="w-3.5 h-3.5" /> {appState.user.email}</span>
					<span class="flex items-center gap-1"><Calendar class="w-3.5 h-3.5" /> Member since {appState.user.joinedDate}</span>
				</p>
			</div>
		</div>

		<div class="flex items-center gap-2">
			<span class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full bg-white/20 text-xs font-semibold backdrop-blur-md">
				<ShieldCheck class="w-4 h-4 text-emerald-300" />
				Active Admin Session
			</span>
		</div>
	</div>

	<!-- Section Tabs -->
	<div class="flex flex-wrap items-center gap-2 border-b border-slate-200 dark:border-slate-800 pb-2">
		<button
			type="button"
			onclick={() => activeSection = 'profile'}
			class={`px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer ${
				activeSection === 'profile'
					? 'bg-sky-400 text-white shadow-sm shadow-sky-300/50'
					: 'text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800'
			}`}
		>
			<User class="w-4 h-4 inline mr-1.5" />
			User Details
		</button>

		<button
			type="button"
			onclick={() => activeSection = 'users'}
			class={`px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer ${
				activeSection === 'users'
					? 'bg-sky-400 text-white shadow-sm shadow-sky-300/50'
					: 'text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800'
			}`}
		>
			<UserPlus class="w-4 h-4 inline mr-1.5" />
			User Management ({appState.usersList.length})
		</button>

		<button
			type="button"
			onclick={() => activeSection = 'history'}
			class={`px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer ${
				activeSection === 'history'
					? 'bg-sky-400 text-white shadow-sm shadow-sky-300/50'
					: 'text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800'
			}`}
		>
			<Clock class="w-4 h-4 inline mr-1.5" />
			Login History ({appState.loginHistory.length})
		</button>

		<button
			type="button"
			onclick={() => activeSection = 'permissions'}
			class={`px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer ${
				activeSection === 'permissions'
					? 'bg-sky-400 text-white shadow-sm shadow-sky-300/50'
					: 'text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800'
			}`}
		>
			<Key class="w-4 h-4 inline mr-1.5" />
			Permissions
		</button>

		<button
			type="button"
			onclick={() => activeSection = 'server'}
			class={`px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer ${
				activeSection === 'server'
					? 'bg-sky-400 text-white shadow-sm shadow-sky-300/50'
					: 'text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800'
			}`}
		>
			<HardDrive class="w-4 h-4 inline mr-1.5" />
			Node Status
		</button>
	</div>

	{#if activeSection === 'profile'}
		<!-- TAB 1: User Details Section -->
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			
			<div class="lg:col-span-2 p-6 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm space-y-4">
				<h3 class="text-base font-bold text-slate-900 dark:text-white border-b border-slate-100 dark:border-slate-800 pb-3 flex items-center gap-2">
					<User class="w-4 h-4 text-sky-500" />
					Account Information
				</h3>

				{#if isSaved}
					<div class="p-3 rounded-xl bg-emerald-50 dark:bg-emerald-950/60 text-emerald-700 dark:text-emerald-300 text-xs font-semibold flex items-center gap-2 border border-emerald-200 dark:border-emerald-800">
						<CheckCircle2 class="w-4 h-4" />
						Profile information updated successfully!
					</div>
				{/if}

				<form onsubmit={saveProfile} class="space-y-4">
					<div class="space-y-1.5">
						<label for="display-name" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Display Name</label>
						<input
							id="display-name"
							type="text"
							bind:value={nameInput}
							class="w-full h-10 px-3.5 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400 focus:bg-white dark:focus:bg-slate-900 focus:ring-2 focus:ring-sky-100 dark:focus:ring-sky-900"
							required
						/>
					</div>

					<div class="space-y-1.5">
						<label for="email-address" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Email Address</label>
						<input
							id="email-address"
							type="email"
							bind:value={emailInput}
							class="w-full h-10 px-3.5 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400 focus:bg-white dark:focus:bg-slate-900 focus:ring-2 focus:ring-sky-100 dark:focus:ring-sky-900"
							required
						/>
					</div>

					<div class="space-y-1.5">
						<label for="user-role" class="text-xs font-semibold text-slate-700 dark:text-slate-300">User Role</label>
						<input
							id="user-role"
							type="text"
							value={appState.user.role}
							disabled
							class="w-full h-10 px-3.5 text-xs rounded-xl bg-slate-100 dark:bg-slate-800/60 border border-slate-200 dark:border-slate-700 text-slate-500 dark:text-slate-400 cursor-not-allowed"
						/>
					</div>

					<button
						type="submit"
						class="px-5 py-2.5 rounded-xl bg-sky-400 hover:bg-sky-500 text-white text-xs font-bold shadow-sm shadow-sky-300/50 transition-all flex items-center gap-2 cursor-pointer"
					>
						<Save class="w-4 h-4" />
						Save Changes
					</button>
				</form>
			</div>

			<div class="p-6 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm space-y-4">
				<h3 class="text-base font-bold text-slate-900 dark:text-white border-b border-slate-100 dark:border-slate-800 pb-3 flex items-center gap-2">
					<Sparkles class="w-4 h-4 text-sky-500" />
					Vault Metrics
				</h3>

				<div class="space-y-3">
					<div class="p-3 rounded-xl bg-sky-50/60 dark:bg-slate-800/60 border border-sky-100 dark:border-slate-800 flex items-center justify-between">
						<span class="text-xs text-slate-600 dark:text-slate-400 font-medium">Photos Stored</span>
						<span class="text-sm font-bold text-sky-600 dark:text-sky-400">{appState.user.storage.photosCount}</span>
					</div>

					<div class="p-3 rounded-xl bg-sky-50/60 dark:bg-slate-800/60 border border-sky-100 dark:border-slate-800 flex items-center justify-between">
						<span class="text-xs text-slate-600 dark:text-slate-400 font-medium">Albums Created</span>
						<span class="text-sm font-bold text-sky-600 dark:text-sky-400">{appState.user.storage.albumsCount}</span>
					</div>

					<div class="p-3 rounded-xl bg-sky-50/60 dark:bg-slate-800/60 border border-sky-100 dark:border-slate-800 flex items-center justify-between">
						<span class="text-xs text-slate-600 dark:text-slate-400 font-medium">Storage Used</span>
						<span class="text-sm font-bold text-sky-600 dark:text-sky-400">{appState.user.storage.used} / {appState.user.storage.total}</span>
					</div>
				</div>
			</div>

		</div>

	{:else if activeSection === 'users'}
		<!-- TAB 2: User Management -->
		<div class="p-6 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm space-y-5">
			
			<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-100 dark:border-slate-800 pb-4">
				<div>
					<h3 class="text-base font-bold text-slate-900 dark:text-white flex items-center gap-2">
						<UserPlus class="w-5 h-5 text-sky-500" />
						Registered Accounts
					</h3>
					<p class="text-xs text-slate-500 dark:text-slate-400">Administrators can invite, create, or remove users</p>
				</div>

				<button
					type="button"
					onclick={() => showAddUserModal = true}
					class="px-4 py-2 rounded-xl bg-sky-400 hover:bg-sky-500 text-white text-xs font-bold shadow-sm shadow-sky-300/50 flex items-center gap-1.5 cursor-pointer transition-all"
				>
					<UserPlus class="w-4 h-4" />
					Add New User
				</button>
			</div>

			<!-- Users Table -->
			<div class="overflow-x-auto">
				<table class="w-full text-left border-collapse">
					<thead>
						<tr class="border-b border-slate-200 dark:border-slate-800 text-[11px] font-bold text-slate-400 uppercase tracking-wider">
							<th class="py-3 px-4">User</th>
							<th class="py-3 px-4">Email</th>
							<th class="py-3 px-4">Role</th>
							<th class="py-3 px-4">Last Active</th>
							<th class="py-3 px-4 text-right">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100 dark:divide-slate-800 text-xs">
						{#each appState.usersList as user}
							<tr class="hover:bg-sky-50/50 dark:hover:bg-slate-800/50 transition-colors">
								<td class="py-3.5 px-4 font-bold text-slate-900 dark:text-white flex items-center gap-3">
									<img src={user.avatar} alt={user.name} class="w-8 h-8 rounded-full object-cover ring-2 ring-sky-200 dark:ring-sky-900" />
									<span>{user.name}</span>
								</td>
								<td class="py-3.5 px-4 text-slate-600 dark:text-slate-300 font-mono">{user.email}</td>
								<td class="py-3.5 px-4">
									<span class={`px-2.5 py-0.5 rounded-full text-[11px] font-semibold ${
										user.role === 'Administrator' ? 'bg-sky-100 dark:bg-sky-950/80 text-sky-700 dark:text-sky-300 border border-sky-200 dark:border-sky-800' :
										user.role === 'Editor' ? 'bg-emerald-100 dark:bg-emerald-950/80 text-emerald-700 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800' :
										'bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700'
									}`}>
										{user.role}
									</span>
								</td>
								<td class="py-3.5 px-4 text-slate-500 dark:text-slate-400">{user.lastLogin}</td>
								<td class="py-3.5 px-4 text-right">
									{#if user.email !== appState.user.email}
										<button
											type="button"
											onclick={() => handleDeleteUser(user.id, user.name)}
											title="Delete User"
											class="p-2 rounded-xl text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/40 border border-transparent hover:border-rose-200 dark:hover:border-rose-900 transition-all cursor-pointer"
										>
											<Trash2 class="w-4 h-4" />
										</button>
									{:else}
										<span class="text-[11px] text-slate-400 italic">Current User</span>
									{/if}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>

		</div>

		{#if showAddUserModal}
			<!-- Add User Modal -->
			<div class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-sm flex items-center justify-center p-4">
				<div class="bg-white dark:bg-slate-900 rounded-3xl p-6 md:p-8 max-w-md w-full border border-slate-200 dark:border-slate-800 shadow-2xl animate-fade-in relative">
					<button
						type="button"
						onclick={() => showAddUserModal = false}
						class="absolute top-5 right-5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1"
					>
						<X class="w-5 h-5" />
					</button>

					<h3 class="text-lg font-bold text-slate-900 dark:text-white mb-1 flex items-center gap-2">
						<UserPlus class="w-5 h-5 text-sky-500" />
						Create New Account
					</h3>
					<p class="text-xs text-slate-500 dark:text-slate-400 mb-6">Add a new user to your self-hosted DeepPhotos node</p>

					<form onsubmit={handleAddUser} class="space-y-4">
						<div class="space-y-1.5">
							<label for="new-name" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Full Name</label>
							<input
								id="new-name"
								type="text"
								bind:value={newUserName}
								placeholder="e.g. John Doe"
								class="w-full h-10 px-3.5 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400 focus:bg-white dark:focus:bg-slate-900"
								required
							/>
						</div>

						<div class="space-y-1.5">
							<label for="new-email" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Email Address</label>
							<input
								id="new-email"
								type="email"
								bind:value={newUserEmail}
								placeholder="john@deepphotos.local"
								class="w-full h-10 px-3.5 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400 focus:bg-white dark:focus:bg-slate-900"
								required
							/>
						</div>

						<div class="space-y-1.5">
							<label for="new-role" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Assign Role</label>
							<select
								id="new-role"
								bind:value={newUserRole}
								class="w-full h-10 px-3 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400 focus:bg-white dark:focus:bg-slate-900"
							>
								<option value="Administrator">Administrator (Full Access)</option>
								<option value="Editor">Editor (Upload & Organize)</option>
								<option value="Viewer">Viewer (Read Only)</option>
							</select>
						</div>

						<div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100 dark:border-slate-800">
							<button
								type="button"
								onclick={() => showAddUserModal = false}
								class="px-4 py-2 rounded-xl text-xs font-semibold text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800"
							>
								Cancel
							</button>
							<button
								type="submit"
								class="px-5 py-2 rounded-xl bg-sky-400 hover:bg-sky-500 text-white text-xs font-bold shadow-sm shadow-sky-300/50"
							>
								Create Account
							</button>
						</div>
					</form>
				</div>
			</div>
		{/if}

	{:else if activeSection === 'history'}
		<!-- TAB 3: User Login History Audit Trail -->
		<div class="p-6 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm space-y-4">
			<div>
				<h3 class="text-base font-bold text-slate-900 dark:text-white flex items-center gap-2">
					<Clock class="w-5 h-5 text-sky-500" />
					Login Activity & Access History
				</h3>
				<p class="text-xs text-slate-500 dark:text-slate-400">Security audit log of recent authentication attempts</p>
			</div>

			<div class="overflow-x-auto">
				<table class="w-full text-left border-collapse">
					<thead>
						<tr class="border-b border-slate-200 dark:border-slate-800 text-[11px] font-bold text-slate-400 uppercase tracking-wider">
							<th class="py-3 px-4">User</th>
							<th class="py-3 px-4">Timestamp</th>
							<th class="py-3 px-4">IP Address</th>
							<th class="py-3 px-4">Device & Client</th>
							<th class="py-3 px-4 text-right">Status</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100 dark:divide-slate-800 text-xs">
						{#each appState.loginHistory as log}
							<tr class="hover:bg-sky-50/50 dark:hover:bg-slate-800/50 transition-colors">
								<td class="py-3.5 px-4 font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2">
									<User class="w-4 h-4 text-sky-500" />
									<span>{log.user}</span>
								</td>
								<td class="py-3.5 px-4 text-slate-600 dark:text-slate-300">{log.timestamp}</td>
								<td class="py-3.5 px-4 font-mono text-slate-500 dark:text-slate-400">{log.ip}</td>
								<td class="py-3.5 px-4 text-slate-600 dark:text-slate-300 flex items-center gap-1.5">
									<Laptop class="w-3.5 h-3.5 text-slate-400" />
									<span>{log.device}</span>
								</td>
								<td class="py-3.5 px-4 text-right">
									{#if log.status === 'Success'}
										<span class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 dark:bg-emerald-950/80 text-emerald-700 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800">
											Success
										</span>
									{:else}
										<span class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-rose-100 dark:bg-rose-950/80 text-rose-700 dark:text-rose-300 border border-rose-200 dark:border-rose-800">
											Failed Attempt
										</span>
									{/if}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>

	{:else if activeSection === 'permissions'}
		<!-- TAB 4: Permissions List -->
		<div class="p-6 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm space-y-4">
			<div>
				<h3 class="text-base font-bold text-slate-900 dark:text-white">User Permissions & Access Control</h3>
				<p class="text-xs text-slate-500 dark:text-slate-400">Active access rights assigned to role: <strong class="text-sky-600 dark:text-sky-400">{appState.user.role}</strong></p>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-3 pt-2">
				{#each appState.user.permissions as perm}
					<div class="p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/60 dark:bg-slate-800/60 flex items-start justify-between gap-3">
						<div class="space-y-1">
							<p class="text-xs font-bold text-slate-900 dark:text-white flex items-center gap-1.5">
								<CheckCircle2 class="w-4 h-4 text-emerald-500 shrink-0" />
								<span>{perm.name}</span>
							</p>
							<p class="text-[11px] text-slate-500 dark:text-slate-400">{perm.desc}</p>
						</div>
						<span class="text-[10px] px-2 py-0.5 rounded-full bg-emerald-100 dark:bg-emerald-950/80 text-emerald-700 dark:text-emerald-300 font-semibold border border-emerald-200 dark:border-emerald-800">
							Granted
						</span>
					</div>
				{/each}
			</div>
		</div>

	{:else if activeSection === 'server'}
		<!-- TAB 5: Server / Node Status -->
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			<div class="p-6 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm space-y-4">
				<h3 class="text-base font-bold text-slate-900 dark:text-white border-b border-slate-100 dark:border-slate-800 pb-3 flex items-center gap-2">
					<Database class="w-4 h-4 text-sky-500" />
					Database & MinIO Engine
				</h3>

				<div class="space-y-2.5 text-xs">
					<div class="flex justify-between py-1.5 border-b border-slate-100 dark:border-slate-800">
						<span class="text-slate-500 dark:text-slate-400">Object Storage:</span>
						<span class="font-bold text-slate-800 dark:text-slate-200">MinIO (S3 Compatible)</span>
					</div>
					<div class="flex justify-between py-1.5 border-b border-slate-100 dark:border-slate-800">
						<span class="text-slate-500 dark:text-slate-400">Metadata Database:</span>
						<span class="font-bold text-slate-800 dark:text-slate-200">SQLite 3 (photos.db)</span>
					</div>
					<div class="flex justify-between py-1.5 border-b border-slate-100 dark:border-slate-800">
						<span class="text-slate-500 dark:text-slate-400">API Endpoint:</span>
						<span class="font-mono text-sky-600 dark:text-sky-400">http://localhost:8080</span>
					</div>
					<div class="flex justify-between py-1.5">
						<span class="text-slate-500 dark:text-slate-400">Thumbnail Engine:</span>
						<span class="font-bold text-emerald-600 dark:text-emerald-400">WebP Async Worker</span>
					</div>
				</div>
			</div>

			<div class="p-6 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm space-y-4">
				<h3 class="text-base font-bold text-slate-900 dark:text-white border-b border-slate-100 dark:border-slate-800 pb-3 flex items-center gap-2">
					<Lock class="w-4 h-4 text-sky-500" />
					Security & Session Info
				</h3>

				<div class="space-y-2.5 text-xs">
					<div class="flex justify-between py-1.5 border-b border-slate-100 dark:border-slate-800">
						<span class="text-slate-500 dark:text-slate-400">Session Type:</span>
						<span class="font-bold text-slate-800 dark:text-slate-200">Secure JWT Bearer Token</span>
					</div>
					<div class="flex justify-between py-1.5 border-b border-slate-100 dark:border-slate-800">
						<span class="text-slate-500 dark:text-slate-400">Locked Vault Protection:</span>
						<span class="font-bold text-emerald-600 dark:text-emerald-400">AES-256 Enabled</span>
					</div>
					<div class="flex justify-between py-1.5">
						<span class="text-slate-500 dark:text-slate-400">Cloud Data Sharing:</span>
						<span class="font-bold text-slate-800 dark:text-slate-200">Disabled (0% External)</span>
					</div>
				</div>
			</div>
		</div>
	{/if}

</div>
