/**
 * Helper to build full API URLs using the environment base
 */
import { PUBLIC_API_BASE_URL } from '$env/static/public';

export function apiRoute(path: string): string {
    // Debugging: This will log to your terminal during server-side loads
    if (!PUBLIC_API_BASE_URL) {
        console.error("❌ ERROR: PUBLIC_API_BASE_URL is not defined in .env");
        // Fallback to help you keep developing while you fix the env issue
        return `http://127.0.0.1:1323/${path.startsWith('/') ? path.slice(1) : path}`;
    }

    const cleanPath = path.startsWith('/') ? path.slice(1) : path;
    return `${PUBLIC_API_BASE_URL}/${cleanPath}`;
}