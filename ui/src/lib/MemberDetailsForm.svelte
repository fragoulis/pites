<script lang="ts">
	import { type Member, type UpdateMemberForm } from '$lib/types';
	import Form from '$lib/Form.svelte';
	import MemberFormFields from './MemberFormFields.svelte';

	export let record: Member;

	let form: UpdateMemberForm = {
		first_name: record.first_name,
		last_name: record.last_name,
		father_name: record.father_name,
		email: record.email,
		mobile: record.mobile,
		phone: record.phone,
		birthdate: record.birthdate,
		id_card_number: record.id_card_number,
		social_security_num: record.social_security_num,
		other_union: record.other_union,
		comments: record.comments,
		education: record.education,
		address_street_id: record.address_street_id,
		address_city_id: record.address_city_id,
		address_street_no: record.address_street_no,
		legacy_address: record.legacy_address,
		legacy_area: record.legacy_area,
		legacy_city: record.legacy_city,
		legacy_post_code: record.legacy_post_code,
		company_id: record.company_id,
		specialty: record.specialty
	};

	let errors: UpdateMemberForm = {};
	const url: string = `/members/${record.id}`;
</script>

<Form {form} method="PATCH" {url} bind:errors on:success on:failure>
	<MemberFormFields bind:form bind:errors {record} />
</Form>
