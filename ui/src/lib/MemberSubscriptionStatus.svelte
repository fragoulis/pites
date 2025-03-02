<script lang="ts">
	import { type Member } from '$lib/types';
	import { Button, Modal } from 'flowbite-svelte';
	import pb from '$lib/pocketbase';
	import { createEventDispatcher } from 'svelte';
	import ToggleField from './ToggleField.svelte';
	const dispatch = createEventDispatcher();

	export let record: Member;

	let activationModal: boolean = false;
	let fee_paid: boolean = false;

	const deactivate = async () => {
		const yes = confirm('Είσαι σίγουρη;');
		if (!yes) return;

		try {
			const res = await pb.send(`/members/${record.id}/subscriptions`, {
				method: 'DELETE'
			});

			dispatch('success');
		} catch (err: any) {
			dispatch('failure', { message: err.message });
		}
	};

	const activate = async () => {
		try {
			const res = await pb.send(`/members/${record.id}/subscriptions`, {
				method: 'POST',
				body: {
					fee_paid: fee_paid
				}
			});

			dispatch('success');
		} catch (err: any) {
			dispatch('failure', { message: err.message });
		}
	};
</script>

{record.subscription_formatted}

{#if record.subscription_active}
	<Button class="ml-5" color="red" on:click={deactivate}>Απενεργοποίηση</Button>
{:else}
	<Button
		class="ml-5"
		color="green"
		on:click={() => {
			activationModal = true;
		}}
	>
		Ενεργοποίηση
	</Button>
	<Modal size="xs" title="Ενεργοποίση" bind:open={activationModal} autoclose outsideclose>
		<ToggleField bind:checked={fee_paid} label="Έχει πληρώσει εγγραφή;" />
		<Button color="green" on:click={activate}>Ενεργοποίηση</Button>
	</Modal>
{/if}
