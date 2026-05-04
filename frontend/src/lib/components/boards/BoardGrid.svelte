<script lang="ts">
	import SquareCell from './SquareCell.svelte';
	import type { BoardPageData } from '$lib/types';
	import { SvelteMap } from 'svelte/reactivity';
	import type { Square } from '$lib/types';

	// Define the shape of the data prop specifically
	let {
		data,
		pendingKeys,
		userColor,
		allWinningCoords,
		onSquareClick
	}: {
		data: BoardPageData;
		pendingKeys: string[];
		userColor: string;
		allWinningCoords: { row: number; col: number; isLocked: boolean; isCurrent: boolean }[];
		onSquareClick: (key: string) => void;
	} = $props();

	const range = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];
	let hasNumbers = $derived(data.board.home_axis_numbers?.length > 0);

	// OPTIMIZATION: Create a Map once when data.squares changes
	// This makes the lookup inside the loop O(1) instead of O(n)
	let squareLookup = $derived.by(() => {
		const map = new SvelteMap<string, Square>();

		// Use the 'unknown' bridge to satisfy the strict compiler
		const grid = data.squares as unknown as Square[][];

		if (grid && Array.isArray(grid)) {
			grid.forEach((row, rowIndex) => {
				// Double check it's actually a row array
				if (Array.isArray(row)) {
					row.forEach((square, colIndex) => {
						map.set(`${rowIndex}-${colIndex}`, square);
					});
				}
			});
		}
		return map;
	});
</script>

<div
	class="rounded-4xl border border-slate-100/50 bg-white p-8 shadow-[0_8px_30px_rgb(0,0,0,0.04)]"
>
	<div
		class="mb-8 ml-19 flex items-center justify-center gap-4 rounded-2xl py-4 text-white shadow-lg"
		style="background: linear-gradient(90deg, {data.homeTeam.primary_color}, {data.homeTeam
			.secondary_color});"
	>
		{#if data.homeTeam.team_logo_url}
			<img src={data.homeTeam.team_logo_url} alt="" class="h-12 w-12 object-contain" />
		{/if}
		<h2 class="text-2xl font-black tracking-widest uppercase">{data.homeTeam.full_name}</h2>
	</div>

	<div class="flex gap-6">
		<div
			class="mt-19 flex rotate-180 items-center justify-center gap-4 rounded-2xl px-4 py-8 text-white shadow-lg [writing-mode:vertical-lr]"
			style="background: linear-gradient(180deg, {data.awayTeam.primary_color}, {data.awayTeam
				.secondary_color});"
		>
			<h2 class="text-2xl font-black tracking-widest uppercase">{data.awayTeam.full_name}</h2>
			{#if data.awayTeam.team_logo_url}
				<img src={data.awayTeam.team_logo_url} alt="" class="h-12 w-12 rotate-180 object-contain" />
			{/if}
		</div>

		<div class="grid grid-cols-[repeat(11,minmax(64px,1fr))] gap-3">
			<div class="aspect-square"></div>
			{#each range as i (i)}
				<div
					class="flex items-center justify-center rounded-xl text-2xl font-black"
					style="color: {data.homeTeam.primary_color}; background-color: {data.homeTeam
						.primary_color}15;"
				>
					{hasNumbers ? data.board.home_axis_numbers[i] : '?'}
				</div>
			{/each}

			{#each range as row (row)}
				<div
					class="flex items-center justify-center rounded-xl text-2xl font-black"
					style="color: {data.awayTeam.primary_color}; background-color: {data.awayTeam
						.primary_color}15;"
				>
					{hasNumbers ? data.board.away_axis_numbers[row] : '?'}
				</div>

				{#each range as col (col)}
					{@const key = `${row}-${col}`}
					{@const square = squareLookup.get(key)}
					{@const winInfo = allWinningCoords.find((c) => c.row === row && c.col === col)}

					<SquareCell
						{square}
						isLocked={data.board.status === 'locked' || !data.user}
						isPending={pendingKeys.includes(key)}
						isWinner={!!winInfo}
						isCurrentWinner={winInfo?.isCurrent ?? false}
						pendingColor={userColor}
						userInitials={data.user?.initials || '??'}
						onclick={() => {
							if (data.user) onSquareClick(key);
						}}
					/>
				{/each}
			{/each}
		</div>
	</div>
</div>
