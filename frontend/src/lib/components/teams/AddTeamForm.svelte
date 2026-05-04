<script lang="ts">
	import { enhance } from '$app/forms';
	import type { ActionData } from '$lib/types';

	const DEFAULT_PRIMARY = '#0f172a';
	const DEFAULT_SECONDARY = '#64748b';

	let { onResult }: { onResult: (data: ActionData) => void } = $props();
	let loading = $state(false);

	// State for live preview
	let city = $state('');
	let mascot = $state('');
	let logoUrl = $state(''); // New state for logo preview
	let primaryColor = $state('#0f172a');
	let secondaryColor = $state('#64748b');
</script>

<form
	method="POST"
	action="?/create"
	use:enhance={() => {
		loading = true;
		return async ({ result, update }) => {
			loading = false;
			if (result.type === 'success' || result.type === 'failure') {
				const actionData = result.data as unknown as ActionData;
				if (actionData) onResult(actionData);
			}
			// Clear the logo state specifically on success
			if (result.type === 'success') {
				city = '';
				mascot = '';
				logoUrl = '';
				primaryColor = DEFAULT_PRIMARY; // Reset to default
				secondaryColor = DEFAULT_SECONDARY; // Reset to default
				const sportSelect = document.querySelector('select[name="sport"]') as HTMLSelectElement;
				if (sportSelect) sportSelect.selectedIndex = 0;
			}
			await update();
		};
	}}
	class="sticky top-6 space-y-6 rounded-4xl border border-slate-200 bg-white p-8 shadow-sm"
>
	<h2 class="text-xl font-bold text-slate-800">Add New Team</h2>
	<div class="space-y-4">
		<div class="grid grid-cols-2 gap-4">
			<div>
				<label
					for="city"
					class="mb-2 block text-xs font-bold tracking-widest text-slate-400 uppercase">City</label
				>
				<input
					name="city"
					bind:value={city}
					placeholder="Seattle"
					class="w-full rounded-xl border border-slate-100 bg-slate-50 px-4 py-3 outline-none focus:bg-white focus:ring-2 focus:ring-blue-500"
					required
				/>
			</div>
			<div>
				<label
					for="mascot"
					class="mb-2 block text-xs font-bold tracking-widest text-slate-400 uppercase"
					>Mascot</label
				>
				<input
					name="mascot"
					bind:value={mascot}
					placeholder="Seahawks"
					class="w-full rounded-xl border border-slate-100 bg-slate-50 px-4 py-3 outline-none focus:bg-white focus:ring-2 focus:ring-blue-500"
					required
				/>
			</div>
		</div>

		<div>
			<label
				for="team_logo_url"
				class="mb-2 block text-xs font-bold tracking-widest text-slate-400 uppercase"
				>Logo URL</label
			>
			<input
				name="team_logo_url"
				bind:value={logoUrl}
				placeholder="https://example.com/logo.png"
				class="w-full rounded-xl border border-slate-100 bg-slate-50 px-4 py-3 outline-none focus:bg-white focus:ring-2 focus:ring-blue-500"
			/>
		</div>

		<div class="grid grid-cols-2 gap-4">
			<div class="space-y-2">
				<label
					for="primary_color"
					class="block text-xs font-bold tracking-widest text-slate-400 uppercase">Primary</label
				>
				<div class="flex ...">
					<input
						type="color"
						id="primary_color"
						name="primary_color"
						bind:value={primaryColor}
						class="h-8 w-10 cursor-pointer border-none bg-transparent"
					/>
					<span class="font-mono text-[10px] text-slate-500 uppercase">{primaryColor}</span>
				</div>
			</div>
			<div class="space-y-2">
				<label
					for="secondary_color"
					class="block text-xs font-bold tracking-widest text-slate-400 uppercase">Secondary</label
				>
				<div class="flex ...">
					<input
						type="color"
						id="secondary_color"
						name="secondary_color"
						bind:value={secondaryColor}
						class="h-8 w-10 cursor-pointer border-none bg-transparent"
					/>
					<span class="font-mono text-[10px] text-slate-500 uppercase">{secondaryColor}</span>
				</div>
			</div>
		</div>

		<div
			class="flex flex-col items-center justify-center rounded-2xl border-2 border-dashed border-slate-100 bg-slate-50/30 py-6"
		>
			<div class="flex items-center gap-3">
				{#if logoUrl}
					<img src={logoUrl} alt="" class="h-10 w-10 object-contain" />
				{/if}
				<div
					class="text-xl font-black tracking-tighter uppercase transition-all"
					style="color: {primaryColor}; border-bottom: 4px solid {secondaryColor}"
				>
					{city || 'Preview'}
					{mascot || ''}
				</div>
			</div>
		</div>

		<div>
			<label
				for="sport"
				class="mb-2 block px-1 text-xs font-bold tracking-widest text-slate-400 uppercase"
				>Sport</label
			>
			<select
				name="sport"
				required
				class="w-full rounded-xl border border-slate-100 bg-slate-50 px-4 py-3 outline-none focus:bg-white"
			>
				<option value="" disabled selected>Select a sport...</option>
				{#each ['football', 'basketball', 'baseball', 'hockey', 'soccer'] as s (s)}
					<option value={s}>{s.toUpperCase()}</option>
				{/each}
			</select>
		</div>
	</div>

	<button
		disabled={loading}
		class="w-full rounded-xl bg-slate-900 py-4 font-bold text-white transition-all hover:bg-blue-600 active:scale-95 disabled:opacity-50"
	>
		{loading ? 'Saving...' : 'Save Team'}
	</button>
</form>
