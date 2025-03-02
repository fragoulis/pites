<script lang="ts">
	import { Breadcrumb, BreadcrumbItem } from 'flowbite-svelte';
	import PaymentCreateForm from '$lib/PaymentCreateForm.svelte';
	import Separator from '$lib/Separator.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import { sharedToast } from '$lib/store';

	const onSuccess = () => {
		$sharedToast.show = true;
		$sharedToast.success = true;
		$sharedToast.message = 'Η πληρωμή περάστηκε επιτυχώς.';
	};

	const onError = async (e: any) => {
		$sharedToast.show = true;
		$sharedToast.success = false;
		$sharedToast.message = e.detail?.message || e.message;
	};
</script>

<PageHeader>Νέα Πληρωμή</PageHeader>
<Breadcrumb>
	<BreadcrumbItem home href="/admin/payments">Πληρωμές</BreadcrumbItem>
	<BreadcrumbItem>Νέα Πληρωμή</BreadcrumbItem>
</Breadcrumb>
<Separator />

<PaymentCreateForm on:success={onSuccess} on:failure={onError} />
