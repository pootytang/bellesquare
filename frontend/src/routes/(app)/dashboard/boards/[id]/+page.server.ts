import { apiRoute } from '$lib/api';
import { error, fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, fetch, locals }) => {
    console.info("********** LOADING THE BOARD **********")

    console.debug("Fetching board id: " + params.id)
    const res = await fetch(apiRoute(`/api/v1/boards/${params.id}`), {
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
        console.error("unable to fetch the board with id: " + params.id)
        throw error(res.status, 'Could not fetch board');
    }

    console.info("board retrieved, grabbing the data")
    const boardData = await res.json();

    console.info("returning board and team data")
    console.debug(`board id ${boardData.board.id}`)
    // console.debug("HOME TEAM DATA")
    // console.table(homeTeamData)
    // console.debug("AWAY TEAM DATA")
    // console.table(awayTeamData)
    // console.debug(boardData.squares)
    // console.table(boardData.board)
    // console.table(boardData.scores.quarters)
    // console.debug("PAYOUTS")
    // console.table(boardData.creator_zelle)
    return {
        board: boardData.board,     // The board metadata
        squares: boardData.squares, // The 100 squares
        homeTeam: boardData.board.home_team,
        awayTeam: boardData.board.away_team,
        scores: boardData.scores ?? { home: 0, away: 0, quarters: [] },
        payouts: boardData.payouts, // The payouts
        square_owners: boardData.square_owners || {},
        isCreator: boardData.board.creator_id === locals.user?.id,
        creator_venmo: boardData.creator_venmo || '',
        creator_zelle: boardData.creator_zelle || '',
        token: locals.token,
        user: locals.user
    };
};

export const actions: Actions = {
    claim: async ({ params, request, fetch, locals }) => {
        console.info("The Square Claim action called")

        console.debug("Checking if there's a user")
        if (!locals.user || !locals.token) {
            console.warn("No user found, redirecting to login")
            throw redirect(303, '/login');
        }

        console.info(`User found: ${locals.user.user_name}. Grabbing form data`)
        const formData = await request.formData();
        const color = formData.get('color') as string;
        const selectionsRaw = formData.get('selections') as string;

        // 1. Parse JSON and handle error in one go
        console.debug("Parsing json")
        let selectionKeys: string[];
        try {
            selectionKeys = JSON.parse(selectionsRaw);
            console.debug(`parsed the selection keys: ${selectionKeys}`)
        } catch (e) {
            console.error("Problem parsing selection keys: ", e);
            return fail(400, { message: 'Invalid selection format' });
        }

        // 2. Validate length
        console.info("Performing validation")
        if (!selectionKeys || selectionKeys.length === 0) {
            console.warn("Validation failed returning a 400")
            return fail(400, { message: 'No squares selected' });
        }

        // 3. Transform
        console.debug("Transorming the selections to something the backend can use")
        const selections = selectionKeys.map((key) => {
            const parts = key.split('-');
            return { row: Number(parts[0]), col: Number(parts[1]) };
        }).filter(s => !isNaN(s.row) && !isNaN(s.col));

        console.info("fetching /api/v1/boards/squares/:id/claim endpoint")
        const res = await fetch(apiRoute(`/api/v1/boards/${params.id}/squares/claim`), {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${locals.token}`
            },
            body: JSON.stringify({
                user_id: locals.user.id,
                user_color: color,
                selections: selections
            })
        });

        if (!res.ok) {
            const errData = await res.json();
            console.error(`bad response from backend: ${errData.error}`)
            return fail(res.status, {                 
                message: errData.error || 'Failed to claim squares' 
            });
        }

        console.info("Successfully updated the selection")
        return { success: true };
    },
    shuffle: async ({ params, fetch, locals }) => {
        console.info(`Shuffling numbers for board ${params.id}...`);
        const res = await fetch(apiRoute(`/api/v1/boards/${params.id}/shuffle`), {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${locals.token}`
            }
        });

        if (!res.ok) {
            console.error(`Failed to shuffle numbers for board ${params.id}`);
            return fail(500, { message: "Could not shuffle numbers" });
        }

        console.info('Numbers shuffled successfully.');
        return { success: true };
    },
    lock: async ({ params, fetch, locals }) => {
        console.info("---------- LOCK ACTION CALLED ----------")
        console.info(`Locking board ${params.id}...`);
        // This endpoint in your Go API sets status to 'locked'
        const res = await fetch(apiRoute(`/api/v1/boards/${params.id}/status`), {
            method: 'PATCH', 
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${locals.token}`
            },
            body: JSON.stringify({ status: 'locked' }) // Explicitly send 'locked'
        });
        
        if (!res.ok) {
            console.error(`Failed to lock board ${params.id}`);
            return fail(500, { message: "Could not lock board" });
        }

        console.info('Board locked successfully.');
        return { success: true };
    },
    unlock: async ({ params, locals, fetch }) => {
        console.info("---------- UNLOCK ACTION CALLED ----------")
        console.info(`unlocking board ${params.id}`)
        // Call the same update status endpoint but pass "open"

        console.debug(`fetching /api/v1/boards/${params.id}/status`)
        const res = await fetch(apiRoute(`/api/v1/boards/${params.id}/status`), {
            method: 'PATCH', // Or POST, depending on your Go route
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${locals.token}`
            },
            body: JSON.stringify({ status: 'open' })
        });

        if (!res.ok) {
            console.error(`Failed to unlock board ${params.id}`)
            return fail(res.status, { message: 'Failed to unlock' })
        };

        console.info("board is unlocked")
        return { success: true };
    },
    savePayouts: async ({ params, request, fetch, locals }) => {
        console.info("Saving Payouts")
        const formData = await request.formData();
        // Gather all Q1-Final amounts from the form
        const periods = ['q1', 'q2', 'q3', 'q4', 'final'];
        const payouts = periods.map(p => ({
            period_name: p,
            amount: parseFloat(formData.get(p) as string || '0')
        }));

        console.debug("Fetching " + apiRoute(`/api/v1/boards/${params.id}/payouts`) + " endpoint");
        const res = await fetch(apiRoute(`/api/v1/boards/${params.id}/payouts`), {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${locals.token}` 
            },
            body: JSON.stringify({ payouts })
        });

        console.info("responding")
        return res.ok ? { success: true } : fail(res.status);
    }
};