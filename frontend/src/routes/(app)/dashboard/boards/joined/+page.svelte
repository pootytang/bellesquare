<script lang="ts">
	import { resolve } from '$app/paths';
	let { data } = $props();
</script>

<div class="mx-auto max-w-4xl p-6">
	<header class="mb-8">
		<h1 class="text-3xl font-black text-slate-900">Joined Games</h1>
		<p class="text-slate-500">Boards where you've claimed squares.</p>
	</header>

	{#if data.boards.length === 0}
		<div class="rounded-2xl border-2 border-dashed border-slate-200 p-12 text-center">
			<p class="text-slate-400">You haven't joined any boards yet.</p>
			<p class="text-sm text-slate-400">Join a game using a link shared by a host!</p>
		</div>
	{:else}
		<div class="grid gap-4">
			{#each data.boards as board (board.id)}
				<a
					href={resolve(`/dashboard/boards/${board.id}`)}
					class="group flex items-center justify-between rounded-2xl border border-slate-200 bg-white p-4 transition-all hover:border-blue-500 hover:shadow-md"
				>
					<div class="flex items-center gap-4">
						<div class="flex -space-x-2">
							<img
								src={board.home_team.team_logo_url}
								alt=""
								class="h-10 w-10 rounded-full border-2 border-white bg-slate-50"
							/>
							<img
								src={board.away_team.team_logo_url}
								alt=""
								class="h-10 w-10 rounded-full border-2 border-white bg-slate-50"
							/>
						</div>
						<div>
							<h3 class="font-bold text-slate-900">{board.title}</h3>
							<p class="text-xs font-semibold text-slate-500 uppercase">
								{board.home_team.mascot} vs {board.away_team.mascot}
							</p>
						</div>
					</div>

					<div class="flex items-center gap-4">
						<span
							class="rounded-full bg-slate-100 px-3 py-1 text-xs font-bold text-slate-600 uppercase"
						>
							{board.status}
						</span>
						<span class="text-blue-600 transition-transform group-hover:translate-x-1">→</span>
					</div>
				</a>
			{/each}
		</div>
	{/if}
</div>
