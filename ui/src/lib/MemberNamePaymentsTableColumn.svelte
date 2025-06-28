<script lang="ts">
	import { A } from 'flowbite-svelte';
	import { type Payment } from '$lib/types';

	export let record: Payment;
</script>

<A href={`/admin/member?id=${record.member_id}`}>
	{record.member_name} | {record.member_no}
</A>
