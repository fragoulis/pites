<script lang="ts">
	import { type Member, type MemberCompanyForm } from '$lib/types';
	import Form from '$lib/Form.svelte';
	import MemberCompanyFormFields from './MemberCompanyFormFields.svelte';

	export let record: Member;

	let form: MemberCompanyForm = {
		company_id: record.company_id
	};
	let errors: MemberCompanyForm = {};

	const url: string = `/members/${record.id}/company`;
</script>

<Form {form} method="POST" {url} bind:errors on:success on:failure>
	<MemberCompanyFormFields company_name={record.company_formatted} bind:form bind:errors />
</Form>
