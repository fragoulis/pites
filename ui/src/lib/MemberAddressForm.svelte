<script lang="ts">
	import { type Member, type MemberAddressForm } from '$lib/types';
	import MemberAddressTabbedFormFields from '$lib/MemberAddressTabbedFormFields.svelte';
	import Form from '$lib/Form.svelte';

	export let record: Member;

	let form: MemberAddressForm = {
		address_street_id: record.address_street_id,
		address_city_id: record.address_city_id,
		address_street_no: record.address_street_no,
		legacy_address: record.legacy_address,
		legacy_area: record.legacy_area,
		legacy_city: record.legacy_city,
		legacy_post_code: record.legacy_post_code
	};
	let errors: MemberAddressForm = {};

	const url: string = `/members/${record.id}/address`;
</script>

<Form {form} method="POST" {url} bind:errors on:success on:failure>
	<MemberAddressTabbedFormFields bind:form bind:errors {record} />
</Form>
