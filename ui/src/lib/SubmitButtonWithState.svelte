<script lang="ts">
	import { Button, Spinner } from 'flowbite-svelte';

	export let loading: boolean = false;
</script>

<Button type="submit" disabled={loading}>
	{#if loading}
		<Spinner class="me-3" size="4" color="white" /><slot />...
	{:else}
		<slot />
	{/if}
</Button>
