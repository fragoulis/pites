<script lang="ts">
	import CompaniesTable from '$lib/CompaniesTable.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import { Button } from 'flowbite-svelte';
	import { PlusOutline } from 'flowbite-svelte-icons';
	import { isAdmin } from '$lib/utils';
</script>

<PageHeader>
	Αναζήτηση εταιρειών
	<svelte:fragment slot="right">
		{#if isAdmin()}
			<Button href="/admin/companies/create">
				<PlusOutline class="w-4 h-4 me-2" />
				<span class="font-bold">Νέα εταιρεία</span>
			</Button>
		{/if}
	</svelte:fragment>
</PageHeader>
<CompaniesTable />
