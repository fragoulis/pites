<script lang="ts">
	export let href: string = '#!';
</script>

<a
	{href}
	class="text-primary transition duration-150 ease-in-out hover:text-primary-600 focus:text-primary-600 active:text-primary-700 dark:text-primary-400 dark:hover:text-primary-500 dark:focus:text-primary-500 dark:active:text-primary-600"
>
	<slot>Primary link</slot>
</a>
