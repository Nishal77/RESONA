import { useState, useCallback } from 'react'
import { Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { Search } from 'lucide-react'
import { exploreApi } from '../api/explore.api'
import { PostCard } from '../components/post/PostCard'
import { LanguageBadge } from '../components/common/LanguageBadge'
import { VRSBadge } from '../components/vrs/VRSBadge'
import { Spinner } from '../components/common/Spinner'
import { SUPPORTED_LANGUAGES } from '../constants/languages'
import type { Language } from '../types/user.types'
import { useAuthStore } from '../store/auth.store'
import { communitiesApi } from '../api/communities.api'
import toast from 'react-hot-toast'

export default function Explore() {
  const { user } = useAuthStore()
  const [lang, setLang] = useState<Language | ''>('')
  const [search, setSearch] = useState('')
  const [searchQuery, setSearchQuery] = useState('')

  const { data: trending } = useQuery({
    queryKey: ['explore', 'trending', lang],
    queryFn: () => exploreApi.trending({ language: lang || undefined, limit: 10 }),
  })

  const { data: tags } = useQuery({
    queryKey: ['explore', 'tags', lang],
    queryFn: () => exploreApi.tags({ language: lang || undefined, limit: 10 }),
  })

  const { data: communities } = useQuery({
    queryKey: ['explore', 'communities', lang],
    queryFn: () => exploreApi.communities({ language: lang || undefined, limit: 6 }),
  })

  const { data: searchResults, isLoading: searching } = useQuery({
    queryKey: ['explore', 'search', searchQuery, lang],
    queryFn: () => exploreApi.search({ q: searchQuery, language: lang || undefined }),
    enabled: searchQuery.length > 1,
  })

  const handleSearch = useCallback((e: React.FormEvent) => {
    e.preventDefault()
    setSearchQuery(search)
  }, [search])

  const handleJoin = async (id: string) => {
    if (!user) { toast.error('Login to join'); return }
    try { await communitiesApi.join(id); toast.success('Joined!') }
    catch { toast.error('Already a member') }
  }

  const trendingPosts = trending?.data?.data ?? []
  const trendingTags = tags?.data?.data ?? []
  const topCommunities = communities?.data?.data ?? []
  const results = searchResults?.data?.data ?? []

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      {/* Search bar */}
      <form onSubmit={handleSearch} className="flex gap-2">
        <div className="flex-1 relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
          <input
            value={search}
            onChange={e => setSearch(e.target.value)}
            placeholder="Search in Bharat…"
            className="w-full pl-9 pr-4 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-resona-saffron"
          />
        </div>
        <select
          value={lang}
          onChange={e => setLang(e.target.value as Language | '')}
          className="border border-gray-200 rounded-xl px-3 py-2 text-sm focus:outline-none"
        >
          <option value="">All Languages</option>
          {SUPPORTED_LANGUAGES.map(l => (
            <option key={l.code} value={l.code}>{l.label}</option>
          ))}
        </select>
      </form>

      {/* Search results */}
      {searchQuery && (
        <div>
          <h2 className="text-sm font-semibold text-gray-600 mb-3">Results for "{searchQuery}"</h2>
          {searching ? <Spinner /> : results.length === 0 ? (
            <p className="text-sm text-gray-500">No posts found</p>
          ) : (
            <div className="space-y-3">
              {results.map(post => <PostCard key={post.id} post={post} />)}
            </div>
          )}
        </div>
      )}

      {!searchQuery && (
        <>
          {/* Trending Tags */}
          <section>
            <h2 className="font-semibold text-gray-900 mb-3">Trending Tags</h2>
            <div className="flex flex-wrap gap-2">
              {trendingTags.map(tag => (
                <Link
                  key={tag.id}
                  to={`/explore/tags/${tag.name}`}
                  className="bg-white border border-gray-200 rounded-full px-3 py-1.5 text-sm hover:border-resona-saffron transition-colors flex items-center gap-2"
                >
                  <span className="text-resona-navy font-medium">#{tag.name}</span>
                  {tag.language && <LanguageBadge language={tag.language} size="xs" />}
                  <span className="text-xs text-gray-400">{tag.usage_count}</span>
                </Link>
              ))}
            </div>
          </section>

          {/* Top Communities */}
          <section>
            <h2 className="font-semibold text-gray-900 mb-3">Top Communities</h2>
            <div className="grid grid-cols-2 gap-3">
              {topCommunities.map(c => (
                <div key={c.id} className="bg-white rounded-xl border border-gray-100 p-3">
                  <div className="flex items-center gap-2 mb-1">
                    <div className="h-8 w-8 rounded-full bg-resona-saffron flex items-center justify-center text-white text-xs font-bold">
                      {c.name.slice(0, 2).toUpperCase()}
                    </div>
                    <div className="flex-1 min-w-0">
                      <Link to={`/c/${c.slug}`} className="text-sm font-medium truncate block hover:underline">
                        {c.name}
                      </Link>
                      <LanguageBadge language={c.primary_language} size="xs" />
                    </div>
                  </div>
                  <p className="text-xs text-gray-500 mb-2">{c.member_count} members</p>
                  <button
                    onClick={() => handleJoin(c.id)}
                    className="w-full text-xs bg-resona-saffron text-white rounded-lg py-1.5 hover:bg-orange-500"
                  >
                    Join
                  </button>
                </div>
              ))}
            </div>
          </section>

          {/* Trending Across Bharat */}
          <section>
            <div className="flex items-center gap-2 mb-3">
              <h2 className="font-semibold text-gray-900">Trending Across Bharat</h2>
              <span className="text-xs bg-resona-saffron text-white px-2 py-0.5 rounded-full">🔥</span>
            </div>
            <div className="space-y-3">
              {trendingPosts.map(post => (
                <div key={post.id} className="bg-white rounded-xl border border-gray-100 p-3 flex items-center gap-3">
                  <div className="flex-1 min-w-0">
                    <Link to={`/posts/${post.id}`} className="text-sm font-medium text-gray-800 line-clamp-2 hover:underline">
                      {post.content_text?.slice(0, 100) ?? 'Media post'}
                    </Link>
                    <div className="flex items-center gap-1.5 mt-1">
                      {post.detected_language && <LanguageBadge language={post.detected_language} size="xs" />}
                      <span className="text-xs text-gray-400">@{post.user.username}</span>
                    </div>
                  </div>
                  <VRSBadge score={post.vrs_score} size="sm" />
                </div>
              ))}
            </div>
          </section>
        </>
      )}
    </div>
  )
}
