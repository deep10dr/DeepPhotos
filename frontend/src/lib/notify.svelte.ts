/**
 * DeepPhotos Reusable Toast & Confirmation Dialog System (Svelte 5)
 *
 * Usage for Toast Notifications:
 *   import { notify } from '$lib/notify.svelte';
 *   notify.success('Photo uploaded!');
 *   notify.error('Failed to save.');
 *
 * Usage for Yes / No Confirmation Dialogs:
 *   import { confirmDialog } from '$lib/notify.svelte';
 *   if (await confirmDialog.ask('Delete this item?')) {
 *       // User clicked Yes
 *   }
 */

export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface ToastMessage {
	id: string;
	type: ToastType;
	title?: string;
	message: string;
	duration?: number;
}

class NotifyState {
	toasts = $state<ToastMessage[]>([]);

	show(message: string, type: ToastType = 'info', title?: string, duration: number = 4000) {
		const id = Math.random().toString(36).substring(2, 9);
		const toast: ToastMessage = { id, type, title, message, duration };

		this.toasts = [...this.toasts, toast];

		if (duration > 0) {
			setTimeout(() => {
				this.remove(id);
			}, duration);
		}

		return id;
	}

	success(message: string, title: string = 'Success', duration: number = 4000) {
		return this.show(message, 'success', title, duration);
	}

	error(message: string, title: string = 'Error', duration: number = 5000) {
		return this.show(message, 'error', title, duration);
	}

	warning(message: string, title: string = 'Warning', duration: number = 4500) {
		return this.show(message, 'warning', title, duration);
	}

	info(message: string, title: string = 'Info', duration: number = 4000) {
		return this.show(message, 'info', title, duration);
	}

	remove(id: string) {
		this.toasts = this.toasts.filter(t => t.id !== id);
	}

	clearAll() {
		this.toasts = [];
	}
}

export const notify = new NotifyState();

// ----------------------------------------------------
// Yes / No Confirmation Dialog System
// ----------------------------------------------------

export interface ConfirmOptions {
	title?: string;
	message: string;
	confirmText?: string;
	cancelText?: string;
	type?: 'danger' | 'warning' | 'info';
}

class ConfirmState {
	active = $state(false);
	options = $state<ConfirmOptions>({ message: '' });
	private resolvePromise: ((value: boolean) => void) | null = null;

	ask(options: ConfirmOptions | string): Promise<boolean> {
		const opts: ConfirmOptions = typeof options === 'string' ? { message: options } : options;
		this.options = {
			title: opts.title || 'Confirm Action',
			message: opts.message,
			confirmText: opts.confirmText || 'Yes',
			cancelText: opts.cancelText || 'No',
			type: opts.type || 'danger'
		};
		this.active = true;

		return new Promise<boolean>((resolve) => {
			this.resolvePromise = resolve;
		});
	}

	confirm() {
		this.active = false;
		this.resolvePromise?.(true);
		this.resolvePromise = null;
	}

	cancel() {
		this.active = false;
		this.resolvePromise?.(false);
		this.resolvePromise = null;
	}
}

export const confirmDialog = new ConfirmState();
