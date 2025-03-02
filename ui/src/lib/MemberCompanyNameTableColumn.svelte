<script lang="ts">
	import { A } from 'flowbite-svelte';
	import { type Member } from '$lib/types';
	export let record: Member;
</script>

{#if record.company_id}
	<A href={`/admin/company/edit?id=${record.company_id}`} target="_blank">
		{record.company_name}
	</A>
{:else}
	{record.company_name}
{/if}
