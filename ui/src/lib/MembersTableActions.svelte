<script lang="ts">
	import { A } from 'flowbite-svelte';
	import { Dropdown, Button } from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	import CreatePaymentTableAction from '$lib/CreatePaymentTableAction.svelte';

	export let createPaymentActionHidden = true;
	export let selectedRows: Set<string>;
</script>

<Button outline pill color="light">Ενέργειες<ChevronDownOutline /></Button>
<Dropdown class="w-56 p-3 space-y-1">
	<li>
		<A
			on:click={(e) => {
				e.preventDefault();
				createPaymentActionHidden = false;
			}}>Δημιουργία είσπραξης</A
		>
	</li>
</Dropdown>

<CreatePaymentTableAction bind:hidden={createPaymentActionHidden} bind:members={selectedRows} />
