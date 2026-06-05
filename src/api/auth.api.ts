import client from './client'
import type { User } from '../types/user.types'

export interface AuthPayload {
  access_token: string
  user: User
}

export const authApi = {
  register: (data: { username: string; email: string; password: string; full_name: string }) =>
    client.post<{ success: boolean; data: AuthPayload }>('/api/v1/auth/register', data),

  login: (data: { email: string; password: string }) =>
    client.post<{ success: boolean; data: AuthPayload }>('/api/v1/auth/login', data),

  googleAuth: (google_token: string) =>
    client.post<{ success: boolean; data: AuthPayload }>('/api/v1/auth/google', { google_token }),

  logout: () =>
    client.post('/api/v1/auth/logout'),
}
