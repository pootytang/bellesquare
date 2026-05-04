import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    default: async ({ cookies }) => {
        console.info("Logout action triggered");
        // 1. Delete the 'at' cookie from the browser
        console.debug("Deleting 'at' cookie to log out user...");
        cookies.delete('at', { path: '/' });

        // 2. Redirect the user to the home page or login
        console.info("User logged out successfully, redirecting to login page...");
        throw redirect(303, '/login');
    }
};