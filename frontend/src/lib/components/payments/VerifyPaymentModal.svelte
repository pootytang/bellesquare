<script lang="ts">
	import type { BoardPlayer } from '$lib/types';
	import { PUBLIC_API_BASE_URL } from '$env/static/public';
	import { fade, scale } from 'svelte/transition';

	let {
		player,
		boardId,
		token,
		onClose
	}: {
		player: BoardPlayer;
		boardId: string;
		token: string;
		onClose: () => void;
	} = $props();

	async function updateStatus(status: number) {
		try {
			const res = await fetch(
				`${PUBLIC_API_BASE_URL}/api/v1/boards/${boardId}/players/${player.user_id}/verify`,
				{
					method: 'PUT',
					headers: {
						'Content-Type': 'application/json',
						Authorization: `Bearer ${token}`
					},
					body: JSON.stringify({ status })
				}
			);

			if (res.ok) {
				onClose();
			}
		} catch (err) {
			console.error('Failed to update status:', err);
		}
	}
</script>

<div
	transition:fade
	class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 p-4 backdrop-blur-sm"
>
	<div
		transition:scale={{ start: 0.95 }}
		class="w-full max-w-sm rounded-3xl bg-white p-8 shadow-2xl"
	>
		<h3 class="text-xl font-black text-slate-800">Verify Payment</h3>

		<p class="mt-3 text-sm leading-relaxed text-slate-500">
			Confirming payment for <span class="font-bold text-slate-900"
				>{player.first_name} {player.last_name}</span
			>.
		</p>

		<div class="mt-8 flex flex-col gap-3">
			<button
				onclick={() => updateStatus(2)}
				class="w-full rounded-2xl bg-green-600 py-4 text-sm font-bold text-white shadow-lg shadow-green-200 transition-transform hover:bg-green-700 active:scale-95"
			>
				Confirm & Mark Paid
			</button>

			<button
				onclick={() => updateStatus(0)}
				class="w-full rounded-2xl bg-red-50 py-4 text-sm font-bold text-red-600 transition-colors hover:bg-red-100"
			>
				Payment Not Received
			</button>

			<button
				onclick={onClose}
				class="mt-2 text-xs font-semibold tracking-widest text-slate-400 uppercase hover:text-slate-600"
			>
				Go Back
			</button>
		</div>
	</div>
</div>
