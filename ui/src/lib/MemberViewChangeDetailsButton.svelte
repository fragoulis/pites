<script lang="ts">
	import { Button } from 'flowbite-svelte';
	import MemberChangeDetailsDrawer from '$lib/MemberChangeDetailsDrawer.svelte';

	export let hideDrawer: boolean = true;
</script>

<Button
	on:click={(e) => {
		e.preventDefault();
		hideDrawer = false;
	}}
>
	<slot name="label" />
</Button>

<MemberChangeDetailsDrawer bind:hidden={hideDrawer}>
	<slot name="title" slot="title" />
	<slot />
</MemberChangeDetailsDrawer>
