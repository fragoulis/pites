<script lang="ts">
	import { onMount, tick } from 'svelte';
	import pb from '$lib/pocketbase';
	import { Select } from 'flowbite-svelte';

	export let value: string = '';
	export let url: string;
	export let sendOptions: Object = {};
	export let optionValue: string = 'id';
	export let optionName: string = 'name';
	export let withDefault: boolean = false;

	let options: any = [];
	if (withDefault) {
		options.push({ value: 'none', name: 'Όλα' });
	}

	onMount(async () => {
		try {
			const res: any[] = await pb.send(url, sendOptions);
			res.forEach((record: any) => {
				options.push({
					value: record[optionValue],
					name: record[optionName]
				});
			});

			// fucking svelte
			options = options;
		} catch (e: any) {
			console.error(e);
		}
	});
</script>

<Select bind:items={options} bind:value />
