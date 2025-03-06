<script lang="ts">
	import { type UpdateMemberForm, type Member } from '$lib/types';
	import AutocompleteAddressFormField from '$lib/AutocompleteAddressFormField.svelte';
	import InputField from '$lib/InputField.svelte';
	import { TabItem, Tabs } from 'flowbite-svelte';

	export let record: Member = {};
	export let form: UpdateMemberForm;
	export let errors: UpdateMemberForm;

	let addressType: 'structured' | 'unstructured' | 'city_only' = 'structured';

	if (form.address_street_id != '' && form.address_city_id != '') {
		addressType = 'structured';
	} else if (form.address_street_id == '' && form.address_city_id != '') {
		addressType = 'city_only';
	} else if (form.legacy_address != '') {
		addressType = 'unstructured';
	}
</script>

<Tabs>
	<TabItem open={addressType == 'structured'} on:click={() => (addressType = 'structured')}>
		<span slot="title">Πλήρης (Αττική)</span>

		<div class="w-full">
			<AutocompleteAddressFormField
				query={record.address_formatted}
				label="Οδός και Δήμος"
				bind:value={form.address_street_id}
				bind:error={errors.address_street_id}
				required={false}
			/>
		</div>

		<div class="w-full">
			<InputField
				id="address_street_no"
				label="Αριθμός Οδού"
				bind:value={form.address_street_no}
				bind:error={errors.address_street_no}
			/>
		</div>
	</TabItem>
	<TabItem open={addressType == 'city_only'} on:click={() => (addressType = 'city_only')}>
		<span slot="title">Μόνο Δήμος (Αττική)</span>

		<div class="w-full">
			<AutocompleteAddressFormField
				query={record.address_city_name}
				label="Δήμος"
				bind:value={form.address_city_id}
				bind:error={errors.address_city_id}
				required={false}
				city_only={true}
			/>
		</div>
	</TabItem>
	<TabItem open={addressType == 'unstructured'} on:click={() => (addressType = 'unstructured')}>
		<span slot="title">Αδόμητη</span>
		<div class="w-full mb-5">
			Μόνο για διευθύνσεις εκτός Αττικής. Αποφύγετε να συμπληρώνετε χύμα τιμές.
		</div>
		<div class="w-full">
			<InputField
				id="legacy_address"
				label="Οδός/Αριθμός"
				bind:value={form.legacy_address}
				bind:error={errors.legacy_address}
			/>
		</div>
		<div class="w-full">
			<InputField
				id="legacy_area"
				label="Δήμος"
				bind:value={form.legacy_area}
				bind:error={errors.legacy_area}
			/>
		</div>
		<div class="w-full">
			<InputField
				id="legacy_city"
				label="Πόλη"
				bind:value={form.legacy_city}
				bind:error={errors.legacy_city}
			/>
		</div>
		<div class="w-full">
			<InputField
				id="legacy_post_code"
				label="ΤΚ"
				bind:value={form.legacy_post_code}
				bind:error={errors.legacy_post_code}
			/>
		</div>
	</TabItem>
</Tabs>
