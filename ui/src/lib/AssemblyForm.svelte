<script lang="ts">
	import { type AssemblyForm } from '$lib/types';
	import TextField from '$lib/TextField.svelte';
	import SubmitButtonWithState from '$lib/SubmitButtonWithState.svelte';
	import { createEventDispatcher } from 'svelte';
	import { objectMap } from '$lib/utils';
	import pb from '$lib/pocketbase';
	import InputField from '$lib/InputField.svelte';

	const dispatch = createEventDispatcher();

	let form: AssemblyForm = {};
	let errors: AssemblyForm = {};
	let saving: boolean = false;

	const submit = async () => {
		errors = {};
		saving = true;
		try {
			const record = await pb.collection('assemblies').create(form);

			dispatch('create:success', { record: record });
		} catch (err: any) {
			errors = objectMap(err.data.data, (v) => {
				return v.message;
			});
			const firstErrorId = Object.keys(errors)[0];
			document.getElementById(firstErrorId)?.focus();
			dispatch('create:failed', { errors: errors });
		}
		saving = false;
	};
</script>

<form on:submit|preventDefault={submit} novalidate>
	<InputField label="Ημερομηνία" type="date" bind:value={form.date} bind:error={errors.date} />
	<TextField label="Σχόλια" bind:value={form.comments} bind:error={errors.comments} />

	<div class="text-left">
		<SubmitButtonWithState loading={saving}>Αποθήκευση</SubmitButtonWithState>
	</div>
</form>
