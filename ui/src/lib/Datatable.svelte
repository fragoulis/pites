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

	$: headers = Object.keys(columns);
	$: fields = Object.values(columns);

	let checked: boolean = false;
</script>

<Table striped={true} hoverable={true}>
	<TableHead>
		<TableHeadCell>#</TableHeadCell>
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
				<DatatableRow index={i + 1} {record} {fields} {checked} />
			{:else}
				@
			{/each}
		{/if}
	</TableBody>
</Table>
