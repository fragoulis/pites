<script lang="ts">
	import { Breadcrumb, BreadcrumbItem } from 'flowbite-svelte';
	import NewMemberForm from '$lib/NewMemberForm.svelte';
	import Separator from '$lib/Separator.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
</script>

<PageHeader>Νέο Μέλος</PageHeader>
<Breadcrumb>
	<BreadcrumbItem home href="/admin/members">Μέλη</BreadcrumbItem>
	<BreadcrumbItem>Νέο Μέλος</BreadcrumbItem>
</Breadcrumb>
<Separator />
<NewMemberForm />
