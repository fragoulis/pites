<script lang="ts">
	import { Alert } from 'flowbite-svelte';
	import { InfoCircleSolid } from 'flowbite-svelte-icons';
	import { type Issue } from '$lib/types';

	export let issues: Issue[];
</script>

{#if issues.length > 0}
	<Alert class="!items-start my-5">
		<span slot="icon">
			<InfoCircleSolid class="w-5 h-5" />
		</span>
		<p class="text-base">Βρήκαμε τα παρακάτω προβλήματα:</p>
		<ul class="mt-1.5 ms-4 list-disc list-inside">
			{#each issues as issue}
				<li>{issue.description}</li>
			{/each}
		</ul>
	</Alert>

	<slot />
{/if}
