<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { Input, Label, Helper, Modal, A, Button } from 'flowbite-svelte';
	import Autocomplete from '$lib/Autocomplete.svelte';
	import AutocompleteCompanyResult from '$lib/AutocompleteCompanyResult.svelte';
	import CompanyForm from '$lib/CompanyForm.svelte';
	import { sharedToast } from '$lib/store';
	import { type Company } from './types';

	const dispatch = createEventDispatcher();

	let label: string = 'Επωνμία Επιχείρησης';
	let name: string = 'company_id';
	export let value: string = '';
	export let query: string = '';
	export let error: string = '';
	export let required: boolean = true;
	export let multi: boolean = false;
	export let help: string = '';

	let open: boolean = false;

	// Used for multiselect capabilities.
	let selections: any = {};

	// Used to pass down to the company form the query for better ux.
	let companyRecord: Company = {};

	$: companyRecord.name = query;

	const select = (e: any) => {
		if (!multi) selections = {};
		selections[e.detail.item.id] = e.detail.item.name_formatted;
		value = Object.keys(selections).join(',');

		query = e.detail.item.name_formatted;
		open = false;
		dispatch('select', { model: e.detail.item });
		dispatch('change', { model: e.detail.item });
	};

	const clear = () => {
		selections = {};
		value = '';
		dispatch('clear');
		dispatch('change', selections);
	};

	let newCompanyModal = false;

	const showNewCompanyModal = (e: any) => {
		e.preventDefault();
		newCompanyModal = true;
		dispatch('new');
	};

	const onCompanyCreated = (e: any) => {
		newCompanyModal = false;
		$sharedToast.show = true;
		$sharedToast.success = true;
		$sharedToast.message = 'Η εταιρεία δημιουργήθηκε.';
		query = e.detail.model.name_formatted;
		value = e.detail.model.id;
		dispatch('company:created', { model: e.detail.model });
	};

	const onCompanyFailure = (e: any) => {
		$sharedToast.show = true;
		$sharedToast.success = false;
		$sharedToast.message = e.detail?.message || e.message;
		dispatch('company:failure');
	};
</script>

<div class="mb-4">
	<Label class="space-y-2" color={error ? 'red' : 'gray'}>
		<span>
			{label}
			{#if required}*{/if}
		</span>
		<Autocomplete
			url="/companies"
			bind:open
			bind:query
			on:select={select}
			on:clear={clear}
			let:record
		>
			<AutocompleteCompanyResult {record} />
			<svelte:fragment slot="empty">
				<A class="text-sky-700 dark:text-sky-200" on:click={showNewCompanyModal}>Δημιουργία</A>
			</svelte:fragment>
		</Autocomplete>
		{#if help}
			<Helper class="mt-2 text-sm">
				{help}
			</Helper>
		{/if}
		{#if error}
			<Helper class="mt-2 text-sm" color="red">
				{error}
			</Helper>
		{/if}
		<Input type="hidden" bind:value {name} />
	</Label>

	{#if multi && value != ''}
		{#each Object.entries(selections) as [id, name]}
			<Button
				class="mr-1 mt-5"
				size="xs"
				outline
				pill
				on:click={() => {
					delete selections[id];
					value = Object.keys(selections).join(',');
				}}
			>
				{name}
			</Button>
		{/each}
	{/if}
</div>

<Modal title="Δημιουργία εταιρείας" bind:open={newCompanyModal}>
	<CompanyForm on:success={onCompanyCreated} on:failure={onCompanyFailure} record={companyRecord} />
</Modal>
