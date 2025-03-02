<script lang="ts">
	import {
		NumberInput,
		Input,
		Label,
		Helper,
		type InputType,
		type FormSizeType,
		Textarea
	} from 'flowbite-svelte';

	export let id: string | undefined = undefined;
	export let label: string = '';
	export let type: InputType | 'textarea' = 'text';
	export let value: string | number | undefined = '';
	export let required: boolean = false;
	export let error: string = '';
	export let placeholder: string | undefined = undefined;
	export let size: FormSizeType = 'md';
	export let min: number | undefined = undefined;
	export let max: number | undefined = undefined;
	export let help: string = '';
</script>

<div class="mb-4">
	{#if type == 'textarea'}
		<Textarea rows="8" placeholder={label} bind:value on:keyup />
		{#if help != ''}
			<Helper class="mt-2 text-sm">
				{help}
			</Helper>
		{/if}
	{:else}
		<Label class="space-y-2" color={error ? 'red' : 'gray'}>
			{#if label !== ''}
				<span>
					{label}
					{#if required}*{/if}
				</span>
			{/if}
			<slot>
				{#if type == 'number'}
					<NumberInput
						{id}
						{placeholder}
						bind:value
						{required}
						{type}
						{size}
						{min}
						{max}
						color={error ? 'red' : 'base'}
						on:keyup
					/>
				{:else}
					<Input
						{id}
						{placeholder}
						bind:value
						{required}
						{type}
						{size}
						color={error ? 'red' : 'base'}
						on:keyup
					/>
				{/if}
			</slot>
			{#if help != ''}
				<Helper class="mt-2 text-sm">
					{help}
				</Helper>
			{/if}
			{#if error}
				<Helper class="mt-2 text-sm" color="red">
					{error}
				</Helper>
			{/if}
		</Label>
	{/if}
</div>
