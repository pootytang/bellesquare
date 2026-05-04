import { redirect, type Handle } from '@sveltejs/kit';
import { apiRoute } from '$lib/api'; // Adjust based on your helper location

export const handle: Handle = async ({ event, resolve }) => {
	console.info("Hooks running")
	
    const token = event.cookies.get('at');
    event.locals.token = token || '';
	console.debug('Token from cookie:', token);

    // 1. If we have a token, hydrate the user object
    if (token) {
        try {
			console.info("requesting 'me' endpoint from HOOKS")
            const res = await event.fetch(apiRoute('/auth/me'), {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (res.ok) {
				console.debug("got a positive response, making user available to the app")
                const userData = await res.json();

				// Simple, reliable initials calculation
                const fName = (userData.first_name || '').trim();
                const lName = (userData.last_name || '').trim();
                console.debug(`HOOKS retrieved first name: ${fName}, last name: ${lName}, preferred_color: ${userData.preferred_color.trim()}`)
                const initials = (
                    (fName.charAt(0) || '') + 
                    (lName.charAt(0) || '')
                ).toUpperCase() || '??';

                event.locals.user = {
                    ...userData,
                    color: userData.preferred_color || '#2563eb',
                    initials
                };
            } else if (res.status === 401) {
				console.warn("token may have expired, deleting it locally")
                // Token is dead, clear it
                event.cookies.delete('at', { path: '/' });
                event.locals.token = '';
            }
        } catch (err) {
            console.error('Hooks Error:', err);
        }
    }

    // 2. Protect Dashboard Routes
    if (event.url.pathname.startsWith('/dashboard')) {
        console.warn("/dashboard/* routes are protected")
        
        if (!event.locals.token || !event.locals.user) {
            console.warn('Unauthorized dashboard access, capturing redirect path...');
            
            // Capture the path they were trying to hit (e.g., /dashboard/boards/uuid)
            const from = event.url.pathname + event.url.search;
            
            // Redirect to register with the 'relay' parameter
            throw redirect(303, `/register?redirectTo=${encodeURIComponent(from)}`);
        }
    }

    return await resolve(event);
};