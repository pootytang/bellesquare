/**
 * Helper to build full API URLs using the environment base
 */
import { PUBLIC_API_BASE_URL } from '$env/static/public';
import { INTERNAL_API_URL } from '$env/static/private';

export function apiRoute(path: string): string {
    const cleanPath = path.startsWith('/') ? path.slice(1) : path;

    // If we are on the server (SSR), use the high-speed internal URL
    if (typeof window === 'undefined' && INTERNAL_API_URL) {
        return `${INTERNAL_API_URL.replace(/\/$/, '')}/${cleanPath}`;
    }

    // Otherwise (client-side), use the public URL
    const cleanBase = PUBLIC_API_BASE_URL?.replace(/\/$/, '') || 'http://127.0.0.1:1323';
    return `${cleanBase}/${cleanPath}`;
}