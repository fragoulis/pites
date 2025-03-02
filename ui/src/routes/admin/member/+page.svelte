<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { page } from '$app/stores';
	import MemberView from '$lib/MemberView.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import AlertWarning from '$lib/AlertWarning.svelte';
	import pb from '$lib/pocketbase';
	import Loading from '$lib/Loading.svelte';
	import { type Member } from '$lib/types';

	let record: Member;
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
			record = await pb.send(`/members/${id}`, {});
		} catch (e: any) {
			error = e;
		}
	};
</script>

{#if record}
	{#key record.name_formatted}
		<PageHeader>{record.name_formatted}</PageHeader>
		<MemberView bind:record />
	{/key}
{:else if error}
	<AlertWarning>{error.message}</AlertWarning>
{:else}
	<Loading />
{/if}
