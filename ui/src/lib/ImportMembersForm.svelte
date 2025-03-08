<script lang="ts">
	import { type ImportMembersForm } from '$lib/types';
	import Dropzone from '$lib/Dropzone.svelte';

	let form: ImportMembersForm = {};
	let errors: ImportMembersForm = {};
	let url: string = '/members/import';
	let method: 'POST' = 'POST';
</script>

<Dropzone />
