<script lang="ts">
	import {
		Tabs,
		TabItem,
		Table,
		TableBody,
		TableBodyRow,
		TableBodyCell,
		TableHead,
		TableHeadCell
	} from 'flowbite-svelte';
	import { UserCircleSolid, ClipboardSolid, AdjustmentsVerticalSolid } from 'flowbite-svelte-icons';

	import { type Member } from '$lib/types';
	import { isComponent } from '$lib/utils';
	import pb from '$lib/pocketbase';
	import { sharedToast } from '$lib/store';

	import MemberViewChangeDetailsButton from '$lib/MemberViewChangeDetailsButton.svelte';
	import MemberDetailsForm from '$lib/MemberDetailsForm.svelte';
	import MemberCommentsAutosave from '$lib/MemberCommentsAutosave.svelte';
	import IssuesAlert from '$lib/IssuesAlert.svelte';
	import MemberPaymentHistory from '$lib/MemberPaymentHistory.svelte';
	import MemberSubscriptionStatus from '$lib/MemberSubscriptionStatus.svelte';
	import PaymentStatusTableColumn from '$lib/PaymentStatusTableColumn.svelte';
	import MemberCompanyFormattedNameColumn from '$lib/MemberCompanyFormattedNameColumn.svelte';

	export let record: Member;

	$: issues = record.issues;

	const profile: any = {
		'Αρ. Μητρώου': (r: Member) => r.member_no,
		Όνομα: (r: Member) => r.first_name,
		Επώνυμο: (r: Member) => r.last_name,
		Πατρώνυμο: (r: Member) => r.father_name,
		Διεύθυνση: (r: Member) => r.address_formatted,
		Εταιρεία: MemberCompanyFormattedNameColumn,
		Ειδικότητα: (r: Member) => r.specialty,
		Σπουδές: (r: Member) => r.education,
		Συνδρομή: MemberSubscriptionStatus,
		Κατάσταση: PaymentStatusTableColumn,
		Email: (r: Member) => r.email,
		Κινητό: (r: Member) => r.mobile,
		Σταθερό: (r: Member) => r.phone,
		'Ημ/νία Γέννησης': (r: Member) => r.birthdate_formatted,
		'Αριθμός Δελτίου Ταυτότητας': (r: Member) => r.id_card_number,
		'Αριθμός Μητρώου Ασφαλισμένου': (r: Member) => r.social_security_num,
		'Μέλος άλλου σωματείου': (r: Member) => (r.other_union ? 'Ναι' : 'Όχι'),
		Σχόλια: MemberCommentsAutosave
	};

	let hideDrawerDetails: boolean = true;

	const onSuccess = async () => {
		$sharedToast.show = true;
		$sharedToast.success = true;
		$sharedToast.message = 'Η αλλαγή αποθηκεύτηκε.';

		try {
			record = await pb.send(`/members/${record.id}`, {});
			hideDrawerDetails = true;
		} catch (e) {
			onError(e);
		}
	};

	const onError = async (e: any) => {
		$sharedToast.show = true;
		$sharedToast.success = false;
		$sharedToast.message = e.detail?.message || e.message;
	};
</script>

<MemberViewChangeDetailsButton bind:hideDrawer={hideDrawerDetails}>
	<svelte:fragment slot="label">Αλλαγή Στοιχείων</svelte:fragment>
	<svelte:fragment slot="title">Αλλαγή στοιχείων για {record.name_formatted}</svelte:fragment>

	<MemberDetailsForm {record} on:success={onSuccess} on:failure={onError} />
</MemberViewChangeDetailsButton>

<IssuesAlert {issues} />

<Tabs tabStyle="underline">
	<TabItem open>
		<div slot="title" class="flex items-center gap-2">
			<UserCircleSolid size="md" />
			Στοιχεία
		</div>

		<Table striped={true} hoverable={true}>
			<TableBody>
				{#each Object.entries(profile) as [header, value]}
					<TableBodyRow>
						<TableBodyCell>{header}</TableBodyCell>
						<TableBodyCell>
							{#if isComponent(value)}
								<svelte:component
									this={value}
									{record}
									on:success={onSuccess}
									on:failure={onError}
								/>
							{:else}
								{value(record)}
							{/if}
						</TableBodyCell>
					</TableBodyRow>
				{/each}
			</TableBody>
		</Table>
	</TabItem>
	<TabItem>
		<div slot="title" class="flex items-center gap-2">
			<AdjustmentsVerticalSolid size="md" />
			Ιστορικό πληρωμών
		</div>

		<MemberPaymentHistory {record} />
	</TabItem>
	<TabItem>
		<div slot="title" class="flex items-center gap-2">
			<ClipboardSolid size="md" />
			Ιστορικό συνδρομών
		</div>

		<Table striped={true} hoverable={true}>
			<TableHead>
				<TableHeadCell>Ενεργή</TableHeadCell>
				<TableHeadCell>Μήνες</TableHeadCell>
				<TableHeadCell>Από</TableHeadCell>
				<TableHeadCell>Μέχρι</TableHeadCell>
				<TableHeadCell>Πληρωμένη Εγγραφή</TableHeadCell>
			</TableHead>
			<TableBody>
				{#each record.subscriptions as subscription}
					<TableBodyRow>
						<TableBodyCell>{subscription.active ? 'Ναι' : 'Όχι'}</TableBodyCell>
						<TableBodyCell>{subscription.months}</TableBodyCell>
						<TableBodyCell>{subscription.start_date_formatted}</TableBodyCell>
						<TableBodyCell>{subscription.end_date_formatted || '-'}</TableBodyCell>
						<TableBodyCell>{subscription.fee_paid ? 'Ναι' : 'Όχι'}</TableBodyCell>
					</TableBodyRow>
				{:else}
					<TableBodyRow>
						<TableBodyCell colspan="5">Δε βρέθηκαν συνδρομές</TableBodyCell>
					</TableBodyRow>
				{/each}
			</TableBody>
		</Table>
	</TabItem>
</Tabs>
