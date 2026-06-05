import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Heart, Send } from 'lucide-react'
import toast from 'react-hot-toast'
import { postsApi } from '../api/posts.api'
import { VRSBadge } from '../components/vrs/VRSBadge'
import { LanguageBadge } from '../components/common/LanguageBadge'
import { UserAvatar } from '../components/user/UserAvatar'
import { Spinner } from '../components/common/Spinner'
import { timeAgo } from '../utils/formatDate'
import { useAuthStore } from '../store/auth.store'

export default function PostDetail() {
  const { id } = useParams<{ id: string }>()
  const qc = useQueryClient()
  const { user } = useAuthStore()
  const [commentText, setCommentText] = useState('')
  const [replyTo, setReplyTo] = useState<string | null>(null)

  const { data: postData, isLoading } = useQuery({
    queryKey: ['post', id],
    queryFn: () => postsApi.getPost(id!),
    enabled: !!id,
  })

  const { data: commentsData, isLoading: loadingComments } = useQuery({
    queryKey: ['comments', id],
    queryFn: () => postsApi.getComments(id!, { limit: 20 }),
    enabled: !!id,
  })

  const addComment = useMutation({
    mutationFn: () =>
      postsApi.createComment(id!, { content: commentText, parent_comment_id: replyTo ?? undefined }),
    onSuccess: () => {
      setCommentText('')
      setReplyTo(null)
      qc.invalidateQueries({ queryKey: ['comments', id] })
      qc.invalidateQueries({ queryKey: ['post', id] })
    },
    onError: () => toast.error('Failed to comment'),
  })

  const post = postData?.data?.data
  const comments = commentsData?.data?.data ?? []

  if (isLoading) return <div className="flex justify-center py-16"><Spinner size="lg" /></div>
  if (!post) return <div className="text-center py-16 text-gray-500">Post not found</div>

  return (
    <div className="max-w-2xl mx-auto">
      <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-5 mb-4">
        {/* VRS Badge — large, prominent */}
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <Link to={`/${post.user.username}`}>
              <UserAvatar avatarUrl={post.user.avatar_url} username={post.user.username} size="md" />
            </Link>
            <div>
              <Link to={`/${post.user.username}`} className="font-semibold text-gray-900 hover:underline">
                {post.user.full_name ?? post.user.username}
              </Link>
              <div className="flex items-center gap-1.5">
                {post.detected_language && <LanguageBadge language={post.detected_language} size="xs" />}
                <span className="text-xs text-gray-400">{timeAgo(post.created_at)}</span>
              </div>
            </div>
          </div>
          <VRSBadge score={post.vrs_score} size="lg" />
        </div>

        {/* Content */}
        {post.content_text && (
          <p className="text-gray-800 leading-relaxed mb-4 whitespace-pre-wrap">{post.content_text}</p>
        )}
        {post.media_url && post.media_type === 'image' && (
          <img src={post.media_url} alt="" className="w-full rounded-xl mb-4 max-h-[500px] object-cover" />
        )}
        {post.media_url && post.media_type === 'video' && (
          <video src={post.media_url} controls className="w-full rounded-xl mb-4 max-h-[500px]" />
        )}

        {/* Tags */}
        {post.tags.length > 0 && (
          <div className="flex flex-wrap gap-1 mb-4">
            {post.tags.map(t => (
              <Link key={t.id} to={`/explore/tags/${t.name}`} className="text-xs text-resona-navy hover:underline">
                #{t.name}
              </Link>
            ))}
          </div>
        )}

        {/* Stats */}
        <div className="flex gap-4 text-sm text-gray-500 border-t border-gray-50 pt-3">
          <span>{post.like_count} likes</span>
          <span>{post.comment_count} comments</span>
          <span>{post.share_count} shares</span>
          <span>{post.view_count} views</span>
        </div>
      </div>

      {/* Comments */}
      <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-5">
        <h2 className="font-semibold text-gray-900 mb-4">Comments</h2>

        {user ? (
          <div className="flex gap-2 mb-5">
            <UserAvatar avatarUrl={user.avatar_url} username={user.username} size="sm" />
            <div className="flex-1 flex gap-2">
              <input
                value={commentText}
                onChange={e => setCommentText(e.target.value)}
                onKeyDown={e => { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); if (commentText.trim()) addComment.mutate() } }}
                placeholder={replyTo ? 'Write a reply…' : 'Add a comment…'}
                className="flex-1 border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-resona-saffron"
              />
              <button
                onClick={() => commentText.trim() && addComment.mutate()}
                disabled={addComment.isPending || !commentText.trim()}
                className="bg-resona-saffron text-white rounded-lg px-3 py-2 disabled:opacity-40"
              >
                <Send className="h-4 w-4" />
              </button>
            </div>
          </div>
        ) : (
          <p className="text-sm text-gray-500 mb-4">
            <Link to="/" className="text-resona-saffron font-medium hover:underline">Login</Link> to comment
          </p>
        )}

        {loadingComments ? <Spinner /> : (
          <div className="space-y-4">
            {comments.map(c => (
              <div key={c.id}>
                <div className="flex gap-2">
                  <UserAvatar avatarUrl={c.user.avatar_url} username={c.user.username} size="sm" />
                  <div className="flex-1">
                    <div className="bg-gray-50 rounded-xl px-3 py-2">
                      <Link to={`/${c.user.username}`} className="text-xs font-semibold text-gray-900 hover:underline">
                        {c.user.username}
                      </Link>
                      <p className="text-sm text-gray-700 mt-0.5">{c.content}</p>
                    </div>
                    <div className="flex items-center gap-3 mt-1 px-1">
                      <span className="text-xs text-gray-400">{timeAgo(c.created_at)}</span>
                      <button
                        onClick={() => postsApi.likeComment(id!, c.id).catch(() => {})}
                        className="text-xs text-gray-400 hover:text-red-500 flex items-center gap-0.5"
                      >
                        <Heart className="h-3 w-3" /> {c.like_count}
                      </button>
                      {user && (
                        <button
                          onClick={() => setReplyTo(c.id)}
                          className="text-xs text-gray-400 hover:text-resona-saffron"
                        >
                          Reply
                        </button>
                      )}
                    </div>
                  </div>
                </div>

                {/* Replies */}
                {c.replies?.map(r => (
                  <div key={r.id} className="flex gap-2 mt-2 ml-10">
                    <UserAvatar avatarUrl={r.user.avatar_url} username={r.user.username} size="xs" />
                    <div className="flex-1">
                      <div className="bg-gray-50 rounded-xl px-3 py-2">
                        <Link to={`/${r.user.username}`} className="text-xs font-semibold hover:underline">{r.user.username}</Link>
                        <p className="text-sm text-gray-700 mt-0.5">{r.content}</p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
