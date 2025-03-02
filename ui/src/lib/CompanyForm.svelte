<script lang="ts">
	import { type CompanyForm, type Company } from '$lib/types';
	import InputField from '$lib/InputField.svelte';
	import AutocompleteAddressFormField from '$lib/AutocompleteAddressFormField.svelte';
	import Form from '$lib/Form.svelte';
	import SelectAsync from '$lib/SelectAsync.svelte';
	import { Label } from 'flowbite-svelte';

	export let record: Company | undefined = undefined;

	let form: CompanyForm = {};
	let errors: CompanyForm = {};
	let url: string = '/companies';
	let method: 'POST' | 'PUT' = 'POST';

	let addressFormatted: string | undefined;

	if (record) {
		form.name = record.name;
		form.email = record.email;
		form.phone = record.phone;
		form.website = record.website;
		form.comments = record.comments;
		form.business_type_id = record.business_type_id;
		form.address_street_id = record.address_street_id;
		form.address_street_no = record.address_street_no;
		addressFormatted = record.address_formatted;

		method = 'PUT';
		url = `/companies/${record.id}`;
	}
</script>

<Form {form} {method} {url} bind:errors on:success on:failure>
	<InputField label="Όνομα" bind:value={form.name} bind:error={errors.name} required={true} />

	<Label class="space-y-2 mb-6" color="gray">
		<div class="mb-2">Ομάδα</div>
		<SelectAsync url="/companies/business_types" bind:value={form.business_type_id} />
	</Label>

	<InputField type="email" label="Email" bind:value={form.email} bind:error={errors.email} />
	<InputField label="Τηλεφωνο" bind:value={form.phone} bind:error={errors.phone} />
	<InputField label="Website" bind:value={form.website} bind:error={errors.website} />
	<AutocompleteAddressFormField
		bind:query={addressFormatted}
		bind:value={form.address_street_id}
		bind:error={errors.address_street_id}
		required={false}
	/>
	<InputField
		label="Αριθμός"
		bind:value={form.address_street_no}
		bind:error={errors.address_street_no}
	/>
	<InputField type="textarea" label="Σχόλια" bind:value={form.comments} />
</Form>
