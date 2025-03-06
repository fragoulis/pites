<script lang="ts">
	import {
		Input,
		InputAddon,
		Dropdown,
		DropdownItem,
		type FormSizeType,
		ButtonGroup,
		Button
	} from 'flowbite-svelte';
	import { SearchOutline, CloseOutline, RefreshOutline } from 'flowbite-svelte-icons';
	import { createEventDispatcher } from 'svelte';
	import pb from '$lib/pocketbase';

	export let id = '';
	export let query = '';
	export let placeholder: string = 'Πληκτρολογήστε για αναζήτηση';
	export let size: FormSizeType = 'md';
	export let url: string;
	export let open: boolean = false;

	const dispatch = createEventDispatcher();

	let loading: boolean = false;
	let records: any = [];
	let highlightedRecordIdx: number = 0;

	$: showClearBtn = query != '';

	const requestKey = `requestkey${url}`;

	const performSearch = async () => {
		highlightedRecordIdx = 0;

		if (query == '') {
			records = [];
		}

		loading = true;

		pb.cancelRequest(requestKey);

		try {
			const res = await pb.send(url, {
				requestKey: requestKey,
				query: { q: query }
			});

			if (res.records) {
				records = res.records;
			} else {
				records = res;
			}
			open = true;
		} catch (err: any) {
			records = [];
		}

		loading = false;
	};

	const clear = (e: any) => {
		e.preventDefault();
		query = '';
		open = false;
		records = [];

		dispatch('clear');
	};

	const onSelect = (i: number) => {
		dispatch('select', { item: records[i] });
	};

	const onKeyUp = (e: any) => {
		if (e.key == 'ArrowDown' || e.key == 'ArrowUp' || e.key == 'Enter') {
			if (e.repeat) return;
			if (!records || records.length == 0) return;

			switch (e.key) {
				case 'ArrowDown':
					highlightedRecordIdx++;
					if (highlightedRecordIdx >= records.length) {
						highlightedRecordIdx = 0;
					}
					break;
				case 'ArrowUp':
					highlightedRecordIdx--;
					if (highlightedRecordIdx < 0) {
						highlightedRecordIdx = records.length - 1;
					}
					break;
				case 'Enter':
					onSelect(highlightedRecordIdx);
					break;
			}
		} else {
			performSearch();
		}
	};
</script>

<ButtonGroup class="w-full" {size}>
	<InputAddon>
		<SearchOutline class="w-6 h-6 text-gray-500 dark:text-gray-400" />
	</InputAddon>
	<Input
		{id}
		{placeholder}
		bind:value={query}
		autocomplete="off"
		on:focus={(e) => e.target?.select()}
		on:keyup={onKeyUp}
	/>
	{#if showClearBtn}
		<Button on:click={clear} color="red"><CloseOutline /></Button>
	{/if}
</ButtonGroup>

<Dropdown bind:open placement="bottom-start">
	{#each records as record, i}
		<DropdownItem
			defaultClass={'font-medium py-2 px-4 text-sm hover:bg-gray-100 ' +
				(i == highlightedRecordIdx ? 'bg-gray-100' : '')}
			on:click={() => onSelect(i)}
		>
			<slot {record} />
		</DropdownItem>
	{:else}
		{#if query != ''}
			<DropdownItem>
				<slot name="empty">Δεν υπάρχουν αποτελέσματα</slot>
			</DropdownItem>
		{/if}
	{/each}
</Dropdown>
