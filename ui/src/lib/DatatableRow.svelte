<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { TableBodyCell, TableBodyRow, Checkbox } from 'flowbite-svelte';

	const dispatch = createEventDispatcher();

	export let index: number;
	export let record: any;
	export let fields: string[];
	export let selected = false;
	export let selectacble = false;

	// Hack to trigger the change event when selected is updated.
	let lastSelected = selected;
	$: if (selected !== lastSelected) {
		lastSelected = selected;
		dispatch('change', {
			id: record.id,
			selected: selected
		});
	}
</script>

<TableBodyRow>
	<TableBodyCell>{index}</TableBodyCell>
	{#if selectacble}
		<TableBodyCell>
			<Checkbox checked={selected} on:change={() => (selected = !selected)} value={record.id} />
		</TableBodyCell>
	{/if}
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
