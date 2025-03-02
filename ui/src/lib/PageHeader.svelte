<script lang="ts">
	let title: HTMLSpanElement;
</script>

<span class="hidden" bind:this={title}><slot /></span>

<svelte:head>
	<title>{title?.innerText || ''}</title>
</svelte:head>

<div class="flex flex-no-wrap w-full justify-between my-5">
	<div class="text-xl mt-2">
		{title?.innerText || ''}
	</div>
	{#if $$slots.right}
		<slot name="right" />
	{/if}
</div>
<hr class="my-5 h-0.5 border-t-0 bg-neutral-100 dark:bg-white/10" />
