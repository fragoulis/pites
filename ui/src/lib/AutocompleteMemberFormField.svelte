<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { Input, Label, Helper } from 'flowbite-svelte';
	import Autocomplete from '$lib/Autocomplete.svelte';
	import AutocompleteMemberResult from '$lib/AutocompleteMemberResult.svelte';

	const dispatch = createEventDispatcher();

	let label: string = 'Μέλος';
	let name: string = 'member_id';
	export let value = '';
	export let id = '';
	export let query = '';
	export let error: string = '';

	let open: boolean = false;

	const select = (e: any) => {
		value = e.detail.item.id;
		query = e.detail.item.name_formatted;
		open = false;
		dispatch('select', { item: e.detail.item });
	};

	const clear = () => {
		value = '';
		dispatch('clear');
	};
</script>

<div class="mb-6">
	<div class="relative">
		<Label class="space-y-2" color={error ? 'red' : 'gray'}>
			<span>{label} *</span>
			<Autocomplete
				url="/members?active_only=false"
				bind:id
				bind:open
				bind:query
				on:select={select}
				on:clear={clear}
				on:new
				let:record
			>
				<AutocompleteMemberResult {record} />
			</Autocomplete>
			{#if error}
				<Helper class="mt-2 text-sm" color="red">
					{error}
				</Helper>
			{/if}
			<Input type="hidden" bind:value {name} />
		</Label>
	</div>
</div>
