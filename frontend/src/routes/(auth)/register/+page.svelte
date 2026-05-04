<script lang="ts">
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';

	let { form } = $props();
	let loading = $state(false);

	// Color picker state
	const presets = ['#2563eb', '#dc2626', '#16a34a', '#ca8a04', '#7c3aed', '#db2777'];
	let selectedColor = $state('#2563eb');

	// Use Svelte 5 derived state to get the redirect path
	let redirectTo = $derived(page.url.searchParams.get('redirectTo'));
</script>

<div class="mx-auto mt-10 max-w-md rounded-2xl border border-slate-100 bg-white p-8 shadow-xl">
	<div class="mb-6">
		<h1 class="text-2xl font-bold text-slate-900">Create Account</h1>
		<p class="text-sm text-slate-500">Join Bellesquare to manage your grids.</p>
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
		class="flex flex-col gap-4"
	>
		<div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
			<div class="space-y-1">
				<label for="user_name" class="text-xs font-bold tracking-wider text-slate-400 uppercase"
					>Username</label
				>
				<input
					name="user_name"
					id="user_name"
					type="text"
					placeholder="johndoe"
					class="w-full rounded-xl border border-slate-200 p-3 transition-all outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20"
					required
				/>
			</div>

			<div class="space-y-1">
				<label for="color" class="text-xs font-bold tracking-wider text-slate-400 uppercase"
					>Square Color</label
				>

				<div class="flex items-start gap-4 rounded-xl border border-slate-100 bg-slate-50/50 p-3">
					<div class="space-y-2">
						<span class="block text-[9px] font-black tracking-tighter text-slate-400 uppercase"
							>Presets</span
						>
						<div class="flex flex-wrap gap-1.5">
							{#each presets as color (color)}
								<button
									type="button"
									onclick={() => (selectedColor = color)}
									class="h-6 w-6 rounded-full border-2 shadow-sm transition-all hover:scale-110 active:scale-95"
									style="background-color: {color}; border-color: {selectedColor === color
										? 'white'
										: 'transparent'}; outline: {selectedColor === color
										? '2px solid #3b82f6'
										: 'none'}"
									aria-label="Select {color}"
								></button>
							{/each}
						</div>
					</div>

					<div class="mt-2 h-10 w-px bg-slate-200"></div>

					<div class="space-y-2">
						<span class="block text-[9px] font-black tracking-tighter text-slate-400 uppercase"
							>Custom</span
						>
						<div
							class="relative h-7 w-7 overflow-hidden rounded-full border-2 border-white shadow-sm ring-1 ring-slate-200"
						>
							<input
								type="color"
								name="preferred_color"
								id="preferred_color"
								bind:value={selectedColor}
								class="absolute -inset-2 h-12 w-12 cursor-pointer border-none bg-transparent"
							/>
						</div>
					</div>
				</div>
				<p class="text-[10px] font-medium text-slate-400 italic">These sets your square color.</p>
			</div>
		</div>

		<div class="grid grid-cols-2 gap-4">
			<div class="space-y-1">
				<label for="first_name" class="text-xs font-bold tracking-wider text-slate-400 uppercase"
					>First Name</label
				>
				<input
					name="first_name"
					id="first_name"
					type="text"
					placeholder="John"
					class="w-full rounded-xl border border-slate-200 p-3 transition-all outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20"
					required
				/>
			</div>
			<div class="space-y-1">
				<label for="last_name" class="text-xs font-bold tracking-wider text-slate-400 uppercase"
					>Last Name</label
				>
				<input
					name="last_name"
					id="last_name"
					type="text"
					placeholder="Doe"
					class="w-full rounded-xl border border-slate-200 p-3 transition-all outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20"
					required
				/>
			</div>
		</div>

		<div class="space-y-1">
			<label for="email" class="text-xs font-bold tracking-wider text-slate-400 uppercase"
				>Email Address</label
			>
			<input
				name="email"
				id="email"
				type="email"
				placeholder="john@example.com"
				class="w-full rounded-xl border border-slate-200 p-3 transition-all outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20"
				required
			/>
		</div>

		<div class="grid grid-cols-2 gap-4 rounded-xl bg-slate-50 p-4">
			<div class="space-y-2">
				<label for="venmo" class="block text-xs font-bold tracking-widest text-blue-500 uppercase">
					Venmo (Optional)
				</label>

				<div class="relative">
					<!-- The Visual Symbol -->
					<span class="absolute inset-y-0 left-0 flex items-center pl-4 font-bold text-slate-400">
						@
					</span>

					<input
						type="text"
						id="venmo"
						name="venmo_handle"
						placeholder="username"
						class="w-full rounded-xl border border-slate-200 bg-white py-3 pr-4 pl-9 font-medium text-slate-900 shadow-sm transition-all focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20"
					/>
				</div>

				<p class="text-[10px] text-slate-600 italic">
					Enter your handle only (e.g., Delane.Jackson)
				</p>
			</div>

			<div class="space-y-2">
				<label
					for="zelle_handle"
					class="block text-xs font-bold tracking-widest text-purple-500 uppercase"
				>
					Zelle (optional)
				</label>
				<div class="relative">
					<!-- The Visual Symbol -->
					<span
						class="absolute inset-y-0 left-0 flex items-center pl-4 text-slate-400 transition-colors group-focus-within:text-purple-500"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							viewBox="0 0 20 20"
							fill="currentColor"
							class="h-4 w-4"
						>
							<path
								d="M3 4a2 2 0 0 0-2 2v1.161l8.441 4.221a1.25 1.25 0 0 0 1.118 0L19 7.162V6a2 2 0 0 0-2-2H3Z"
							/>
							<path
								d="m19 8.839-7.97 3.985a2.75 2.75 0 0 1-2.46 0L1 8.839V14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V8.839Z"
							/>
						</svg>
					</span>
					<input
						name="zelle_handle"
						id="zelle_handle"
						type="text"
						placeholder="Phone or Email"
						class="w-full rounded-xl border border-slate-200 bg-white py-3 pr-4 pl-9 font-medium text-slate-900 shadow-sm transition-all focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20"
					/>
				</div>
			</div>
		</div>

		{#if !redirectTo?.includes('/board')}
			<div class="space-y-1">
				<label for="password" class="text-xs font-bold text-slate-400 uppercase">Password</label>
				<input name="password" id="password" type="password" required />
			</div>
			<!-- ... confirm password ... -->
		{:else}
			<!-- Hidden inputs for guests -->
			<input type="hidden" name="password" value="GUEST_NO_PASSWORD" />
			<input type="hidden" name="confirm_password" value="GUEST_NO_PASSWORD" />
			<input type="hidden" name="is_guest" value="true" />

			<div class="mb-4 rounded-xl border border-blue-100 bg-blue-50 p-4">
				<p class="text-xs font-medium text-blue-700">
					✨ Joining as a guest. You won't need a password to claim squares.
				</p>
			</div>
		{/if}

		{#if form?.message}
			<p class="rounded-lg bg-red-50 p-2 text-center text-xs font-bold text-red-500">
				{form.message}
			</p>
		{/if}

		<button
			type="submit"
			disabled={loading}
			class="mt-2 rounded-xl bg-blue-600 p-3.5 font-bold text-white transition-all hover:bg-blue-700 disabled:opacity-50"
		>
			{loading ? 'Creating Account...' : 'Register'}
		</button>
	</form>

	<p class="mt-6 text-center text-sm text-slate-500">
		Already have an account?
		<a
			href={resolve(redirectTo ? `/login?redirectTo=${encodeURIComponent(redirectTo)}` : '/login')}
			class="font-bold text-blue-600 hover:underline"
		>
			Sign In
		</a>
	</p>
</div>
