<script lang="ts">
	import { TableBodyCell, TableBodyRow, Checkbox } from 'flowbite-svelte';

	export let index: number;
	export let record: any;
	export let fields: string[];
	export let checked = false;

	const selectRow = (e: MouseEvent) => {
		checked = !checked;
	};
</script>

<TableBodyRow>
	<TableBodyCell>{index}</TableBodyCell>
	{#each fields as field}
		<TableBodyCell class="break-works whitespace-normal">
			{#if typeof field === 'function'}
				<svelte:component this={field} {record} />
			{:else if typeof field === 'string'}
				{record[field]}
			{/if}
		</TableBodyCell>
	{/each}
</TableBodyRow>
