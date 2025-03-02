<script lang="ts">
	import { Toggle } from 'flowbite-svelte';

	export let name: string = '';
	export let label: string;
	export let checked: boolean = false;
</script>

<div class="mb-3">
	<Toggle {name} bind:checked on:change><div class="text-sm">{label}</div></Toggle>
</div>
