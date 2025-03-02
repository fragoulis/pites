<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import InputField from '$lib/InputField.svelte';
	import SubmitButtonWithState from '$lib/SubmitButtonWithState.svelte';
	import AlertWarning from '$lib/AlertWarning.svelte';
	import { isLoggedIn } from '$lib/utils';
	import pb from '$lib/pocketbase';

	let email: string = '';
	let password: string = '';
	let isLoading = false;
	let error: string = '';

	onMount(async () => {
		if (isLoggedIn()) {
			goto('/admin');
		}
	});

	const login = async () => {
		if (isLoading) {
			return;
		}

		isLoading = true;

		loginUser()
			.then(() => {
				error = '';
				goto('/');
			})
			.catch((e) => {
				console.log(e);
				error = 'Invalid login credentials.';
				document.getElementById('email')?.focus();
			})
			.finally(() => {
				isLoading = false;
			});
	};

	const loginUser = () => {
		return pb.collection('users').authWithPassword(email, password);
	};
</script>

<svelte:head>
	<title>Login</title>
</svelte:head>

<section class="h-screen">
	<div class="h-full">
		<!-- Left column container with background-->
		<div class="g-6 flex h-full flex-wrap items-center justify-center">
			<div
				class="shrink-1 mb-12 grow-0 basis-auto md:mb-0 md:w-9/12 md:shrink-0 lg:w-6/12 xl:w-6/12"
			>
				<img src="/login.webp" class="w-full" alt="hero" />
			</div>

			<!-- Right column container -->
			<div class="mb-12 md:mb-0 md:w-8/12 lg:w-5/12 xl:w-5/12">
				{#if error !== ''}
					<AlertWarning>{error}</AlertWarning>
				{/if}
				<form on:submit|preventDefault={login}>
					<div class="text-left">
						<InputField
							type="email"
							id="email"
							label="Email address"
							required={true}
							bind:value={email}
						/>
						<InputField
							type="password"
							id="password"
							label="Password"
							required={true}
							bind:value={password}
						/>
						<SubmitButtonWithState loading={isLoading}>Σύνδεση</SubmitButtonWithState>
					</div>
				</form>
			</div>
		</div>
	</div>
</section>
