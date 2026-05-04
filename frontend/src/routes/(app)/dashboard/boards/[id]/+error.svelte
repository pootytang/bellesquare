<script>
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
</script>

<div class="flex min-h-[60vh] flex-col items-center justify-center p-6 text-center">
	<h1 class="text-4xl font-black text-slate-900">404</h1>
	<p class="mt-2 font-bold text-slate-500">{page.error?.message || 'Board not found'}</p>
	<a
		href={resolve('/dashboard')}
		class="mt-6 rounded-xl bg-slate-900 px-6 py-3 font-bold text-white"
	>
		Back to Dashboard
	</a>
</div>
