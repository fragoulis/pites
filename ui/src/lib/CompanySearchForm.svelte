<script lang="ts">
	import { Label } from 'flowbite-svelte';
	import { type DatatableSearchForm, type CompanyBusinessType } from '$lib/types';
	import MultiSelectAsync from '$lib/MultiSelectAsync.svelte';

	export let form: DatatableSearchForm;

	export const reset = () => {
		form.business_type_ids = undefined;
	};
</script>

<div class="w-full">
	<Label class="space-y-2" color="gray">
		<div class="mb-2">Ομάδα</div>
		<MultiSelectAsync url="/companies/business_types" bind:values={form.business_type_ids} />
	</Label>
</div>
