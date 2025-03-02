<script lang="ts">
	import { Alert } from 'flowbite-svelte';
	import { InfoCircleSolid } from 'flowbite-svelte-icons';
</script>

<Alert color="yellow" border class="mb-5">
	<InfoCircleSolid slot="icon" class="w-5 h-5" />
	<slot />
</Alert>
