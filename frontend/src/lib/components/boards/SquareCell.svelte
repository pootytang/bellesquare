<script lang="ts">
	import type { Square } from '$lib/types';

	let {
		square,
		isLocked,
		isPending,
		isWinner,
		isCurrentWinner,
		pendingColor,
		userInitials,
		onclick
	}: {
		square: Square | undefined;
		isLocked: boolean;
		isPending: boolean;
		isWinner: boolean;
		isCurrentWinner: boolean;
		pendingColor: string;
		userInitials: string;
		onclick: () => void;
	} = $props();

	// Determine if text should be white or black based on background brightness
	function getContrastColor(hex: string) {
		if (!hex) return 'white';
		const r = parseInt(hex.slice(1, 3), 16);
		const g = parseInt(hex.slice(3, 5), 16);
		const b = parseInt(hex.slice(5, 7), 16);
		const yiq = (r * 299 + g * 587 + b * 114) / 1000;
		return yiq >= 128 ? '#0f172a' : 'white';
	}

	let bgColor = $derived.by(() => {
		if (square?.user_id) return square.user_color || '#2563eb';
		if (isPending) return pendingColor;
		return '#f8fafc'; // slate-50
	});

	let textColor = $derived(getContrastColor(bgColor));
</script>

<button
	type="button"
	disabled={isLocked || !!square?.user_id}
	{onclick}
	style:background-color={bgColor}
	style:color={textColor}
	class="group relative flex aspect-square items-center justify-center rounded-xl border-2 transition-all
           {square?.user_id
		? 'border-transparent shadow-sm'
		: 'border-slate-100 hover:border-blue-400'}
           {isPending ? 'z-10 scale-110 shadow-xl ring-2 ring-blue-500 ring-offset-2' : ''}
           
           /* WINNER LOGIC */
           {isWinner ? 'z-20 ring-4' : ''}
           {isWinner && isCurrentWinner
		? 'scale-110 animate-pulse shadow-[0_0_25px_rgba(250,204,21,0.8)] ring-yellow-400'
		: ''}
           {isWinner && !isCurrentWinner
		? 'scale-100 opacity-95 shadow-md ring-emerald-400/70'
		: ''}"
>
	{#if square?.user_id}
		<span class="text-[10px] font-black tracking-tighter uppercase opacity-90">
			{square.user_initials || 'Taken'}
		</span>
	{:else if isPending}
		<span class="animate-pulse text-[10px] font-black tracking-tighter uppercase">
			{userInitials}
		</span>
	{:else if !isLocked}
		<span class="text-xl font-light text-slate-300 opacity-0 group-hover:opacity-100">+</span>
	{/if}

	{#if isWinner && isCurrentWinner}
		<span class="absolute -top-3 -right-3 text-lg drop-shadow-md">👑</span>
	{:else if isWinner}
		<span class="absolute -top-2 -right-2 text-xs opacity-70 grayscale">🏁</span>
	{/if}
</button>

<!-- <button
	type="button"
	disabled={isLocked || !!square?.user_id}
	{onclick}
	style:background-color={bgColor}
	style:color={textColor}
	class="group relative flex aspect-square items-center justify-center rounded-xl border-2 transition-all
           {square?.user_id
		? 'border-transparent shadow-sm'
		: 'border-slate-100 hover:border-blue-400'}
           {isPending ? 'z-10 scale-110 shadow-xl ring-2 ring-blue-500 ring-offset-2' : ''}
           {isWinner
		? 'z-20 scale-105 animate-pulse shadow-[0_0_20px_rgba(250,204,21,0.6)] ring-4 ring-yellow-400'
		: ''}"
>
	{#if square?.user_id}
		<span class="text-[10px] font-black tracking-tighter uppercase opacity-90">
			{square.user_initials || 'Taken'}
		</span>
	{:else if isPending}
		<span class="animate-pulse text-[10px] font-black tracking-tighter uppercase">
			{userInitials}
		</span>
	{:else if !isLocked}
		<span class="text-xl font-light text-slate-300 opacity-0 group-hover:opacity-100">+</span>
	{/if}

	{#if isWinner}
		<span class="absolute -top-2 -right-2 text-xs">👑</span>
	{/if}
</button> -->

<!-- <script lang="ts">
	let { square, isLocked, isPending, pendingColor, userInitials, onclick } = $props();

	// Determine if text should be white or black based on background brightness
	function getContrastColor(hex: string) {
		if (!hex) return 'white';
		const r = parseInt(hex.slice(1, 3), 16);
		const g = parseInt(hex.slice(3, 5), 16);
		const b = parseInt(hex.slice(5, 7), 16);
		const yiq = (r * 299 + g * 587 + b * 114) / 1000;
		return yiq >= 128 ? '#0f172a' : 'white';
	}

	let bgColor = $derived.by(() => {
		if (square?.user_id) return square.user_color || '#2563eb';
		if (isPending) return pendingColor;
		return '#f8fafc'; // slate-50
	});

	let textColor = $derived(getContrastColor(bgColor));
</script>

<button
	type="button"
	disabled={isLocked || !!square?.user_id}
	{onclick}
	style:background-color={bgColor}
	style:color={textColor}
	class="group relative flex aspect-square items-center justify-center rounded-xl border-2 transition-all
           {square?.user_id
		? 'border-transparent shadow-sm'
		: 'border-slate-100 hover:border-blue-400'}
           {isPending ? 'z-10 scale-110 shadow-xl ring-2 ring-blue-500 ring-offset-2' : ''}"
>
	{#if square?.user_id}
		<span class="text-[10px] font-black tracking-tighter uppercase opacity-90">
			{square.user_initials || 'Taken'}
		</span>
	{:else if isPending}
		<span class="animate-pulse text-[10px] font-black tracking-tighter uppercase">
			{userInitials}
		</span>
	{:else if !isLocked}
		<span class="text-xl font-light text-slate-300 opacity-0 group-hover:opacity-100">+</span>
	{/if}
</button> -->
