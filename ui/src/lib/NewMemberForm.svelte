<script lang="ts">
	import { goto } from '$app/navigation';
	import pb from '$lib/pocketbase';
	import { onMount } from 'svelte';
	import { type NewMemberForm, type Member } from '$lib/types';
	import InputField from '$lib/InputField.svelte';
	import ToggleField from '$lib/ToggleField.svelte';
	import MemberAddressTabbedFormFields from '$lib/MemberAddressTabbedFormFields.svelte';
	import Form from '$lib/Form.svelte';
	import InputGroup from '$lib/InputGroup.svelte';
	import { sharedToast } from '$lib/store';
	import { Table, TableBodyCell, TableBodyRow, Card, A } from 'flowbite-svelte';
	import MemberCompanyFormFields from '$lib/MemberCompanyFormFields.svelte';
	import MemberFormFields from '$lib/MemberFormFields.svelte';

	let form: NewMemberForm = {};
	let errors: NewMemberForm = {};

	const onSuccess = async () => {
		$sharedToast.show = true;
		$sharedToast.success = true;
		$sharedToast.message = 'Το νέο μέλος δημιουργήθηκε.';
		await goto(`/admin/members?q=${form.member_no}`);
	};

	const onError = async (e: any) => {
		$sharedToast.show = true;
		$sharedToast.success = false;
		$sharedToast.message = e.detail?.message || e.message;
	};

	let loading: boolean = false;
	let records: Member[] = [];

	const findSimilar = async () => {
		loading = true;

		pb.cancelAllRequests();
		let query = [
			form?.last_name?.trim(),
			form?.email?.trim(),
			form?.mobile?.trim(),
			form?.phone?.trim()
		]
			.filter(Boolean)
			.join(' ');

		if (query == '') {
			records = [];
			return;
		}

		try {
			const res = await pb.send(`/members`, {
				query: {
					name: form.last_name,
					email: form.email,
					mobile: form.mobile,
					phone: form.phone,
					active_only: false,
					limit: 5
				}
			});
			records = res.records;
		} catch (err: any) {
			records = [];
		}
		loading = false;

		return records;
	};

	onMount(async () => {
		form.member_no = await pb.send('/members/next', {});
	});
</script>

<Form {form} url="/members" bind:errors on:success={onSuccess} on:failure={onError}>
	<div class="grid gap-4 grid-cols-1 mb-10">
		<div class="w-full">
			<InputField
				type="number"
				id="member_no"
				label="Αριθμός Μέλους"
				bind:value={form.member_no}
				bind:error={errors.member_no}
				required={true}
				min={1}
				max={10000}
			/>
			<div>Ο αριθμός μέλους συμπληρώνεται αυτόματα βάσει των υπαρχόντων μελών.</div>
		</div>
	</div>
	<div class="grid gap-4 grid-cols-2">
		<div class="col-span-1">
			<MemberFormFields bind:form bind:errors on:keyup={findSimilar} />

			<InputGroup legend="Συνδρομή">
				<div class="w-full">
					<InputField
						id="start_date"
						type="date"
						label="Ημ/νία έναρξης συνδρομής"
						bind:value={form.start_date}
						bind:error={errors.start_date}
					/>
				</div>

				<div class="w-full">
					<ToggleField bind:checked={form.fee_paid} label="Έχει εξοφλήσει εγγραφή" />
				</div>
			</InputGroup>
		</div>

		<Card class="col-span-1" size="xl">
			<h1 class="underline text-center">Παρόμοια αποτελέσματα</h1>
			<Table>
				{#each records as record}
					<TableBodyRow>
						<TableBodyCell class="break-works whitespace-normal">{record.member_no}</TableBodyCell>
						<TableBodyCell class="break-works whitespace-normal">
							<A href={`/admin/member?id=${record.id}`} target="_blank">
								{record.name_formatted}
							</A>
						</TableBodyCell>
						<TableBodyCell class="break-works whitespace-normal">{record.email}</TableBodyCell>
						<TableBodyCell class="break-works whitespace-normal">{record.mobile}</TableBodyCell>
					</TableBodyRow>
				{/each}
			</Table>
		</Card>
	</div>
</Form>
