<script lang="ts">
	import { enhance } from '$app/forms';
	import type { Square } from '$lib/types';
	import { page } from '$app/state';
	import { resolve } from '$app/paths';

	let { data, isFull, hasNumbers, isLocked } = $props();

	// Calculate claimed count from the live grid
	let claimedCount = $derived(data.squares.flat().filter((s: Square) => !!s.user_id).length);

	// Invite Logic
	let copied = $state(false);
	const shareUrl = $derived(page.url.href);

	function copyToClipboard() {
		navigator.clipboard.writeText(shareUrl);
		copied = true;
		setTimeout(() => (copied = false), 2000);
	}
</script>

<div
	class="mb-12 grid grid-cols-3 items-center rounded-3xl border border-slate-200 bg-slate-50 p-6"
>
	<div class="flex flex-col gap-1">
		<p class="text-xs font-bold tracking-widest text-slate-500 uppercase">
			Status: <span class="text-slate-900">{data.board.status}</span>
		</p>
		<div class="flex items-center gap-2">
			<span
				class="text-xs font-black tracking-tighter {isFull
					? 'text-green-600'
					: 'text-slate-400'} uppercase"
			>
				{isFull ? 'Board Full!' : `${claimedCount}/100 Squares`}
			</span>
			{#if !hasNumbers}
				<span
					class="rounded-full bg-amber-100 px-2 py-0.5 text-[10px] font-bold text-amber-700 uppercase"
				>
					Numbers Pending
				</span>
			{/if}
		</div>
	</div>

	<div class="text-center">
		<h1 class="text-3xl font-black tracking-tight text-slate-900">{data.board.title}</h1>
	</div>

	{#if data.user}
		{#if data.user.is_guest}
			<!-- User is logged in as a guest: Bring them to UPGRADE -->
			<div class="flex justify-end">
				<a
					href={resolve('/dashboard/upgrade')}
					class="animate-pulse rounded-lg bg-amber-400 px-4 py-2 font-bold text-white shadow-lg shadow-amber-200 transition-colors hover:bg-amber-500"
				>
					Unlock Hosting 🚀
				</a>
			</div>
		{/if}
	{:else}
		<!-- User is NOT logged in: Bring them to REGISTER -->
		<div class="flex justify-end">
			<a
				href={resolve('/register')}
				class="rounded-lg bg-blue-600 px-4 py-2 font-bold text-white transition-colors hover:bg-blue-700"
			>
				Create Account to Host
			</a>
		</div>
	{/if}

	<div class="flex justify-end">
		{#if data.isCreator && !isLocked}
			<div class="flex gap-3">
				<button
					onclick={copyToClipboard}
					class="rounded-xl border border-slate-200 bg-white px-5 py-2.5 text-sm font-bold text-slate-700 transition-all hover:bg-slate-50 active:scale-95"
				>
					{copied ? 'Copied!' : 'Invite Players'}
				</button>

				<form method="POST" action="?/shuffle" use:enhance>
					<button
						class="rounded-xl border border-slate-200 bg-white px-5 py-2.5 text-sm font-bold text-slate-700 transition-all hover:bg-slate-50 active:scale-95"
					>
						{hasNumbers ? 'Re-shuffle' : 'Generate Numbers'}
					</button>
				</form>

				{#if hasNumbers}
					<form method="POST" action="?/lock" use:enhance>
						<button
							class="rounded-xl bg-red-600 px-5 py-2.5 text-sm font-bold text-white shadow-lg shadow-red-100 transition-all hover:bg-red-700 active:scale-95"
						>
							Lock Board
						</button>
					</form>
				{/if}
			</div>
		{:else if isLocked}
			<div
				class="flex items-center gap-2 rounded-lg bg-green-100 px-4 py-2 font-bold text-green-700"
			>
				<span>🔒 Locked</span>
			</div>

			{#if data.isCreator}
				<form method="POST" action="?/unlock" use:enhance>
					<button
						type="submit"
						class="rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-xs font-bold text-slate-500 transition-all hover:bg-slate-100 hover:text-slate-900 active:scale-95"
					>
						Unlock
					</button>
				</form>
			{/if}
		{/if}
	</div>
</div>
