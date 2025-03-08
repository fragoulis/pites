<script lang="ts">
	import { Breadcrumb, BreadcrumbItem } from 'flowbite-svelte';
	import ImportMembersForm from '$lib/ImportMembersForm.svelte';
	import Separator from '$lib/Separator.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
</script>

<PageHeader>Εισαγωγή</PageHeader>
<Breadcrumb>
	<BreadcrumbItem home href="/admin/members">Μέλη</BreadcrumbItem>
	<BreadcrumbItem>Εισαγωγή</BreadcrumbItem>
</Breadcrumb>
<Separator />
<ImportMembersForm />
