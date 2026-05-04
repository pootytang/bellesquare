import { apiRoute } from '$lib/api';
import { redirect, fail } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';

export const load: PageServerLoad = async ({ url, fetch, locals }) => {
    console.info('Grabbing teams for the board');
    const sport = url.searchParams.get('sport');

    // If the user tries to access /new without ?sport=..., send them back
    if (!sport) {
        console.warn("Missing team parameter")
        throw redirect(303, '/dashboard/boards');
    }

    console.debug("fetching the teams")
    const res = await fetch(apiRoute(`/api/v1/teams?sport=${sport}`), {
        headers: {
            'Authorization': `Bearer ${locals.token}`
        }
    });

    if (!res.ok) {
        // Handle API errors (e.g., token expired or backend down)
        console.error("Fetch failed with a status of: " + res.status)
        return { teams: [], activeSport: sport, error: 'Could not load teams' };
    }

    const teams = await res.json();
    
    console.info("Successfully retrieved teams: " + teams)
    return {
        teams: teams || [],
        noTeams: !teams || teams.length === 0,
        activeSport: sport
    };
};

export const actions: Actions = {
    create: async ({ request, fetch, locals }) => {
        console.info("Creating a board")
        const formData = await request.formData();
        
        const payload = {
            title: formData.get('name'), // maps to board.Title
            sport: formData.get('sport'),
            price_per_square: parseFloat(formData.get('price_per_square') as string) || 0,
            home_team_id: formData.get('home_team_id'),
            away_team_id: formData.get('away_team_id'),
            status: 'draft', // Initial models.BoardStatus
            // Payout amounts mapping
            payouts: {
                "q1": parseFloat(formData.get('payout_q1') as string) || 0,
                "q2": parseFloat(formData.get('payout_q2') as string) || 0,
                "q3": parseFloat(formData.get('payout_q3') as string) || 0,
                "q4": parseFloat(formData.get('payout_q4') as string) || 0,
                "final": parseFloat(formData.get('payout_final') as string) || 0
            }
        };

        console.info("Requesting the boards endpoint")
        const res = await fetch(apiRoute('/api/v1/boards/new'), {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${locals.token}`
            },
            body: JSON.stringify(payload)
        });

        if (!res.ok) {
            console.error("Fetch failed with a status of " + res.status)
            const err = await res.json();
            return fail(res.status, { success: false, message: err.error });
        }

        console.info("setting up the board")
        const newBoard = await res.json();
        
        // Redirect to the newly created board
        console.info("Redirecting to the newly created board")
        console.debug(`boardId: ${newBoard.id}`)
        // throw redirect(303, `/dashboard/boards/${newBoard.id}`);
        return { success: true, boardId: newBoard.id };
    }
};