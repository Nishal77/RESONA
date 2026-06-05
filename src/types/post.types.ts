import type { User } from './user.types'
import type { Community } from './community.types'

export interface Tag {
  id: string
  name: string
  language: string | null
  usage_count: number
  trending_score: number
}

export interface Post {
  id: string
  user_id: string
  user: User
  content_text: string | null
  media_url: string | null
  media_type: 'image' | 'video' | null
  community_id: string | null
  community: Community | null
  detected_language: string | null
  language_confidence: number | null
  language_locality_score: number
  vrs_score: number
  like_count: number
  comment_count: number
  share_count: number
  view_count: number
  save_count: number
  tags: Tag[]
  created_at: string
  updated_at: string
}

export interface Comment {
  id: string
  post_id: string
  user_id: string
  user: User
  content: string
  parent_comment_id: string | null
  replies: Comment[]
  like_count: number
  created_at: string
}
