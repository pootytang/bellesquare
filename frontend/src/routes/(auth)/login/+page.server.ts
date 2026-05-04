import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { apiRoute } from '$lib/api';

export const load: PageServerLoad = async ({ parent }) => {
    console.info("Login page load: checking if user is already authenticated...");
    const { user } = await parent();
    if (user) {
        console.debug("Login->PageServerLoad: User is already authenticated, redirecting to dashboard...");
        throw redirect(303, '/dashboard');
    }

    console.debug("Login->PageServerLoad: No authenticated user, rendering login page");
    return {};
};

export const actions: Actions = {
    default: async ({ request, fetch, cookies, url }) => {
        console.info("Login action triggered, processing form data...");

        const formData = await request.formData();
        const redirectTo = url.searchParams.get('redirectTo');

        console.debug(`redirectTo: ${redirectTo}`)

        const email = formData.get('email');
        const password = formData.get('password');

        console.debug("Validating email and password:");
        if (!email || !password) {
            console.error("Login failed: Missing email or password");
            return fail(400, { message: 'Email and password are required' });
        }

        // Convert FormData to URLSearchParams for compatibility
        console.debug("Converting FormData to URLSearchParams for backend compatibility...");
        const params = new URLSearchParams();
        for (const [key, value] of formData) {
            params.append(key, value.toString());
        }

        // Forward the login request to Go
        console.info("Sending login request to Go backend...");
        const res = await fetch(apiRoute('/auth/login'), {
            method: 'POST',
            headers: {'Content-Type': 'application/x-www-form-urlencoded'},
            body: params.toString()
        });

        console.debug("Received response from Go backend with status:", res.status);
        if (res.status === 401) {
            console.warn("Login failed: 401 Unauthorized return from backend");
            return fail(401, { message: 'Invalid email or password' });
        }

        if (!res.ok) {
            console.error("Login failed: server returned status", res.status);
            return fail(500, { message: 'Server error. Please try again later.' });
        }

        console.debug("Retrieving token from response body...");
        const data = await res.json(); // { "token": "..." }
        const token = data.token;

        console.info("Received token:", token);
        console.debug("Setting 'at' cookie in the browser with the received token...");
        if (token) {
            cookies.set('at', token, {
                path: '/',
                httpOnly: true,
                sameSite: 'lax',
                // In SvelteKit, 'dev' is a built-in import to check environment
                // import { dev } from '$app/environment';
                secure: process.env.NODE_ENV === 'production', 
                maxAge: 60 * 60 * 24 * 3 // 3 days
            });
        } else {
            console.error("Login failed: Token not received");
            return fail(500, { message: 'Unable to retrieve authentication token.' });
        }

        console.info("Login successful, redirecting to dashboard...");
        throw redirect(302, redirectTo || '/dashboard');
    }
};