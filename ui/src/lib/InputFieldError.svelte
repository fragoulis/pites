<script lang="ts">
	export let error: string;
</script>

{#if error}
	<span
		class="mb-2 text-sm text-red-500 peer-[&:not(:placeholder-shown):not(:focus):invalid]:block"
	>
		{error}
	</span>
{/if}
