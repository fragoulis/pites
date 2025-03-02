import { writable } from 'svelte/store';

export const activePayment = writable({
	receipt_block_no: 0,
	receipt_no: 0,
	issued_at: ''
});
export const sharedToast = writable({ message: '', show: false, success: true });
