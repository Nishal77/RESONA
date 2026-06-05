import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { usersApi } from '../api/users.api'
import { PostCard } from '../components/post/PostCard'
import { LanguageBadge } from '../components/common/LanguageBadge'
import { UserAvatar } from '../components/user/UserAvatar'
import { VRSBadge } from '../components/vrs/VRSBadge'
import { Spinner } from '../components/common/Spinner'
import { useAuthStore } from '../store/auth.store'
import { MapPin } from 'lucide-react'

export default function Profile() {
  const { username } = useParams<{ username: string }>()
  const qc = useQueryClient()
  const { user: me } = useAuthStore()
  const isOwnProfile = me?.username === username
  const [following, setFollowing] = useState(false)

  const { data, isLoading } = useQuery({
    queryKey: ['user', username],
    queryFn: () => usersApi.getUser(username!),
    enabled: !!username,
  })

  const { data: postsData } = useQuery({
    queryKey: ['user-posts', data?.data?.data?.id],
    queryFn: () => usersApi.getUserPosts(data!.data.data.id, { limit: 20 }),
    enabled: !!data?.data?.data?.id,
  })

  const followMutation = useMutation({
    mutationFn: () =>
      following ? usersApi.unfollow(data!.data.data.id) : usersApi.follow(data!.data.data.id),
    onSuccess: () => {
      setFollowing(f => !f)
      qc.invalidateQueries({ queryKey: ['user', username] })
    },
    onError: () => toast.error('Action failed'),
  })

  const user = data?.data?.data
  const posts = postsData?.data?.data ?? []

  if (isLoading) return <div className="flex justify-center py-16"><Spinner size="lg" /></div>
  if (!user) return <div className="text-center py-16 text-gray-500">User not found</div>

  return (
    <div className="max-w-2xl mx-auto">
      <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-5 mb-4">
        <div className="flex items-start gap-4">
          <UserAvatar avatarUrl={user.avatar_url} username={user.username} size="lg" />
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <h1 className="text-xl font-bold text-gray-900">
                {user.full_name ?? user.username}
              </h1>
              <LanguageBadge language={user.primary_language} />
            </div>
            <p className="text-sm text-gray-500 mb-1">@{user.username}</p>
            {user.bio && <p className="text-sm text-gray-700 mb-2">{user.bio}</p>}
            {(user.state || user.city) && (
              <div className="flex items-center gap-1 text-xs text-gray-400 mb-2">
                <MapPin className="h-3 w-3" />
                {[user.city, user.state].filter(Boolean).join(', ')}
              </div>
            )}

            {/* Stats row */}
            <div className="flex items-center gap-4 text-sm">
              <span><b>{user.follower_count}</b> <span className="text-gray-500">followers</span></span>
              <span><b>{user.following_count}</b> <span className="text-gray-500">following</span></span>
              <div className="ml-auto">
                <VRSBadge score={user.vrs_total} size="md" />
              </div>
            </div>
          </div>
        </div>

        {!isOwnProfile && me && (
          <button
            onClick={() => followMutation.mutate()}
            disabled={followMutation.isPending}
            className={`mt-4 w-full rounded-lg py-2 text-sm font-medium transition-colors ${
              following
                ? 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                : 'bg-resona-saffron text-white hover:bg-orange-500'
            }`}
          >
            {following ? 'Unfollow' : 'Follow'}
          </button>
        )}
      </div>

      {/* Posts */}
      <h2 className="font-semibold text-gray-700 mb-3 text-sm">Posts</h2>
      {posts.length === 0 ? (
        <div className="text-center py-8 text-gray-400 text-sm">No posts yet</div>
      ) : (
        <div className="space-y-3">
          {posts.map(post => <PostCard key={post.id} post={post} />)}
        </div>
      )}
    </div>
  )
}
