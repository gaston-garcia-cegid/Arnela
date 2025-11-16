import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';

// User interface following project conventions (camelCase for props)
export interface User {
  id: string;
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

      login: (token, user) =>
        set({
          token,
          user,
          isAuthenticated: true,
        }),

      logout: () =>
        set({
          user: null,
          token: null,
          isAuthenticated: false,
        }),

      setUser: (user) =>
        set({ user }),
    }),
    {
      name: 'auth-storage',
    }
  )
);
