import type { PageServerLoad } from './$types';
import { apiRoute } from '$lib/api';

export const load: PageServerLoad = async ({ fetch, locals }) => {
    console.info("********** LOADING DASHBOARD - MANAGED BOARDS, JOINED BOARDS, AND COUNTS **********");
    // Trigger all three requests in parallel
    console.debug("Initiating parallel fetches for managed boards, joined boards, and team counts");
    const [managedRes, joinedRes, countsRes] = await Promise.all([
        fetch(apiRoute('/api/v1/boards'), { // Calls GetBoardsByUserID
            headers: { 'Authorization': `Bearer ${locals.token}` }
        }),
        fetch(apiRoute('/api/v1/boards/joined'), { // Calls GetJoinedBoards
            headers: { 'Authorization': `Bearer ${locals.token}` }
        }),
        fetch(apiRoute('/api/v1/teams/counts'), {
            headers: { 'Authorization': `Bearer ${locals.token}` }
        })
    ]);

    // 2. Parse responses
    console.debug("Parsing responses for managed boards, joined boards, and team counts");
    const managedBoards = managedRes.ok ? await managedRes.json() : [];
    const joinedBoards = joinedRes.ok ? await joinedRes.json() : [];
    const teamCounts = countsRes.ok ? await countsRes.json() : {};

    console.info("********** DASHBOARD LOADING COMPLETE **********");
    return { 
        boards: managedBoards || [], 
        joined: joinedBoards || [], 
        teamCounts 
    };
};

// TODO: REMOVE THIS OLD CODE ONCE THE NEW DASHBOARD IS VERIFIED TO BE WORKING. KEEPING FOR REFERENCE DURING DEVELOPMENT
// import type { PageServerLoad } from './$types';
// import { apiRoute } from '$lib/api';

// export const load: PageServerLoad = async ({ fetch, locals }) => {
//     console.info("Grabbing all boards and team counts")
//     const [boardsRes, countsRes] = await Promise.all([
//         fetch(apiRoute('/api/v1/boards'), {
//             headers: { 'Authorization': `Bearer ${locals.token}` }
//         }),
//         // New endpoint: returns {"football": 2, "basketball": 0...}
//         fetch(apiRoute('/api/v1/teams/counts'), {
//             headers: { 'Authorization': `Bearer ${locals.token}` }
//         })
//     ]);

//     console.debug("Checking the response")
//     const boards = boardsRes.ok ? await boardsRes.json() : [];
//     const teamCounts = countsRes.ok ? await countsRes.json() : {};

    
//     console.info(`Boards Retrieved ${boards.length}`)
//     console.info("--- TEAM COUNTS ---")
//     console.table(teamCounts)
//     return { 
//         boards: boards || [],
//         teamCounts // Pass the whole map to the frontend
//     };
// };