<script lang="ts">
	import { type UpdateMemberForm } from '$lib/types';
	import AutocompleteCompanyFormField from '$lib/AutocompleteCompanyFormField.svelte';
	import InputField from '$lib/InputField.svelte';

	export let form: UpdateMemberForm;
	export let errors: UpdateMemberForm;
	export let company_name: string = '';

	const onSuccess = (e: any) => {
		form.company_id = e.detail.model.id;
	};
</script>

<AutocompleteCompanyFormField
	query={company_name}
	bind:value={form.company_id}
	bind:error={errors.company_id}
	on:company:created={onSuccess}
	help="Αφήστε κενό εάν είναι άνεργος/η"
/>

<div class="w-full">
	<InputField label="Ειδικότητα" bind:value={form.specialty} bind:error={errors.specialty} />
</div>
