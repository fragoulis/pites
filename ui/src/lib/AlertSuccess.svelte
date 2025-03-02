<script lang="ts">
	import { Alert } from 'flowbite-svelte';
	import { InfoCircleSolid } from 'flowbite-svelte-icons';
</script>

<Alert color="green">
	<InfoCircleSolid slot="icon" class="w-5 h-5" />
	<slot />
</Alert>
