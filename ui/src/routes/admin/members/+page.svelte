<script lang="ts">
	import MembersTable from '$lib/MembersTable.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import { Button } from 'flowbite-svelte';
	import { PlusOutline, TabletOutline } from 'flowbite-svelte-icons';
	import { isAdmin } from '$lib/utils';
</script>

<PageHeader>
	Αναζήτηση μελών
	<svelte:fragment slot="right">
		{#if isAdmin()}
			<Button href="/admin/members/create">
				<PlusOutline class="w-4 h-4 me-2" />
				<span class="font-bold">Νέο μέλος</span>
			</Button>
		{/if}
	</svelte:fragment>
</PageHeader>
<MembersTable />
