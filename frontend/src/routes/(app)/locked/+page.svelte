<script lang="ts">
	import { appState } from '$lib/state.svelte';
	import { Lock, Unlock, Key, ShieldCheck, FolderClosed, Plus, X, Trash2, User, Clock, AlertCircle } from 'lucide-svelte';

	interface LockedFolder {
		id: string;
		user_id: string;
		user_name: string;
		name: string;
		description: string;
		photos_count: number;
		created_at: string;
	}

	let folders = $state<LockedFolder[]>([]);
	let isLoading = $state(true);

	// Create Locked Folder Modal State (2-Step Validation)
	let showCreateModal = $state(false);
	let step = $state<1 | 2>(1);
	let folderName = $state('');
	let folderDesc = $state('');
	let passcode = $state('');
	let passcodeConfirm = $state('');
	let passcodeError = $state('');

	// Unlock Folder State
	let selectedFolder = $state<LockedFolder | null>(null);
	let unlockPasscode = $state('');
	let isUnlocked = $state(false);
	let unlockError = $state('');

	async function fetchLockedFolders() {
		isLoading = true;
		try {
			const res = await fetch(`${appState.apiBaseUrl}/api/locked-folders`);
			if (res.ok) {
				folders = await res.json();
			}
		} catch (e) {
			console.warn('API error fetching locked folders:', e);
		} finally {
			isLoading = false;
		}
	}

	$effect(() => {
		fetchLockedFolders();
	});

	function openCreateModal() {
		step = 1;
		folderName = '';
		folderDesc = '';
		passcode = '';
		passcodeConfirm = '';
		passcodeError = '';
		showCreateModal = true;
	}

	function goToStep2(e: Event) {
		e.preventDefault();
		if (!folderName) return;
		step = 2;
	}

	async function handleCreateLockedFolder(e: Event) {
		e.preventDefault();
		passcodeError = '';

		if (passcode.length < 4) {
			passcodeError = 'Passcode must be at least 4 digits';
			return;
		}

		if (passcode !== passcodeConfirm) {
			passcodeError = 'Passcodes do not match! Please confirm your passcode.';
			return;
		}

		try {
			const res = await fetch(`${appState.apiBaseUrl}/api/locked-folders`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					name: folderName,
					description: folderDesc,
					passcode: passcode
				})
			});

			if (res.ok) {
				showCreateModal = false;
				await fetchLockedFolders();
			} else {
				passcodeError = 'Failed to create locked folder.';
			}
		} catch (err) {
			console.error('Error creating locked folder:', err);
			passcodeError = 'Network error while creating folder.';
		}
	}

	async function handleUnlockFolder(e: Event) {
		e.preventDefault();
		if (!selectedFolder) return;
		unlockError = '';

		try {
			const res = await fetch(`${appState.apiBaseUrl}/api/locked-folders/${selectedFolder.id}/verify`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ passcode: unlockPasscode })
			});

			if (res.ok) {
				isUnlocked = true;
				unlockPasscode = '';
			} else {
				unlockError = 'Incorrect passcode! Access denied.';
			}
		} catch (err) {
			console.error('Error verifying passcode:', err);
			unlockError = 'Network error.';
		}
	}

	async function handleDeleteFolder(id: string, name: string, e: Event) {
		e.stopPropagation();
		if (!confirm(`Delete locked folder "${name}"?`)) return;

		try {
			const res = await fetch(`${appState.apiBaseUrl}/api/locked-folders/${id}`, { method: 'DELETE' });
			if (res.ok) {
				await fetchLockedFolders();
			}
		} catch (err) {
			console.error('Error deleting locked folder:', err);
		}
	}
</script>

<div class="space-y-6 animate-fade-in">
	<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
		<div>
			<h1 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
				<Lock class="w-6 h-6 text-sky-500" />
				Locked Vault & Folders
			</h1>
			<p class="text-xs text-slate-500 dark:text-slate-400">Passcode-protected AES encrypted storage (`lockedfolder/...` in MinIO)</p>
		</div>

		<button
			type="button"
			onclick={openCreateModal}
			class="px-4 py-2 rounded-xl bg-sky-400 text-white text-xs font-bold flex items-center gap-1.5 shadow-sm shadow-sky-300/50 hover:bg-sky-500 transition-all cursor-pointer"
		>
			<Plus class="w-4 h-4" />
			New Locked Folder
		</button>
	</div>

	{#if !isLoading && folders.length === 0}
		<div class="p-12 text-center rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm max-w-md mx-auto my-12 space-y-3">
			<div class="w-12 h-12 rounded-2xl bg-sky-100 dark:bg-sky-950 text-sky-500 flex items-center justify-center mx-auto">
				<Lock class="w-6 h-6" />
			</div>
			<h3 class="text-base font-bold text-slate-900 dark:text-white">No locked folders created</h3>
			<p class="text-xs text-slate-500 dark:text-slate-400">Click <strong class="text-sky-500">New Locked Folder</strong> to create a passcode-protected private folder.</p>
		</div>
	{:else}
		<!-- Folders Grid -->
		<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
			{#each folders as folder}
				<div
					role="button"
					tabindex="0"
					onclick={() => { selectedFolder = folder; isUnlocked = false; unlockPasscode = ''; unlockError = ''; }}
					onkeydown={(e) => e.key === 'Enter' && (selectedFolder = folder)}
					class="p-5 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:shadow-md transition-all group cursor-pointer relative"
				>
					<div class="flex items-center justify-between mb-3">
						<div class="w-10 h-10 rounded-xl bg-sky-100 dark:bg-sky-950 text-sky-500 flex items-center justify-center">
							<Lock class="w-5 h-5" />
						</div>

						<button
							type="button"
							onclick={(e) => handleDeleteFolder(folder.id, folder.name, e)}
							class="p-1.5 rounded-lg text-slate-400 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/40 transition-colors"
							title="Delete Locked Folder"
						>
							<Trash2 class="w-4 h-4" />
						</button>
					</div>

					<h3 class="text-sm font-bold text-slate-900 dark:text-white">{folder.name}</h3>
					<p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">{folder.photos_count || 0} items • {folder.description || 'Private Vault'}</p>

					<!-- Admin Audit Oversight Info -->
					<div class="mt-4 pt-3 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between text-[10px] text-slate-400">
						<span class="flex items-center gap-1"><User class="w-3 h-3 text-sky-500" /> {folder.user_name || 'Admin'}</span>
						<span class="flex items-center gap-1"><Clock class="w-3 h-3" /> {folder.created_at || 'Recently'}</span>
					</div>
				</div>
			{/each}
		</div>
	{/if}

	<!-- CREATE LOCKED FOLDER MODAL (2-STEP VALIDATION) -->
	{#if showCreateModal}
		<div class="fixed inset-0 z-50 bg-slate-950/60 backdrop-blur-sm flex items-center justify-center p-4">
			<div class="bg-white dark:bg-slate-900 rounded-3xl p-6 md:p-8 max-w-md w-full border border-slate-200 dark:border-slate-800 shadow-2xl animate-fade-in relative">
				<button
					type="button"
					onclick={() => showCreateModal = false}
					class="absolute top-5 right-5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1"
				>
					<X class="w-5 h-5" />
				</button>

				<div class="flex items-center gap-2 mb-1">
					<Lock class="w-5 h-5 text-sky-500" />
					<h3 class="text-lg font-bold text-slate-900 dark:text-white">New Locked Folder (Step {step} of 2)</h3>
				</div>
				<p class="text-xs text-slate-500 dark:text-slate-400 mb-6">
					{step === 1 ? 'Set folder name and optional description' : 'Set 4-digit passcode & confirmation'}
				</p>

				{#if passcodeError}
					<div class="p-3 mb-4 rounded-xl bg-rose-50 dark:bg-rose-950/60 text-rose-600 dark:text-rose-400 text-xs font-semibold flex items-center gap-2 border border-rose-200 dark:border-rose-800">
						<AlertCircle class="w-4 h-4 shrink-0" />
						<span>{passcodeError}</span>
					</div>
				{/if}

				{#if step === 1}
					<form onsubmit={goToStep2} class="space-y-4">
						<div class="space-y-1.5">
							<label for="folder-name" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Folder Name</label>
							<input
								id="folder-name"
								type="text"
								bind:value={folderName}
								placeholder="e.g. Confidential Financial Docs"
								class="w-full h-10 px-3.5 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400"
								required
							/>
						</div>

						<div class="space-y-1.5">
							<label for="folder-desc" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Description</label>
							<input
								id="folder-desc"
								type="text"
								bind:value={folderDesc}
								placeholder="e.g. Private personal documents"
								class="w-full h-10 px-3.5 text-xs rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400"
							/>
						</div>

						<div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100 dark:border-slate-800">
							<button
								type="button"
								onclick={() => showCreateModal = false}
								class="px-4 py-2 rounded-xl text-xs font-semibold text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800"
							>
								Cancel
							</button>
							<button
								type="submit"
								class="px-5 py-2 rounded-xl bg-sky-400 hover:bg-sky-500 text-white text-xs font-bold shadow-sm shadow-sky-300/50"
							>
								Next: Set Passcode →
							</button>
						</div>
					</form>
				{:else}
					<form onsubmit={handleCreateLockedFolder} class="space-y-4">
						<div class="space-y-1.5">
							<label for="passcode-1" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Set 4-Digit Passcode</label>
							<input
								id="passcode-1"
								type="password"
								maxLength={4}
								bind:value={passcode}
								placeholder="••••"
								class="w-full h-11 text-center font-mono text-lg tracking-widest rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400"
								required
							/>
						</div>

						<div class="space-y-1.5">
							<label for="passcode-2" class="text-xs font-semibold text-slate-700 dark:text-slate-300">Confirm 4-Digit Passcode</label>
							<input
								id="passcode-2"
								type="password"
								maxLength={4}
								bind:value={passcodeConfirm}
								placeholder="••••"
								class="w-full h-11 text-center font-mono text-lg tracking-widest rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400"
								required
							/>
						</div>

						<div class="flex items-center justify-between pt-4 border-t border-slate-100 dark:border-slate-800">
							<button
								type="button"
								onclick={() => step = 1}
								class="px-4 py-2 rounded-xl text-xs font-semibold text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800"
							>
								← Back
							</button>
							<button
								type="submit"
								class="px-5 py-2 rounded-xl bg-sky-400 hover:bg-sky-500 text-white text-xs font-bold shadow-sm shadow-sky-300/50"
							>
								Create Locked Folder
							</button>
						</div>
					</form>
				{/if}
			</div>
		</div>
	{/if}

	<!-- UNLOCK FOLDER MODAL -->
	{#if selectedFolder}
		<div class="fixed inset-0 z-50 bg-slate-950/70 backdrop-blur-md flex items-center justify-center p-4">
			<div class="bg-white dark:bg-slate-900 rounded-3xl p-6 md:p-8 max-w-md w-full border border-slate-200 dark:border-slate-800 shadow-2xl animate-fade-in relative text-center">
				<button
					type="button"
					onclick={() => selectedFolder = null}
					class="absolute top-5 right-5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1"
				>
					<X class="w-5 h-5" />
				</button>

				{#if !isUnlocked}
					<div class="w-14 h-14 rounded-2xl bg-sky-100 dark:bg-sky-950 text-sky-500 flex items-center justify-center mx-auto mb-3">
						<Lock class="w-7 h-7" />
					</div>
					<h3 class="text-lg font-bold text-slate-900 dark:text-white">{selectedFolder.name}</h3>
					<p class="text-xs text-slate-500 dark:text-slate-400 mt-1 mb-6">Enter 4-digit passcode to unlock private contents</p>

					{#if unlockError}
						<div class="p-3 mb-4 rounded-xl bg-rose-50 dark:bg-rose-950/60 text-rose-600 dark:text-rose-400 text-xs font-semibold flex items-center justify-center gap-2 border border-rose-200 dark:border-rose-800">
							<AlertCircle class="w-4 h-4 shrink-0" />
							<span>{unlockError}</span>
						</div>
					{/if}

					<form onsubmit={handleUnlockFolder} class="space-y-4">
						<input
							type="password"
							maxLength={4}
							bind:value={unlockPasscode}
							placeholder="••••"
							class="w-44 h-12 mx-auto text-center font-mono text-xl tracking-widest rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-slate-100 focus:outline-none focus:border-sky-400"
							required
						/>

						<div class="pt-2">
							<button
								type="submit"
								class="w-full h-11 rounded-xl bg-sky-400 hover:bg-sky-500 text-white font-bold text-xs shadow-md shadow-sky-300/50"
							>
								Unlock Folder
							</button>
						</div>
					</form>
				{:else}
					<!-- Unlocked Folder Contents -->
					<div class="w-14 h-14 rounded-2xl bg-emerald-100 dark:bg-emerald-950 text-emerald-500 flex items-center justify-center mx-auto mb-3">
						<Unlock class="w-7 h-7" />
					</div>
					<h3 class="text-lg font-bold text-slate-900 dark:text-white">{selectedFolder.name} (Unlocked)</h3>
					<p class="text-xs text-emerald-600 dark:text-emerald-400 mt-1 mb-6">AES-256 decrypted access granted</p>

					<div class="p-6 rounded-2xl bg-slate-50 dark:bg-slate-800 text-xs text-slate-600 dark:text-slate-300 space-y-2 text-left">
						<p class="font-bold text-slate-800 dark:text-white">Folder Metadata Audit:</p>
						<p>• Created By: {selectedFolder.user_name || 'Admin'}</p>
						<p>• Created Date: {selectedFolder.created_at || 'Recently'}</p>
						<p>• Total Protected Items: {selectedFolder.photos_count || 0}</p>
					</div>

					<div class="pt-6">
						<button
							type="button"
							onclick={() => selectedFolder = null}
							class="px-6 py-2.5 rounded-xl bg-slate-200 dark:bg-slate-800 text-slate-800 dark:text-slate-200 font-bold text-xs"
						>
							Lock & Close Folder
						</button>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
