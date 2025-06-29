<script lang="ts">
	import ToggleField from '$lib/ToggleField.svelte';
	import AutocompleteCompanyFormField from '$lib/AutocompleteCompanyFormField.svelte';
	import { type DatatableSearchForm, type Chapter } from '$lib/types';
	import InputField from '$lib/InputField.svelte';
	import MultiSelectAsyncFormField from '$lib/MultiSelectAsyncFormField.svelte';
	import SelectAsyncFormField from './SelectAsyncFormField.svelte';
	import { Tabs, TabItem } from 'flowbite-svelte';

	export let form: DatatableSearchForm;

	let companyQuery: string = '';
	let addressFilterType: 'address' | 'chapter' = 'address';

	export const reset = () => {
		companyQuery = '';

		form.active_only = true;
		form.with_comments = false;
		form.company_id = '';
		form.business_type_ids = [];
		form.address_city_ids = [];
		form.legacy_area = '';
		form.chapter_id = '';
		form.with_fixed_monthly_payment = undefined;
	};
</script>

<div class="w-full">
	<ToggleField name="active_only" bind:checked={form.active_only} label="Μόνο τα ενεργά μέλη" />
</div>
<div class="w-full">
	<ToggleField name="with_comments" bind:checked={form.with_comments} label="Με σχόλια" />
</div>
<div class="w-full">
	<ToggleField
		name="with_fixed_monthly_payment"
		bind:checked={form.with_fixed_monthly_payment}
		label="Πληρώνουν με πάγια εντολή"
	/>
</div>

<div class="w-full grid grid-cols-2 gap-4 px-4 pt-4 bg-gray-100 rounded-lg border">
	<AutocompleteCompanyFormField
		multi={true}
		required={false}
		bind:value={form.company_id}
		bind:query={companyQuery}
	/>

	<MultiSelectAsyncFormField
		url="/companies/business_types"
		label="Ομάδα"
		bind:values={form.business_type_ids}
	/>
</div>

<div class="w-full px-4 bg-gray-100 rounded-lg border">
	<Tabs tabStyle="underline" contentClass="p-4 rounded-lg" divider={false}>
		<TabItem open on:click={() => (addressFilterType = 'address')}>
			<span slot="title">Με διεύθυνση</span>
			<div class="w-full grid grid-cols-2 gap-4">
				<MultiSelectAsyncFormField
					url="/address_cities"
					label="Δήμος"
					help="Φιλτράρει μέλη βάσει Δήμου. Χρειάζεται το μέλος να έχει δομημένη διεύθυνση. Μπορεί να συνδυαστεί με την Περιοχή (Παλιά βάση)."
					bind:values={form.address_city_ids}
				/>

				<InputField
					label="Περιοχή (Παλιά βάση)"
					help="Φιλτράρει μέλη βάσει περιοχής όπως ήταν περασμένη στην παλιά βάση. Υποστηρίζει πολλά φίλτρα χωρισμένα με κόμα."
					bind:value={form.legacy_area}
					placeholder="ν. ιωνια, νεα ιωνια"
				/>
			</div>
		</TabItem>
		<TabItem on:click={() => (addressFilterType = 'chapter')}>
			<span slot="title">Με παράρτημα</span>
			<div class="w-full grid grid-cols-2 gap-4">
				<SelectAsyncFormField
					url="/chapters"
					label="Παράρτημα"
					bind:value={form.chapter_id}
					help="Η αναζήτηση βάσει παραρτήματος είναι πρακτικά μια προεπιλεγμένη αναζήτηση βάσει διεύθυνσης."
				/>
			</div>
		</TabItem>
	</Tabs>
</div>
