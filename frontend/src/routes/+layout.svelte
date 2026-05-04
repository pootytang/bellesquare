<script lang="ts">
	import './layout.css';
	import Navbar from '$lib/components/Navbar.svelte';
	import { auth } from '$lib/auth.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state'; // Svelte 5 reactive page state
	import { resolve } from '$app/paths';
	import { untrack } from 'svelte';

	let { data, children } = $props();

	// 2. Reactive Guard - This was working
	$effect(() => {
		// We want to watch 'data.user'. Whenever SvelteKit
		// provides new data, this block triggers.
		const newUser = data?.user ?? null;

		// We re-sync in case data.user changes during client-side nav
		// auth.setUser(data?.user ?? null);
		untrack(() => {
			auth.setUser(newUser);
		});

		// The "Bouncer" Logic:
		// If we have finished initializing, but have no user,
		// and we are sitting on a dashboard page... kick them out.
		const isDashboard = page.url.pathname.startsWith(resolve('/dashboard'));
		if (auth.initialized && !auth.user && isDashboard) {
			goto(resolve('/login'));
		}
	});
</script>

<div class="min-h-screen bg-white">
	<Navbar serverUser={data?.user ?? null} />

	<main>
		{@render children()}
	</main>

	<footer class="mt-20 border-t border-slate-100 py-10 text-center text-sm text-slate-400">
		<p>&copy; 2026 Bellesquare. Built with Go & Svelte 5.</p>
	</footer>
</div>
