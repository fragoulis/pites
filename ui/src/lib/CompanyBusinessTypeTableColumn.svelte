<script lang="ts">
	import { type Company } from '$lib/types';
	export let record: Company;
</script>

{record.business_type?.name || ''}
