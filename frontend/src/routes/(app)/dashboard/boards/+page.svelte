<script lang="ts">
	import type { PageData } from './$types';
	import type { BoardSummary } from '$lib/types';
	import { resolve } from '$app/paths';

	let { data }: { data: PageData } = $props();

	const hasEnoughTeams = $derived(
		Object.values(data.teamCounts).some((count) => (count as number) >= 2)
	);
</script>

<div class="mx-auto max-w-6xl p-6">
	<header class="mb-10 flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-black tracking-tight text-slate-900">Your Boards</h1>
			<p class="text-slate-500">Manage your active games.</p>
		</div>

		{#if hasEnoughTeams}
			<a
				href={resolve('/dashboard/boards/new?sport=football')}
				class="rounded-xl bg-blue-600 px-6 py-3 font-bold text-white shadow-lg shadow-blue-100 transition-all hover:bg-blue-700"
			>
				+ Create New Board
			</a>
		{:else}
			<div class="group relative">
				<button
					disabled
					class="cursor-not-allowed rounded-xl border border-slate-200 bg-slate-100 px-6 py-3 font-bold text-slate-400"
				>
					Not Enough Teams
				</button>
				<div
					class="absolute top-full right-0 z-10 mt-2 hidden w-64 rounded-xl bg-slate-900 p-3 text-center text-xs text-white shadow-xl group-hover:block"
				>
					<p class="mb-1 font-bold">Teams Required</p>
					<p class="opacity-70">You need at least 2 teams in a single sport to create a board.</p>
				</div>
			</div>
		{/if}
	</header>

	<div class="space-y-12">
		<!-- SECTION: MANAGED BOARDS -->
		<section>
			<div class="mb-6 flex items-center gap-3">
				<h2 class="text-xs font-black tracking-widest text-slate-400 uppercase">Managed Boards</h2>
				<div class="h-px flex-1 bg-slate-100"></div>
			</div>

			<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
				{#each data.boards as board (board.id)}
					{@render BoardCard(board, true)}
				{:else}
					<div
						class="col-span-full rounded-3xl border-2 border-dashed border-slate-100 py-12 text-center text-slate-400"
					>
						You haven't created any boards yet.
					</div>
				{/each}
			</div>
		</section>

		<!-- SECTION: JOINED BOARDS -->
		<section>
			<div class="mb-6 flex items-center gap-3">
				<h2 class="text-xs font-black tracking-widest text-slate-400 uppercase">Joined Boards</h2>
				<div class="h-px flex-1 bg-slate-100"></div>
			</div>

			<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
				{#each data.joined as board (board.id)}
					{@render BoardCard(board, false)}
				{:else}
					<div
						class="col-span-full rounded-3xl border-2 border-dashed border-slate-100 py-12 text-center text-slate-400"
					>
						You haven't joined any boards yet.
					</div>
				{/each}
			</div>
		</section>
	</div>
</div>

{#snippet BoardCard(board: BoardSummary, isOwner: boolean)}
	<a
		href={resolve(isOwner ? `/dashboard/boards/${board.id}` : `/dashboard/boards/${board.id}`)}
		class="group relative overflow-hidden rounded-3xl border border-slate-200 bg-white p-6 shadow-sm transition-all hover:border-blue-200 hover:shadow-xl"
	>
		<!-- The Dynamic Gradient Bar -->
		<div
			class="absolute top-0 left-0 h-1.5 w-full"
			style="background: linear-gradient(90deg, {board.home_team.secondary_color}, {board.away_team
				.secondary_color})"
		></div>

		<div class="mb-4 flex items-center justify-between">
			<span
				class="rounded-full px-3 py-1 text-[10px] font-black tracking-widest uppercase {isOwner
					? 'bg-blue-50 text-blue-600'
					: 'bg-slate-100 text-slate-500'}"
			>
				{isOwner ? 'Admin' : board.status}
			</span>
			<span class="text-xs font-bold text-slate-400">${board.price_per_square}/sq</span>
		</div>

		<h2 class="mb-6 text-xl font-black text-slate-900 group-hover:text-blue-600">
			{board.title}
		</h2>

		<div class="mt-4 flex items-center gap-3">
			<div class="flex -space-x-3">
				<img
					src={board.home_team.team_logo_url}
					alt={board.home_team.mascot}
					class="h-10 w-10 rounded-full border-2 border-white bg-white object-contain shadow-sm"
				/>
				<img
					src={board.away_team.team_logo_url}
					alt={board.away_team.mascot}
					class="h-10 w-10 rounded-full border-2 border-white bg-white object-contain shadow-sm"
				/>
			</div>
			<div class="flex flex-col">
				<span class="text-[10px] font-black tracking-tighter text-slate-400 uppercase">Matchup</span
				>
				<span class="text-xs font-bold text-slate-600">
					{board.home_team.mascot} vs {board.away_team.mascot}
				</span>
			</div>
		</div>

		<div class="mt-4 flex items-center justify-between border-t border-slate-50 pt-4">
			<span class="text-xs font-bold text-slate-500">
				{isOwner ? 'Manage Settings' : 'View Your Squares'} →
			</span>
		</div>
	</a>
{/snippet}

<!-- <script lang="ts">
	import type { PageData } from './$types';
	import { resolve } from '$app/paths';

	let { data }: { data: PageData } = $props();

	// Check if at least one sport has 2+ teams
	const hasEnoughTeams = $derived(
		data.teamCounts.football >= 2 ||
			data.teamCounts.basketball >= 2 ||
			data.teamCounts.baseball >= 2 ||
			data.teamCounts.soccer >= 2 ||
			data.teamCounts.hockey >= 2
	);
</script>

<div class="mx-auto max-w-6xl p-6">
	<header class="mb-10 flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-black tracking-tight text-slate-900">Your Boards</h1>
			<p class="text-slate-500">Manage your active games.</p>
		</div>

		{#if hasEnoughTeams}
			<a
				href={resolve('/dashboard/boards/new?sport=football')}
				class="rounded-xl bg-blue-600 px-6 py-3 font-bold text-white shadow-lg shadow-blue-100 transition-all hover:bg-blue-700"
			>
				+ Create New Board
			</a>
		{:else}
			<div class="group relative">
				<button
					disabled
					class="cursor-not-allowed rounded-xl border border-slate-200 bg-slate-100 px-6 py-3 font-bold text-slate-400"
				>
					Not Enough Teams
				</button>
				<div
					class="absolute top-full right-0 z-10 mt-2 hidden w-64 rounded-xl bg-slate-900 p-3 text-center text-xs text-white shadow-xl group-hover:block"
				>
					<p class="mb-1 font-bold">Teams Required</p>
					<p class="opacity-70">You need at least 2 teams in a single sport to create a board.</p>
				</div>
			</div>
		{/if}
	</header>

	<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
		{#each data.boards as board (board.id)}
			<a
				href={resolve(`/dashboard/boards/${board.id}`)}
				class="group relative overflow-hidden rounded-3xl border border-slate-200 bg-white p-6 shadow-sm transition-all hover:border-blue-200 hover:shadow-xl"
			>
				<div
					class="absolute top-0 left-0 h-1.5 w-full"
					style="background: linear-gradient(90deg, 
                    {board.home_team?.secondary_color || '#cbd5e1'}, 
                    {board.away_team?.secondary_color || '#94a3b8'}
                )"
				></div>
				<div class="mb-4 flex items-center justify-between">
					<span
						class="rounded-full bg-slate-100 px-3 py-1 text-[10px] font-black tracking-widest text-slate-500 uppercase"
					>
						{board.status}
					</span>
					<span class="text-xs font-bold text-slate-400">
						${board.price_per_square}/sq
					</span>
				</div>

				<h2 class="mb-6 text-xl font-black text-slate-900 group-hover:text-blue-600">
					{board.title}
				</h2>

				<div class="mt-4 flex items-center gap-3">
					<div class="flex -space-x-3">
						{#if board.home_team?.team_logo_url}
							<img
								src={board.home_team.team_logo_url}
								alt=""
								class="h-10 w-10 rounded-full border-2 border-white bg-white object-contain shadow-sm"
							/>
						{/if}
						{#if board.away_team?.team_logo_url}
							<img
								src={board.away_team.team_logo_url}
								alt=""
								class="h-10 w-10 rounded-full border-2 border-white bg-white object-contain shadow-sm"
							/>
						{/if}
					</div>
					<div class="flex flex-col">
						<span class="text-[10px] font-black tracking-tighter text-slate-400 uppercase"
							>Matchup</span
						>
						<span class="text-xs font-bold text-slate-600">
							{board.home_team?.mascot} vs {board.away_team?.mascot}
						</span>
					</div>
				</div>

				<div class="flex items-center justify-between border-t border-slate-50 pt-4">
					<div class="flex -space-x-2">
						<div class="h-8 w-8 rounded-full border-2 border-white bg-slate-200"></div>
						<div class="h-8 w-8 rounded-full border-2 border-white bg-slate-300"></div>
					</div>
					<span class="text-xs font-bold text-slate-500"> View Details → </span>
				</div>
			</a>
		{:else}
			<div
				class="col-span-full flex flex-col items-center justify-center rounded-4xl border-2 border-dashed border-slate-200 py-20"
			>
				<p class="text-slate-400 italic">No boards found. Ready to start a game?</p>
			</div>
		{/each}
	</div>
</div> -->
