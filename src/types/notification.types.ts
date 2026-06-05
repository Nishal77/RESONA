import type { User } from './user.types'

export type NotificationType = 'like' | 'comment' | 'follow' | 'trending'

export interface Notification {
  id: string
  user_id: string
  type: NotificationType
  actor_id: string | null
  actor: User | null
  post_id: string | null
  message: string | null
  read: boolean
  created_at: string
}
