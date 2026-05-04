<script lang="ts">
	import type { Team } from '$lib/types';
	let { team, onclick, ondelete }: { team: Team; onclick: () => void; ondelete: () => void } =
		$props();
</script>

<div class="group relative">
	<button
		type="button"
		{onclick}
		class="flex w-full cursor-pointer items-center gap-4 overflow-hidden rounded-2xl border border-slate-200 bg-white p-4 text-left shadow-sm transition-all hover:border-slate-300 hover:shadow-md active:scale-[0.98]"
	>
		<div
			class="absolute top-0 left-0 h-1 w-full"
			style="background: linear-gradient(90deg, {team.primary_color}, {team.secondary_color})"
		></div>

		{#if team.team_logo_url}
			<div
				class="flex h-14 w-14 shrink-0 items-center justify-center rounded-xl p-2 transition-transform group-hover:scale-110"
				style="background-color: {team.primary_color}15"
			>
				<img
					src={team.team_logo_url}
					alt={team.full_name}
					class="h-full w-full object-contain drop-shadow-sm"
				/>
			</div>
		{:else}
			<div
				class="flex h-14 w-14 shrink-0 items-center justify-center rounded-xl font-black text-white uppercase shadow-inner"
				style="background: linear-gradient(135deg, {team.primary_color}, {team.secondary_color})"
			>
				{team.city[0]}
			</div>
		{/if}

		<div class="flex-1">
			<h3 class="font-black tracking-tight text-slate-900" style="color: {team.primary_color}">
				{team.full_name}
			</h3>
			<p class="text-[10px] font-black tracking-widest text-slate-400 uppercase">
				{team.sport}
			</p>
		</div>

		<div class="text-slate-300 group-hover:text-slate-500">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="h-4 w-4"
				fill="none"
				viewBox="0 0 24 24"
				stroke="currentColor"
			>
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M9 5l7 7-7 7" />
			</svg>
		</div>
	</button>

	<button
		type="button"
		onclick={(e) => {
			e.stopPropagation();
			ondelete();
		}}
		class="absolute -top-2 -right-2 z-10 flex h-8 w-8 cursor-pointer items-center justify-center rounded-full border border-red-100 bg-white text-red-500 shadow-lg transition-all hover:bg-red-500 hover:text-white lg:opacity-0 lg:group-hover:opacity-100"
		aria-label="Delete team"
	>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			class="h-4 w-4"
			fill="none"
			viewBox="0 0 24 24"
			stroke="currentColor"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M6 18L18 6M6 6l12 12"
			/>
		</svg>
	</button>
</div>

<!-- <script lang="ts">
	import type { Team } from '$lib/types';
	let { team, onclick, ondelete }: { team: Team; onclick: () => void; ondelete: () => void } =
		$props();
</script>

<div class="group relative">
	<button
		type="button"
		{onclick}
		class="flex w-full cursor-pointer items-center gap-4 rounded-2xl border border-slate-200 bg-white p-4 text-left shadow-sm transition-all hover:border-blue-500 hover:shadow-md"
	>
		{#if team.team_logo_url}
			<img
				src={team.team_logo_url}
				alt=""
				class="h-12 w-12 rounded-full bg-slate-50 object-contain p-1"
			/>
		{:else}
			<div
				class="flex h-12 w-12 items-center justify-center rounded-full bg-slate-100 font-bold text-slate-400 uppercase"
			>
				{team.city[0]}
			</div>
		{/if}

		<div class="flex-1">
			<h3 class="font-bold text-slate-900">{team.full_name}</h3>
			<p class="text-[10px] font-black tracking-widest text-slate-400 uppercase">{team.sport}</p>
		</div>
	</button>

	<button
		type="button"
		onclick={(e) => {
			e.stopPropagation();
			ondelete();
		}}
		class="absolute -top-2 -right-2 flex h-8 w-8 cursor-pointer items-center justify-center rounded-full border border-red-100 bg-white text-red-500 shadow-sm transition-all hover:bg-red-500 hover:text-white lg:opacity-0 lg:group-hover:opacity-100"
		aria-label="Delete team"
	>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			class="h-4 w-4"
			fill="none"
			viewBox="0 0 24 24"
			stroke="currentColor"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M6 18L18 6M6 6l12 12"
			/>
		</svg>
	</button>
</div> -->
