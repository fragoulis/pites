<script lang="ts">
	import { Label, Helper, type InputType, type FormSizeType, Textarea } from 'flowbite-svelte';

	export let id: string = '';
	export let label: string = '';
	export let type: InputType = 'text';
	export let value: string = '';
	export let required: boolean = false;
	export let error: string = '';
	export let placeholder: string = '';
	export let size: FormSizeType = 'md';
</script>

<div class="mb-6">
	<div class="relative">
		<Label class="space-y-2">
			{#if label !== ''}
				<span>
					{label}
					{#if required}*{/if}
				</span>
			{/if}
			<slot>
				<Textarea {id} {type} {placeholder} {size} bind:value {required} />
			</slot>
			{#if error}
				<Helper class="mt-2" color="red">
					{error}
				</Helper>
			{/if}
		</Label>
	</div>
</div>
