<script lang="ts">
	import { Input, Label, Helper } from 'flowbite-svelte';
	import Autocomplete from '$lib/Autocomplete.svelte';
	import AutocompleteAddressResult from '$lib/AutocompleteAddressResult.svelte';

	export let label: string = 'Διεύθυνση';
	export let value = '';
	export let query = '';
	export let error: string = '';
	export let required: boolean = true;
	export let city_only: boolean = false;

	let open: boolean = false;

	let url: string = '/address';
	if (city_only) url = 'address_cities';

	const select = (e: any) => {
		if (!e.detail.item) return;

		value = e.detail.item.id;
		query = e.detail.item.name;
		open = false;
	};

	const clear = () => {
		value = '';
	};
</script>

<div class="mb-6">
	<div class="relative">
		<Label class="space-y-2" color={error ? 'red' : 'gray'}>
			<span>
				{label}
				{#if required}*{/if}
			</span>
			<Autocomplete
				{url}
				bind:open
				bind:query
				on:select={select}
				on:clear={clear}
				on:new
				let:record
			>
				<AutocompleteAddressResult {record} />
			</Autocomplete>
			{#if error}
				<Helper class="mt-2 text-sm" color="red">
					{error}
				</Helper>
			{/if}
			<Input type="hidden" bind:value />
		</Label>
	</div>
</div>
