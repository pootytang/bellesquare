<script lang="ts">
	import { enhance } from '$app/forms';
	import type { Team, ActionData } from '$lib/types';

	let {
		team,
		onResult,
		onclose
	}: { team: Team; onResult: (data: ActionData) => void; onclose: () => void } = $props();
	let loading = $state(false);
</script>

<div class="fixed inset-0 z-110 flex items-center justify-center p-4">
	<button
		aria-label="Close Modal"
		class="absolute inset-0 cursor-default bg-slate-900/60 backdrop-blur-sm"
		onclick={onclose}
	>
	</button>

	<div class="animate-in zoom-in-95 relative w-full max-w-sm rounded-3xl bg-white p-8 shadow-2xl">
		<div class="mb-6 text-center">
			<div
				class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-red-100 text-red-600"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-8 w-8"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
					/>
				</svg>
			</div>
			<h2 class="text-xl font-black text-slate-900">Delete Team?</h2>
			<p class="mt-2 text-sm text-slate-500">
				Are you sure you want to remove <span class="font-bold text-slate-800"
					>{team.full_name}</span
				>? This action cannot be undone.
			</p>
		</div>

		<form
			method="POST"
			action="?/delete"
			use:enhance={() => {
				loading = true;
				return async ({ result, update }) => {
					loading = false;
					if (result.type === 'success' || result.type === 'failure') {
						onResult(result.data as unknown as ActionData);
						if (result.type === 'success') onclose();
					}
					await update();
				};
			}}
		>
			<input type="hidden" name="id" value={team.id} />
			<div class="flex flex-col gap-2">
				<button
					disabled={loading}
					class="w-full cursor-pointer rounded-xl bg-red-600 py-4 font-bold text-white transition-all hover:bg-red-700 disabled:opacity-50"
				>
					{loading ? 'Deleting...' : 'Yes, Delete Team'}
				</button>
				<button
					type="button"
					onclick={onclose}
					class="w-full cursor-pointer py-2 text-sm font-bold text-slate-400 hover:text-slate-600"
				>
					Cancel
				</button>
			</div>
		</form>
	</div>
</div>
