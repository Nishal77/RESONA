import client from './client'
import type { Post, Comment } from '../types/post.types'
import type { PaginatedResponse } from '../types/api.types'

export const postsApi = {
  getFeed: (params: { page?: number; limit?: number; language?: string }) =>
    client.get<PaginatedResponse<Post>>('/api/v1/posts', { params }),

  getPost: (id: string) =>
    client.get<{ success: boolean; data: Post }>(`/api/v1/posts/${id}`),

  createPost: (data: {
    content_text?: string
    media_url?: string
    media_type?: string
    community_id?: string
    tag_ids?: string[]
    manual_language?: string
  }) => client.post<{ success: boolean; data: Post }>('/api/v1/posts', data),

  deletePost: (id: string) =>
    client.delete(`/api/v1/posts/${id}`),

  likePost: (id: string) =>
    client.post(`/api/v1/posts/${id}/like`),

  unlikePost: (id: string) =>
    client.delete(`/api/v1/posts/${id}/like`),

  sharePost: (id: string) =>
    client.post(`/api/v1/posts/${id}/share`),

  viewPost: (id: string) =>
    client.post(`/api/v1/posts/${id}/view`),

  savePost: (id: string) =>
    client.post(`/api/v1/posts/${id}/save`),

  unsavePost: (id: string) =>
    client.delete(`/api/v1/posts/${id}/save`),

  getComments: (postId: string, params?: { page?: number; limit?: number }) =>
    client.get<PaginatedResponse<Comment>>(`/api/v1/posts/${postId}/comments`, { params }),

  createComment: (postId: string, data: { content: string; parent_comment_id?: string }) =>
    client.post<{ success: boolean; data: Comment }>(`/api/v1/posts/${postId}/comments`, data),

  deleteComment: (postId: string, commentId: string) =>
    client.delete(`/api/v1/posts/${postId}/comments/${commentId}`),

  likeComment: (postId: string, commentId: string) =>
    client.post(`/api/v1/posts/${postId}/comments/${commentId}/like`),
}
