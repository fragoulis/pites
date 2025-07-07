<script lang="ts">
	import PaymentSearchForm from '$lib/payment/SearchForm.svelte';
	import MemberNamePaymentsTableColumn from '$lib/MemberNamePaymentsTableColumn.svelte';
	import IssueDatePaymentsTableColumn from '$lib/IssueDatePaymentsTableColumn.svelte';
	import ReceiptPaymentsTableColumn from '$lib/payment/ReceiptPaymentsTableColumn.svelte';
	import Datatable from '$lib/Datatable.svelte';
	import DatatableSearchForm from '$lib/DatatableSearchForm.svelte';
	import { type DatatableColumns, type Payment } from '$lib/types';
	import PaymentsTableActions from '$lib/payment/TableActions.svelte';

	const availableColumns: DatatableColumns = {
		'Ημ/νία': IssueDatePaymentsTableColumn,
		Μέλος: MemberNamePaymentsTableColumn,
		Ποσό: 'amount',
		Απόδειξη: ReceiptPaymentsTableColumn,
		Μήνες: 'months'
	};
	let selectedColumns: {} = {
		'Ημ/νία': IssueDatePaymentsTableColumn,
		Μέλος: MemberNamePaymentsTableColumn,
		Ποσό: 'amount',
		Απόδειξη: ReceiptPaymentsTableColumn,
		Μήνες: 'months'
	};
	let records: any[] = [];
	let selectedRows: Map<string, any> = new Map<string, any>();
</script>

<div class="my-5">
	<DatatableSearchForm
		{availableColumns}
		bind:selectedColumns
		bind:records
		bind:selectedRows
		searchForm={PaymentSearchForm}
		actions={PaymentsTableActions}
		collection="payments"
		placeholder="Αναζήτηση βάσει μέλους (ονομα, email, αρ. μητρώου, τηλεφωνο)"
	/>
</div>

<Datatable bind:records bind:columns={selectedColumns} bind:selectedRows selectacble={true}
></Datatable>
