<script lang="ts">
	import { Drawer, CloseButton } from 'flowbite-svelte';
	import { UserCircleOutline } from 'flowbite-svelte-icons';
	import { sineIn } from 'svelte/easing';
	import ReceiptCreateFormForBatchPayments from '$lib/payment/ReceiptCreateFormForBatchPayments.svelte';
	import { sharedToast } from '$lib/store';

	export let hidden = true;
	export let payments: Map<string, any>;

	let transitionParams = {
		x: 320,
		duration: 200,
		easing: sineIn
	};

	const onSuccess = () => {
		$sharedToast.show = true;
		$sharedToast.success = true;
		$sharedToast.message = 'Οι αποδείξεις δημιουργήθηκαν επιτυχώς.';

		hidden = true;
	};

	const onError = async (e: any) => {
		$sharedToast.show = true;
		$sharedToast.success = false;
		$sharedToast.message = e.detail?.message || e.message;
	};
</script>

<Drawer placement="right" transitionType="fly" {transitionParams} bind:hidden width="w-2/3">
	<div class="flex items-center">
		<h5
			id="drawer-label"
			class="inline-flex items-center mb-4 text-xl font-bold text-gray-800 dark:text-gray-400"
		>
			<UserCircleOutline class="w-5 h-5 me-2.5" />Δημιουργία αποδείξεων για τις επιλεγμένες πληρωμές
		</h5>
		<CloseButton on:click={() => (hidden = true)} class="mb-4 dark:text-white" />
	</div>

	<ReceiptCreateFormForBatchPayments bind:payments on:success={onSuccess} on:failure={onError} />
</Drawer>
