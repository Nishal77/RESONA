import client from './client'
import type { Notification } from '../types/notification.types'
import type { PaginatedResponse } from '../types/api.types'

export const notificationsApi = {
  list: (params?: { page?: number; limit?: number }) =>
    client.get<PaginatedResponse<Notification>>('/api/v1/notifications', { params }),

  unreadCount: () =>
    client.get<{ success: boolean; data: { count: number } }>('/api/v1/notifications/unread-count'),

  markRead: (id: string) =>
    client.put(`/api/v1/notifications/${id}/read`),

  markAllRead: () =>
    client.put('/api/v1/notifications/read-all'),
}
