<script lang="ts">
	import {
		Table,
		TableHead,
		TableBody,
		TableHeadCell,
		TableBodyCell,
		TableBodyRow
	} from 'flowbite-svelte';
	import { onMount } from 'svelte';
	import pb from '$lib/pocketbase';
	import { loggedInUserID } from '$lib/utils';
	import { type Payment } from '$lib/types';
	import { EditSolid } from 'flowbite-svelte-icons';

	let records: Payment[] = [];

	export const refresh = async () => {
		try {
			records = await pb.send('/payments', {
				query: { user_ids: [loggedInUserID()] }
			});
		} catch (e: any) {
			console.error(e);
		}
	};

	onMount(refresh);
</script>

<Table class="text-sm">
	<TableHead>
		<TableHeadCell>Απόδ.</TableHeadCell>
		<TableHeadCell>Μέλος</TableHeadCell>
		<TableHeadCell>Ποσό</TableHeadCell>
		<TableHeadCell>Επεξ.</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each records as record}
			<TableBodyRow>
				<TableBodyCell tdClass="px-2 py-2">#{record.receipt_no}</TableBodyCell>
				<TableBodyCell tdClass="px-2 py-2">
					<a
						class="flex grow items-center text-sky-500 hover:underline after:content-['_↗']"
						target="_blank"
						href={`/admin/member?id=${record.member_id}`}
					>
						{record.member_name}
						<br />
						{record.member_no}
					</a>
				</TableBodyCell>
				<TableBodyCell tdClass="px-2 py-2">{record.amount}€</TableBodyCell>
				<TableBodyCell tdClass="px-2 py-2">
					<a
						class="flex grow items-center text-green-500 hover:underline"
						target="_blank"
						href={`/admin/payment/edit?id=${record.id}`}
					>
						<EditSolid />
					</a>
				</TableBodyCell>
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
