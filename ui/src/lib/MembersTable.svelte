<script lang="ts">
	import MemberNameTableColumn from '$lib/MemberNameTableColumn.svelte';
	import MemberCompanyNameTableColumn from '$lib/MemberCompanyNameTableColumn.svelte';
	import PaymentStatusTableColumn from '$lib/PaymentStatusTableColumn.svelte';
	import Datatable from '$lib/Datatable.svelte';
	import MemberSearchForm from '$lib/MemberSearchForm.svelte';
	import DatatableSearchForm from '$lib/DatatableSearchForm.svelte';
	import { type DatatableColumns } from '$lib/types';

	const availableColumns: DatatableColumns = {
		'Α/Μ': 'member_no',
		Όνομα: MemberNameTableColumn,
		Διεύθυνση: 'address_formatted',
		Κινητό: 'mobile',
		Email: 'email',
		Συνδρομή: 'subscription_formatted',
		Οικονομικά: PaymentStatusTableColumn,
		Ομάδα: 'business_type_name',
		Εταιρεία: MemberCompanyNameTableColumn,
		Παράρτημα: 'company_branch_name',
		'Δ/ση Εταιρείας': 'company_address',
		ΑΔΤ: 'id_card_number'
	};
	let selectedColumns: {} = {
		'Α/Μ': 'member_no',
		Όνομα: MemberNameTableColumn,
		Διεύθυνση: 'address_formatted',
		Συνδρομή: 'subscription_formatted',
		Οικονομικά: PaymentStatusTableColumn,
		Εταιρεία: MemberCompanyNameTableColumn
	};
	let records: any[] = [];
</script>

<div class="my-5">
	<DatatableSearchForm
		{availableColumns}
		bind:selectedColumns
		bind:records
		collection="members"
		placeholder="Αναζήτηση βάσει ονόματος, αριθμό μητρώου, email, τηλεφώνου"
		searchForm={MemberSearchForm}
		showExport={true}
	/>
</div>

<Datatable bind:records bind:columns={selectedColumns}></Datatable>
