import type { Language } from './user.types'

export interface Community {
  id: string
  name: string
  slug: string
  description: string | null
  primary_language: Language
  secondary_languages: Language[]
  avatar_url: string | null
  banner_url: string | null
  member_count: number
  post_count: number
  snap_of_week_post_id: string | null
  created_by: string | null
  created_at: string
}

export interface CommunityMember {
  id: string
  community_id: string
  user_id: string
  role: 'member' | 'moderator' | 'admin'
  joined_at: string
}
