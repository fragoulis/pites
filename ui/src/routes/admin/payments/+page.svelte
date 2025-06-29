<script lang="ts">
	import PaymentsTable from '$lib/PaymentsTable.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import { Button } from 'flowbite-svelte';
	import { PlusOutline } from 'flowbite-svelte-icons';
	import { isAdmin } from '$lib/utils';

	const headers: string[] = ['Μέλος', 'Ποσό', 'Ημερομηνία', 'Χρήστης'];
	const fields: string[] = ['member_id', 'amount_in_euros', 'issued_at', 'created_by_user_id'];
</script>

<PageHeader>
	Πληρωμές
	<svelte:fragment slot="right">
		{#if isAdmin()}
			<Button href="/admin/payments/create">
				<PlusOutline class="w-4 h-4 me-2" />
				<span class="font-bold">Νέα πληρωμή</span>
			</Button>
		{/if}
	</svelte:fragment>
</PageHeader>

<PaymentsTable />
