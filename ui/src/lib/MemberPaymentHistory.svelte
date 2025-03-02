<script lang="ts">
	import {
		Table,
		TableBody,
		TableBodyRow,
		TableBodyCell,
		TableHead,
		TableHeadCell,
		A
	} from 'flowbite-svelte';
	import { type Member } from '$lib/types';

	export let record: Member;
	export let short: boolean = false;

	let tdClass: string = short ? 'px-6 py-2' : '';
</script>

<Table striped={true} hoverable={true}>
	{#if !short}
		<TableHead>
			<TableHeadCell>#</TableHeadCell>
			<TableHeadCell>Πότε</TableHeadCell>
			<TableHeadCell>Πόσο</TableHeadCell>
			<TableHeadCell>Μήνες</TableHeadCell>
			<TableHeadCell>Αρ. Μπλοκ</TableHeadCell>
			<TableHeadCell>Απόδειξη</TableHeadCell>
			<TableHeadCell>Σχόλια</TableHeadCell>
			<TableHeadCell>Εώς</TableHeadCell>
			<TableHeadCell>Ταμίας</TableHeadCell>
		</TableHead>
	{/if}
	<TableBody>
		{#each record.payments as payment, i (i)}
			{#if !short || payment.amount > 0}
				<TableBodyRow>
					<TableBodyCell tdClass="">{i + 1}.</TableBodyCell>
					<TableBodyCell {tdClass}>
						<A href={`/admin/payment/edit?id=${payment.id}`} target="_blank">
							{payment.issued_at_formatted}
						</A>
					</TableBodyCell>
					<TableBodyCell {tdClass}>{payment.amount} €</TableBodyCell>
					{#if !short}
						<TableBodyCell>{payment.months}</TableBodyCell>
						<TableBodyCell>{payment.receipt_block_no}</TableBodyCell>
						<TableBodyCell>{payment.receipt_no}</TableBodyCell>
						<TableBodyCell>{payment.comments}</TableBodyCell>
						<TableBodyCell>{payment.legacy_to_formatted}</TableBodyCell>
						<TableBodyCell>{payment.created_by_user?.username || ''}</TableBodyCell>
					{/if}
				</TableBodyRow>
			{/if}
		{:else}
			<TableBodyRow>
				<TableBodyCell colspan="5">Δε βρέθηκαν προηγούμενες πληρωμές</TableBodyCell>
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
