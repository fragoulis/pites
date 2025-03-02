<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { page } from '$app/stores';
	import PaymentUpdateForm from '$lib/PaymentUpdateForm.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import AlertWarning from '$lib/AlertWarning.svelte';
	import pb from '$lib/pocketbase';
	import Loading from '$lib/Loading.svelte';
	import { type Payment } from '$lib/types';
	import { sharedToast } from '$lib/store';

	let loading = false;
	let error: any;

	let record: Payment;
	let id: string | null | undefined;
	let name: string;

	onMount(async () => {
		id = $page.url.searchParams.get('id');
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
			record = await pb.send<Payment>(`/payments`, { query: { id: id } });
		} catch (e: any) {
			error = e;
		}
		loading = false;
	};

	const onSuccess = async (e: any) => {
		$sharedToast.show = true;
		$sharedToast.success = true;
		$sharedToast.message = 'Η πληρωμή ανανεώθηκε.';
	};

	const onError = (e: any) => {
		$sharedToast.show = true;
		$sharedToast.success = false;
		$sharedToast.message = e.detail?.message || e.message;
	};
</script>

{#if record}
	<PageHeader>Επεξεργασία πληρωμής για {record.member_name}</PageHeader>
	<PaymentUpdateForm bind:record on:success={onSuccess} on:failure={onError} />
{:else if error}
	<AlertWarning>{error.message}</AlertWarning>
{:else}
	<Loading />
{/if}
