<script lang="ts">
	import { resolve } from '$app/paths';
	let { onclose }: { onclose: () => void } = $props();

	const sports = [
		{ name: 'football', icon: '🏈', color: 'hover:border-blue-500 hover:bg-blue-50' },
		{ name: 'basketball', icon: '🏀', color: 'hover:border-orange-500 hover:bg-orange-50' },
		{ name: 'baseball', icon: '⚾', color: 'hover:border-red-500 hover:bg-red-50' },
		{ name: 'soccer', icon: '⚽', color: 'hover:border-green-500 hover:bg-green-50' }
	];
</script>

<div class="fixed inset-0 z-100 flex items-center justify-center p-4">
	<button
		aria-label="Close modal"
		title="Close"
		class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm"
		onclick={onclose}
	></button>

	<div class="animate-in zoom-in-95 relative w-full max-w-lg rounded-4xl bg-white p-8 shadow-2xl">
		<header class="mb-8 flex items-center justify-between">
			<div>
				<h2 class="text-2xl font-black text-slate-900">Squares Boards</h2>
				<p class="text-sm text-slate-500">Pick a sport to start or view history.</p>
			</div>
			<button
				aria-label="Close modal"
				title="Close"
				onclick={onclose}
				class="rounded-full bg-slate-100 p-2 text-slate-400 hover:text-slate-600"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-5 w-5"
					viewBox="0 0 20 20"
					fill="currentColor"
				>
					<path
						fill-rule="evenodd"
						d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
						clip-rule="evenodd"
					/>
				</svg>
			</button>
		</header>

		<div class="space-y-6">
			<a
				href={resolve('/dashboard/boards')}
				class="flex w-full items-center justify-between rounded-2xl border-2 border-slate-100 p-4 transition-all hover:border-purple-500 hover:bg-purple-50"
			>
				<div class="flex items-center gap-4">
					<div
						class="flex h-10 w-10 items-center justify-center rounded-xl bg-purple-100 text-purple-600"
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
								d="M4 6h16M4 10h16M4 14h16M4 18h16"
							/>
						</svg>
					</div>
					<span class="font-bold text-slate-700">View My Existing Boards</span>
				</div>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-5 w-5 text-slate-300"
					viewBox="0 0 20 20"
					fill="currentColor"
				>
					<path
						fill-rule="evenodd"
						d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
						clip-rule="evenodd"
					/>
				</svg>
			</a>

			<div class="relative py-2">
				<div class="absolute inset-0 flex items-center">
					<span class="w-full border-t border-slate-100"></span>
				</div>
				<div class="relative flex justify-center text-xs tracking-widest text-slate-400 uppercase">
					<span class="bg-white px-2 italic">or Create New</span>
				</div>
			</div>

			<div class="grid grid-cols-2 gap-3">
				{#each sports as sport (sport)}
					<a
						href={resolve(`/dashboard/boards/new?sport=${sport.name}`)}
						class="flex flex-col items-center gap-2 rounded-2xl border-2 border-slate-100 p-4 transition-all {sport.color}"
					>
						<span class="text-2xl">{sport.icon}</span>
						<span class="text-xs font-black tracking-widest text-slate-600 uppercase"
							>{sport.name}</span
						>
					</a>
				{/each}
			</div>
		</div>
	</div>
</div>
