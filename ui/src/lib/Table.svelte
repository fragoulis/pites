<script lang="ts">
	import { onMount, tick } from 'svelte';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell
	} from 'flowbite-svelte';
	import pb from '$lib/pocketbase';
	import { isComponent } from '$lib/utils';

	export let collection: string;
	export let headers: string[];
	export let fields: string[];
	export let sortBy: string = '-created';

	let loading: boolean = false;
	let records: any = [];

	onMount(async () => {
		await tick();
		await performSearch();
	});

	const performSearch = async () => {
		loading = true;

		pb.cancelAllRequests();

		try {
			records = await pb.collection(collection).getFullList({
				sort: sortBy
			});
		} catch (err: any) {
			console.error(err);
			records = [];
		}

		loading = false;

		return records;
	};
</script>

<Table striped={true} hoverable={true}>
	<TableHead>
		<TableHeadCell>#</TableHeadCell>
		{#each headers as header}
			<TableHeadCell>{header}</TableHeadCell>
		{/each}
	</TableHead>
	<TableBody tableBodyClass="divide-y">
		{#if records.length === 0 && !loading}
			<TableBodyRow>
				<TableBodyCell colspan={headers.length + 1} class="text-center py-5">
					Δεν υπάρχουν αποτελέσματα
				</TableBodyCell>
			</TableBodyRow>
		{:else}
			{#each records as record, i (i)}
				<TableBodyRow>
					<TableBodyCell>{i + 1}</TableBodyCell>
					{#each fields as field}
						<TableBodyCell>
							{#if isComponent(field)}
								<svelte:component this={field} {record} />
							{:else if typeof field === 'function'}
								{field(record)}
							{:else if typeof record[field] === 'boolean'}
								{#if record[field]}
									Ναι
								{:else}
									Όχι
								{/if}
							{:else}
								{record[field]}
							{/if}
						</TableBodyCell>
					{/each}
				</TableBodyRow>
			{/each}
		{/if}
	</TableBody>
</Table>
