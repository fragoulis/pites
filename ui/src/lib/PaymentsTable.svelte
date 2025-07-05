<script lang="ts">
	import PaymentSearchForm from '$lib/payment/SearchForm.svelte';
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
	let selectedRows: Set<string> = new Set<string>();
</script>

<div class="my-5">
	<DatatableSearchForm
		{availableColumns}
		bind:selectedColumns
		bind:records
		bind:selectedRows
		searchForm={PaymentSearchForm}
		collection="payments"
		placeholder="Αναζήτηση βάσει μέλους (ονομα, email, αρ. μητρώου, τηλεφωνο)"
	/>
</div>

<Datatable bind:records bind:columns={selectedColumns} bind:selectedRows selectacble={true}
></Datatable>
