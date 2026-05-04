<script lang="ts">
	import { enhance } from '$app/forms';
	import { auth } from '$lib/auth.svelte';
	let { form } = $props();
	let loading = $state(false);
</script>

<div class="mx-auto mt-16 max-w-md">
	<div class="rounded-3xl border border-slate-100 bg-white p-8 shadow-2xl shadow-slate-200/50">
		<header class="mb-8 text-center">
			<div
				class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-amber-100 text-amber-600"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2"
					stroke="currentColor"
					class="h-8 w-8"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M15.362 5.214A8.252 8.252 0 0 1 12 21 8.25 8.25 0 0 1 6.038 7.047 8.287 8.287 0 0 0 9 9.601a8.983 8.983 0 0 1 3.361-6.867 8.21 8.21 0 0 0 3 2.48Z"
					/>
				</svg>
			</div>
			<h1 class="text-2xl font-black text-slate-900">Unlock Full Access</h1>
			<p class="mt-2 text-sm font-medium text-slate-500">
				Choose a password to finish your account and start hosting your own grids.
			</p>
		</header>

		<form
			method="POST"
			use:enhance={() => {
				loading = true;
				return async ({ result, update }) => {
					if (result.type === 'redirect' || result.type === 'success') {
						// Manually update the local store so the UI reacts instantly
						if (auth.user) {
							auth.user.is_guest = false;
						}
					}
					loading = false;
					await update();
				};
			}}
			class="flex flex-col gap-5"
		>
			<div class="space-y-1">
				<label for="password" class="text-xs font-bold tracking-widest text-slate-400 uppercase"
					>Set Password</label
				>
				<input
					name="password"
					id="password"
					type="password"
					placeholder="••••••••"
					class="w-full rounded-xl border border-slate-200 p-3 transition-all outline-none focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10"
					required
				/>
			</div>

			<div class="space-y-1">
				<label
					for="confirm_password"
					class="text-xs font-bold tracking-widest text-slate-400 uppercase">Confirm Password</label
				>
				<input
					name="confirm_password"
					id="confirm_password"
					type="password"
					placeholder="••••••••"
					class="w-full rounded-xl border border-slate-200 p-3 transition-all outline-none focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10"
					required
				/>
			</div>

			{#if form?.message}
				<div class="rounded-lg bg-red-50 p-3 text-center text-xs font-bold text-red-500">
					{form.message}
				</div>
			{/if}

			<button
				type="submit"
				disabled={loading}
				class="mt-2 rounded-2xl bg-slate-900 py-4 font-black text-white transition-all hover:bg-slate-800 hover:shadow-xl active:scale-[0.98] disabled:opacity-50"
			>
				{loading ? 'Securing Account...' : 'Upgrade Now'}
			</button>
		</form>
	</div>

	<p class="mt-6 text-center text-xs font-bold tracking-widest text-slate-400 uppercase">
		Bellesquare Pro • Unlimited Boards
	</p>
</div>
