<script lang="ts">
	import { A } from 'flowbite-svelte';
	import { Dropdown, Button } from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	import CreateReceiptTableAction from '$lib/payment/CreateReceiptTableAction.svelte';

	export let createReceiptActionHidden = true;
	export let selectedRows: Map<string, any>;
</script>

<Button outline pill color="light">Ενέργειες<ChevronDownOutline /></Button>
<Dropdown class="w-56 p-3 space-y-1">
	<li>
		<A
			on:click={(e) => {
				e.preventDefault();
				createReceiptActionHidden = false;
			}}>Δημιουργία απόδειξης</A
		>
	</li>
</Dropdown>

<CreateReceiptTableAction bind:hidden={createReceiptActionHidden} bind:payments={selectedRows} />
