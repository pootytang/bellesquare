<script lang="ts">
	import { PUBLIC_API_BASE_URL } from '$env/static/public';

	let { scores, homeTeam, awayTeam, isCreator, boardId, token } = $props();

	const LOCK = '🔒';
	const UNLOCK = '🔓';

	async function updateScore(i: number) {
		// Force Svelte to look at the latest state of the proxy
		const qData = scores.quarters[i];

		const payload = {
			quarter: i + 1,
			home_score: Number(qData.home_score),
			away_score: Number(qData.away_score),
			is_locked: qData.is_locked // This is now updated by toggleLock before calling here
		};

		// console.log('Sending Score Update:', payload);

		const response = await fetch(`${PUBLIC_API_BASE_URL}/api/v1/boards/${boardId}/scores`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify(payload)
		});

		if (response.ok) console.log('Score saved!');
	}

	function toggleLock(i: number) {
		// 1. Update the state object property directly
		scores.quarters[i].is_locked = !scores.quarters[i].is_locked;

		// 2. Local propagation logic
		if (scores.quarters[i].is_locked && i < 4 && !scores.quarters[4].is_locked) {
			scores.quarters[4].home_score = scores.quarters[i].home_score;
			scores.quarters[4].away_score = scores.quarters[i].away_score;
		}

		// 3. Save the new lock state to the DB
		updateScore(i);
	}
</script>

<div class="overflow-hidden rounded-xl border border-gray-700 bg-gray-900 shadow-xl">
	<table class="w-full table-auto text-left text-sm text-gray-300">
		<thead class="bg-gray-800/50 text-xs font-bold tracking-wider text-gray-400 uppercase">
			<tr>
				<th class="px-4 py-3">Period</th>
				<th class="px-2 py-3 text-center">{homeTeam.mascot}</th>
				<th class="px-2 py-3 text-center">{awayTeam.mascot}</th>
				{#if isCreator}
					<th class="w-20 px-4 py-3 text-center">Status</th>
				{/if}
			</tr>
		</thead>
		<tbody class="divide-y divide-gray-800">
			{#each scores.quarters as qData, i (qData.quarter)}
				<tr
					class="transition-colors hover:bg-gray-800/30
                    {qData.is_locked ? 'bg-gray-950/40 text-gray-500' : ''}
                    {i === 4 ? 'border-t border-blue-500/30 bg-blue-900/20' : ''}"
				>
					<td
						class="px-4 py-4 font-semibold whitespace-nowrap italic {i === 4
							? 'text-blue-400'
							: ''}"
					>
						{i === 4 ? 'Final' : `Q${i + 1}`}
					</td>

					<!-- Home Score Input -->
					<td class="px-2 py-4 text-center">
						{#if isCreator}
							<input
								type="number"
								value={qData.home_score}
								disabled={qData.is_locked}
								oninput={(e) => {
									// Update the state immediately as you type
									const val = Number(e.currentTarget.value);
									scores.quarters[i].home_score = val;

									// Auto-update final row locally
									if (i < 4 && !scores.quarters[4].is_locked) {
										scores.quarters[4].home_score = val;
									}
								}}
								onblur={() => updateScore(i)}
								class="w-14 rounded-md border-none bg-gray-800 px-1 py-1.5 text-center text-white ring-1 ring-gray-700 focus:ring-2 focus:ring-blue-500 disabled:opacity-30 {i ===
								4
									? 'ring-blue-500/50'
									: ''}"
							/>
						{:else}
							<span class="text-lg font-bold {i === 4 ? 'text-white' : ''}">{qData.home_score}</span
							>
						{/if}
					</td>

					<!-- Away Score Input -->
					<td class="px-2 py-4 text-center">
						{#if isCreator}
							<input
								type="number"
								value={qData.away_score}
								disabled={qData.is_locked}
								oninput={(e) => {
									// Update the state immediately as you type
									const val = Number(e.currentTarget.value);
									scores.quarters[i].away_score = val;

									// Auto-update final row locally
									if (i < 4 && !scores.quarters[4].is_locked) {
										scores.quarters[4].away_score = val;
									}
								}}
								onblur={() => updateScore(i)}
								class="w-14 rounded-md border-none bg-gray-800 px-1 py-1.5 text-center text-white ring-1 ring-gray-700 focus:ring-2 focus:ring-blue-500 disabled:opacity-30 {i ===
								4
									? 'ring-blue-500/50'
									: ''}"
							/>
						{:else}
							<span class="text-lg font-bold {i === 4 ? 'text-white' : ''}">{qData.away_score}</span
							>
						{/if}
					</td>

					{#if isCreator}
						<td class="px-4 py-4 text-center">
							<button
								type="button"
								onclick={() => toggleLock(i)}
								class="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg border border-gray-700 bg-gray-800 text-lg transition-all hover:border-gray-500 hover:bg-gray-700 active:scale-90 {i ===
								4
									? 'border-blue-500/50'
									: ''}"
							>
								<span class="block -translate-y-px">
									{qData.is_locked ? LOCK : UNLOCK}
								</span>
							</button>
						</td>
					{/if}
				</tr>
			{/each}
		</tbody>
	</table>
</div>
