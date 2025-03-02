<script lang="ts">
	import SuccessToast from '$lib/SuccessToast.svelte';
	import ErrorToast from '$lib/ErrorToast.svelte';
	import { sharedToast } from '$lib/store';

	// Show toast on either event
	$: showSuccess = $sharedToast.show && $sharedToast.success;
	$: showError = $sharedToast.show && !$sharedToast.success;

	// On change, make sure toast disappears after a few seconds.
	$: showSuccess && setTimeout(() => (showSuccess = false), 4000);
	$: showError && setTimeout(() => (showError = false), 2000);
</script>

<SuccessToast bind:open={showSuccess}>{$sharedToast.message}</SuccessToast>
<ErrorToast bind:open={showError}>{$sharedToast.message}</ErrorToast>
