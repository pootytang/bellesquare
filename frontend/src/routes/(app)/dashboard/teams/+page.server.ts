import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import { apiRoute } from '$lib/api';

export const load: PageServerLoad = async ({ fetch, locals, url }) => {
    console.info('*****Fetching teams for selection*****');
    // Pull 'sport' from the browser URL (e.g., /dashboard/teams?sport=football)
    // Default to 'football' if no query param is present
    const sport = url.searchParams.get('sport') || 'football';
    console.debug(`Calling teams endpoint: Sport = ${sport} Token = ${locals.token}`);

    const res = await fetch(apiRoute(`/api/v1/teams?sport=${sport}`), {
        headers: {
            'Authorization': `Bearer ${locals.token}`
        }
    });
    
    console.debug('Received a response status of:', res.status);
    if (res.status === 401) {
		// If the token is expired/invalid, the Go backend returns 401.
		// We can handle that here or let the hook catch it on the next refresh.
		throw redirect(303, '/login');
	}
    
    if (!res.ok) {
        console.error('Failed to fetch teams');
        return { teams: [] }; // Return empty array on error to prevent crashes
    }

    console.info('Teams fetched successfully. Parsing response...');
    const teams = await res.json();

    console.info('Teams: ', teams);
    return { 
        teams, // This makes 'teams' available on the data prop
        activeSport: sport // This makes data.activeSport available in the Svelte file
    };
};

export const actions: Actions = {
    create: async ({ request, fetch, locals }) => {
        console.info('Processing new team creation...');
        const formData = await request.formData();
        
        const teamData = {
            city: formData.get('city'),
            mascot: formData.get('mascot'),
            sport: formData.get('sport'),
            team_logo_url: formData.get('team_logo_url') || null,
            primary_color: formData.get('primary_color'),
            secondary_color: formData.get('secondary_color')
        };

        console.debug('Team data to create:', teamData);
        const res = await fetch(apiRoute('/api/v1/teams/new'), {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${locals.token}`
            },
            body: JSON.stringify(teamData)
        });

        const result = await res.json();

        if (!res.ok) {
            console.error('Failed to save team', res.status);
            return fail(res.status, { 
                success: false,
                message: result.error || 'Failed to save team'
            });
        }

        console.info('Team created successfully');
        return { 
            success: true, 
            team: result,
        };
    },
    update: async ({ request, fetch, locals }) => {
        console.info('Processing team update...');
        const formData = await request.formData();
        const id = formData.get('id'); // Pass the UUID in a hidden input

        console.debug('Team ID to update:', id);
        const payload = {
            city: formData.get('city'),
            mascot: formData.get('mascot'),
            sport: formData.get('sport'),
            team_logo_url: formData.get('team_logo_url') || null,
            primary_color: formData.get('primary_color'),
            secondary_color: formData.get('secondary_color')
        };

        console.debug('Team update payload:', payload);
        const res = await fetch(apiRoute(`/api/v1/teams/${id}`), {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${locals.token}`
            },
            body: JSON.stringify(payload)
        });

        if (!res.ok) {
            const error = await res.json();
            console.error('Failed to update team', res.status, error);
            return fail(res.status, { success: false, message: error.error });
        }

        console.info('Team updated successfully');
        return { success: true, message: 'Team updated!' };
    },
    delete: async ({ request, fetch, locals }) => {
        console.info('Delete action called.');

        console.debug("Grabbing form data")
        const formData = await request.formData();
        const id = formData.get('id');

        console.warn("Attempting to delete team with id: " + id)
        const res = await fetch(apiRoute(`/api/v1/teams/${id}`), {
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${locals.token}` }
        });

        if (!res.ok) {
            console.error("Problem deleting team")
            return fail(res.status, { success: false, message: 'Failed to delete team.' });
        }

        console.info("successfully deleted team")
        return { success: true, message: 'Team removed successfully.' };
    }
};