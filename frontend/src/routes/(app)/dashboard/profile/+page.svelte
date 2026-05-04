<script lang="ts">
	import { enhance } from '$app/forms';
	import { auth } from '$lib/auth.svelte';

	let { data, form } = $props();

	// Initialize local form state
	let userName = $state('');
	let firstName = $state('');
	let lastName = $state('');
	let venmo = $state('');
	let zelle = $state('');
	let preferredColor = $state('#2563eb');
	let saving = $state(false);

	const presets = ['#2563eb', '#dc2626', '#16a34a', '#ca8a04', '#7c3aed', '#db2777'];

	// Sync local state whenever 'data' changes
	$effect(() => {
		userName = data.user?.user_name ?? '';
		firstName = data.user?.first_name ?? '';
		lastName = data.user?.last_name ?? '';
		venmo = data.user?.venmo_handle ?? '';
		zelle = data.user?.zelle_handle ?? '';
		preferredColor = data.user?.preferred_color ?? '#2563eb';
	});

	let isUserInitialized = $derived(auth.initialized);

	// Track if anything has changed to enable/disable save button
	let hasChanges = $derived(
		userName !== (data.user?.user_name ?? '') ||
			firstName !== (data.user?.first_name ?? '') ||
			lastName !== (data.user?.last_name ?? '') ||
			venmo !== (data.user?.venmo_handle ?? '') ||
			zelle !== (data.user?.zelle_handle ?? '') ||
			preferredColor !== (data.user?.preferred_color ?? '#2563eb')
	);
</script>

<div class="mx-auto mt-10 max-w-xl rounded-2xl border border-slate-200 bg-white p-8 shadow-sm">
	{#if isUserInitialized}
		<header class="mb-8">
			<h1 class="text-2xl font-bold text-slate-900">Profile Settings</h1>
			<p class="text-slate-500">Manage your Bellesquare identity.</p>
		</header>

		<form
			method="POST"
			action="?/update"
			use:enhance={() => {
				saving = true;
				return async ({ update }) => {
					await update();
					saving = false;
				};
			}}
		>
			<div class="space-y-6">
				<div class="flex flex-col gap-1.5">
					<label for="username" class="text-sm font-semibold text-slate-700">Username</label>
					<input
						id="username"
						name="username"
						bind:value={userName}
						disabled
						class="w-full cursor-not-allowed rounded-lg border border-slate-200 bg-slate-50 px-4 py-2 text-slate-500 outline-none"
					/>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div class="flex flex-col gap-1.5">
						<label for="first_name" class="text-sm font-semibold text-slate-700">First Name</label>
						<input
							id="first_name"
							name="first_name"
							bind:value={firstName}
							class="rounded-lg border border-slate-300 px-4 py-2 transition outline-none focus:ring-2 focus:ring-blue-500"
						/>
					</div>
					<div class="flex flex-col gap-1.5">
						<label for="last_name" class="text-sm font-semibold text-slate-700">Last Name</label>
						<input
							id="last_name"
							name="last_name"
							bind:value={lastName}
							class="rounded-lg border border-slate-300 px-4 py-2 transition outline-none focus:ring-2 focus:ring-blue-500"
						/>
					</div>
				</div>

				<div class="flex flex-col gap-2">
					<label for="preferred_square_color_button" class="text-sm font-semibold text-slate-700"
						>Preferred Square Color</label
					>
					<div class="flex items-start gap-4 rounded-xl border border-slate-100 bg-slate-50/50 p-4">
						<div class="space-y-2">
							<span class="block text-[9px] font-black tracking-tighter text-slate-400 uppercase"
								>Presets</span
							>
							<div class="flex flex-wrap gap-1.5">
								{#each presets as color (color)}
									<button
										type="button"
										aria-label="Preferred Color"
										onclick={() => (preferredColor = color)}
										class="h-6 w-6 rounded-full border-2 shadow-sm transition-all hover:scale-110 active:scale-95"
										style="background-color: {color}; border-color: {preferredColor === color
											? 'white'
											: 'transparent'}; outline: {preferredColor === color
											? '2px solid #3b82f6'
											: 'none'}"
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
									bind:value={preferredColor}
									class="absolute -inset-2 h-12 w-12 cursor-pointer border-none bg-transparent"
								/>
							</div>
						</div>
					</div>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div class="flex flex-col gap-1.5">
						<label for="venmo_handle" class="text-sm font-semibold text-slate-700"
							>Venmo Handle</label
						>
						<input
							id="venmo_handle"
							name="venmo_handle"
							bind:value={venmo}
							placeholder="@username"
							class="rounded-lg border border-slate-300 px-4 py-2 transition outline-none focus:ring-2 focus:ring-blue-500"
						/>
					</div>
					<div class="flex flex-col gap-1.5">
						<label for="zelle_handle" class="text-sm font-semibold text-slate-700"
							>Zelle Handle</label
						>
						<input
							id="zelle_handle"
							name="zelle_handle"
							bind:value={zelle}
							placeholder="Phone or Email"
							class="rounded-lg border border-slate-300 px-4 py-2 transition outline-none focus:ring-2 focus:ring-blue-500"
						/>
					</div>
				</div>

				<button
					type="submit"
					disabled={!hasChanges || saving}
					class="w-full rounded-xl bg-slate-900 px-4 py-4 font-bold text-white shadow-lg transition-all hover:bg-blue-600 active:scale-[0.98] disabled:bg-slate-400 disabled:opacity-30"
				>
					{saving ? 'Saving...' : hasChanges ? 'Save Changes' : 'No Changes'}
				</button>
			</div>
		</form>

		{#if form?.success}
			<div
				class="mt-6 rounded-lg border border-green-200 bg-green-50 p-3 text-center text-sm font-bold text-green-700"
			>
				Profile updated successfully! 🎉
			</div>
		{/if}
	{:else}
		<div class="flex flex-col items-center justify-center py-10">
			<div class="h-8 w-8 animate-spin rounded-full border-b-2 border-blue-600"></div>
			<p class="mt-4 text-center text-slate-500">Loading profile...</p>
		</div>
	{/if}
</div>
