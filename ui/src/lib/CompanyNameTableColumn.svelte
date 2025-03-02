<script lang="ts">
	import { A } from 'flowbite-svelte';
	import { type Company } from '$lib/types';
	export let record: Company;
</script>

<A
	href={`/admin/company/edit?id=${record.id}`}
	class="text-sky-700 dark:text-sky-200"
	target="_blank"
>
	{record.name}
</A>
