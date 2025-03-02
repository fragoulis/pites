<script lang="ts">
	import { A } from 'flowbite-svelte';
	import { type Member } from '$lib/types';
	import MemberViewDrawer from './MemberViewDrawer.svelte';

	export let record: Member;
	export let hidden = true;
</script>

<A
	href={`/admin/member?id=${record.id}`}
	on:click={(e) => {
		e.preventDefault();
		hidden = false;
	}}
>
	{record.name_formatted}
</A>

<MemberViewDrawer {record} bind:hidden />
