<script lang="ts">
	import { Card } from 'flowbite-svelte';
</script>

<div class="grid grid-cols-4 gap-4 my-10 mx-5">
	<Card>
		<a href="/">
			<h5 class="mb-2 text-2xl font-semibold tracking-tight text-gray-900 dark:text-white">
				Πληρωμές
			</h5>
		</a>
		<a
			href="/admin/payments/create"
			class="inline-flex items-center text-primary-600 hover:underline"
		>
			Πατήστε εδώ για να μεταβείτε στη φόρμα πληρωμών.
		</a>
	</Card>
</div>
