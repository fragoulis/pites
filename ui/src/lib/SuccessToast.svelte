<script lang="ts">
	import { Toast } from 'flowbite-svelte';
	import { CheckCircleSolid } from 'flowbite-svelte-icons';
	import { slide } from 'svelte/transition';

	export let open: boolean = false;
</script>

<Toast
	dismissable={false}
	position="top-right"
	color="green"
	transition={slide}
	bind:open
	divClass="z-[60] w-full max-w-xs p-4 text-gray-500 bg-green-100 shadow dark:text-gray-400 dark:bg-gray-800 gap-3"
>
	<svelte:fragment slot="icon">
		<CheckCircleSolid slot="icon" class="w-5 h-5" />
	</svelte:fragment>
	<slot />
</Toast>
