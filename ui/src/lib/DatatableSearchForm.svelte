<script lang="ts">
	import InputGroup from '$lib/InputGroup.svelte';
	import { onMount, tick } from 'svelte';
	import { page } from '$app/stores';
	import { Search, Button, Dropdown, Toggle } from 'flowbite-svelte';
	import pb from '$lib/pocketbase';
	import { browser } from '$app/environment';
	import { type DatatableSearchForm } from '$lib/types';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';

	export let collection: string;
	export let placeholder: string = 'Αναζήτηση';
	export let loading: boolean = false;
	export let records: any[] = [];
	export let searchForm: any = null;
	export let actions: any = null;
	export let showExport: boolean = false;
	export let availableColumns: {};
	export let selectedColumns: {};
	export let selectedRows: Set<string>;

	let form: DatatableSearchForm = { active_only: true };
	let total: number = 0;
	let requestKey: string = 'datatablesearchform';

	$: performSearch(form);

	onMount(async () => {
		form.q = $page.url.searchParams.get('q') || undefined;

		await performSearch(form);
	});

	const performSearch = async (data: DatatableSearchForm) => {
		if (!browser) {
			return;
		}

		await tick();

		Object.keys(data).forEach((key) => {
			if (data[key] === undefined) {
				delete data[key];
				return;
			}

			if (key == 'name' && data[key] == '') {
				delete data[key];
				return;
			}
		});

		pb.cancelRequest(requestKey);

		loading = true;
		try {
			records = [];

			const res = await pb.send(`/${collection}`, {
				query: data,
				requestKey: requestKey,
				cache: 'no-cache'
			});

			if (res.records) {
				records = res.records;
				total = res.total;
			} else {
				// Weird error breaks the datatable
				records = [];
				total = 0;
			}
		} catch (err: any) {
			records = [];
		}
		loading = false;
	};

	let searchFormRef: any;

	const onReset = async () => {
		form.q = '';
		searchFormRef?.reset();
	};

	const performExport = async () => {
		let data: DatatableSearchForm = form;

		Object.keys(data).forEach((key) => {
			if (data[key] === undefined) {
				delete data[key];
				return;
			}

			if (key == 'name' && data[key] == '') {
				delete data[key];
				return;
			}
		});

		data = Object.assign({}, data, { columns: Object.keys(selectedColumns) });

		pb.cancelRequest(requestKey);

		loading = true;
		try {
			const res = await pb.send(`/${collection}/export`, {
				method: 'POST',
				body: data,
				requestKey: requestKey,
				fetch: async (url, config) => {
					const response = await fetch(url, config);
					const blob = await response.blob();

					const windowURL = window.URL || window.webkitURL;
					const downloadUrl = windowURL.createObjectURL(blob);

					const link = document.createElement('a');
					link.href = downloadUrl;
					link.download = 'export.xlsx';
					document.body.append(link);
					link.click();

					setTimeout(() => {
						document.body.removeChild(link);
						windowURL.revokeObjectURL(downloadUrl);
					}, 100);
				}
			});
		} catch (err: any) {}
		loading = false;
	};

	const toggleColumn = (name: string) => {
		if (name in selectedColumns) {
			delete selectedColumns[name];
		} else {
			// Keep the order in which the columns appear steady.

			// Get currently shown columns.
			const selectedColumnsCopy = selectedColumns;

			// Empty it.
			selectedColumns = {};

			// Iterate over available columns because they have the
			// fixed order.
			for (const [k, v] of Object.entries(availableColumns)) {
				if (k in selectedColumnsCopy || k == name) {
					selectedColumns[k] = v;
				}
			}
		}

		selectedColumns = selectedColumns;
	};
</script>

<form
	id="datatable-search-form"
	action={`/${collection}`}
	autocomplete="off"
	on:submit|preventDefault
>
	<InputGroup legend="Φίλτρα">
		<Search
			autofocus
			{placeholder}
			class="bg-gray-100"
			name="q"
			bind:value={form.q}
			on:focus={(e) => e.target.select()}
		/>
		{#if searchForm}
			<svelte:component this={searchForm} bind:this={searchFormRef} bind:form />
		{/if}
		<div class="flex items-center space-x-4 justify-between">
			<div>Αποτελέσματα: {total}</div>
			<div class="flex items-center space-x-4 justify-between">
				{#if actions}
					<svelte:component this={actions} {selectedRows} />
				{/if}
				<Button outline pill color="light">Στήλες<ChevronDownOutline /></Button>
				<Dropdown class="w-56 p-3 space-y-1">
					{#each Object.keys(availableColumns) as name}
						<li>
							<Toggle
								class="rounded p-2 hover:bg-gray-100"
								checked={name in selectedColumns}
								on:change={() => toggleColumn(name)}
							>
								{name}
							</Toggle>
						</li>
					{/each}
				</Dropdown>

				{#if showExport}
					<Button disabled={loading} color="yellow" on:click={performExport}>Εξαγωγή</Button>
				{/if}
				<Button color="light" on:click={onReset}>Εκκαθάρριση</Button>
			</div>
		</div>
	</InputGroup>
</form>
