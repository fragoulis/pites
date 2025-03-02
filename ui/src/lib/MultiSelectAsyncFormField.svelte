<script lang="ts">
	import { Input, Label, Helper } from 'flowbite-svelte';
	import MultiSelectAsync from '$lib/MultiSelectAsync.svelte';

	export let url: string;
	export let label: string = '';
	export let help: string = '';
	export let values: string[];
</script>

<div class="mb-4">
	<Label class="space-y-2" color="gray">
		{#if label}
			<div class="mb-2">{label}</div>
		{/if}

		<MultiSelectAsync {url} bind:values />

		{#if help}
			<Helper class="mt-2 text-sm">{help}</Helper>
		{/if}
	</Label>
</div>
