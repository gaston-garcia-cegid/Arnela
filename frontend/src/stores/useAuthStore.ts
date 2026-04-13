import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';

// User interface following project conventions (camelCase for props)
// Note: clientId is NOT stored here. Backend derives it automatically from userId when needed.
export interface User {
  id: string; // This is the USER ID (not client ID)
  email: string;
  firstName: string;
  lastName: string;
  isActive: boolean;
  role: 'admin' | 'employee' | 'client';
}

// AuthState interface
interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (token: string, user: User) => void;
  logout: () => void;
  setUser: (user: User) => void;
}

// Create store with persist middleware
export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,

      login: (token, user) => {
        document.cookie = `auth-token=${token}; path=/; max-age=${60 * 60 * 24 * 7}; SameSite=Lax`;
        set({
          token,
          user,
          isAuthenticated: true,
        });
      },

      logout: () => {
        document.cookie = 'auth-token=; path=/; max-age=0';
        set({
          user: null,
          token: null,
          isAuthenticated: false,
        });
      },

      setUser: (user) =>
        set({ user }),
    }),
    {
      name: 'auth-storage',
    }
  )
);
