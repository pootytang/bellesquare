import type { User } from './types';

class AuthState {
    user = $state<User | null>(null);
    initialized = $state(false);

    // setUser(userData: User | null) {
    //     this.user = userData;
    //     this.initialized = true;
    // }
    setUser(userData: User | null) {
    // Only update if the data is actually different
    if (this.user?.id !== userData?.id) {
        this.user = userData;
    }
    this.initialized = true;
}

    logout() {
        this.user = null;
        // this.initialized = false;
    }
}

export const auth = new AuthState();