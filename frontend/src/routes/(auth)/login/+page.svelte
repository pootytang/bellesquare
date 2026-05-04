<script lang="ts">
	import { resolve } from '$app/paths';
	import { enhance } from '$app/forms';
	import { page } from '$app/state';

	let { form } = $props();
	let loading = $state(false);

	// 1. Determine if they are specifically trying to join a board
	let isGuestPath = $derived(page.url.searchParams.get('redirectTo')?.includes('/board'));

	// 2. Track if the user wants to log in as a full member
	// If it's not a guest path, we'll let them toggle.
	let showPassword = $state(false);

	// If it's a guest path (joining a board), we always hide password.
	// If not, we show/hide based on the toggle.
	let effectiveShowPassword = $derived(!isGuestPath && showPassword);
</script>

<div class="mx-auto mt-20 max-w-sm rounded-2xl border border-slate-100 bg-white p-8 shadow-xl">
	<div class="mb-8">
		<h1 class="text-2xl font-bold text-slate-900">
			{isGuestPath ? 'Join the Game' : 'Sign In'}
		</h1>
		<p class="text-sm text-slate-500">
			{isGuestPath || !showPassword
				? 'Enter your email to continue.'
				: 'Welcome back, enter your credentials.'}
		</p>
	</div>

	<form
		method="POST"
		use:enhance={() => {
			loading = true;
			return async ({ update }) => {
				loading = false;
				await update();
			};
		}}
		class="flex flex-col gap-5"
	>
		<div class="space-y-1">
			<label for="email" class="text-xs font-bold text-slate-400 uppercase">Email</label>
			<input
				name="email"
				id="email"
				type="email"
				placeholder="name@example.com"
				class="w-full rounded-xl border border-slate-200 p-3 transition-all outline-none focus:ring-2 focus:ring-blue-500"
				required
			/>
		</div>

		{#if effectiveShowPassword}
			<div class="space-y-1">
				<div class="flex items-center justify-between">
					<label for="password" class="text-xs font-bold text-slate-400 uppercase">Password</label>
					<button
						type="button"
						onclick={() => (showPassword = false)}
						class="text-[10px] font-bold text-blue-600 uppercase hover:underline"
					>
						Guest Login?
					</button>
				</div>
				<input
					name="password"
					id="password"
					type="password"
					placeholder="••••••••"
					class="w-full rounded-xl border border-slate-200 p-3 transition-all outline-none focus:ring-2 focus:ring-blue-500"
					required
				/>
			</div>
		{:else}
			<!-- Hidden field to satisfy backend validation if password isn't visible -->
			<input type="hidden" name="password" value="GUEST_SESSION" />

			{#if !isGuestPath}
				<button
					type="button"
					onclick={() => (showPassword = true)}
					class="text-left text-xs font-bold text-blue-600 hover:underline"
				>
					Are you a Member? Sign in with password
				</button>
			{/if}
		{/if}

		{#if form?.message}
			<div class="rounded-lg border border-red-100 bg-red-50 p-3">
				<p class="text-center text-xs font-semibold text-red-600">{form.message}</p>
			</div>
		{/if}

		<button
			type="submit"
			disabled={loading}
			class="rounded-xl bg-blue-600 p-4 font-bold text-white shadow-lg shadow-blue-100 transition-all hover:bg-blue-700 disabled:opacity-50"
		>
			{loading
				? 'Authenticating...'
				: isGuestPath || !showPassword
					? 'Continue as Guest'
					: 'Sign In'}
		</button>
	</form>

	<p class="mt-8 text-center text-sm text-slate-500">
		New here? <a
			href={resolve(
				isGuestPath
					? `/register?redirectTo=${encodeURIComponent(page.url.searchParams.get('redirectTo') || '')}`
					: '/register'
			)}
			class="font-bold text-blue-600 hover:underline">Create account</a
		>
	</p>
</div>
