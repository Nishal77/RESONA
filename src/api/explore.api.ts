import client from './client'
import type { Post, Tag } from '../types/post.types'
import type { Community } from '../types/community.types'
import type { PaginatedResponse } from '../types/api.types'

export const exploreApi = {
  trending: (params?: { language?: string; limit?: number }) =>
    client.get<{ success: boolean; data: Post[] }>('/api/v1/explore/trending', { params }),

  tags: (params?: { language?: string; limit?: number }) =>
    client.get<{ success: boolean; data: Tag[] }>('/api/v1/explore/tags', { params }),

  tagPosts: (tagName: string, params?: { page?: number; limit?: number; language?: string }) =>
    client.get<PaginatedResponse<Post>>(`/api/v1/explore/tags/${tagName}/posts`, { params }),

  communities: (params?: { language?: string; limit?: number }) =>
    client.get<{ success: boolean; data: Community[] }>('/api/v1/explore/communities', { params }),

  search: (params: { q: string; language?: string; page?: number; limit?: number }) =>
    client.get<PaginatedResponse<Post>>('/api/v1/explore/search', { params }),
}
