<script lang="ts">
	import { Breadcrumb, BreadcrumbItem } from 'flowbite-svelte';
	import CompanyForm from '$lib/CompanyForm.svelte';
	import Separator from '$lib/Separator.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
	import { sharedToast } from '$lib/store';

	const onSuccess = () => {
		$sharedToast.show = true;
		$sharedToast.success = true;
		$sharedToast.message = 'Η εταιρεία δημιουργήθηκε επιτυχώς.';
	};

	const onError = async (e: any) => {
		$sharedToast.show = true;
		$sharedToast.success = false;
		$sharedToast.message = e.detail?.message || e.message;
	};
</script>

<PageHeader>Νέα Εταιρεία</PageHeader>
<Breadcrumb>
	<BreadcrumbItem home href="/admin/companies">Εταιρείες</BreadcrumbItem>
	<BreadcrumbItem>Νέα Εταιρεία</BreadcrumbItem>
</Breadcrumb>
<Separator />

<CompanyForm on:success={onSuccess} on:failure={onError} />
