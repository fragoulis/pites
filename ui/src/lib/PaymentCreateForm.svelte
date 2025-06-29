<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { Alert, Card } from 'flowbite-svelte';
	import { InfoCircleSolid } from 'flowbite-svelte-icons';
	import { type CreatePaymentForm, type PaymentDetails, type Issue, type Member } from '$lib/types';
	import InputField from '$lib/InputField.svelte';
	import ToggleField from '$lib/ToggleField.svelte';
	import AutocompleteMemberFormField from '$lib/AutocompleteMemberFormField.svelte';
	import Form from '$lib/Form.svelte';
	import { activePayment } from '$lib/store';
	import { todayStr } from '$lib/utils';
	import InputGroup from '$lib/InputGroup.svelte';
	import IssuesAlert from '$lib/IssuesAlert.svelte';
	import { Tabs, TabItem, Input } from 'flowbite-svelte';
	import { onMount } from 'svelte';
	import MemberPaymentHistory from '$lib/MemberPaymentHistory.svelte';
	import LatestCashierPayments from '$lib/LatestCashierPayments.svelte';

	const dispatch = createEventDispatcher();

	let form: CreatePaymentForm = {
		amount: 2,
		months: 0,
		receipt_no: $activePayment.receipt_no,
		receipt_block_no: $activePayment.receipt_block_no,
		issued_at: $activePayment.issued_at == '' ? todayStr() : $activePayment.issued_at,
		without_receipt: false
	};

	let errors: CreatePaymentForm = {};
	let paymentDetails: PaymentDetails;
	let showPaymentDetails: boolean = false;
	let selectedMember: Member | undefined;
	let paymentMethod: 'amount' | 'months' = 'amount';
	let canMemberPay: boolean = true;
	let latestCashierPayments: any;

	const onSuccess = async (e: any) => {
		// Reset form.
		form.amount = 2;
		form.months = 0;
		form.comments = '';
		form.without_receipt = false;
		paymentMethod = 'amount';

		if (form.receipt_no && form.receipt_no >= 1 && form.receipt_no < 50) {
			form.receipt_no += 1;
		} else {
			form.receipt_no = 0;
		}

		if (form.receipt_block_no && form.receipt_block_no >= 1) {
		} else {
			form.receipt_block_no = 0;
		}

		paymentDetails = e.detail.model;
		showPaymentDetails = true;

		$activePayment = {
			receipt_block_no: form.receipt_block_no > 0 ? form.receipt_block_no : 0,
			receipt_no: form.receipt_no > 0 ? form.receipt_no + 1 : 0,
			issued_at: form.issued_at === undefined ? '' : form.issued_at
		};

		selectedMember = undefined;
		latestCashierPayments.refresh();
		dispatch('success');
	};

	let issues: Issue[] = [];

	const onSelectMember = (e: any) => {
		issues = e.detail.item.issues;
		selectedMember = e.detail.item;
		if (selectedMember) {
			form.member_id = selectedMember.id;
			canMemberPay = !selectedMember.payment_status.is_payment_disabled;
		}
		showPaymentDetails = false;
	};

	onMount(async () => {
		document.getElementById('member_dropdown')?.focus();
	});
</script>

<div class="grid gap-4 grid-cols-3">
	<div class="col-span-2">
		{#if showPaymentDetails}
			<div class="mb-10">
				<Alert border color="dark">
					<InfoCircleSolid slot="icon" class="w-5 h-5" />
					<div class="text-lg mb-5">Στοιχεία πληρωμής</div>
					<ul class="text-base">
						<li>
							Ημερομηνία: <span class="font-bold">{paymentDetails.issued_at_formatted}</span>
						</li>
						<li>
							Απόδειξη: <span class="font-bold">{paymentDetails.receipt_no}</span>
						</li>
						<li>
							Ποσό: <span class="font-bold">{paymentDetails.amount} €</span>
						</li>
						<li>
							Ονομα: <span class="font-bold">{paymentDetails.member_name}</span>
						</li>
						<li>
							Αριθμός μέλους: <span class="font-bold">{paymentDetails.member_no}</span>
						</li>
						<li>
							Κατάσταση: <span class="font-bold">{paymentDetails.status}</span>
						</li>
					</ul>
				</Alert>
			</div>
		{/if}

		{#if !selectedMember}
			<InputGroup>
				<div class="w-full">
					<AutocompleteMemberFormField
						id="member_dropdown"
						bind:value={form.member_id}
						bind:error={errors.member_id}
						on:select={onSelectMember}
					/>
				</div>
			</InputGroup>
		{/if}

		{#if selectedMember}
			<InputGroup>
				<h2 class="w-full flex">
					<a
						class="flex grow items-center text-sky-500 hover:underline after:content-['_↗']"
						target="_blank"
						href={`/admin/member?id=${form.member_id}`}
					>
						{selectedMember.name_formatted}
					</a>

					<a class="text-red-500" href={'#'} on:click={() => (selectedMember = undefined)}>
						<svg
							class="w-5 h-5"
							aria-hidden="true"
							xmlns="http://www.w3.org/2000/svg"
							width="24"
							height="24"
							fill="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								fill-rule="evenodd"
								d="M2 12C2 6.477 6.477 2 12 2s10 4.477 10 10-4.477 10-10 10S2 17.523 2 12Zm7.707-3.707a1 1 0 0 0-1.414 1.414L10.586 12l-2.293 2.293a1 1 0 1 0 1.414 1.414L12 13.414l2.293 2.293a1 1 0 0 0 1.414-1.414L13.414 12l2.293-2.293a1 1 0 0 0-1.414-1.414L12 10.586 9.707 8.293Z"
								clip-rule="evenodd"
							/>
						</svg>
					</a>
				</h2>

				<IssuesAlert {issues}>
					<p class="text-sm w-full">
						Μεταβείτε στη
						<a class="text-sky-500" target="_blank" href={`/admin/member?id=${form.member_id}`}>
							σελίδα του μέλους
						</a>
						για να λύσετε.
					</p>
				</IssuesAlert>

				<div
					class={selectedMember.payment_status.ok
						? 'text-sm p-5 bg-green-300'
						: 'text-sm p-5 bg-red-300'}
				>
					{selectedMember.payment_status.formatted}
				</div>
			</InputGroup>

			{#if canMemberPay}
				<div class="grid gap-2 grid-cols-2">
					<div class="col-span-1">
						<Form bind:form url="/payments" bind:errors on:success={onSuccess} on:failure>
							<Input
								id="member_id"
								bind:value={form.member_id}
								name="member_id"
								required="true"
								type="hidden"
							/>

							<div class="w-full mb-5">
								<Tabs>
									<TabItem open on:click={() => (paymentMethod = 'amount')}>
										<span slot="title">Ποσό €</span>

										<InputGroup legend="Είσπραξη">
											<InputField
												id="amount_in_euros"
												type="number"
												min={0}
												max={1000}
												bind:value={form.amount}
												bind:error={errors.amount}
												help="Το ποσό πρέπει να είναι τουλάχιστον 2€ που αντιστοιχεί στη μηνιαία συνδρομή. Αν το ποσό ειναι μονός αριθμός, η στρογγυλοποίηση γίνεται προς τα κάτω (πχ 5€ = 2 μήνες)."
											/>
											<div class="w-full">
												<ToggleField
													bind:checked={form.contains_registration_fee}
													label="Περιέχει εγγραφή (4€)"
												/>
											</div>
										</InputGroup>

										<InputGroup legend="Απόδειξη">
											<div class="w-full">
												<ToggleField bind:checked={form.without_receipt} label="Χωρίς απόδειξη" />
											</div>

											<div class="w-full" class:hidden={form.without_receipt}>
												<InputField
													id="receipt_no"
													type="number"
													label="Απόδειξη"
													bind:value={form.receipt_no}
													bind:error={errors.receipt_no}
													min={0}
													max={50}
												/>
											</div>
											<div class="w-full" class:hidden={form.without_receipt}>
												<InputField
													id="receipt_block_no"
													type="number"
													label="Μπλοκ αποδείξεων"
													bind:value={form.receipt_block_no}
													bind:error={errors.receipt_block_no}
													min={0}
													max={1000}
												/>
											</div>
										</InputGroup>
									</TabItem>
									<TabItem on:click={() => (paymentMethod = 'months')}>
										<span slot="title">Μήνες</span>
										<div class="text-sm mb-3">
											Ειδικού τύπου πληρωμή για ανέργους και για διόρθωση λαθών.
										</div>
										<InputField
											id="months"
											type="number"
											min={0}
											max={24}
											bind:value={form.months}
											bind:error={errors.months}
											help="Αν η τιμή είναι μεγαλύτερη του μηδέν (0), τότε η τιμή στο πεδίο 'Ποσό' αγνοείται."
										/>
									</TabItem>
								</Tabs>
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

							<InputField id="comments" type="textarea" label="Σχόλια" bind:value={form.comments} />
						</Form>
					</div>
					<div class="col-span-1">
						<Card class="mb-10">
							<div class="underline mb-5">Πληροφορίες υπολογισμού χρωστούμενων</div>
							<div class="text-sm text-black grid grid-cols-2 gap-4">
								<div>Η τελευταία πληρωμή ήταν μέχρι και:</div>
								<div>{selectedMember.payment_status.last_payment_until_formatted}</div>
								<div>Εγγραφή:</div>
								<div>{selectedMember.payment_status.registered_at_formatted}</div>
								<div>Βάση υπολογισμού για όλους:</div>
								<div>{selectedMember.payment_status.cutoff_date_formatted}</div>
							</div>
						</Card>
						<Card>
							<div class="underline">Προηγούμενες πληρωμές του μέλους</div>
							<MemberPaymentHistory record={selectedMember} short={true} />
						</Card>
					</div>
				</div>
			{/if}
		{/if}
	</div>
	<div class="col-span-1">
		<Card padding="sm" size="xl" class="h-96">
			<div class="underline mb-5">Τελευταίες εισπράξεις του ταμία</div>
			{#key paymentDetails}
				<LatestCashierPayments bind:this={latestCashierPayments} />
			{/key}
		</Card>
	</div>
</div>
