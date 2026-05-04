import type { PageServerLoad } from './$types';
import { error } from '@sveltejs/kit';
import { apiRoute } from '$lib/api';

export const load: PageServerLoad = async ({ fetch, locals }) => {
    console.info("Grabbing all joined boards")
    const res = await fetch(apiRoute("/api/v1/boards/joined"), {
            headers: {
                'Authorization': `Bearer ${locals.token}`
            }
        });
    
        // Handle 404 or other errors immediately
        if (res.status === 404) {
            throw error(404, {
                message: 'This board does not exist or has been deleted.'
            });
        }
    
        if (!res.ok) {
            console.error("unable to fetch the boards the user has joined")
            throw error(res.status, 'Could not fetch board');
        }
    
        console.info("board retrieved, grabbing the data")
        const boardData = await res.json();

        console.info("returning boards")
        return {
        boards: boardData
        };
};