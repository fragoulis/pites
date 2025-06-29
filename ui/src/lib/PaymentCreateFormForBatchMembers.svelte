<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { type CreatePaymentFormForBatchMembers } from '$lib/types';
	import InputField from '$lib/InputField.svelte';
	import Form from '$lib/Form.svelte';
	import { activePayment } from '$lib/store';
	import { todayStr } from '$lib/utils';
	import InputGroup from '$lib/InputGroup.svelte';
	import { Input } from 'flowbite-svelte';

	const dispatch = createEventDispatcher();

	export let members: Set<string>;

	let form: CreatePaymentFormForBatchMembers = {
		member_ids: Array.from(members),
		amount: 2,
		issued_at: $activePayment.issued_at == '' ? todayStr() : $activePayment.issued_at,
		comments: ''
	};

	let errors: CreatePaymentFormForBatchMembers = {};
</script>

<Form bind:form url="/payments/batch" bind:errors on:success on:failure>
	<div class="w-full mb-5">
		<InputGroup legend="Είσπραξη">
			<InputField
				type="number"
				min={0}
				max={1000}
				bind:value={form.amount}
				bind:error={errors.amount}
			/>
		</InputGroup>
	</div>

	<div class="w-full">
		<InputField
			id="issued_at"
			type="date"
			label="Ημ/νία Καταχώρισης"
			bind:value={form.issued_at}
			bind:error={errors.issued_at}
		/>
	</div>

	<InputField type="textarea" label="Σχόλια" bind:value={form.comments} />
</Form>
