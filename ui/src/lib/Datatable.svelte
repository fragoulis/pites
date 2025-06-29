<script lang="ts">
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Checkbox
	} from 'flowbite-svelte';
	import DatatableRow from '$lib/DatatableRow.svelte';
	import { type DatatableColumns } from '$lib/types';

	export let columns: DatatableColumns;
	export let records: any[] = [];
	export let loading: boolean = false;
	export let selectacble: boolean = false;
	export let selectedRows: Set<string>;

	let selectAll: boolean = false;

	$: headers = Object.keys(columns);
	$: fields = Object.values(columns);

	// Hack to reset selected rows when records change.
	let lastRecords = records;
	$: if (records !== lastRecords) {
		lastRecords = records;
		selectedRows.clear();
		selectAll = false;
	}

	const onSelectableChangeState = (e: any) => {
		if (e.detail.selected) {
			selectedRows.add(e.detail.id);
		} else {
			selectedRows.delete(e.detail.id);
		}
	};
</script>

<Table striped={true} hoverable={true}>
	<TableHead>
		<TableHeadCell>#</TableHeadCell>
		{#if selectacble}
			<TableHeadCell>
				<Checkbox bind:checked={selectAll} />
			</TableHeadCell>
		{/if}
		{#each headers as header}
			<TableHeadCell>{header}</TableHeadCell>
		{/each}
	</TableHead>
	<TableBody tableBodyClass="divide-y">
		{#if !loading && records.length === 0}
			<TableBodyRow>
				<TableBodyCell colspan={headers.length + 1} class="text-center py-5">
					Ξεκινήστε μία αναζήτηση για να δείτε αποτελέσματα
				</TableBodyCell>
			</TableBodyRow>
		{:else if records.length === 0 && !loading}
			<TableBodyRow>
				<TableBodyCell colspan={headers.length + 1} class="text-center py-5">
					Δεν υπάρχουν αποτελέσματα
				</TableBodyCell>
			</TableBodyRow>
		{:else}
			{#each records as record, i (i)}
				<DatatableRow
					index={i + 1}
					{record}
					{fields}
					selected={selectAll}
					{selectacble}
					on:change={onSelectableChangeState}
				/>
			{:else}
				@
			{/each}
		{/if}
	</TableBody>
</Table>
