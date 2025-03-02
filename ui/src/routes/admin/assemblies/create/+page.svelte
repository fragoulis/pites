<script lang="ts">
	import { Breadcrumb, BreadcrumbItem } from 'flowbite-svelte';
	import AssemblyForm from '$lib/AssemblyForm.svelte';
	import Separator from '$lib/Separator.svelte';
	import PageHeader from '$lib/PageHeader.svelte';
</script>

<PageHeader>Νέα Συνέλευση</PageHeader>
<Breadcrumb>
	<BreadcrumbItem home href="/admin/assemblies">Συνελεύσεις</BreadcrumbItem>
	<BreadcrumbItem>Νέα Συνέλευση</BreadcrumbItem>
</Breadcrumb>
<Separator />
<AssemblyForm />
