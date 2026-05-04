<script lang="ts">
	import { getVenmoLink } from '$lib/utils/payments';
	import { fade } from 'svelte/transition';
	import { onMount } from 'svelte';
	import { PUBLIC_API_BASE_URL } from '$env/static/public';

	let {
		unpaidCount,
		paymentStatus,
		pricePerSquare,
		hostVenmo,
		hostZelle,
		boardTitle,
		boardId, // Passed from +page.svelte
		token // Passed from +page.svelte
	}: {
		unpaidCount: number;
		paymentStatus: number;
		pricePerSquare: number;
		hostVenmo: string;
		hostZelle: string;
		boardTitle: string;
		boardId: string;
		token: string;
	} = $props();

	let totalOwed = $derived(unpaidCount * pricePerSquare);
	let copied = $state(false);
	let isMobile = $state(false);
	let isNotifying = $state(false);

	let isPending = $derived(paymentStatus === 1);

	onMount(() => {
		isMobile = /iPhone|iPad|iPod|Android/i.test(navigator.userAgent);
	});

	let venmoUrl = $derived(getVenmoLink(hostVenmo, totalOwed, boardTitle, isMobile));

	function copyZelle() {
		if (!hostZelle || isPending) return;
		navigator.clipboard.writeText(hostZelle);
		copied = true;
		setTimeout(() => (copied = false), 2000);
	}

	async function notifyPaymentSent() {
		isNotifying = true;
		try {
			const res = await fetch(`${PUBLIC_API_BASE_URL}/api/v1/boards/${boardId}/payment-sent`, {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${token}`
				}
			});
			if (res.ok) {
				// Success! The WS will broadcast the update and this
				// PaymentHub will disappear (because unpaidCount becomes 0 in square_status)
				// or the status will visually change.
			}
		} finally {
			isNotifying = false;
		}
	}
</script>

<div transition:fade class="rounded-3xl border border-amber-100 bg-amber-50 p-6 shadow-sm">
	<!-- Header Section -->
	<div class="mb-4 flex items-center justify-between">
		<div>
			<h3 class="text-xs font-black tracking-widest text-amber-800 uppercase">
				{isPending ? 'Payment Status' : 'Unpaid Squares'}
			</h3>
			<p class="text-2xl font-black text-amber-900">
				{isPending ? 'Pending Review' : `$${totalOwed.toFixed(2)}`}
			</p>
		</div>
		<div
			class="flex h-10 w-10 items-center justify-center rounded-full bg-amber-200 text-amber-700"
		>
			{#if isPending}
				<!-- Clock Icon for Pending -->
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2"
					stroke="currentColor"
					class="h-6 w-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
					/>
				</svg>
			{:else}
				<!-- Cash Icon -->
				<svg
					xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 24 24"
					fill="currentColor"
					class="h-6 w-6"
				>
					<path d="M12 7.5a2.25 2.25 0 1 0 0 4.5 2.25 2.25 0 0 0 0-4.5Z" />
					<path
						fill-rule="evenodd"
						d="M1.5 4.875C1.5 3.839 2.34 3 3.375 3h17.25c1.035 0 1.875.84 1.875 1.875v14.25c0 1.036-.84 1.875-1.875 1.875H3.375A1.875 1.875 0 0 1 1.5 19.125V4.875ZM12 15.75a3 3 0 1 0 0-6 3 3 0 0 0 0 6Z"
						clip-rule="evenodd"
					/>
				</svg>
			{/if}
		</div>
	</div>

	<div class="space-y-3">
		{#if hostVenmo}
			<div class="flex flex-col gap-1">
				<a
					href={isPending ? 'javascript:void(0)' : venmoUrl}
					rel="external noreferrer noopener"
					class="flex w-full items-center justify-center gap-2 rounded-xl py-3 text-xs font-bold text-white transition-all
                    {isPending
						? 'cursor-not-allowed bg-slate-300'
						: 'bg-[#3d95ce] hover:scale-[1.02]'}"
				>
					{isMobile ? 'Open Venmo App' : 'Pay via Venmo'}
				</a>
				<span class="text-center text-[10px] font-bold tracking-tighter text-slate-400 uppercase">
					Recipient: @{hostVenmo.replace('@', '')}
				</span>
			</div>
		{/if}

		{#if hostZelle}
			<button
				onclick={copyZelle}
				disabled={isPending}
				class="w-full rounded-xl border py-3 text-[10px] font-bold transition-colors
                {isPending
					? 'cursor-not-allowed border-slate-200 bg-slate-100 text-slate-400'
					: 'border-amber-200 bg-white text-amber-800 hover:bg-amber-100'}"
			>
				{isPending ? 'Zelle Info Locked' : copied ? 'Copied!' : `Copy Zelle: ${hostZelle}`}
			</button>
		{/if}

		<div class="border-t border-amber-200/50 pt-2">
			{#if isPending}
				<div
					class="flex w-full items-center justify-center gap-2 rounded-xl bg-amber-100 py-3 text-[10px] font-black tracking-widest text-amber-700 uppercase"
				>
					<span class="relative flex h-2 w-2">
						<span
							class="absolute inline-flex h-full w-full animate-ping rounded-full bg-amber-400 opacity-75"
						></span>
						<span class="relative inline-flex h-2 w-2 rounded-full bg-amber-500"></span>
					</span>
					Waiting for Creator Review
				</div>
			{:else}
				<button
					onclick={notifyPaymentSent}
					disabled={isNotifying}
					class="w-full rounded-xl bg-slate-900 py-3 text-[10px] font-black tracking-widest text-white uppercase shadow-lg shadow-slate-200 transition-all hover:bg-slate-800 active:scale-95 disabled:opacity-50"
				>
					{isNotifying ? 'Notifying Host...' : "I've Sent Payment"}
				</button>
			{/if}
		</div>
	</div>
</div>
