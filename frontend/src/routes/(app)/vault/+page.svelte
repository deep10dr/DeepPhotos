<script lang="ts">
	import { onMount } from 'svelte';
	import { appState } from '$lib/state.svelte';
	import { apiFetch, getMediaUrl } from '$lib/api';
	import { notify, confirmDialog } from '$lib/notify.svelte';
	import { Lock, Unlock, Key, ShieldCheck, FolderClosed, Plus, X, Trash2, User, Clock, AlertCircle, FileText, Image as ImageIcon, Video, Play, Download } from 'lucide-svelte';

	interface LockedFolder {
		id: string;
		user_id: string;
		user_name: string;
		name: string;
		description: string;
		photos_count: number;
		created_at: string;
	}

	interface MediaItem {
		id: string;
		title: string;
		filename: string;
		mime_type: string;
		file_type: string;
		size: number;
		url?: string;
		thumbnail_url?: string;
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
	let lockedMediaItems = $state<MediaItem[]>([]);
	let isFetchingMedia = $state(false);

	async function fetchLockedFolders() {
		isLoading = true;
		try {
			const res = await apiFetch('/api/locked-folders');
			if (res.ok) {
				folders = await res.json();
			}
		} catch (e) {
			console.warn('API error fetching locked folders:', e);
		} finally {
			isLoading = false;
		}
	}

	async function fetchLockedFolderMedia(folderId: string) {
		isFetchingMedia = true;
		try {
			const res = await apiFetch(`/api/photos?locked_folder_id=${folderId}`);
			if (res.ok) {
				const data = await res.json();
				lockedMediaItems = data.map((item: any) => ({
					...item,
					url: getMediaUrl(`/api/photos/${item.id}/file`),
					thumbnail_url: getMediaUrl(`/api/photos/${item.id}/thumbnail`)
				}));
			}
		} catch (e) {
			console.warn('API error fetching locked media:', e);
		} finally {
			isFetchingMedia = false;
		}
	}

	onMount(() => {
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
			const res = await apiFetch('/api/locked-folders', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					name: folderName,
					description: folderDesc,
					passcode: passcode
				})
			});

			if (res.ok) {
				const newFolder = await res.json();
				folders = [...folders, newFolder];
				showCreateModal = false;
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
			const res = await apiFetch(`/api/locked-folders/${selectedFolder.id}/verify`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ passcode: unlockPasscode })
			});

			if (res.ok) {
				isUnlocked = true;
				unlockPasscode = '';
				await fetchLockedFolderMedia(selectedFolder.id);
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
		const confirmed = await confirmDialog.ask({
			title: 'Delete Locked Vault Folder',
			message: `Delete locked folder "${name}"?`,
			confirmText: 'Yes, Delete Folder',
			cancelText: 'Cancel',
			type: 'danger'
		});
		if (!confirmed) return;

		try {
			const res = await apiFetch(`/api/locked-folders/${id}`, { method: 'DELETE' });
			if (res.ok) {
				notify.success(`Locked folder "${name}" deleted.`);
				await fetchLockedFolders();
			} else {
				notify.error('Failed to delete locked folder.');
			}
		} catch (err) {
			console.error('Error deleting locked folder:', err);
			notify.error('Network error deleting locked folder.');
		}
	}

	function formatBytes(bytes: number): string {
		if (!bytes) return '0 B';
		const k = 1024;
		const sizes = ['B', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
	}
</script>

<div class="h-full flex flex-col gap-6 animate-fade-in">
	<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 shrink-0">
		<div>
			<h1 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
				<Lock class="w-6 h-6 text-sky-500" />
				Vault
			</h1>
		</div>

		<button
			type="button"
			onclick={openCreateModal}
			class="px-4 py-2 rounded-xl bg-sky-400 text-white text-xs font-bold flex items-center gap-1.5 shadow-sm shadow-sky-300/50 hover:bg-sky-500 transition-all cursor-pointer shrink-0"
		>
			<Plus class="w-4 h-4" />
			New Locked Folder
		</button>
	</div>

	<!-- Scrollable Content Area — only this scrolls, not the whole page -->
	<div class="flex-1 overflow-y-auto min-h-0">

	<!-- EMPTY STATE WITH EMOJI -->
	{#if !isLoading && folders.length === 0}
		<div class="p-12 text-center rounded-3xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm max-w-md mx-auto my-12 space-y-3">
			<img src="/empty-folder.png" alt="Empty" class="w-16 h-16 mx-auto mb-3 opacity-60 dark:opacity-40 select-none pointer-events-none drop-shadow-sm" />
			<h3 class="text-base font-bold text-slate-900 dark:text-white">This locked vault is empty</h3>
			<p class="text-xs text-slate-500 dark:text-slate-400">Click <strong class="text-sky-500">New Locked Folder</strong> to create a passcode-protected private folder.</p>
		</div>
	{:else}
		<!-- Folders Grid -->
		<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
			{#each folders as folder}
				<div
					role="button"
					tabindex="0"
					onclick={() => { selectedFolder = folder; isUnlocked = false; unlockPasscode = ''; unlockError = ''; lockedMediaItems = []; }}
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

	</div><!-- end scrollable content -->

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

	<!-- UNLOCK FOLDER & VIEW CONTENTS MODAL -->
	{#if selectedFolder}
		<div class="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-md flex items-center justify-center p-4">
			<div class="bg-white dark:bg-slate-900 rounded-3xl p-6 md:p-8 max-w-3xl w-full border border-slate-200 dark:border-slate-800 shadow-2xl animate-fade-in relative text-center flex flex-col max-h-[85vh]">
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
					<!-- UNLOCKED FOLDER MEDIA DISPLAY & EMOJI EMPTY STATE -->
					<div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-4 mb-4">
						<div class="flex items-center gap-3">
							<div class="w-10 h-10 rounded-xl bg-emerald-100 dark:bg-emerald-950 text-emerald-500 flex items-center justify-center">
								<Unlock class="w-5 h-5" />
							</div>
							<div class="text-left">
								<h3 class="text-base font-bold text-slate-900 dark:text-white">{selectedFolder.name} (Unlocked)</h3>
								<p class="text-xs text-slate-500 dark:text-slate-400">{lockedMediaItems.length} protected items stored</p>
							</div>
						</div>

						<button
							type="button"
							onclick={() => selectedFolder = null}
							class="px-4 py-2 rounded-xl bg-slate-200 dark:bg-slate-800 text-slate-800 dark:text-slate-200 font-bold text-xs"
						>
							Lock & Close
						</button>
					</div>

					<!-- UNLOCKED FOLDER EMPTY STATE WITH EMOJI -->
					{#if lockedMediaItems.length === 0}
						<div class="p-8 text-center space-y-3 my-auto">
							<img src="/empty-folder.png" alt="Empty" class="w-24 h-24 mx-auto mb-4 opacity-60 dark:opacity-40 select-none pointer-events-none drop-shadow-sm" />
							<h4 class="text-base font-bold text-slate-900 dark:text-white">This locked folder is empty</h4>
							<p class="text-xs text-slate-500 dark:text-slate-400">Open any photo in gallery and click <strong class="text-sky-500">Move to Vault</strong> to add items to this folder.</p>
						</div>
					{:else}
						<!-- MEDIA GRID INSIDE LOCKED FOLDER -->
						<div class="flex-1 overflow-y-auto grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 text-left p-1">
							{#each lockedMediaItems as media}
								<div class="p-3 rounded-2xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 flex flex-col justify-between">
									<div class="aspect-4/3 rounded-xl overflow-hidden bg-slate-950 mb-2 relative flex items-center justify-center">
										{#if media.file_type === 'video' || media.mime_type.startsWith('video/')}
											<video src={media.url} controls class="w-full h-full object-cover"></video>
										{:else if media.file_type === 'document'}
											<FileText class="w-8 h-8 text-sky-400" />
										{:else}
											<img src={media.thumbnail_url || media.url} alt={media.title} class="w-full h-full object-cover" />
										{/if}
									</div>
									<div class="flex items-center justify-between text-xs">
										<span class="font-bold text-slate-800 dark:text-slate-100 truncate">{media.title}</span>
										<a href={media.url} download={media.filename} target="_blank" class="text-sky-500 p-1 hover:bg-sky-50 dark:hover:bg-slate-700 rounded-lg">
											<Download class="w-3.5 h-3.5" />
										</a>
									</div>
								</div>
							{/each}
						</div>
					{/if}
				{/if}
			</div>
		</div>
	{/if}
</div>
