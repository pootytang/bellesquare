import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = ({ locals }) => {
    console.info("Layout load: passing authenticated user from locals to client");

    // The hook has already done the fetch and validation.
    // We just return it so it's available via $page.data.user
    return {
        user: locals.user
    };
};