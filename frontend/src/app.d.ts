// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user: User | null; // This allows you to set the user in hooks
			token: string | null; // Store the auth token for API requests
		}
		// interface PageData {
		// 	user: User | null; // This ensures $page.data.user is typed
		// }
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
