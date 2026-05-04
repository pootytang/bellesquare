<script lang="ts">
	import type { BoardPlayer } from '$lib/types';
	import VerifyPaymentModal from '../payments/VerifyPaymentModal.svelte'; // Ensure path is correct

	let {
		players,
		isCreator,
		boardId,
		token,
		currentUserId
	}: {
		players: BoardPlayer[];
		isCreator: boolean;
		boardId: string;
		token: string;
		currentUserId?: string;
	} = $props();

	// 1. Track which player is currently being "verified"
	let activePlayer = $state<BoardPlayer | null>(null);

	// 2. Open modal function
	function handlePaymentClick(player: BoardPlayer) {
		if (!isCreator) return;
		activePlayer = player;
	}
</script>

<div class="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
	<div class="border-b border-slate-200 bg-slate-50 px-5 py-3">
		<h3 class="text-sm font-bold text-slate-800">Participants ({players.length})</h3>
	</div>

	<div class="max-h-64 divide-y divide-slate-100 overflow-y-auto">
		{#each players as player (player.user_id)}
			<div class="flex items-center justify-between px-5 py-3">
				<div class="flex items-center gap-3">
					<div
						class="h-3 w-3 rounded-full border border-black/10"
						style="background-color: {player.color}"
					></div>

					<div class="flex flex-col">
						<span class="text-sm font-medium text-slate-700">
							{player.first_name}
							{player.last_name}
						</span>
						<span class="text-[10px] font-semibold tracking-wider text-slate-400 uppercase">
							{player.square_count}
							{player.square_count === 1 ? 'Square' : 'Squares'}
						</span>
					</div>
				</div>

				<div class="flex items-center gap-2">
					{#if isCreator}
						<!-- 3. Clicking this now opens the modal instead of calling an API -->
						<button
							onclick={() => handlePaymentClick(player)}
							class="rounded border px-2 py-1 text-[10px] font-bold uppercase transition-all
							{player.payment_status === 2
								? 'border-green-200 bg-green-100 text-green-700 hover:bg-green-200'
								: ''}
							{player.payment_status === 1
								? 'animate-pulse border-amber-200 bg-amber-100 text-amber-700 hover:bg-amber-200'
								: ''}
							{player.payment_status === 0 ? 'border-red-100 bg-red-50 text-red-600 hover:bg-red-100' : ''}
							cursor-pointer hover:scale-105 active:scale-95"
						>
							{#if player.payment_status === 2}
								Paid
							{/if}
							{#if player.payment_status === 1}
								Review
							{/if}
							{#if player.payment_status === 0}
								Unpaid
							{/if}
						</button>
					{:else if currentUserId && player.user_id === currentUserId}
						<span
							class="text-[10px] font-bold uppercase {player.payment_status === 2
								? 'text-green-500'
								: 'text-slate-400'}"
						>
							{player.payment_status === 2
								? 'Paid'
								: player.payment_status === 1
									? 'Pending'
									: 'Unpaid'}
						</span>
					{:else}
						<div class="w-4"></div>
					{/if}
				</div>
			</div>
		{:else}
			<div class="px-5 py-8 text-center text-slate-400 text-sm">
				No players have claimed squares yet.
			</div>
		{/each}
	</div>
</div>

<!-- Mount the Modal component -->
{#if activePlayer}
	<VerifyPaymentModal
		player={activePlayer}
		{boardId}
		{token}
		onClose={() => (activePlayer = null)}
	/>
{/if}

<!-- TODO: REMOVE THIS: THIS WAS WORKING FOR THE TOGGLING OF IS_PAID BUT MOVING TO PAYMENT_STATUS
 <script lang="ts">
	import type { BoardPlayer } from '$lib/types';
	import { PUBLIC_API_BASE_URL } from '$env/static/public';

	let {
		players,
		isCreator,
		boardId,
		token
	}: { players: BoardPlayer[]; isCreator: boolean; boardId: string; token: string } = $props();

	async function togglePaid(player: BoardPlayer) {
		if (!isCreator) return;
		const res = await fetch(
			`${PUBLIC_API_BASE_URL}/api/v1/boards/${boardId}/players/${player.user_id}/toggle-paid`,
			{
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${token}`
				}
				// Include this if you're using HttpOnly cookies for auth
				// credentials: 'include'
			}
		);
		if (res.ok) console.log('Payment status updated');
	}
</script>

<div class="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
	<div class="border-b border-slate-200 bg-slate-50 px-5 py-3">
		<h3 class="text-sm font-bold text-slate-800">Participants ({players.length})</h3>
	</div>

	<div class="max-h-64 divide-y divide-slate-100 overflow-y-auto">
		{#each players as player (player.user_id)}
			<div class="flex items-center justify-between px-5 py-3">
				<div class="flex items-center gap-3">
					<div
						class="h-3 w-3 rounded-full border border-black/10"
						style="background-color: {player.color}"
					></div>

					<div class="flex flex-col">
						<span class="text-sm font-medium text-slate-700">
							{player.first_name}
							{player.last_name}
						</span>
						<span class="text-[10px] font-semibold tracking-wider text-slate-400 uppercase">
							{player.square_count}
							{player.square_count === 1 ? 'Square' : 'Squares'}
						</span>
					</div>
				</div>
				{#if isCreator}
					<button
						onclick={() => togglePaid(player)}
						disabled={!isCreator}
						class="rounded px-2 py-1 text-[10px] font-bold uppercase transition-all
                    {player.is_paid
							? 'border border-green-200 bg-green-100 text-green-700 hover:bg-green-200'
							: 'border border-red-100 bg-red-50 text-red-600 hover:bg-red-100'} 
                    {!isCreator
							? 'cursor-default'
							: 'cursor-pointer hover:scale-105 active:scale-95'}"
					>
						{player.is_paid ? 'Paid' : 'Unpaid'}
					</button>
				{/if}
			</div>
		{:else}
			<div class="px-5 py-8 text-center text-slate-400 text-sm">
				No players have claimed squares yet.
			</div>
		{/each}
	</div>
</div> -->
