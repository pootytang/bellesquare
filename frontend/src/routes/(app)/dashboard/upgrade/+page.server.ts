import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { apiRoute } from '$lib/api';

export const load: PageServerLoad = async ({ parent }) => {
    console.info('********** UPGRADE LOAD FUNCTION **********');
    const { user } = await parent();
    // If they are already a member, they don't need to be here
    if (user && !user.is_guest) {
        console.info('User is already a member, redirecting to dashboard');
        throw redirect(303, '/dashboard');
    }

    console.info('User is a guest, showing upgrade page');
    return {};
};

export const actions: Actions = {
    default: async ({ request, fetch, locals }) => {
        console.info('********** UPGRADE ACTION **********');

        console.info('Processing upgrade form submission');
        const formData = await request.formData();
        const password = formData.get('password');
        const confirmPassword = formData.get('confirm_password');

        // Basic client-side validation check
        console.debug('Validating form data');
        if (password !== confirmPassword) {
            console.error('Password validation failed: Passwords do not match.');
            return fail(400, { message: 'Passwords do not match.' });
        }

        if (!password || password.toString().length < 8) {
            console.error('Password validation failed: Password must be at least 8 characters.');
            return fail(400, { message: 'Password must be at least 8 characters.' });
        }

        // Forward to the Go Upgrade handler
        // Using JSON because your new UpgradeRequest DTO expects JSON
        const route = apiRoute('/api/v1/users/upgrade');
        console.info(`Sending upgrade request to API at ${route}`);
        const res = await fetch(route, {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${locals.token}`
            },
            body: JSON.stringify({ password })
        });

        if (!res.ok) {
            const errorData = await res.json();
            console.error('Upgrade failed:', errorData);
            return fail(res.status, { 
                message: errorData.error || 'Upgrade failed. Please try again.' 
            });
        }

        // Success! Redirect to dashboard. 
        // The user object in the 'parent' layout will refresh on the next load.
        console.info('Upgrade successful, redirecting to dashboard');
        throw redirect(303, '/dashboard');
    }
};