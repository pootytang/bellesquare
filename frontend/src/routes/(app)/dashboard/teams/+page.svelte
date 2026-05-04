<script lang="ts">
	import type { PageData } from './$types';
	import type { Team, ActionData } from '$lib/types';
	import { resolve } from '$app/paths';
	import AddTeamForm from '$lib/components/teams/AddTeamForm.svelte';
	import TeamCard from '$lib/components/teams/TeamCard.svelte';
	import EditTeamModal from '$lib/components/teams/EditTeamModal.svelte';
	import DeleteConfirmModal from '$lib/components/teams/DeleteConfirmModal.svelte';

	let teamToDelete = $state<Team | null>(null);

	let { data }: { data: PageData } = $props();

	let showNotification = $state(false);
	let localFormResult = $state<ActionData | null>(null);
	let selectedTeam = $state<Team | null>(null);

	function handleResult(result: ActionData) {
		localFormResult = result;
		showNotification = true;
		setTimeout(() => (showNotification = false), 5000);
	}
</script>

<div class="mx-auto max-w-6xl p-6">
	{#if showNotification && localFormResult}
		<div
			class="animate-in fade-in slide-in-from-top-4 fixed top-10 left-1/2 z-50 w-full max-w-sm -translate-x-1/2 transform px-4 duration-300"
		>
			<div
				class="flex items-center gap-3 rounded-2xl border bg-white p-4 shadow-xl {localFormResult.success
					? 'border-green-200 text-green-800'
					: 'border-red-200 text-red-800'}"
			>
				<div
					class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full {localFormResult.success
						? 'bg-green-500'
						: 'bg-red-500'} text-white"
				>
					{#if localFormResult.success}
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-5 w-5"
							viewBox="0 0 20 20"
							fill="currentColor"
						>
							<path
								fill-rule="evenodd"
								d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
								clip-rule="evenodd"
							/>
						</svg>
					{:else}
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-5 w-5"
							viewBox="0 0 20 20"
							fill="currentColor"
						>
							<path
								fill-rule="evenodd"
								d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
								clip-rule="evenodd"
							/>
						</svg>
					{/if}
				</div>
				<div class="flex flex-col">
					<span class="text-sm font-bold"
						>{localFormResult.success
							? localFormResult.team
								? `Saved ${localFormResult.team.full_name}!`
								: 'Success!'
							: 'Error'}</span
					>
					<span class="text-xs opacity-80"
						>{localFormResult.message ||
							(localFormResult.success ? 'The roster has been updated.' : 'Check your data.')}</span
					>
				</div>
			</div>
		</div>
	{/if}

	<header class="mb-10">
		<h1 class="text-3xl font-black tracking-tight text-slate-900">Team Management</h1>
		<p class="text-slate-500">Add new teams or view your current teams list.</p>
	</header>

	<div class="grid grid-cols-1 items-start gap-12 lg:grid-cols-12">
		<div class="lg:col-span-4">
			<AddTeamForm onResult={handleResult} />
		</div>

		<div class="lg:col-span-8">
			<div class="mb-6 flex items-center justify-between">
				<h2 class="text-xl font-bold text-slate-800">Current Teams</h2>

				<div class="flex gap-2 rounded-2xl bg-slate-100 p-1">
					{#each ['football', 'basketball', 'baseball', 'hockey', 'soccer'] as sport (sport)}
						<a
							href={resolve(`/dashboard/teams?sport=${sport}`)}
							class="rounded-xl px-4 py-1.5 text-xs font-black tracking-widest uppercase transition-all
                            {data.activeSport === sport
								? 'bg-white text-blue-600 shadow-sm'
								: 'text-slate-400 hover:text-slate-600'}"
						>
							{sport}
						</a>
					{/each}
				</div>
			</div>

			<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
				{#each data.teams as team (team.id)}
					<TeamCard
						{team}
						onclick={() => (selectedTeam = team)}
						ondelete={() => (teamToDelete = team)}
					/>
				{:else}
					<div
						class="col-span-2 flex flex-col items-center justify-center rounded-4xl border border-dashed border-slate-200 bg-slate-50 py-20"
					>
						<div class="mb-3 rounded-full bg-slate-100 p-4 text-slate-300">
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
									d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
								/>
							</svg>
						</div>
						<p class="italic text-slate-400">No {data.activeSport} teams found.</p>
					</div>
				{/each}
				{#if teamToDelete}
					<DeleteConfirmModal
						team={teamToDelete}
						onResult={handleResult}
						onclose={() => (teamToDelete = null)}
					/>
				{/if}
			</div>
		</div>
	</div>
</div>

{#if selectedTeam}
	<EditTeamModal
		team={selectedTeam}
		onResult={handleResult}
		onclose={() => (selectedTeam = null)}
	/>
{/if}
