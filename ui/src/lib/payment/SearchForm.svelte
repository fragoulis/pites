<script lang="ts">
	import { type PaymentSearchForm } from '$lib/types';
	import InputField from '$lib/InputField.svelte';
	import { Select, Label, Button } from 'flowbite-svelte';
	import { todayStr, weekAgoStr, monthAgoStr } from '$lib/utils';

	export let form: PaymentSearchForm;

	const receiptStateOptions = [
		{ value: '', name: '' },
		{ value: 'with', name: 'Με απόδειξη' },
		{ value: 'without', name: 'Χωρίς απόδειξη' }
	];

	export const reset = () => {
		form.receipt_state = '';
		form.issue_date_from = '';
		form.issue_date_to = '';
	};
</script>

<div class="w-full grid grid-cols-4 gap-4 p-4 bg-gray-100 rounded-lg border">
	<div>
		<Label class="space-y-2" color="gray">
			<div class="mb-2">Απόδειξη</div>

			<Select items={receiptStateOptions} bind:value={form.receipt_state} />
		</Label>
	</div>

	<div>
		<InputField type="date" label="Ημ/νία έκδοσης (από)" bind:value={form.issue_date_from} />
		<Button
			pill={true}
			size="xs"
			on:click={() => {
				form.issue_date_from = monthAgoStr();
			}}>Μήνα πριν</Button
		>
		<Button
			pill={true}
			size="xs"
			on:click={() => {
				form.issue_date_from = weekAgoStr();
			}}>Βδομάδα πριν</Button
		>
		<Button
			pill={true}
			size="xs"
			on:click={() => {
				form.issue_date_from = todayStr();
			}}>Σήμερα</Button
		>
	</div>
	<div>
		<InputField type="date" label="Ημ/νία έκδοσης (μέχρι)" bind:value={form.issue_date_to} />
		<Button
			pill={true}
			size="xs"
			on:click={() => {
				form.issue_date_to = monthAgoStr();
			}}>Μήνα πριν</Button
		>
		<Button
			pill={true}
			size="xs"
			on:click={() => {
				form.issue_date_to = weekAgoStr();
			}}>Βδομάδα πριν</Button
		>
		<Button
			pill={true}
			size="xs"
			on:click={() => {
				form.issue_date_to = todayStr();
			}}>Σήμερα</Button
		>
	</div>
</div>
