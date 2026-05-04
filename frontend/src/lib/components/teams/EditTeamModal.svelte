<script lang="ts">
	import { enhance } from '$app/forms';
	import type { Team, ActionData } from '$lib/types';

	let {
		team,
		onResult,
		onclose
	}: { team: Team; onResult: (data: ActionData) => void; onclose: () => void } = $props();

	let loading = $state(false);

	// Reactive state for the form fields
	let primary = $state('#0f172a');
	let secondary = $state('#64748b');
	let logoUrl = $state('');

	// Sync state whenever the "team" prop changes
	$effect(() => {
		primary = team.primary_color || '#0f172a';
		secondary = team.secondary_color || '#64748b';
		logoUrl = team.team_logo_url || '';
	});
</script>

<div class="fixed inset-0 z-100 flex items-center justify-center p-4">
	<button
		aria-label="Edit Modal"
		class="absolute inset-0 cursor-default bg-slate-900/60 backdrop-blur-sm"
		onclick={onclose}
	></button>

	<div class="animate-in zoom-in-95 relative w-full max-w-lg rounded-3xl bg-white p-8 shadow-2xl">
		<h2 class="mb-6 text-2xl font-black text-slate-900">Edit {team.full_name}</h2>

		<form
			method="POST"
			action="?/update"
			use:enhance={() => {
				loading = true;
				return async ({ result, update }) => {
					loading = false;
					if (result.type === 'success' || result.type === 'failure') {
						const actionData = result.data as unknown as ActionData;
						if (actionData) onResult(actionData);
					}
					if (result.type === 'success') onclose();
					await update();
				};
			}}
			class="space-y-6"
		>
			<input type="hidden" name="id" value={team.id} />

			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-1">
					<label
						for="city"
						class="px-1 text-[10px] font-bold tracking-widest text-slate-400 uppercase">City</label
					>
					<input
						name="city"
						id="city"
						value={team.city}
						class="w-full rounded-xl border border-slate-200 px-4 py-3 outline-none focus:ring-2 focus:ring-blue-500"
						required
					/>
				</div>
				<div class="space-y-1">
					<label
						for="mascot"
						class="px-1 text-[10px] font-bold tracking-widest text-slate-400 uppercase"
						>Mascot</label
					>
					<input
						name="mascot"
						id="mascot"
						value={team.mascot}
						class="w-full rounded-xl border border-slate-200 px-4 py-3 outline-none focus:ring-2 focus:ring-blue-500"
						required
					/>
				</div>
			</div>

			<div class="space-y-1">
				<label
					for="team_logo_url"
					class="px-1 text-[10px] font-bold tracking-widest text-slate-400 uppercase"
					>Logo URL</label
				>
				<input
					name="team_logo_url"
					id="team_logo_url"
					bind:value={logoUrl}
					placeholder="https://example.com/logo.png"
					class="w-full rounded-xl border border-slate-200 px-4 py-3 outline-none focus:ring-2 focus:ring-blue-500"
				/>
			</div>

			<div class="grid grid-cols-2 gap-4">
				<div class="flex items-center gap-3 rounded-xl border border-slate-100 bg-slate-50 p-2">
					<input
						type="color"
						name="primary_color"
						bind:value={primary}
						class="h-10 w-10 cursor-pointer border-none bg-transparent"
					/>
					<span class="font-mono text-xs text-slate-600 uppercase">{primary}</span>
				</div>
				<div class="flex items-center gap-3 rounded-xl border border-slate-100 bg-slate-50 p-2">
					<input
						type="color"
						name="secondary_color"
						bind:value={secondary}
						class="h-10 w-10 cursor-pointer border-none bg-transparent"
					/>
					<span class="font-mono text-xs text-slate-600 uppercase">{secondary}</span>
				</div>
			</div>

			<div
				class="flex flex-col items-center justify-center rounded-2xl border-2 border-dashed border-slate-100 bg-slate-50/50 py-8"
			>
				<div class="flex items-center gap-4">
					{#if logoUrl}
						<img src={logoUrl} alt="" class="h-12 w-12 object-contain" />
					{/if}
					<div
						class="text-2xl font-black tracking-tighter uppercase transition-all"
						style="color: {primary}; border-bottom: 4px solid {secondary}"
					>
						{team.city}
						{team.mascot}
					</div>
				</div>
			</div>

			<div class="space-y-1">
				<label
					for="sport"
					class="px-1 text-[10px] font-bold tracking-widest text-slate-400 uppercase"
				>
					Sport
				</label>
				<select
					name="sport"
					value={team.sport}
					required
					class="w-full rounded-xl border border-slate-200 bg-white px-4 py-3 outline-none focus:ring-2 focus:ring-blue-500"
				>
					{#each ['football', 'basketball', 'baseball', 'hockey', 'soccer'] as s (s)}
						<option value={s}>{s.toUpperCase()}</option>
					{/each}
				</select>
			</div>

			<div class="flex gap-3 pt-2">
				<button
					type="button"
					onclick={onclose}
					class="flex-1 rounded-xl bg-slate-100 py-4 font-bold text-slate-600 transition-colors hover:bg-slate-200"
				>
					Cancel
				</button>
				<button
					disabled={loading}
					class="flex-2 rounded-xl bg-blue-600 py-4 font-bold text-white shadow-lg shadow-blue-100 transition-all hover:bg-blue-700 active:scale-95 disabled:opacity-50"
				>
					{loading ? 'Updating...' : 'Save Changes'}
				</button>
			</div>
		</form>
	</div>
</div>

<!-- <script lang="ts">
	import { enhance } from '$app/forms';
	import type { Team, ActionData } from '$lib/types';

	let {
		team,
		onResult,
		onclose
	}: { team: Team; onResult: (data: ActionData) => void; onclose: () => void } = $props();
	let loading = $state(false);

	// Initialize with a default
	let primary = $state('#0f172a');
	let secondary = $state('#64748b');

	// Sync state whenever the "team" prop changes
	$effect(() => {
		primary = team.primary_color || '#0f172a';
		secondary = team.secondary_color || '#64748b';
	});
</script>

<div class="fixed inset-0 z-100 flex items-center justify-center p-4">
	<button
		aria-label="Edit Modal"
		class="absolute inset-0 cursor-default bg-slate-900/60 backdrop-blur-sm"
		onclick={onclose}
	></button>

	<div class="animate-in zoom-in-95 relative w-full max-w-lg rounded-3xl bg-white p-8 shadow-2xl">
		<h2 class="mb-6 text-2xl font-black text-slate-900">Edit {team.full_name}</h2>

		<form
			method="POST"
			action="?/update"
			use:enhance={() => {
				loading = true;
				return async ({ result, update }) => {
					loading = false;
					if (result.type === 'success' || result.type === 'failure') {
						const actionData = result.data as unknown as ActionData;
						if (actionData) onResult(actionData);
					}
					if (result.type === 'success') onclose();
					await update();
				};
			}}
			class="space-y-6"
		>
			<input type="hidden" name="id" value={team.id} />

			<div class="grid grid-cols-2 gap-4">
				<input
					name="city"
					value={team.city}
					class="rounded-xl border border-slate-200 px-4 py-3"
					required
				/>
				<input
					name="mascot"
					value={team.mascot}
					class="rounded-xl border border-slate-200 px-4 py-3"
					required
				/>
			</div>

			<div class="grid grid-cols-2 gap-4">
				<div class="flex items-center gap-3 rounded-xl border border-slate-100 bg-slate-50 p-2">
					<input
						type="color"
						name="primary_color"
						bind:value={primary}
						class="h-10 w-10 cursor-pointer border-none bg-transparent"
					/>
					<span class="font-mono text-xs uppercase">{primary}</span>
				</div>
				<div class="flex items-center gap-3 rounded-xl border border-slate-100 bg-slate-50 p-2">
					<input
						type="color"
						name="secondary_color"
						bind:value={secondary}
						class="h-10 w-10 cursor-pointer border-none bg-transparent"
					/>
					<span class="font-mono text-xs uppercase">{secondary}</span>
				</div>
			</div>

			<div class="flex gap-3 pt-4">
				<div
					class="flex flex-col items-center justify-center rounded-2xl border-2 border-dashed border-slate-100 py-6"
				>
					<div
						class="text-xl font-black tracking-tighter uppercase transition-all"
						style="color: {primary}; border-bottom: 4px solid {secondary}"
					>
						{team.city}
						{team.mascot}
					</div>
				</div>
				<button
					type="button"
					onclick={onclose}
					class="flex-1 rounded-xl bg-slate-100 py-4 font-bold">Cancel</button
				>
				<button
					disabled={loading}
					class="flex-2 rounded-xl bg-blue-600 py-4 font-bold text-white transition-all"
				>
					{loading ? 'Updating...' : 'Save Changes'}
				</button>
			</div>
		</form>
	</div>
</div> -->
