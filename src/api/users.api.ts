import client from './client'
import type { User } from '../types/user.types'
import type { Post } from '../types/post.types'
import type { PaginatedResponse } from '../types/api.types'

export const usersApi = {
  getMe: () =>
    client.get<{ success: boolean; data: User }>('/api/v1/users/me'),

  updateMe: (data: Partial<User>) =>
    client.put<{ success: boolean; data: User }>('/api/v1/users/me', data),

  getUser: (username: string) =>
    client.get<{ success: boolean; data: User }>(`/api/v1/users/${username}`),

  getUserPosts: (id: string, params?: { page?: number; limit?: number }) =>
    client.get<PaginatedResponse<Post>>(`/api/v1/users/${id}/posts`, { params }),

  getFollowers: (id: string, params?: { page?: number; limit?: number }) =>
    client.get<PaginatedResponse<User>>(`/api/v1/users/${id}/followers`, { params }),

  getFollowing: (id: string, params?: { page?: number; limit?: number }) =>
    client.get<PaginatedResponse<User>>(`/api/v1/users/${id}/following`, { params }),

  follow: (id: string) =>
    client.post(`/api/v1/users/${id}/follow`),

  unfollow: (id: string) =>
    client.delete(`/api/v1/users/${id}/follow`),
}
