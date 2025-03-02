<script lang="ts">
	import { Input, Label, Helper } from 'flowbite-svelte';
	import SelectAsync from '$lib/SelectAsync.svelte';

	export let url: string;
	export let label: string = '';
	export let help: string = '';
	export let value: string;
</script>

<Label class="space-y-2" color="gray">
	{#if label}
		<div class="mb-2">{label}</div>
	{/if}

	<SelectAsync {url} bind:value />

	{#if help}
		<Helper class="mt-2 text-sm">{help}</Helper>
	{/if}
</Label>
