<script lang="ts">
	import { onMount } from 'svelte';
	import { appState, type RegisteredUser, type LoginLog } from '$lib/state.svelte';
	import { apiFetch } from '$lib/api';
	import { notify, confirmDialog } from '$lib/notify.svelte';
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
		X,
		KeyRound,
		ShieldAlert
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

	// Admin Password Reset Modal State
	let showResetPassModal = $state(false);
	let selectedUserForPass = $state<RegisteredUser | null>(null);
	let newPasswordInput = $state('');

	let usersList = $state<RegisteredUser[]>([]);
	let loginHistory = $state<LoginLog[]>([]);

	async function fetchUsers() {
		try {
			const res = await apiFetch('/api/users');
			if (res.ok) {
				usersList = await res.json();
			}
		} catch (e) {
			console.warn('API error fetching users:', e);
		}
	}

	async function fetchAuditLogs() {
		try {
			const res = await apiFetch('/api/audit-logs');
			if (res.ok) {
				loginHistory = await res.json();
			}
		} catch (e) {
			console.warn('API error fetching audit logs:', e);
		}
	}

	onMount(() => {
		if (appState.user.role === 'Administrator') {
			fetchUsers();
			fetchAuditLogs();
		}
	});

	function saveProfile(e: Event) {
		e.preventDefault();
		appState.user.name = nameInput;
		appState.user.email = emailInput;
		isSaved = true;
		setTimeout(() => isSaved = false, 2500);
	}

	async function handleAddUser(e: Event) {
		e.preventDefault();
		if (!newUserName || !newUserEmail) return;

		await appState.addUser({
			name: newUserName,
			email: newUserEmail,
			role: newUserRole,
			avatar: 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=150&auto=format&fit=crop&q=80',
			status: 'Active'
		});

		newUserName = '';
		newUserEmail = '';
		showAddUserModal = false;
		fetchUsers();
	}

	async function handleResetPassword(e: Event) {
		e.preventDefault();
		if (!selectedUserForPass || !newPasswordInput) return;

		try {
			const res = await apiFetch(`/api/users/${selectedUserForPass.id}/password`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ new_password: newPasswordInput })
			});
			if (res.ok) {
				alert(`Password for ${selectedUserForPass.name} successfully updated!`);
				showResetPassModal = false;
				newPasswordInput = '';
				selectedUserForPass = null;
			}
		} catch (err) {
			console.error('Error resetting password:', err);
		}
	}

	async function handleRoleChange(user: RegisteredUser, newRole: string) {
		user.role = newRole as any;
		try {
			await apiFetch(`/api/users/${user.id}/role`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ role: newRole })
			});
		} catch (err) {
			console.error('Error changing role:', err);
		}
	}

	async function handleDeleteUser(userId: string, userName: string) {
		const confirmed = await confirmDialog.ask({
			title: 'Delete User Account',
			message: `Are you sure you want to delete user account "${userName}"? This cannot be undone.`,
			confirmText: 'Yes, Delete User',
			cancelText: 'Cancel',
			type: 'danger'
		});
		if (!confirmed) return;

		try {
			await appState.deleteUser(userId);
			notify.success(`User "${userName}" deleted.`);
			fetchUsers();
		} catch (err) {
			console.error('Error deleting user:', err);
			notify.error('Failed to delete user account.');
		}
	}
</script>

<div class="h-full flex flex-col gap-6 max-w-5xl animate-fade-in">
	
	<!-- Header Banner with User Summary -->
	<div class="shrink-0 p-6 md:p-8 rounded-3xl bg-gradient-to-r from-sky-400 via-sky-500 to-blue-600 text-white shadow-lg shadow-sky-300/30 dark:shadow-sky-900/30 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
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

	<div class="flex-1 overflow-y-auto min-h-0">
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

		{#if appState.user.role === 'Administrator'}
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
				User Management ({usersList.length})
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
				Login History ({loginHistory.length})
			</button>
		{/if}

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
		<div class="p-6 rounded-2xl bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 shadow-sm space-y-5">
			
			<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-100 dark:border-slate-800 pb-4">
				<div>
					<h3 class="text-base font-bold text-slate-900 dark:text-white flex items-center gap-2">
						<UserPlus class="w-5 h-5 text-sky-500" />
						Registered Accounts
					</h3>
					<p class="text-xs text-slate-500 dark:text-slate-400">Administrators can create users, change roles, or reset passwords</p>
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
			<div class="overflow-x-auto rounded-xl border border-slate-100 dark:border-slate-800">
				<table class="w-full text-left border-collapse">
					<thead>
						<tr class="bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-800 text-[11px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider">
							<th class="py-3 px-4">User</th>
							<th class="py-3 px-4">Email</th>
							<th class="py-3 px-4">Role</th>
							<th class="py-3 px-4">Last Active</th>
							<th class="py-3 px-4 text-right">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100 dark:divide-slate-800 text-xs">
						{#each usersList as user}
							<tr class="hover:bg-slate-50 dark:hover:bg-slate-900/60 transition-colors">
								<td class="py-3.5 px-4 font-semibold text-slate-900 dark:text-slate-100">
									<div class="flex items-center gap-3">
										<img src={user.avatar || 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=150&auto=format&fit=crop&q=80'} alt={user.name} class="w-8 h-8 rounded-full object-cover" />
										<span>{user.name}</span>
									</div>
								</td>
								<td class="py-3.5 px-4 text-slate-500 dark:text-slate-400 font-mono">{user.email}</td>
								<td class="py-3.5 px-4">
									<select
										value={user.role}
										onchange={(e) => handleRoleChange(user, (e.target as HTMLSelectElement).value)}
										class="px-2.5 py-1 rounded-lg text-[11px] font-semibold bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-200 cursor-pointer"
									>
										<option value="Administrator">Administrator</option>
										<option value="Editor">Editor</option>
										<option value="Viewer">Viewer</option>
									</select>
								</td>
								<td class="py-3.5 px-4 text-slate-400 dark:text-slate-500">{user.lastLogin || 'Recently'}</td>
								<td class="py-3.5 px-4 text-right">
									<div class="flex items-center justify-end gap-1.5">
										<button
											type="button"
											onclick={() => { selectedUserForPass = user; showResetPassModal = true; }}
											title="Admin Reset Password"
											class="p-1.5 rounded-lg text-amber-500 hover:bg-amber-50 dark:hover:bg-amber-950/30 transition-all cursor-pointer"
										>
											<KeyRound class="w-4 h-4" />
										</button>

										{#if user.email !== appState.user.email}
											<button
												type="button"
												onclick={() => handleDeleteUser(user.id, user.name)}
												title="Delete User"
												class="p-1.5 rounded-lg text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/30 transition-all cursor-pointer"
											>
												<Trash2 class="w-4 h-4" />
											</button>
										{/if}
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>

		</div>

		{#if showResetPassModal && selectedUserForPass}
			<!-- Admin Password Reset Modal -->
			<div class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-sm flex items-center justify-center p-4">
				<div class="bg-white dark:bg-slate-900 rounded-3xl p-6 md:p-8 max-w-md w-full border border-slate-200 dark:border-slate-800 shadow-2xl animate-fade-in relative">
					<button
						type="button"
						onclick={() => showResetPassModal = false}
						class="absolute top-5 right-5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1"
					>
						<X class="w-5 h-5" />
					</button>

					<h3 class="text-lg font-bold text-slate-900 dark:text-white mb-1 flex items-center gap-2">
						<KeyRound class="w-5 h-5 text-amber-500" />
						Admin Password Reset
					</h3>
					<p class="text-xs text-slate-500 dark:text-slate-400 mb-6">Reset password for <strong class="text-slate-800 dark:text-slate-200">{selectedUserForPass.name}</strong> ({selectedUserForPass.email})</p>

					<form onsubmit={handleResetPassword} class="space-y-4">
						<div class="space-y-1.5">
							<label for="new-pass" class="text-xs font-semibold text-slate-700 dark:text-slate-300">New Password</label>
							<input
								id="new-pass"
								type="password"
								bind:value={newPasswordInput}
								placeholder="Enter new password"
								class="w-full h-10 px-3.5 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-amber-400"
								required
							/>
						</div>

						<div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100 dark:border-slate-800">
							<button
								type="button"
								onclick={() => showResetPassModal = false}
								class="px-4 py-2 rounded-xl text-xs font-semibold text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800"
							>
								Cancel
							</button>
							<button
								type="submit"
								class="px-5 py-2 rounded-xl bg-amber-500 hover:bg-amber-600 text-white text-xs font-bold shadow-sm shadow-amber-300/50"
							>
								Update Password
							</button>
						</div>
					</form>
				</div>
			</div>
		{/if}

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
		<div class="p-6 rounded-2xl bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 shadow-sm space-y-4">
			<div>
				<h3 class="text-base font-bold text-slate-900 dark:text-white flex items-center gap-2">
					<Clock class="w-5 h-5 text-sky-500" />
					Login Activity & Access History
				</h3>
				<p class="text-xs text-slate-500 dark:text-slate-400">Security audit log of recent authentication attempts</p>
			</div>

			<div class="overflow-x-auto rounded-xl border border-slate-100 dark:border-slate-800">
				<table class="w-full text-left border-collapse">
					<thead>
						<tr class="bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-800 text-[11px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider">
							<th class="py-3 px-4">User</th>
							<th class="py-3 px-4">Timestamp</th>
							<th class="py-3 px-4">IP Address</th>
							<th class="py-3 px-4">Device & Client</th>
							<th class="py-3 px-4 text-right">Status</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100 dark:divide-slate-800 text-xs">
						{#if loginHistory.length === 0}
							<tr>
								<td colspan="5" class="py-10 text-center text-xs text-slate-400 dark:text-slate-500">
									<Clock class="w-8 h-8 mx-auto mb-2 opacity-30" />
									No login records found
								</td>
							</tr>
						{:else}
							{#each loginHistory as log}
								<tr class="hover:bg-slate-50 dark:hover:bg-slate-900/60 transition-colors">
									<td class="py-3.5 px-4 font-semibold text-slate-800 dark:text-slate-100 flex items-center gap-2">
										<User class="w-4 h-4 text-sky-500 shrink-0" />
										<span>{log.user}</span>
									</td>
									<td class="py-3.5 px-4 text-slate-500 dark:text-slate-400">{log.timestamp}</td>
									<td class="py-3.5 px-4 font-mono text-slate-500 dark:text-slate-400">{log.ip}</td>
									<td class="py-3.5 px-4 text-slate-500 dark:text-slate-400">
										<span class="flex items-center gap-1.5">
											<Laptop class="w-3.5 h-3.5 shrink-0" />
											{log.device}
										</span>
									</td>
									<td class="py-3.5 px-4 text-right">
										{#if log.status === 'Success'}
											<span class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 dark:bg-emerald-950/40 text-emerald-600 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-900">
												Success
											</span>
										{:else}
											<span class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-rose-50 dark:bg-rose-950/40 text-rose-600 dark:text-rose-400 border border-rose-200 dark:border-rose-900">
												Failed
											</span>
										{/if}
									</td>
								</tr>
							{/each}
						{/if}
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
						<span class="text-slate-500 dark:text-slate-400">MinIO Folder Format:</span>
						<span class="font-mono text-sky-600 dark:text-sky-400">category/uuid1/uuid2/uuid3/file</span>
					</div>
					<div class="flex justify-between py-1.5 border-b border-slate-100 dark:border-slate-800">
						<span class="text-slate-500 dark:text-slate-400">Metadata Database:</span>
						<span class="font-bold text-slate-800 dark:text-slate-200">SQLite 3 (photos.db)</span>
					</div>
					<div class="flex justify-between py-1.5 border-b border-slate-100 dark:border-slate-800">
						<span class="text-slate-500 dark:text-slate-400">API Endpoint:</span>
						<span class="font-mono text-sky-600 dark:text-sky-400">http://localhost:8080</span>
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
						<span class="font-bold text-emerald-600 dark:text-emerald-400">AES-256 Passcode Protection</span>
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
</div>
