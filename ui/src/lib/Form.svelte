<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import pb from '$lib/pocketbase';
	import { objectMap } from '$lib/utils';
	import SubmitButtonWithState from '$lib/SubmitButtonWithState.svelte';
	import { type MethodOptions } from '$lib/types';

	const dispatch = createEventDispatcher();

	export let form: any;
	export let errors: any;
	export let url: string;
	export let method: MethodOptions = 'POST';

	let saving: boolean = false;

	const submit = async () => {
		if (!dispatch('validate', '', { cancelable: true })) {
			return;
		}

		saving = true;
		try {
			dispatch('beforeSend');

			const record = await pb.send(url, { method: method, body: form });

			dispatch('success', { model: record });
		} catch (err: any) {
			if (err.originalError) {
				console.debug(err.originalError);
			}

			if (err.data) {
				errors = {};

				let errData = err.data.data;

				for (let field in errData) {
					if ('message' in errData[field]) {
						// single error
						errors[field] = errData[field].message;
					} else {
						// array of errors
						errors[field] = [];
						for (let index in errData[field]) {
							errors[field][index] = errData[field][index].message;
						}
					}
				}

				const firstErrorId = Object.keys(errors)[0];
				document.getElementById(firstErrorId)?.focus();
				dispatch('failure', { message: err.message, errors: errors });
			} else {
				console.error(err);
			}
		}
		saving = false;
	};
</script>

<form class="mb-10" on:submit|preventDefault={submit} novalidate>
	<slot />

	<div class="text-left">
		<SubmitButtonWithState loading={saving}>Αποθήκευση</SubmitButtonWithState>
	</div>
</form>
