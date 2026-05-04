<script lang="ts">
	// import { enhance } from '$app/forms';
	// import { goto } from '$app/navigation';
	import type { PageData } from './$types';
	import { resolve } from '$app/paths';
	// import { applyAction, enhance } from '$app/forms';
	import { goto } from '$app/navigation';
	import { enhance } from '$app/forms';

	let { data }: { data: PageData } = $props();
	let loading = $state(false);

	let homeTeamId = $state('');
	let awayTeamId = $state('');
	let pricePerSquare = $state(0);

	// Explicitly track all 5 payout states
	let q1 = $state(0);
	let q2 = $state(0);
	let q3 = $state(0);
	let q4 = $state(0);
	let final = $state(0);

	let isSameTeam = $derived(homeTeamId && awayTeamId ? homeTeamId === awayTeamId : false);

	// Total Pot for 100 squares
	let totalPot = $derived(pricePerSquare * 100);

	// Sum of all FIVE fields
	let totalPayout = $derived(q1 + q2 + q3 + q4 + final);

	// Balanced check (using a small epsilon for float math safety)
	let isPotBalanced = $derived(Math.abs(totalPot - totalPayout) < 0.01 && totalPot > 0);

	function splitPotEvenly() {
		const share = totalPot / 5;
		q1 = q2 = q3 = q4 = final = share;
	}
</script>

<div class="mx-auto max-w-2xl p-6">
	<header class="mb-8">
		<a href={resolve('/dashboard/boards')} class="text-sm font-bold text-blue-600 hover:underline">
			← Back to Boards
		</a>
		<h1 class="mt-2 text-3xl font-black tracking-tight text-slate-900 uppercase">
			New {data.activeSport} Matchup
		</h1>
	</header>

	{#if !data.teams || data.teams.length === 0}
		<div
			class="flex flex-col items-center justify-center rounded-4xl border-2 border-dashed border-slate-200 bg-white px-8 py-16 text-center shadow-sm"
		>
			<div
				class="mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-orange-50 text-orange-500"
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
						d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"
					/>
				</svg>
			</div>
			<h2 class="text-xl font-black tracking-tight text-slate-900 uppercase">No Teams Available</h2>
			<p class="mt-2 max-w-xs text-sm text-slate-500">
				You need to create at least two teams for the <strong>{data.activeSport}</strong> category before
				you can start a board.
			</p>
			<div class="mt-8 flex gap-3">
				<a
					href={resolve('/dashboard/teams')}
					class="rounded-xl bg-slate-900 px-6 py-3 text-sm font-bold text-white transition-all hover:bg-blue-600"
				>
					Go to Team Manager
				</a>
			</div>
		</div>
	{:else}
		<form
			method="POST"
			action="?/create"
			use:enhance={() => {
				loading = true;
				return async ({ result }) => {
					console.log(`RESULT TYPE: ${result.type}`);
					console.log(`RESULT STATUS: ${result.status}`);
					// console.log(`DATA: ${result.data?.boardId}`);
					// result.data contains the { success, boardId } we returned above
					if (result.type === 'success' && result.data?.boardId) {
						const path = resolve(`/dashboard/boards/${result.data.boardId}`);
						await goto(path, { invalidateAll: true });
					} else {
						// If it's a redirect (thrown by SvelteKit) or failure
						loading = false;
					}
				};
			}}
			class="space-y-6 rounded-3xl border border-slate-200 bg-white p-8 shadow-sm"
		>
			<input type="hidden" name="sport" value={data.activeSport} />

			<div class="space-y-4">
				<div>
					<label
						for="name"
						class="mb-2 block text-xs font-bold tracking-widest text-slate-400 uppercase"
					>
						Board Name / Event
					</label>
					<input
						name="name"
						placeholder="e.g. Super Bowl LXI"
						class="w-full rounded-xl border border-slate-100 bg-slate-50 px-4 py-3 outline-none focus:bg-white focus:ring-2 focus:ring-blue-500"
						required
					/>
				</div>

				<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
					<div>
						<label
							for="home_team_id"
							class="mb-2 block text-xs font-bold tracking-widest text-slate-400 uppercase"
						>
							Home Team
						</label>
						<select
							bind:value={homeTeamId}
							name="home_team_id"
							required
							class="w-full rounded-xl border border-slate-100 bg-slate-50 px-4 py-3 outline-none focus:bg-white focus:ring-2 focus:ring-blue-500"
						>
							<option value="" disabled selected>Select Home...</option>
							{#each data.teams as team (team.id)}
								<option value={team.id}>{team.full_name}</option>
							{/each}
						</select>
					</div>

					<div>
						<label
							for="away_team_id"
							class="mb-2 block text-xs font-bold tracking-widest text-slate-400 uppercase"
						>
							Away Team
						</label>
						<select
							bind:value={awayTeamId}
							name="away_team_id"
							required
							class="w-full rounded-xl border border-slate-100 bg-slate-50 px-4 py-3 outline-none focus:bg-white focus:ring-2 focus:ring-blue-500"
						>
							<option value="" disabled selected>Select Away...</option>
							{#each data.teams as team (team.id)}
								<option value={team.id}>{team.full_name}</option>
							{/each}
						</select>
					</div>
				</div>
			</div>

			<div class="space-y-4 border-t border-slate-100 pt-6">
				<h2 class="text-xs font-bold tracking-widest text-slate-400 uppercase">Financials</h2>

				<div class="flex flex-col gap-2">
					<label for="price_per_square" class="block text-xs font-bold text-slate-500">
						Price per Square ($)
					</label>
					<input
						type="number"
						step="0.01"
						bind:value={pricePerSquare}
						name="price_per_square"
						class="w-full rounded-xl border border-slate-100 bg-slate-50 px-4 py-3 outline-none focus:bg-white focus:ring-2 focus:ring-blue-500"
						required
					/>
					<button
						type="button"
						onclick={splitPotEvenly}
						class="w-fit text-[10px] font-black text-blue-600 uppercase transition-colors hover:text-blue-800"
					>
						⚡ Auto-Split Pot 5 Ways (${totalPot / 5} each)
					</button>
				</div>

				<div class="grid grid-cols-2 gap-3 sm:grid-cols-5">
					<div>
						<label for="payout_q1" class="mb-1 block text-[10px] font-bold text-slate-400 uppercase"
							>1st Qtr</label
						>
						<input
							type="number"
							step="0.01"
							bind:value={q1}
							name="payout_q1"
							id="payout_q1"
							class="w-full rounded-xl border bg-slate-50 p-2 text-sm"
						/>
					</div>
					<div>
						<label for="payout_q2" class="mb-1 block text-[10px] font-bold text-slate-400 uppercase"
							>2nd Qtr</label
						>
						<input
							type="number"
							step="0.01"
							bind:value={q2}
							name="payout_q2"
							id="payout_q2"
							class="w-full rounded-xl border bg-slate-50 p-2 text-sm"
						/>
					</div>
					<div>
						<label for="payout_q3" class="mb-1 block text-[10px] font-bold text-slate-400 uppercase"
							>3rd Qtr</label
						>
						<input
							type="number"
							step="0.01"
							bind:value={q3}
							name="payout_q3"
							id="payout_q3"
							class="w-full rounded-xl border bg-slate-50 p-2 text-sm"
						/>
					</div>
					<div>
						<label for="payout_q4" class="mb-1 block text-[10px] font-bold text-slate-400 uppercase"
							>4th Qtr</label
						>
						<input
							type="number"
							step="0.01"
							bind:value={q4}
							name="payout_q4"
							id="payout_q4"
							class="w-full rounded-xl border bg-slate-50 p-2 text-sm"
						/>
					</div>
					<div>
						<label
							for="payout_final"
							class="mb-1 block text-[10px] font-bold text-blue-500 uppercase">Final Score</label
						>
						<input
							type="number"
							step="0.01"
							bind:value={final}
							name="payout_final"
							id="payout_final"
							class="w-full rounded-xl border border-blue-100 bg-blue-50 p-2 text-sm font-bold text-blue-700"
						/>
					</div>
				</div>

				{#if pricePerSquare > 0 && !isPotBalanced}
					<div
						class="flex items-center justify-between rounded-lg border border-orange-100 bg-orange-50 p-3"
					>
						<p class="text-xs font-bold text-orange-700">
							Mismatch: Payouts (${totalPayout}) ≠ Pot (${totalPot})
						</p>
					</div>
				{/if}
			</div>

			{#if isSameTeam}
				<p class="text-center text-sm font-bold text-red-500">
					Home and Away teams must be different!
				</p>
			{/if}

			<button
				disabled={loading || isSameTeam}
				class="w-full rounded-xl bg-slate-900 py-4 font-bold text-white transition-all hover:bg-blue-600 disabled:opacity-30"
			>
				{loading ? 'Creating...' : 'Create Board'}
			</button>
		</form>
	{/if}
</div>
