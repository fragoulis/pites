<script lang="ts">
	import CompanyBusinessTypeTableColumn from '$lib/CompanyBusinessTypeTableColumn.svelte';
	import CompanyNameTableColumn from '$lib/CompanyNameTableColumn.svelte';
	import CompanySearchForm from '$lib/CompanySearchForm.svelte';
	import Datatable from '$lib/Datatable.svelte';
	import DatatableSearchForm from '$lib/DatatableSearchForm.svelte';
	import { type DatatableColumns } from '$lib/types';

	let availableColumns: DatatableColumns = {
		Όνομα: CompanyNameTableColumn,
		Παράρτημα: 'branch',
		Ομάδα: CompanyBusinessTypeTableColumn,
		Διέυθυνση: 'address_formatted',
		Email: 'email',
		Τηλέφωνο: 'phone',
		Ιστοσελίδα: 'website'
	};
	let selectedColumns: {} = {
		Όνομα: CompanyNameTableColumn,
		Παράρτημα: 'branch',
		Ομάδα: CompanyBusinessTypeTableColumn
	};
	let records: any[] = [];
</script>

<div class="m-5">
	<DatatableSearchForm
		bind:availableColumns
		bind:selectedColumns
		bind:records
		collection="companies"
		searchForm={CompanySearchForm}
	/>
</div>

<Datatable bind:records bind:columns={selectedColumns}></Datatable>
