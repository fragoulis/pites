<script lang="ts">
	import { onMount, tick } from 'svelte';
	import pb from '$lib/pocketbase';
	import { MultiSelect } from 'flowbite-svelte';

	export let values: any = [];
	export let url: string;
	export let sendOptions: Object = {};
	export let optionValue: string = 'id';
	export let optionName: string = 'name';

	let options: any = [];

	onMount(async () => {
		try {
			const res: any[] = await pb.send(url, sendOptions);
			res?.forEach((record: any) => {
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

<MultiSelect bind:items={options} bind:value={values} />
