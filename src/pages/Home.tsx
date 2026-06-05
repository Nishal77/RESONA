import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useInfiniteQuery } from '@tanstack/react-query'
import { PlusCircle } from 'lucide-react'
import { postsApi } from '../api/posts.api'
import { PostCard } from '../components/post/PostCard'
import { LanguageFilter } from '../components/common/LanguageFilter'
import { Spinner } from '../components/common/Spinner'
import { EmptyState } from '../components/common/EmptyState'
import { useAuthStore } from '../store/auth.store'
import type { Language } from '../types/user.types'

export default function Home() {
  const { user } = useAuthStore()
  const [language, setLanguage] = useState<Language | ''>('')

  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading } = useInfiniteQuery({
    queryKey: ['feed', language],
    queryFn: ({ pageParam = 1 }) =>
      postsApi.getFeed({ page: pageParam as number, limit: 20, language: language || undefined }),
    getNextPageParam: (last) => {
      const meta = last.data.meta
      return meta.has_more ? meta.page + 1 : undefined
    },
    initialPageParam: 1,
  })

  const posts = data?.pages.flatMap(p => p.data.data) ?? []

  return (
    <div className="max-w-2xl mx-auto">
      {/* Language filter */}
      <div className="mb-4">
        <LanguageFilter value={language} onChange={setLanguage} />
      </div>

      {isLoading ? (
        <div className="flex justify-center py-16"><Spinner size="lg" /></div>
      ) : posts.length === 0 ? (
        <EmptyState
          title="No posts yet"
          description="Follow creators or join communities to see posts here"
        />
      ) : (
        <div className="space-y-3">
          {posts.map(post => (
            <PostCard key={post.id} post={post} />
          ))}
        </div>
      )}

      {hasNextPage && (
        <div className="flex justify-center py-6">
          <button
            onClick={() => fetchNextPage()}
            disabled={isFetchingNextPage}
            className="bg-white border border-gray-200 rounded-full px-5 py-2 text-sm text-gray-600 hover:bg-gray-50 disabled:opacity-50"
          >
            {isFetchingNextPage ? <Spinner size="sm" /> : 'Load more'}
          </button>
        </div>
      )}

      {/* Floating create button — mobile */}
      <Link
        to="/create"
        className="fixed bottom-20 right-4 sm:hidden bg-resona-saffron text-white rounded-full p-3.5 shadow-lg hover:bg-orange-500 transition-colors"
      >
        <PlusCircle className="h-6 w-6" />
      </Link>
    </div>
  )
}
