import { fail, redirect } from '@sveltejs/kit';
import { apiRoute } from '$lib/api';
import type { Actions } from './$types';

export const actions: Actions = {
    update: async ({ request, fetch, cookies }) => {
        console.info('Processing profile update...');
        const formData = await request.formData();
        const token = cookies.get('at');

        // 1. Validate token presence
        if (!token) {
            console.warn('No authentication token found. User may not be logged in.');
            throw redirect(303, '/login');
        }
        console.debug('Authentication token found:', token);
        // for (const pair of formData.entries()) {
        //     console.log(pair[0], pair[1]);
        // }

        // 2. Prepare the data for the backend
        const profileData = {
            first_name: formData.get('first_name'),
            last_name: formData.get('last_name'),
            venmo_handle: formData.get('venmo_handle'),
            zelle_handle: formData.get('zelle_handle'),
            preferred_color: formData.get('preferred_color'),
        };

        console.debug('Profile data to update:', profileData);
        const res = await fetch(apiRoute('/api/v1/users/profile'), {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(profileData)
        });

        console.debug('Backend response status:', res.status);
        if (!res.ok) {
            console.error('Profile update failed:', res.status);
            return fail(res.status, { message: 'Update failed' });
        }

        console.info('Profile updated successfully.');
        return { success: true };
    }
};