<script lang="ts">
    import { auth } from '$lib/auth.svelte';
    import { resolve } from '$app/paths';
    import { enhance } from '$app/forms';
    import { invalidateAll } from '$app/navigation';

    // Receive the user from the layout as a fallback
    let { serverUser } = $props();

    // Use a derived value: check the global auth state first, 
    // then fall back to the server data
    let user = $derived(auth.user ?? serverUser ?? null);
</script>

<nav class="border-b border-slate-100 bg-white/80 backdrop-blur-md sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
        <a href={resolve('/')} class="text-xl font-black tracking-tighter text-slate-900">
            BELLE<span class="text-blue-600">SQUARE</span>
        </a>

        <div class="flex items-center gap-6">
            {#if user}
                <div class="flex items-center gap-5">
                    <div class="flex flex-col items-end">
                        <!-- <span class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Account</span> -->
                        <a href={resolve('/dashboard/profile')} class="text-sm font-semibold text-slate-700 hover:text-blue-600 transition-colors">
                            {user.user_name}
                        </a>
                    </div>

                    <div class="h-6 w-px bg-slate-200"></div>

                    <a href={resolve('/dashboard')} class="text-sm font-bold text-slate-900 hover:text-blue-600 transition-colors">
                        Dashboard
                    </a>

                    <form action="/logout" method="POST" use:enhance={() => {
                        return async ({update}) => {
                            auth.logout();
                            await invalidateAll();
                            await update();
                        };
                    }}>
                        <button type="submit" class="text-sm font-semibold text-red-600 hover:text-red-700 transition-colors">
                            Logout
                        </button>
                    </form>
                </div>
            {:else}
                <a href={resolve('/login')} class="text-sm font-bold text-slate-600 hover:text-slate-900">
                    Sign In
                </a>
                <a href={resolve('/register')} class="bg-slate-900 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-slate-800 transition-all">
                    Get Started
                </a>
            {/if}
        </div>
    </div>
</nav>

<!-- THIS WAS WORKING NICE BUT ADDING THE DERIVED USER TO TEST
 <nav class="border-b border-slate-100 bg-white/80 backdrop-blur-md sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
        <a href={resolve('/')} class="text-xl font-black tracking-tighter text-slate-900">
            BELLE<span class="text-blue-600">SQUARE</span>
        </a>

        <div class="flex items-center gap-6">
            {#if auth.user}
                <div class="flex items-center gap-5">
                    <div class="flex flex-col items-end">
                        <span class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Account</span>
                        <a href={resolve('/dashboard/profile')} class="text-sm font-semibold text-slate-700 hover:text-blue-600 transition-colors">
                            {auth.user.user_name}
                        </a>
                    </div>

                    <div class="h-6 w-px bg-slate-200"></div>

                    <a href={resolve('/dashboard')} class="text-sm font-bold text-slate-900 hover:text-blue-600 transition-colors">
                        Dashboard
                    </a>

                    <form action="/logout" method="POST" use:enhance={() => {
                        return async ({update}) => {
                            auth.logout();
                            await invalidateAll();
                            await update();
                        };
                    }}>
                        <button type="submit" class="text-sm font-semibold text-red-600 hover:text-red-700 transition-colors">
                            Logout
                        </button>
                    </form>
                </div>
            {:else}
                <a href={resolve('/login')} class="text-sm font-bold text-slate-600 hover:text-slate-900">
                    Sign In
                </a>
                <a href={resolve('/register')} class="bg-slate-900 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-slate-800 transition-all">
                    Get Started
                </a>
            {/if}
        </div>
    </div>
</nav> -->
