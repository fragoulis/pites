<script lang="ts">
	import {
		Helper,
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell
	} from 'flowbite-svelte';
	import { type CreateReceiptFormForBatchPayments } from '$lib/types';
	import InputField from '$lib/InputField.svelte';
	import AlertWarning from '$lib/AlertWarning.svelte';
	import Form from '$lib/Form.svelte';
	import { todayStr } from '$lib/utils';

	export let payments: Map<string, any>;

	let form: CreateReceiptFormForBatchPayments = {
		payment_ids: [...payments.keys()],
		receipt_nos: [],
		block_nos: [],
		issued_at: todayStr(),
		comments: ''
	};

	for (let i = 0; i < form.payment_ids.length; i++) {
		form.receipt_nos[i] = 0;
		form.block_nos[i] = 0;
	}

	type Errors = {
		payment_ids: string[];
		receipt_nos: string[];
		block_nos: string[];
		issued_at: string;
	};

	let errors: Errors = {
		payment_ids: [],
		receipt_nos: [],
		block_nos: []
	};
</script>

<Form bind:form url="/receipts/batch" bind:errors on:success on:failure>
	<AlertWarning>
		Η ενέργεια θα δημιουργήσει {payments.size} αποδείξεις.
	</AlertWarning>

	<Table>
		<TableHead>
			<TableHeadCell>#</TableHeadCell>
			<TableHeadCell>Μέλος</TableHeadCell>
			<TableHeadCell>Ποσό</TableHeadCell>
			<TableHeadCell>Αριθμός απόδειξης</TableHeadCell>
			<TableHeadCell>Μπλοκ</TableHeadCell>
		</TableHead>
		<TableBody>
			{#each payments.values() as payment, i}
				<TableBodyRow>
					<TableBodyCell
						tdClass="px-1 py-1 whitespace-nowrap font-medium text-gray-900 dark:text-white"
						>{i + 1}</TableBodyCell
					>
					<TableBodyCell
						tdClass="px-1 py-1 whitespace-nowrap font-medium text-gray-900 dark:text-white"
					>
						{payment.member_name}
						{#if errors.payment_ids[i]}
							<Helper class="mt-2 text-sm" color="red">
								{errors.payment_ids[i]}
							</Helper>
						{/if}
					</TableBodyCell>
					<TableBodyCell
						tdClass="px-1 py-1 whitespace-nowrap font-medium text-gray-900 dark:text-white"
						>{payment.amount}€</TableBodyCell
					>
					<TableBodyCell
						tdClass="px-1 py-1 whitespace-nowrap font-medium text-gray-900 dark:text-white"
					>
						<InputField
							type="number"
							size="sm"
							divClass=""
							bind:value={form.receipt_nos[i]}
							bind:error={errors.receipt_nos[i]}
						/>
					</TableBodyCell>
					<TableBodyCell
						tdClass="px-1 py-1 whitespace-nowrap font-medium text-gray-900 dark:text-white"
					>
						<InputField
							type="number"
							size="sm"
							divClass=""
							bind:value={form.block_nos[i]}
							bind:error={errors.block_nos[i]}
						/>
					</TableBodyCell>
				</TableBodyRow>
			{/each}
		</TableBody>
	</Table>

	<div class="w-full">
		<InputField
			type="date"
			label="Ημ/νία απόδειξης"
			bind:value={form.issued_at}
			bind:error={errors.issued_at}
		/>
	</div>

	<InputField type="textarea" label="Σχόλια" bind:value={form.comments} />
</Form>
