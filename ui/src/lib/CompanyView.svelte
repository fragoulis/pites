<script lang="ts">
	import { Tabs, TabItem, Table, TableBody, TableBodyRow, TableBodyCell } from 'flowbite-svelte';
	import { UserCircleSolid } from 'flowbite-svelte-icons';

	import { type Company } from '$lib/types';

	export let record: Company;

	const profile: any = {
		Όνομα: (r: Company) => r.name,
		Διεύθυνση: (r: Company) => r.address_formatted,
		Παράρτημα: (r: Company) => r.branch
	};
</script>

<Tabs tabStyle="underline">
	<TabItem open>
		<div slot="title" class="flex items-center gap-2">
			<UserCircleSolid size="md" />
			Στοιχεία
		</div>

		<Table striped={true} hoverable={true}>
			<TableBody>
				{#each Object.entries(profile) as [header, value]}
					<TableBodyRow>
						<TableBodyCell>{header}</TableBodyCell>
						<TableBodyCell>{value(record)}</TableBodyCell>
					</TableBodyRow>
				{/each}
			</TableBody>
		</Table>
	</TabItem>
</Tabs>
