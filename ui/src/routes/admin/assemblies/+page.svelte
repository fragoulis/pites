<script lang="ts">
	import PageHeader from '$lib/PageHeader.svelte';
	import { Button } from 'flowbite-svelte';
	import { PlusOutline } from 'flowbite-svelte-icons';
	import { isAdmin } from '$lib/utils';
	import Table from '$lib/Table.svelte';

	const headers: string[] = ['Ημερομηνία', 'Σχόλια', 'Ενεργή'];
	const fields: string[] = ['date_formatted', 'comments', 'active'];
</script>

<PageHeader>
	Συνελεύσεις
	<svelte:fragment slot="right">
		{#if isAdmin()}
			<Button href="/admin/assemblies/create">
				<PlusOutline class="w-4 h-4 me-2" />
				<span class="font-bold">Νέα συνέλευση</span>
			</Button>
		{/if}
	</svelte:fragment>
</PageHeader>
<Table {headers} {fields} collection="assemblies" sortBy="-date" />
