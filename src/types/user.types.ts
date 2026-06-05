export type Language = 'kannada' | 'tamil' | 'telugu' | 'malayalam' | 'hindi' | 'english'

export interface User {
  id: string
  username: string
  email: string
  full_name: string | null
  avatar_url: string | null
  bio: string | null
  primary_language: Language
  state: string | null
  city: string | null
  vrs_total: number
  follower_count: number
  following_count: number
  onboarding_completed: boolean
  created_at: string
  updated_at: string
}
