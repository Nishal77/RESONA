import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Heart, MessageCircle, Share2, Bookmark } from 'lucide-react'
import toast from 'react-hot-toast'
import type { Post } from '../../types/post.types'
import { VRSBadge } from '../vrs/VRSBadge'
import { LanguageBadge } from '../common/LanguageBadge'
import { UserAvatar } from '../user/UserAvatar'
import { timeAgo } from '../../utils/formatDate'
import { postsApi } from '../../api/posts.api'
import { useAuthStore } from '../../store/auth.store'

interface Props {
  post: Post
  onUpdate?: (updated: Post) => void
}

export function PostCard({ post, onUpdate: _onUpdate }: Props) {
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const [liked, setLiked] = useState(false)
  const [saved, setSaved] = useState(false)
  const [counts, setCounts] = useState({
    likes: post.like_count,
    comments: post.comment_count,
    shares: post.share_count,
  })

  const requireAuth = (fn: () => void) => {
    if (!user) { navigate('/'); toast.error('Login to interact'); return }
    fn()
  }

  const handleLike = () => requireAuth(async () => {
    try {
      if (liked) {
        await postsApi.unlikePost(post.id)
        setCounts(c => ({ ...c, likes: c.likes - 1 }))
      } else {
        await postsApi.likePost(post.id)
        setCounts(c => ({ ...c, likes: c.likes + 1 }))
      }
      setLiked(!liked)
    } catch { toast.error('Action failed') }
  })

  const handleShare = async () => {
    await navigator.clipboard.writeText(`${window.location.origin}/posts/${post.id}`)
    toast.success('Link copied!')
    try { await postsApi.sharePost(post.id) } catch { /* non-fatal */ }
  }

  const handleSave = () => requireAuth(async () => {
    try {
      if (saved) { await postsApi.unsavePost(post.id) }
      else { await postsApi.savePost(post.id) }
      setSaved(!saved)
    } catch { toast.error('Action failed') }
  })

  const truncated = post.content_text && post.content_text.length > 280
    ? post.content_text.slice(0, 280) + '…'
    : post.content_text

  return (
    <article className="bg-white rounded-xl border border-gray-100 shadow-sm hover:shadow-md transition-shadow p-4">
      {/* Header */}
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center gap-2">
          <Link to={`/${post.user.username}`}>
            <UserAvatar avatarUrl={post.user.avatar_url} username={post.user.username} size="sm" />
          </Link>
          <div>
            <Link to={`/${post.user.username}`} className="font-medium text-gray-900 text-sm hover:underline">
              {post.user.full_name ?? post.user.username}
            </Link>
            <div className="flex items-center gap-1 mt-0.5">
              {post.detected_language && (
                <LanguageBadge language={post.detected_language} size="xs" />
              )}
              <span className="text-xs text-gray-400">{timeAgo(post.created_at)}</span>
            </div>
          </div>
        </div>
        {/* VRS Badge — always visible, no exceptions */}
        <VRSBadge score={post.vrs_score} size="sm" />
      </div>

      {/* Community */}
      {post.community && (
        <Link
          to={`/c/${post.community.slug}`}
          className="text-xs text-resona-saffron font-medium mb-2 block hover:underline"
        >
          c/{post.community.name}
        </Link>
      )}

      {/* Content */}
      <Link to={`/posts/${post.id}`} className="block">
        {post.content_text && (
          <p className="text-gray-800 text-sm leading-relaxed mb-3">{truncated}</p>
        )}
        {post.media_url && post.media_type === 'image' && (
          <img
            src={post.media_url}
            alt="post media"
            className="w-full rounded-lg object-cover max-h-72 mb-3"
            loading="lazy"
          />
        )}
        {post.media_url && post.media_type === 'video' && (
          <video
            src={post.media_url}
            className="w-full rounded-lg max-h-72 mb-3"
            controls
          />
        )}
      </Link>

      {/* Tags */}
      {post.tags && post.tags.length > 0 && (
        <div className="flex flex-wrap gap-1 mb-3">
          {post.tags.map(tag => (
            <Link
              key={tag.id}
              to={`/explore/tags/${tag.name}`}
              className="text-xs text-resona-navy hover:underline"
            >
              #{tag.name}
            </Link>
          ))}
        </div>
      )}

      {/* Actions */}
      <div className="flex items-center gap-4 pt-2 border-t border-gray-50">
        <button
          onClick={handleLike}
          className={`flex items-center gap-1.5 text-sm transition-colors ${liked ? 'text-red-500' : 'text-gray-500 hover:text-red-500'}`}
        >
          <Heart className={`h-4 w-4 ${liked ? 'fill-current' : ''}`} />
          <span>{counts.likes}</span>
        </button>

        <Link
          to={`/posts/${post.id}`}
          className="flex items-center gap-1.5 text-sm text-gray-500 hover:text-blue-500 transition-colors"
        >
          <MessageCircle className="h-4 w-4" />
          <span>{counts.comments}</span>
        </Link>

        <button
          onClick={handleShare}
          className="flex items-center gap-1.5 text-sm text-gray-500 hover:text-green-500 transition-colors"
        >
          <Share2 className="h-4 w-4" />
          <span>{counts.shares}</span>
        </button>

        <button
          onClick={handleSave}
          className={`ml-auto flex items-center gap-1.5 text-sm transition-colors ${saved ? 'text-resona-saffron' : 'text-gray-400 hover:text-resona-saffron'}`}
        >
          <Bookmark className={`h-4 w-4 ${saved ? 'fill-current' : ''}`} />
        </button>
      </div>
    </article>
  )
}
