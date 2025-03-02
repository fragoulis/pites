<script lang="ts">
	import { Drawer, CloseButton, A } from 'flowbite-svelte';
	import { UserCircleOutline } from 'flowbite-svelte-icons';
	import { sineIn } from 'svelte/easing';
	import MemberView from '$lib/MemberView.svelte';
	import { type Member } from '$lib/types';

	export let record: Member;
	export let hidden = true;

	let transitionParams = {
		x: 320,
		duration: 200,
		easing: sineIn
	};
</script>

<Drawer placement="right" transitionType="fly" {transitionParams} bind:hidden width="w-2/3">
	<div class="flex items-center">
		<h5
			id="drawer-label"
			class="inline-flex items-center mb-4 text-xl font-bold text-gray-800 dark:text-gray-400"
		>
			{#key record.name_formatted}
				<UserCircleOutline class="w-5 h-5 me-2.5" />{record.name_formatted}
			{/key}
		</h5>
		<CloseButton on:click={() => (hidden = true)} class="mb-4 dark:text-white" />
	</div>
	<MemberView bind:record />
</Drawer>
