<script lang="ts">
	import type {
		BoardPageData,
		Square,
		SquareUpdate,
		AxisUpdate,
		BoardPlayer,
		QuarterScore
	} from '$lib/types';
	import BoardHeader from '$lib/components/boards/BoardHeader.svelte';
	import BoardGrid from '$lib/components/boards/BoardGrid.svelte';
	import SelectionSidebar from '$lib/components/boards/SelectionSidebar.svelte';
	import PayoutsCard from '$lib/components/boards/PayoutsCard.svelte';
	import PlayersList from '$lib/components/boards/PlayersList.svelte';
	import Scoreboard from '$lib/components/boards/Scoreboard.svelte';
	import PaymentHub from '$lib/components/payments/PaymentHub.svelte';
	import { PUBLIC_API_BASE_URL } from '$env/static/public';
	import { onMount, untrack } from 'svelte';
	import { SvelteMap } from 'svelte/reactivity';
	// import { page } from '$app/state';

	let { data }: { data: BoardPageData } = $props();

	// --- Reactive State (The "Live" Overrides) ---
	// We initialize these to null/empty so we can fallback to 'data' reactively
	let liveUpdates = $state<Record<string, SquareUpdate>>({});
	let liveStatus = $state<string | null>(null);
	let liveHomeAxis = $state<number[] | null>(null);
	let liveAwayAxis = $state<number[] | null>(null);

	// ********** SCOREBOARD DATA ********** //
	// Track independent quarters
	let liveQuarterScores = $state([
		{ quarter: 1, home_score: 0, away_score: 0, is_locked: false },
		{ quarter: 2, home_score: 0, away_score: 0, is_locked: false },
		{ quarter: 3, home_score: 0, away_score: 0, is_locked: false },
		{ quarter: 4, home_score: 0, away_score: 0, is_locked: false },
		{ quarter: 5, home_score: 0, away_score: 0, is_locked: false }
	]);

	// ********** SCOREBOARD SYNC & AUTO-FINAL LOGIC ********** //
	$effect(() => {
		// We want this effect to run when 'data' changes (server/load)
		// but NOT when 'liveQuarterScores' changes (typing)
		if (data.scores?.quarters) {
			untrack(() => {
				data.scores.quarters.forEach((q: QuarterScore) => {
					const index = q.quarter - 1;
					const local = liveQuarterScores[index];
					if (local) {
						if (local.home_score !== q.home_score) local.home_score = q.home_score;
						if (local.away_score !== q.away_score) local.away_score = q.away_score;
						if (local.is_locked !== q.is_locked) local.is_locked = q.is_locked;
					}
				});
			});
		}
	});

	// --- Derived Logic (The "Brain") ---
	const grid = $derived(data.squares as Square[][]);
	const totalHome = $derived(
		[...liveQuarterScores].reverse().find((q) => q.home_score > 0 || q.away_score > 0)
			?.home_score ?? 0
	);
	const totalAway = $derived(
		[...liveQuarterScores].reverse().find((q) => q.home_score > 0 || q.away_score > 0)
			?.away_score ?? 0
	);
	// Calculate the squares once and use in multiple places
	const processedSquares = $derived(
		grid.map((row, rIdx) =>
			row.map((sq, cIdx) => {
				const key = `${rIdx}-${cIdx}`;
				const update = liveUpdates[key];
				// Merge live updates. Crucially, ensure we keep first_name/last_name from sq
				return update ? { ...sq, ...update } : sq;
			})
		)
	);
	// ****************************************************** //

	// ********** HIGHLIGHTING THE WINNING CELL ********** //
	const allWinningCoords = $derived.by(() => {
		const lastDigit = (n: number) => Math.abs(n % 10);
		const lastActiveIdx = liveQuarterScores.findLastIndex(
			(qs) => qs.home_score > 0 || qs.away_score > 0
		);

		return liveQuarterScores.flatMap((q, idx) => {
			// Skip quarters with no score yet
			if (q.home_score === 0 && q.away_score === 0 && idx !== 0) return [];

			const rowIdx = data.board.away_axis_numbers?.indexOf(lastDigit(q.away_score));
			const colIdx = data.board.home_axis_numbers?.indexOf(lastDigit(q.home_score));

			if (rowIdx === -1 || colIdx === -1 || rowIdx === undefined) return [];

			return [
				{
					row: rowIdx,
					col: colIdx,
					isLocked: q.is_locked,
					isCurrent: idx === lastActiveIdx,
					quarter: idx + 1
				}
			];
		});
	});
	// ****************************************************** //

	// ********** PAYOUT WINNER LOGIC ********** //
	const winningPlayersByQuarter = $derived.by(() => {
		const lastDigit = (n: number) => Math.abs(n % 10);

		// Identify the latest quarter that has a score entered
		const lastActiveIdx = liveQuarterScores.findLastIndex(
			(qs) => qs.home_score > 0 || qs.away_score > 0
		);

		return liveQuarterScores.map((q, idx) => {
			// A quarter is eligible for display if it's locked OR it's the current live quarter
			const isLive = (q.home_score > 0 || q.away_score > 0) && idx === lastActiveIdx;

			// Final Score (Quarter 5) often remains unlocked until the very end,
			// so we treat it as "viewable" if it has a score.
			if (!q.is_locked && !isLive && q.quarter !== 5) return null;

			const rowIdx = data.board.away_axis_numbers?.indexOf(lastDigit(q.away_score));
			const colIdx = data.board.home_axis_numbers?.indexOf(lastDigit(q.home_score));

			if (rowIdx === -1 || colIdx === -1 || rowIdx === undefined || colIdx === undefined) {
				return 'TBD';
			}

			const key = `${rowIdx}-${colIdx}`;
			const owner = data.square_owners?.[key];

			if (owner && owner.first_name) {
				return `${owner.first_name} ${owner.last_name}`;
			}

			// If coordinates exist but no owner is found in the map
			return 'House Wins (Unsold)';
		});
	});
	// ****************************************************** //

	// ********** LIVE DATA / WEB SOCKET ********** //
	// This creates our final source of truth by merging initial data with live overrides
	let liveData = $derived({
		...data,
		board: {
			...data.board,
			status: (liveStatus ?? data.board.status) as string,
			home_axis_numbers: liveHomeAxis ?? data.board.home_axis_numbers,
			away_axis_numbers: liveAwayAxis ?? data.board.away_axis_numbers,
			payouts: data.payouts ?? []
		},
		squares: processedSquares, // Use the pre-calculated squares
		scores: {
			home: totalHome,
			away: totalAway,
			quarters: liveQuarterScores
		},
		players: (() => {
			const playerMap = new SvelteMap<string, BoardPlayer>();

			// console.log('Owner Map:', data.square_owners);
			processedSquares.flat().forEach((sq) => {
				if (sq.user_id) {
					const key = `${sq.row_index}-${sq.col_index}`;

					// 1. Get the owner details from the lookup map we just added
					const owner = data.square_owners?.[key];

					const existing = playerMap.get(sq.user_id);
					if (existing) {
						existing.square_count += 1;
						// This handles cases where some squares might be marked paid before others
						if (sq.payment_status > existing.payment_status) {
							existing.payment_status = sq.payment_status;
						}
					} else {
						playerMap.set(sq.user_id, {
							user_id: sq.user_id,
							first_name: owner?.first_name || sq.user_initials || 'User',
							last_name: owner?.last_name || '',
							color: owner?.color || sq.user_color || '#cccccc',
							square_count: 1,
							// ADD THESE TWO:
							is_paid: sq.is_paid ?? false,
							payment_status: sq.payment_status ?? 0 // Default to Unpaid (0)
						});
					}
				}
			});

			return Array.from(playerMap.values()).sort((a, b) => b.square_count - a.square_count);
		})()
	});

	onMount(() => {
		const wsUrl = PUBLIC_API_BASE_URL.replace('http', 'ws');
		const socket = new WebSocket(`${wsUrl}/api/v1/boards/${data.board.id}/ws`);

		socket.onmessage = (event) => {
			const rawPayload = JSON.parse(event.data);

			// Handle Status Updates
			if (rawPayload?.type === 'STATUS_UPDATE') {
				liveStatus = rawPayload.status;
			}

			// Handle Axis Updates
			if (rawPayload?.type === 'AXIS_UPDATE') {
				const update = rawPayload as AxisUpdate;
				liveHomeAxis = update.home;
				liveAwayAxis = update.away;
			}

			// Handle Square Updates (Refined)
			if (rawPayload?.type === 'SQUARE_UPDATE' && Array.isArray(rawPayload.squares)) {
				const batch = rawPayload.squares as SquareUpdate[];
				batch.forEach((update) => {
					// Update key to match your new SquareUpdate structure
					const key = `${update.row}-${update.col}`;
					liveUpdates[key] = update;
				});
			}

			if (rawPayload?.type === 'SCORE_UPDATE') {
				// console.table(rawPayload);
				const { quarter, home_score, away_score, is_locked } = rawPayload;
				// quarter is 1-indexed from Go (1, 2, 3, 4)
				const index = quarter - 1;

				if (liveQuarterScores[index]) {
					liveQuarterScores[index].home_score = home_score;
					liveQuarterScores[index].away_score = away_score;
					liveQuarterScores[index].is_locked = is_locked;
				}
			}

			// Legacy/Initial Claim Handling
			// Handle Square Updates (Array of updates)
			else if (Array.isArray(rawPayload)) {
				const batch = rawPayload as SquareUpdate[];
				batch.forEach((update) => {
					const key = `${update.row}-${update.col}`;
					liveUpdates[key] = update;
				});
			}
		};

		return () => socket.close();
	});
	// ****************************************************** //

	// --- 4. UI State & Helpers ---
	let pendingKeys = $state<string[]>([]);
	let selectedColor = $state<string | null>(null);
	let userColor = $derived(selectedColor ?? data.user?.color ?? '#2563eb');

	let isLocked = $derived(liveData.board.status === 'locked');
	let hasNumbers = $derived((liveData.board.home_axis_numbers?.length ?? 0) > 0);

	let dbFilledCount = $derived(liveData.squares.flat().filter((s) => !!s.user_id).length);
	let isFull = $derived(dbFilledCount + pendingKeys.length >= 100);

	function handleSquareClick(key: string) {
		if (isLocked) return;
		if (pendingKeys.includes(key)) {
			pendingKeys = pendingKeys.filter((k) => k !== key);
		} else if (!isFull) {
			pendingKeys = [...pendingKeys, key];
		}
	}

	// console.log(`ZELLE: ${data.creator_zelle}, VENMO: ${data.creator_venmo}`);
</script>

<div class="min-h-screen bg-slate-50 pb-20">
	<div class="mb-8 w-full">
		<BoardHeader data={liveData} {isFull} {hasNumbers} {isLocked} />
	</div>

	<!-- 
      MAIN WRAPPER: 
      - items-start: Ensures all columns pin to the top 
      - justify-between: Spreads columns to the edges
    -->
	<div
		class="mx-auto flex w-full max-w-425 flex-col gap-6 px-4 lg:flex-row lg:items-start lg:justify-between"
	>
		<!-- 
          LEFT ASIDE: Reference Info 
          - Order 2 on mobile, 1 on desktop
        -->
		<aside
			class="order-2 flex w-full shrink-0 flex-col gap-6 lg:sticky lg:top-8 lg:order-1 lg:w-64"
		>
			<PayoutsCard
				payouts={liveData.board.payouts}
				quarterWinners={winningPlayersByQuarter}
				isCreator={data.isCreator}
			/>

			<PlayersList
				players={liveData.players ?? []}
				isCreator={data.isCreator}
				boardId={data.board.id}
				token={data.token}
				currentUserId={data.user?.id}
			/>
		</aside>

		<!-- 
          CENTER MAIN: The Game Board 
          - Order 1 on mobile (Board first), 2 on desktop
        -->
		<main class="order-1 min-w-0 flex-1 lg:order-2 lg:-mt-4">
			<!-- pb-10 added for scrollbar clearance -->
			<div class="custom-scrollbar w-full overflow-x-auto pb-10">
				<div class="inline-block align-top">
					<BoardGrid
						data={liveData}
						{pendingKeys}
						{userColor}
						{allWinningCoords}
						onSquareClick={handleSquareClick}
					/>
				</div>
			</div>
		</main>

		<!-- 
          RIGHT ASIDE: Interactive Actions 
          - Order 3 on mobile and desktop
        -->
		<aside class="order-3 flex w-full shrink-0 flex-col gap-6 lg:sticky lg:top-8 lg:w-80">
			{#if data.user && !data.isCreator}
				{@const unpaidCount = liveData.squares
					.flat()
					.filter((s) => s.user_id === data.user?.id && !s.is_paid).length}

				{#if unpaidCount > 0}
					{@const currentPlayer = liveData.players.find((p) => p.user_id === data.user?.id)}
					<PaymentHub
						{unpaidCount}
						paymentStatus={currentPlayer?.payment_status ?? 0}
						pricePerSquare={data.board.price_per_square}
						hostVenmo={data.creator_venmo}
						hostZelle={data.creator_zelle}
						boardTitle={data.board.title}
						boardId={data.board.id}
						token={data.token}
					/>
				{/if}
			{/if}
			<SelectionSidebar
				bind:selectedKeys={pendingKeys}
				pricePerSquare={data.board.price_per_square}
				bind:userColor
				{isFull}
				{dbFilledCount}
				user={liveData.user}
			/>

			<Scoreboard
				scores={{ quarters: liveQuarterScores }}
				homeTeam={data.homeTeam}
				awayTeam={data.awayTeam}
				isCreator={data.isCreator}
				boardId={data.board.id}
				token={data.token}
			/>
		</aside>
	</div>
</div>

<style>
	/* Prevent the main body from horizontal scrolling, only the grid should */
	:global(body) {
		overflow-x: hidden;
	}
	.custom-scrollbar::-webkit-scrollbar {
		height: 8px;
	}
	.custom-scrollbar::-webkit-scrollbar-thumb {
		background: #cbd5e1;
		border-radius: 10px;
	}
</style>
