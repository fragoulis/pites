<script lang="ts">
	import { page } from '$app/stores';
	import {
		Avatar,
		Sidebar,
		SidebarBrand,
		SidebarGroup,
		SidebarItem,
		SidebarWrapper
	} from 'flowbite-svelte';
	import {
		ChartPieSolid,
		UserSolid,
		ArrowLeftToBracketOutline,
		BuildingSolid,
		UsersGroupSolid,
		EuroOutline
	} from 'flowbite-svelte-icons';
	import { logOutUser } from '$lib/utils';
	import { isAdmin } from '$lib/utils';

	let spanClass = 'flex-1 ms-3 whitespace-nowrap';

	let site = {
		name: 'ΣΕΤΗΠ',
		href: '/',
		img: '/setip.png'
	};

	$: activeUrl = $page.url.pathname;

	export let loggedInUser;
</script>

<Sidebar
	{activeUrl}
	asideClass="fixed top-0 left-0 z-40 w-72 h-screen border-r transition-transform -translate-x-full sm:translate-x-0"
>
	<SidebarWrapper>
		<SidebarGroup>
			<SidebarBrand {site} />
			<div class="flex">
				<div class="w-8">
					<Avatar border />
				</div>
				<SidebarItem label={loggedInUser?.email} />
			</div>
		</SidebarGroup>
		<SidebarGroup border>
			<SidebarItem label="Πίνακας Ελέγχου" href="/admin">
				<svelte:fragment slot="icon">
					<ChartPieSolid
						class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
					/>
				</svelte:fragment>
			</SidebarItem>
			<SidebarItem label="Μέλη" {spanClass} href="/admin/members">
				<svelte:fragment slot="icon">
					<UserSolid
						class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
					/>
				</svelte:fragment>
			</SidebarItem>
			<!-- <SidebarItem label="Πληρωμές" href="/admin/payments">
				<svelte:fragment slot="icon">
					<EuroOutline
						class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
					/>
				</svelte:fragment>
			</SidebarItem> -->
			{#if isAdmin()}
				<SidebarItem label="Εταιρείες" href="/admin/companies">
					<svelte:fragment slot="icon">
						<BuildingSolid
							class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
						/>
					</svelte:fragment>
				</SidebarItem>
				<SidebarItem label="Συνελεύσεις" href="/admin/assemblies">
					<svelte:fragment slot="icon">
						<UsersGroupSolid
							class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
						/>
					</svelte:fragment>
				</SidebarItem>
			{/if}
		</SidebarGroup>

		<SidebarGroup border>
			<SidebarItem label="Νέα Πληρωμή" href="/admin/payments/create">
				<svelte:fragment slot="icon">
					<EuroOutline
						class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
					/>
				</svelte:fragment>
			</SidebarItem>
		</SidebarGroup>

		<SidebarGroup border>
			<SidebarItem label="Αποσύνδεση" on:click={logOutUser}>
				<svelte:fragment slot="icon">
					<ArrowLeftToBracketOutline
						class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
					/>
				</svelte:fragment>
			</SidebarItem>
		</SidebarGroup>
	</SidebarWrapper>
</Sidebar>
