<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { page } from '$app/stores';
	import CompanyForm from '$lib/CompanyForm.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import AlertWarning from '$lib/AlertWarning.svelte';
	import PrimaryLink from '$lib/PrimaryLink.svelte';
	import pb from '$lib/pocketbase';
	import Loading from '$lib/Loading.svelte';
	import type { Company } from '$lib/types';

	let loading = false;
	let record: Company;
	let error: Company;

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
		loading = true;
		pb.cancelAllRequests();
		try {
			record = await pb.send(`/companies/${id}`, {});
		} catch (e: any) {
			error = e;
		}
		loading = false;
	};
</script>

{#if record}
	<PageHeader>
		Επεξεργασία της "{record.name}"
	</PageHeader>
	<PrimaryLink href="/admin/company?id={record.id}">Πίσω</PrimaryLink>
	<CompanyForm {record} />
{:else if error}
	<AlertWarning>{error.message}</AlertWarning>
{:else}
	<Loading />
{/if}
