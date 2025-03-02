<script lang="ts">
	export let legend: string = '';
</script>

<div class="grid gap-2 grid-cols-1 rounded-xl mb-5">
	{#if legend != ''}
		<div class="col-span-1 text-xl">
			{legend}
		</div>
	{/if}
	<slot />
</div>
