<script lang="ts">
	let { boardId, homeTeam, awayTeam } = $props();
	let homeScore = $state(0);
	let awayScore = $state(0);

	async function updateScore() {
		// This should hit your Go backend and broadcast via WebSocket
		await fetch(`/api/v1/boards/${boardId}/scores`, {
			method: 'POST',
			body: JSON.stringify({ homeScore, awayScore })
		});
	}
</script>

<div class="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm">
	<h3 class="mb-4 text-center text-xs font-bold tracking-widest text-slate-400 uppercase">
		Live Score Update
	</h3>
	<div class="flex items-center justify-around gap-4">
		<div class="text-center">
			<p class="mb-1 text-[10px] font-bold uppercase">{homeTeam}</p>
			<input
				type="number"
				bind:value={homeScore}
				class="w-16 rounded-lg border-slate-200 text-center text-xl font-black"
			/>
		</div>
		<div class="text-xl font-black text-slate-300">:</div>
		<div class="text-center">
			<p class="mb-1 text-[10px] font-bold uppercase">{awayTeam}</p>
			<input
				type="number"
				bind:value={awayScore}
				class="w-16 rounded-lg border-slate-200 text-center text-xl font-black"
			/>
		</div>
	</div>
	<button
		onclick={updateScore}
		class="mt-4 w-full rounded-lg bg-slate-900 py-2 text-xs font-bold text-white transition-colors hover:bg-slate-800"
	>
		Update Score
	</button>
</div>
