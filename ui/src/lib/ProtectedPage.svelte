<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { isLoggedIn } from '$lib/utils';

	onMount(async () => {
		if (!isLoggedIn()) {
			goto('/login');
		}
	});
</script>

<slot />
