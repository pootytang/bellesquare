<script lang="ts">
	import { enhance } from '$app/forms';
	import type { User } from '$lib/types';
	import { fade } from 'svelte/transition';

	let {
		selectedKeys = $bindable(),
		pricePerSquare,
		userColor = $bindable(),
		isFull,
		dbFilledCount,
		user
	}: {
		selectedKeys: string[];
		pricePerSquare: number;
		userColor: string;
		isFull: boolean;
		dbFilledCount: number;
		user: User | null;
	} = $props();

	let total = $derived(selectedKeys.length * pricePerSquare);
	let showSuccess = $state(false);

	// Standard professional presets
	const colorPresets = ['#2563eb', '#dc2626', '#16a34a', '#ca8a04', '#7c3aed', '#db2777'];
</script>

<div
	class="sticky top-6 w-full max-w-sm space-y-6 rounded-3xl border border-slate-200 bg-white p-8 shadow-sm"
>
	<h3 class="text-xl font-black tracking-tight text-slate-900 uppercase">Claim Squares</h3>

	{#if showSuccess}
		<div transition:fade class="rounded-2xl border border-green-100 bg-green-50 p-4 text-center">
			<p class="text-sm font-bold text-green-700">Squares Claimed Successfully! 🎉</p>
		</div>
	{/if}

	<div class="space-y-6">
		<div class="space-y-3">
			<label
				for="color"
				class="block text-[10px] font-bold tracking-widest text-slate-400 uppercase"
			>
				Square Color
			</label>

			<div class="flex items-center gap-3">
				<div class="flex flex-wrap gap-2">
					{#each colorPresets as color (color)}
						<button
							type="button"
							onclick={() => (userColor = color)}
							class="h-7 w-7 rounded-full border-2 shadow-sm transition-transform hover:scale-110 active:scale-95"
							style="background-color: {color}; border-color: {userColor === color
								? 'white'
								: 'transparent'}; outline: {userColor === color ? '2px solid #3b82f6' : 'none'}"
							aria-label="Select {color}"
						></button>
					{/each}
				</div>

				<div class="h-6 w-px bg-slate-200"></div>

				<div class="flex flex-col items-center gap-1">
					<div
						class="relative h-8 w-8 overflow-hidden rounded-full border-2 shadow-sm transition-transform hover:scale-105"
						style="border-color: {!colorPresets.includes(userColor)
							? 'white'
							: '#e2e8f0'}; outline: {!colorPresets.includes(userColor)
							? '2px solid #3b82f6'
							: 'none'}"
					>
						<input
							type="color"
							id="user_color"
							bind:value={userColor}
							class="absolute -inset-2 h-12 w-12 cursor-pointer border-none bg-transparent"
						/>
					</div>
					<span class="text-[8px] font-black tracking-tighter text-slate-400 uppercase">Custom</span
					>
				</div>
			</div>
			<p class="text-[10px] text-slate-400 italic">Select a preset or choose a custom color.</p>
		</div>

		<div class="border-t border-slate-100 pt-6">
			<div class="mb-4 flex items-end justify-between">
				<span class="text-xs font-bold text-slate-500 uppercase">
					Total for {selectedKeys.length} squares
				</span>
				<span class="text-2xl font-black text-blue-600">${total.toFixed(2)}</span>
			</div>

			{#if !user}
				<div class="rounded-2xl border border-amber-100 bg-amber-50 p-4 text-center">
					<p class="text-xs font-bold text-amber-700">Please sign in to claim squares</p>
				</div>
			{:else}
				<form
					method="POST"
					action="?/claim"
					use:enhance={() => {
						return async ({ result, update }) => {
							if (result.type === 'success') {
								selectedKeys = [];
								showSuccess = true;
								setTimeout(() => (showSuccess = false), 4000);
							}
							await update();
						};
					}}
				>
					<input type="hidden" name="color" value={userColor} />
					<input type="hidden" name="selections" value={JSON.stringify(selectedKeys)} />

					<button
						type="submit"
						disabled={selectedKeys.length === 0}
						class="w-full rounded-2xl bg-slate-900 py-4 font-black text-white shadow-lg shadow-slate-200 transition-all hover:bg-blue-600 hover:shadow-blue-200 active:scale-95 disabled:bg-slate-100 disabled:text-slate-400"
					>
						{isFull && dbFilledCount < 100 ? 'Claim Final Squares' : 'Claim & Pay'}
					</button>
				</form>
			{/if}
		</div>
	</div>
</div>
