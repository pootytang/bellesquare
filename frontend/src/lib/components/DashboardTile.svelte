<script lang="ts">
	import type { Snippet } from 'svelte';
	import { resolve } from '$app/paths';
	import type { Pathname } from '$app/types';

	type TileColor = 'blue' | 'green' | 'slate' | 'purple' | 'red';

	interface Props {
		href?: Pathname;
		title: string;
		description: string;
		iconColor?: TileColor;
		children?: Snippet;
	}

	let { href, title, description, iconColor = 'blue', children }: Props = $props();

	const colorClasses: Record<TileColor, string> = {
		blue: 'bg-blue-50 group-hover:bg-blue-100 text-blue-600',
		green: 'bg-green-50 group-hover:bg-green-100 text-green-600',
		slate: 'bg-slate-50 group-hover:bg-slate-100 text-slate-600',
		purple: 'bg-purple-50 group-hover:bg-purple-100 text-purple-600',
		red: 'bg-red-50 group-hover:bg-red-100 text-red-600'
	};

	// Helper to determine the hover border class
	const borderHoverClasses: Record<TileColor, string> = {
		blue: 'hover:border-blue-500',
		green: 'hover:border-green-500',
		slate: 'hover:border-slate-500',
		purple: 'hover:border-purple-500',
		red: 'hover:border-red-500'
	};

	// Only resolve if href is actually passed
	const resolvedHref = $derived(href ? resolve(href) : undefined);
</script>

<svelte:element
	this={href ? 'a' : 'div'}
	href={resolvedHref}
	class="group rounded-2xl border border-slate-200 bg-white p-6 shadow-sm {borderHoverClasses[
		iconColor
	]} text-left transition-all hover:shadow-md"
>
	<div class="mb-4 flex items-center justify-between">
		<!-- Add the if-guard here -->
		{#if children}
			<div class="rounded-lg p-3 transition-colors {colorClasses[iconColor]}">
				{@render children()}
			</div>
		{/if}

		{#if href}
			<span class="text-{iconColor}-600 opacity-0 transition-opacity group-hover:opacity-100"
				>Manage →</span
			>
		{/if}
	</div>

	<h3 class="text-lg font-bold text-slate-900">{title}</h3>
	<p class="text-sm text-slate-500">{description}</p>
</svelte:element>
