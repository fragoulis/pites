<script lang="ts">
	import { type Payment, type UpdatePaymentForm } from '$lib/types';
	import InputField from '$lib/InputField.svelte';
	import Form from '$lib/Form.svelte';

	export let record: Payment;

	let form: UpdatePaymentForm = {
		amount: record.amount,
		receipt_block_no: record.receipt_block_no,
		receipt_no: record.receipt_no,
		comments: record.comments
	};
	let errors: UpdatePaymentForm = {};
</script>

<div class="grid gap-4 grid-cols-1">
	<div class="col-span-1">
		<Form
			bind:form
			method="PATCH"
			url={`/payments?id=${record.id}`}
			bind:errors
			on:success
			on:failure
		>
			<div class="w-full mb-5">
				<InputField
					type="number"
					label="Ποσό €"
					min={0}
					max={1000}
					bind:value={form.amount}
					bind:error={errors.amount}
				/>
			</div>

			<div class="w-full">
				<InputField
					type="number"
					label="Απόδειξη"
					bind:value={form.receipt_no}
					bind:error={errors.receipt_no}
					min={0}
					max={50}
				/>
			</div>
			<div class="w-full">
				<InputField
					type="number"
					label="Μπλοκ αποδείξεων"
					bind:value={form.receipt_block_no}
					bind:error={errors.receipt_block_no}
					min={0}
					max={1000}
				/>
			</div>

			<InputField type="textarea" label="Σχόλια" bind:value={form.comments} />
		</Form>
	</div>
</div>
