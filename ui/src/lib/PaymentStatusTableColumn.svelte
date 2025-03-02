<script lang="ts">
	import { Indicator, type IndicatorColorType } from 'flowbite-svelte';
	import { type Member, type MemberPaymentStatus } from '$lib/types';

	export let record: Member;
	const status: MemberPaymentStatus = record.payment_status;
	const color: IndicatorColorType = status?.ok ? 'green' : 'red';
</script>

<span class="flex items-center">
	<Indicator size="lg" class="me-1.5" {color} />
	{status?.formatted}
</span>
