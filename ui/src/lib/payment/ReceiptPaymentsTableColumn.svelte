<script lang="ts">
	import { A } from 'flowbite-svelte';
	import { type Payment } from '$lib/types';

	export let record: Payment;
</script>

{#if record.receipt_id == ''}
	Χωρίς απόδειξη
{:else}
	#{record.receipt_no} ({record.receipt_block_no})
{/if}
