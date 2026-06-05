import client from './client'
import type { Community } from '../types/community.types'
import type { Post } from '../types/post.types'
import type { PaginatedResponse } from '../types/api.types'

export const communitiesApi = {
  list: (params?: { page?: number; limit?: number; language?: string }) =>
    client.get<PaginatedResponse<Community>>('/api/v1/communities', { params }),

  getCommunity: (slug: string) =>
    client.get<{ success: boolean; data: Community }>(`/api/v1/communities/${slug}`),

  createCommunity: (data: {
    name: string
    description?: string
    primary_language: string
    secondary_languages?: string[]
    avatar_url?: string
    banner_url?: string
  }) => client.post<{ success: boolean; data: Community }>('/api/v1/communities', data),

  getPosts: (id: string, params?: { page?: number; limit?: number; sort?: string }) =>
    client.get<PaginatedResponse<Post>>(`/api/v1/communities/${id}/posts`, { params }),

  join: (id: string) =>
    client.post(`/api/v1/communities/${id}/join`),

  leave: (id: string) =>
    client.delete(`/api/v1/communities/${id}/join`),
}
