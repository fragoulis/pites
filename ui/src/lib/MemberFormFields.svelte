<script lang="ts">
	import InputField from '$lib/InputField.svelte';
	import ToggleField from '$lib/ToggleField.svelte';
	import MemberAddressTabbedFormFields from '$lib/MemberAddressTabbedFormFields.svelte';
	import InputGroup from '$lib/InputGroup.svelte';
	import MemberCompanyFormFields from '$lib/MemberCompanyFormFields.svelte';
	import { type UpdateMemberForm, type Member } from '$lib/types';

	export let record: Member = {};
	export let form: UpdateMemberForm = {};
	export let errors: UpdateMemberForm = {};
</script>

<InputGroup legend="Προσωπικά">
	<div class="w-full">
		<InputField
			id="last_name"
			label="Επώνυμο"
			bind:value={form.last_name}
			bind:error={errors.last_name}
			required={true}
			on:keyup
		/>
	</div>
	<div class="w-full">
		<InputField
			id="first_name"
			label="Όνομα"
			bind:value={form.first_name}
			bind:error={errors.first_name}
			required={true}
			on:keyup
		/>
	</div>
	<div class="w-full">
		<InputField
			id="father_name"
			label="Όνομα πατρός"
			bind:value={form.father_name}
			bind:error={errors.father_name}
		/>
	</div>
	<div class="w-full">
		<InputField
			id="birthdate"
			type="date"
			label="Ημ/νία Γέννησης"
			bind:value={form.birthdate}
			bind:error={errors.birthdate}
		/>
	</div>
</InputGroup>

<InputGroup legend="Διεύθυνση Κατοικίας">
	<MemberAddressTabbedFormFields bind:form bind:errors {record} />
</InputGroup>

<InputGroup legend="Επικοινωνία">
	<div class="w-full">
		<InputField
			id="phone"
			label="Τηλέφωνο Οικίας"
			bind:value={form.phone}
			bind:error={errors.phone}
			on:keyup
		/>
	</div>
	<div class="w-full">
		<InputField
			id="mobile"
			label="Κινητό"
			bind:value={form.mobile}
			bind:error={errors.mobile}
			on:keyup
		/>
	</div>
	<div class="w-full">
		<InputField
			id="email"
			label="Email"
			type="email"
			bind:value={form.email}
			bind:error={errors.email}
			on:keyup
		/>
	</div>
</InputGroup>

<InputGroup legend="Ταυτότητα / Ασφάλιση">
	<div class="w-full">
		<InputField
			id="id_card_number"
			label="Αριθμός Αστυνομικής Ταυτότητας"
			bind:value={form.id_card_number}
			bind:error={errors.id_card_number}
		/>
	</div>

	<div class="w-full">
		<InputField
			id="social_security_num"
			label="Αριθμός Βιβλιαρίου Ασφάλισης"
			bind:value={form.social_security_num}
			bind:error={errors.social_security_num}
		/>
	</div>

	<div class="w-full">
		<InputField label="Επίπεδο σπουδών" bind:value={form.education} bind:error={errors.education} />
	</div>
</InputGroup>

<InputGroup legend="Εργασία">
	<MemberCompanyFormFields bind:form bind:errors company_name={record.company_name} />
</InputGroup>

<div class="w-full">
	<ToggleField bind:checked={form.other_union} label="Μέλος άλλους σωματείου" />
</div>

<InputGroup legend="Γενικά σχόλια">
	<div class="w-full">
		<InputField id="comments" type="textarea" label="Σχόλια" bind:value={form.comments} />
	</div>
</InputGroup>
