import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { apiRoute } from '$lib/api';

export const actions: Actions = {
    default: async ({ request, fetch, url }) => {
        console.debug("Register action triggered, processing form data...");

        const formData = await request.formData();
        
        // Extract fields for validation
        const username = formData.get('user_name');
        const firstName = formData.get('first_name');
        const lastName = formData.get('last_name');
        const email = formData.get('email');
        const password = formData.get('password');
        const confirmPassword = formData.get('confirm_password');

        // Basic client-side validation check
        if (!username || !firstName || !lastName || !email || !password) {
            console.error("Registration failed: All fields are required");
            return fail(400, { message: 'All fields are required' });
        }

        // Client-side validation in the server action
        console.debug("Validating form password...");
        if (password !== confirmPassword) {
            console.error("Registration failed: Passwords do not match");
            return fail(400, { message: 'Passwords do not match' });
        }

        // Convert FormData to URLSearchParams for compatibility
        const params = new URLSearchParams();
        for (const [key, value] of formData) {
            params.append(key, value.toString());
        }

        // Send all fields (including preferred_color) to the backend
        console.debug("Sending registration request to Go backend...");
        const res = await fetch(apiRoute('/auth/register'), {
            method: 'POST',
            headers: {'Content-Type': 'application/x-www-form-urlencoded'},
            body: params.toString()
        });

        if (!res.ok) {
            const errorData = await res.json();
            console.error("Registration failed:", errorData.error);
            return fail(res.status, { 
                message: errorData.error || 'Registration failed' 
            });
        }

        // After successful registration, send them to login
        // Alternatively, if your Go backend logs them in immediately, 
        // you would handle the cookie here.
        console.info("Registration successful, redirecting to login...");
        const redirectTo = url.searchParams.get('redirectTo') || '/dashboard';
        
        console.debug(`RedirectTo from register: ${redirectTo}`)
        const loginUrl = redirectTo 
            ? `/login?redirectTo=${encodeURIComponent(redirectTo)}` 
            : '/login';
        
        throw redirect(303, loginUrl);
    }
};