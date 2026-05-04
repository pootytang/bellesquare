import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ parent }) => {
    // 1. Get the layout data which contains the user from the backend
    console.info('Loading dashboard page... Checking user authentication.');
    const { user } = await parent();

    // 2. Security Check: No user? Go back to login.
    if (!user) {
        console.warn('No authenticated user found. Redirecting to login page.');
        throw redirect(303, '/login');
    }

    // 3. User exists, proceed to load dashboard-specific data
    console.info('User authenticated. Loading dashboard data.');
    return {
        user
    };
};