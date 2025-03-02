<script lang="ts">
	import ProtectedPage from '$lib/ProtectedPage.svelte';
	import SideNav from '$lib/SideNav.svelte';
	import SharedEventToast from '$lib/SharedEventToast.svelte';
	import '../../app.css';

	export let data;
</script>

<ProtectedPage>
	<div>
		<SideNav loggedInUser={data.user} />
		<div class="!pl-[310px] !pr-[10px]" id="content">
			<div class="container mx-auto">
				<SharedEventToast />

				<slot />
			</div>
		</div>
	</div>
</ProtectedPage>
