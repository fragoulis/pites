<script lang="ts">
	import { onMount } from 'svelte';
	import { Spinner } from 'flowbite-svelte';
	import { isAdmin } from '$lib/utils';
	import AlertWarning from '$lib/AlertWarning.svelte';

	let hasPermissions: boolean = false;
	let loading: boolean = true;

	onMount(async () => {
		hasPermissions = isAdmin();
		loading = false;
	});
</script>

{#if hasPermissions}
	<slot />
{:else if loading}
	<section class="m-5">
		<Spinner class="me-3" size="10" color="primary" />
	</section>
{:else}
	<section class="m-5">
		<AlertWarning>You do not have permissions to view this page.</AlertWarning>
	</section>
{/if}
