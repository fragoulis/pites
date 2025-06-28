<script lang="ts">
	import { A } from 'flowbite-svelte';
	import { type Payment } from '$lib/types';

	export let record: Payment;
</script>

<A href={`/admin/payment/edit?id=${record.id}`}>
	{record.issued_at_formatted}
</A>
