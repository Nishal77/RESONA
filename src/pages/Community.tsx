import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { communitiesApi } from '../api/communities.api'
import { PostCard } from '../components/post/PostCard'
import { LanguageBadge } from '../components/common/LanguageBadge'
import { VRSBadge } from '../components/vrs/VRSBadge'
import { Spinner } from '../components/common/Spinner'
import { useAuthStore } from '../store/auth.store'

export default function Community() {
  const { slug } = useParams<{ slug: string }>()
  const { user } = useAuthStore()
  const [sort, setSort] = useState<'vrs' | 'latest' | 'top'>('vrs')
  const [joined, setJoined] = useState(false)

  const { data, isLoading } = useQuery({
    queryKey: ['community', slug],
    queryFn: () => communitiesApi.getCommunity(slug!),
    enabled: !!slug,
  })

  const { data: postsData, isLoading: loadingPosts } = useQuery({
    queryKey: ['community-posts', data?.data?.data?.id, sort],
    queryFn: () => communitiesApi.getPosts(data!.data.data.id, { sort, limit: 20 }),
    enabled: !!data?.data?.data?.id,
  })

  const community = data?.data?.data
  const posts = postsData?.data?.data ?? []

  const handleJoin = async () => {
    if (!user) { toast.error('Login to join'); return }
    if (!community) return
    try {
      if (joined) { await communitiesApi.leave(community.id); setJoined(false); toast.success('Left') }
      else { await communitiesApi.join(community.id); setJoined(true); toast.success('Joined!') }
    } catch (e: any) { toast.error(e.response?.data?.error ?? 'Action failed') }
  }

  if (isLoading) return <div className="flex justify-center py-16"><Spinner size="lg" /></div>
  if (!community) return <div className="text-center py-16 text-gray-500">Community not found</div>

  return (
    <div className="max-w-2xl mx-auto">
      {/* Banner */}
      {community.banner_url && (
        <img src={community.banner_url} alt="" className="w-full h-32 object-cover rounded-xl mb-4" />
      )}

      {/* Community header */}
      <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-4 mb-4">
        <div className="flex items-start gap-3">
          <div className="h-12 w-12 rounded-full bg-resona-saffron flex items-center justify-center text-white font-bold text-lg shrink-0">
            {community.name.slice(0, 2).toUpperCase()}
          </div>
          <div className="flex-1">
            <div className="flex items-center gap-2 flex-wrap">
              <h1 className="text-lg font-bold text-gray-900">{community.name}</h1>
              <LanguageBadge language={community.primary_language} />
            </div>
            {community.description && (
              <p className="text-sm text-gray-600 mt-1">{community.description}</p>
            )}
            <p className="text-xs text-gray-400 mt-1">{community.member_count} members</p>
          </div>
          <button
            onClick={handleJoin}
            className={`shrink-0 rounded-full px-4 py-1.5 text-sm font-medium transition-colors ${
              joined ? 'bg-gray-100 text-gray-600 hover:bg-gray-200' : 'bg-resona-saffron text-white hover:bg-orange-500'
            }`}
          >
            {joined ? 'Leave' : 'Join'}
          </button>
        </div>

        {/* Snap of the Week */}
        {community.snap_of_week_post_id && (
          <div className="mt-4 bg-gradient-to-r from-orange-50 to-yellow-50 rounded-lg p-3 border border-orange-100">
            <div className="text-xs font-semibold text-resona-saffron mb-1">⚡ Snap of the Week</div>
            <a href={`/posts/${community.snap_of_week_post_id}`} className="text-sm text-gray-700 hover:underline">
              View this week's top post →
            </a>
          </div>
        )}
      </div>

      {/* Sort tabs */}
      <div className="flex gap-1 mb-4 bg-white rounded-xl border border-gray-100 p-1">
        {(['vrs', 'latest', 'top'] as const).map(s => (
          <button
            key={s}
            onClick={() => setSort(s)}
            className={`flex-1 rounded-lg py-1.5 text-sm font-medium transition-colors ${
              sort === s ? 'bg-resona-saffron text-white' : 'text-gray-500 hover:bg-gray-50'
            }`}
          >
            {s === 'vrs' ? '🔥 VRS' : s === 'latest' ? '🕐 Latest' : '⬆ Top'}
          </button>
        ))}
      </div>

      {loadingPosts ? (
        <div className="flex justify-center py-8"><Spinner /></div>
      ) : (
        <div className="space-y-3">
          {posts.map(post => <PostCard key={post.id} post={post} />)}
        </div>
      )}
    </div>
  )
}
