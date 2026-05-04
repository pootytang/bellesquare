<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
</script>

<div class="min-h-[80vh] flex flex-col items-center justify-center px-4 text-center">
	<h1 class="text-5xl md:text-7xl font-black text-slate-900 tracking-tight mb-4">
		BELLE<span class="text-blue-600">SQUARE</span>
	</h1>
	
	<p class="text-lg md:text-xl text-slate-500 max-w-2xl mb-10 leading-relaxed">
		The ultimate 10x10 grid management system. Claim your squares, 
		track your wins, and manage your boards with ease.
	</p>

	<div class="flex flex-col sm:flex-row gap-4 w-full max-w-md justify-center">
		{#if auth.user}
			<a 
				href={resolve("/dashboard")} 
				class="bg-slate-900 text-white px-8 py-4 rounded-2xl font-bold text-lg hover:bg-slate-800 transition-all shadow-xl shadow-slate-200"
			>
				Go to Dashboard
			</a>
		{:else}
			<a 
				href={resolve("/login")} 
				class="bg-blue-600 text-white px-8 py-4 rounded-2xl font-bold text-lg hover:bg-blue-700 transition-all shadow-xl shadow-blue-100"
			>
				Sign In
			</a>
			<a 
				href={resolve("/register")} 
				class="bg-white text-slate-900 border-2 border-slate-200 px-8 py-4 rounded-2xl font-bold text-lg hover:border-slate-900 transition-all"
			>
				Create Account
			</a>
		{/if}
	</div>

	<div class="mt-20 grid grid-cols-1 md:grid-cols-3 gap-8 text-left max-w-5xl">
		<div class="p-6 bg-slate-50 rounded-2xl">
			<div class="text-blue-600 mb-2">●</div>
			<h3 class="font-bold text-slate-900 mb-1">Fast Performance</h3>
			<p class="text-sm text-slate-500">Built with Svelte 5 and Go for near-instant interaction.</p>
		</div>
		<div class="p-6 bg-slate-50 rounded-2xl">
			<div class="text-blue-600 mb-2">●</div>
			<h3 class="font-bold text-slate-900 mb-1">Secure Auth</h3>
			<p class="text-sm text-slate-500">Industrial-grade Argon2 hashing and HttpOnly JWT cookies.</p>
		</div>
		<div class="p-6 bg-slate-50 rounded-2xl">
			<div class="text-blue-600 mb-2">●</div>
			<h3 class="font-bold text-slate-900 mb-1">Real-time Updates</h3>
			<p class="text-sm text-slate-500">Keep track of every square claim the moment it happens.</p>
		</div>
	</div>
</div>
