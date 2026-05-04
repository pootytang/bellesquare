<script lang="ts">
	import type { Payouts } from '$lib/types';

	let {
		payouts,
		quarterWinners, // Pass winningPlayersByQuarter here
		isCreator,
		onEdit
	}: {
		payouts: Payouts[];
		quarterWinners: (string | null)[]; // Properly typed array of names/nulls
		isCreator: boolean;
		onEdit?: () => void;
	} = $props();

	// Maps period_name to the index in liveQuarterScores
	const periodMap: Record<string, number> = {
		q1: 0,
		q2: 1,
		q3: 2,
		q4: 3,
		final: 4
	};

	const periodLabels: Record<string, string> = {
		q1: '1st Quarter',
		q2: '2nd Quarter',
		q3: '3rd Quarter',
		q4: '4th Quarter',
		final: 'Final Score'
	};
</script>

<div class="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
	<div class="flex items-center justify-between border-b border-slate-200 bg-slate-50 px-5 py-3">
		<h3 class="font-bold text-slate-800">Payout Schedule</h3>
		{#if isCreator}
			<button onclick={onEdit} class="text-xs font-semibold text-blue-600 hover:text-blue-700">
				Edit
			</button>
		{/if}
	</div>

	<div class="divide-y divide-slate-100">
		{#each payouts as payout, i (i)}
			{@const quarterIndex = periodMap[payout.period_name]}
			{@const winnerName = quarterWinners[quarterIndex]}

			<div class="flex items-center justify-between px-5 py-3 {winnerName ? 'bg-green-50/50' : ''}">
				<div>
					<p class="text-sm font-medium text-slate-600">
						{periodLabels[payout.period_name] ?? payout.period_name}
					</p>
					{#if winnerName}
						<div class="animate-in fade-in slide-in-from-left-2 mt-1 flex items-center gap-1">
							<span class="text-xs text-amber-500">🏆</span>
							<p class="text-xs font-bold text-green-700">
								Winner: {winnerName}
							</p>
						</div>
					{:else}
						<p class="mt-1 text-[10px] tracking-tighter text-slate-400 uppercase">TBD</p>
					{/if}
				</div>
				<span class="font-mono text-lg font-bold text-slate-900">
					${payout.amount.toLocaleString()}
				</span>
			</div>
		{/each}
	</div>

	<div class="flex items-center justify-between bg-blue-50 px-5 py-3">
		<span class="text-xs font-bold tracking-wider text-blue-700 uppercase">Total Prize Pool</span>
		<span class="font-mono text-lg font-black text-blue-900">
			${payouts.reduce((sum, p) => sum + p.amount, 0).toLocaleString()}
		</span>
	</div>
</div>

<!-- <script lang="ts">
	import type { Payouts } from '$lib/types';

	let {
		payouts,
		isCreator,
		onEdit
	}: {
		payouts: Payouts[];
		isCreator: boolean;
		onEdit?: () => void;
	} = $props();

	const periodLabels: Record<string, string> = {
		q1: '1st Quarter',
		q2: '2nd Quarter',
		q3: '3rd Quarter',
		q4: '4th Quarter',
		final: 'Final Score'
	};
</script>

<div class="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
	<div class="flex items-center justify-between border-b border-slate-200 bg-slate-50 px-5 py-3">
		<h3 class="font-bold text-slate-800">Payout Schedule</h3>
		{#if isCreator}
			<button onclick={onEdit} class="text-xs font-semibold text-blue-600 hover:text-blue-700">
				Edit
			</button>
		{/if}
	</div>

	<div class="divide-y divide-slate-100">
		{#each payouts as payout, i (i)}
			<div class="flex items-center justify-between px-5 py-3">
				<div>
					<p class="text-sm font-medium text-slate-600">
						{periodLabels[payout.period_name] ?? payout.period_name}
					</p>
					{#if payout.winner_name}
						<div class="mt-1 flex items-center gap-1">
							<span class="text-amber-500">🏆</span>
							<p class="text-xs font-bold text-green-600">
								Winner: {payout.winner_name}
							</p>
						</div>
					{:else}
						<p class="mt-1 text-[10px] tracking-tighter text-slate-400 uppercase">TBD</p>
					{/if}
				</div>
				<span class="font-mono text-lg font-bold text-slate-900">
					${payout.amount.toLocaleString()}
				</span>
			</div>
		{/each}
	</div>

	<div class="flex items-center justify-between bg-blue-50 px-5 py-3">
		<span class="text-xs font-bold tracking-wider text-blue-700 uppercase">Total Prize Pool</span>
		<span class="font-mono text-lg font-black text-blue-900">
			${payouts.reduce((sum, p) => sum + p.amount, 0).toLocaleString()}
		</span>
	</div>
</div> -->
