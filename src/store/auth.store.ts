import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '../types/user.types'

interface AuthState {
  user: User | null
  accessToken: string | null
  setAuth: (user: User, token: string) => void
  setAccessToken: (token: string) => void
  setUser: (user: User) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,

      setAuth: (user, accessToken) => set({ user, accessToken }),
      setAccessToken: (accessToken) => set({ accessToken }),
      setUser: (user) => set({ user }),
      logout: () => set({ user: null, accessToken: null }),
    }),
    {
      name: 'resona-auth',
      partialize: (state) => ({ user: state.user, accessToken: state.accessToken }),
    },
  ),
)
