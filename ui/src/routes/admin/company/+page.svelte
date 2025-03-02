<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { page } from '$app/stores';
	import { Button } from 'flowbite-svelte';
	import CompanyView from '$lib/CompanyView.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import AlertWarning from '$lib/AlertWarning.svelte';
	import pb from '$lib/pocketbase';
	import Loading from '$lib/Loading.svelte';
	import { type Company } from '$lib/types';

	let record: Company;
	let error: any;

	onMount(async () => {
		const id: string | null = $page.url.searchParams.get('id');
		if (id == null) {
			return;
		}

		await tick();
		await fetchRecordByID(id);
	});

	const fetchRecordByID = async (id: string) => {
		error = undefined;
		pb.cancelAllRequests();
		try {
			record = await pb.send(`/companies/${id}`, {});
		} catch (e: any) {
			error = e;
		}
	};
</script>

{#if record}
	<PageHeader>
		{record.name_formatted}
		<svelte:fragment slot="right">
			<Button href={`/admin/company/edit/?id=${record.id}`}>
				<span class="font-bold">Επεξεργασία</span>
			</Button>
		</svelte:fragment>
	</PageHeader>
	<CompanyView {record} />
{:else if error}
	<AlertWarning>{error.message}</AlertWarning>
{:else}
	<Loading />
{/if}
