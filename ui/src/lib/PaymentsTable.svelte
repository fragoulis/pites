<script lang="ts">
	import MemberNamePaymentsTableColumn from '$lib/MemberNamePaymentsTableColumn.svelte';
	import IssueDatePaymentsTableColumn from '$lib/IssueDatePaymentsTableColumn.svelte';
	import Datatable from '$lib/Datatable.svelte';
	import DatatableSearchForm from '$lib/DatatableSearchForm.svelte';
	import { type DatatableColumns, type Payment } from '$lib/types';

	const availableColumns: DatatableColumns = {
		'Ημ/νία': IssueDatePaymentsTableColumn,
		Μέλος: MemberNamePaymentsTableColumn,
		Ποσό: 'amount',
		Μπλοκ: 'receipt_block_no',
		Απόδειξη: 'receipt_no',
		Μήνες: 'months'
	};
	let selectedColumns: {} = {
		'Ημ/νία': IssueDatePaymentsTableColumn,
		Μέλος: MemberNamePaymentsTableColumn,
		Ποσό: 'amount',
		Μπλοκ: 'receipt_block_no',
		Απόδειξη: 'receipt_no',
		Μήνες: 'months'
	};
	let records: any[] = [];
</script>

<div class="my-5">
	<DatatableSearchForm
		{availableColumns}
		bind:selectedColumns
		bind:records
		collection="payments"
		placeholder="Αναζήτηση βάσει μέλους (ονομα, email, αρ. μητρώου, τηλεφωνο), ημερομηνίας (πχ 2024-11-08)"
	/>
</div>

<Datatable bind:records bind:columns={selectedColumns}></Datatable>
