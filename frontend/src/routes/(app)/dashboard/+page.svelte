<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
	import DashboardTile from '$lib/components/DashboardTile.svelte';
	import BoardActionModal from '$lib/components/boards/BoardActionModal.svelte';

	let displayName = $derived(auth.user?.full_name || auth.user?.user_name || 'User');
	let isBoardModalOpen = $state(false);

	// Identify guest status from your auth state
	let isGuest = $derived(auth.user?.is_guest);
</script>

<div class="mx-auto max-w-5xl p-6">
	<header class="mb-10">
		<h1 class="text-4xl font-black tracking-tight text-slate-900">Welcome, {displayName}</h1>
		<p class="text-slate-600">What would you like to manage today?</p>
	</header>

	<!-- Marketing Banner for Guest Users -->
	{#if isGuest}
		<div
			class="mb-10 overflow-hidden rounded-3xl bg-linear-to-br from-indigo-600 to-blue-700 p-1 shadow-xl"
		>
			<div
				class="flex flex-col items-center justify-between gap-6 rounded-[calc(1.5rem-1px)] bg-white/5 p-8 text-white md:flex-row"
			>
				<div class="flex-1 space-y-2 text-center md:text-left">
					<h2 class="text-2xl font-bold">Host your own Big Game board!</h2>
					<p class="text-indigo-100">
						Upgrade to a full account to create custom grids, set payouts, and manage your own
						players.
					</p>
				</div>
				<a
					href={resolve('/dashboard/upgrade')}
					class="rounded-xl bg-white px-8 py-4 font-black whitespace-nowrap text-indigo-600 shadow-lg transition-all hover:scale-105 active:scale-95"
				>
					UPGRADE NOW
				</a>
			</div>
		</div>
	{/if}

	<div
		class="grid grid-cols-1 gap-6 md:grid-cols-2 {isGuest
			? 'lg:flex lg:justify-center'
			: 'lg:grid-cols-3'}"
	>
		{#if isGuest}
			<!-- GUEST VIEW: Just show Squares and potentially a "Joined Boards" tile -->
			<div class="relative opacity-80">
				<div
					aria-hidden="true"
					class="absolute -top-2 -right-2 z-10 rounded-lg bg-amber-400 px-2 py-1 text-[10px] font-black tracking-widest text-white uppercase shadow-md"
				>
					Member Only
				</div>
				<a
					href={resolve('/dashboard/upgrade')}
					class="block transition-transform hover:scale-[1.02]"
				>
					<DashboardTile
						title="Create a Board"
						description="Upgrade your account to start hosting your own grids."
						iconColor="slate"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-6 w-6"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
							/>
						</svg>
					</DashboardTile>
				</a>
			</div>

			<DashboardTile
				href="/dashboard/boards/joined"
				title="Joined Boards"
				description="View the grids you are currently playing in."
				iconColor="blue"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-6 w-6"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
					/>
				</svg>
			</DashboardTile>
		{:else}
			<!-- MEMBER VIEW: Full access, no extra wrappers to break the layout -->
			<button
				onclick={() => (isBoardModalOpen = true)}
				class="text-left transition-transform outline-none hover:scale-[1.02] active:scale-95"
			>
				<DashboardTile
					title="Squares Boards"
					description="Create new grids or manage your existing games."
					iconColor="purple"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-6 w-6"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"
						/>
					</svg>
				</DashboardTile>
			</button>

			<DashboardTile
				href="/dashboard/teams"
				title="Teams"
				description="Add and update mascots, cities, and logos."
				iconColor="green"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-6 w-6"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
					/>
				</svg>
			</DashboardTile>

			<DashboardTile
				href="/dashboard/profile"
				title="Profile"
				description="Manage your account and payment handles."
				iconColor="blue"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-6 w-6"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
					/>
				</svg>
			</DashboardTile>
		{/if}
	</div>

	{#if isBoardModalOpen}
		<BoardActionModal onclose={() => (isBoardModalOpen = false)} />
	{/if}
</div>
