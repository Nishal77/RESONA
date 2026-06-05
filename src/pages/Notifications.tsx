import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { notificationsApi } from '../api/notifications.api'
import { UserAvatar } from '../components/user/UserAvatar'
import { Spinner } from '../components/common/Spinner'
import { timeAgo } from '../utils/formatDate'

export default function Notifications() {
  const qc = useQueryClient()

  const { data, isLoading } = useQuery({
    queryKey: ['notifications'],
    queryFn: () => notificationsApi.list({ limit: 50 }),
  })

  const markAll = useMutation({
    mutationFn: () => notificationsApi.markAllRead(),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['notifications'] }),
  })

  const notifications = data?.data?.data ?? []
  const unread = notifications.filter(n => !n.read).length

  return (
    <div className="max-w-2xl mx-auto">
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-xl font-bold text-gray-900">Notifications</h1>
        {unread > 0 && (
          <button
            onClick={() => markAll.mutate()}
            className="text-sm text-resona-saffron hover:underline"
          >
            Mark all read
          </button>
        )}
      </div>

      {isLoading ? (
        <div className="flex justify-center py-12"><Spinner size="lg" /></div>
      ) : notifications.length === 0 ? (
        <div className="text-center py-12 text-gray-400">No notifications yet</div>
      ) : (
        <div className="space-y-1">
          {notifications.map(n => {
            const href = n.post_id ? `/posts/${n.post_id}` : n.actor ? `/${n.actor.username}` : '#'
            return (
              <Link
                key={n.id}
                to={href}
                onClick={() => !n.read && notificationsApi.markRead(n.id).then(() => qc.invalidateQueries({ queryKey: ['notifications'] }))}
                className={`flex items-center gap-3 p-3 rounded-xl transition-colors ${
                  n.read ? 'bg-white hover:bg-gray-50' : 'bg-orange-50 hover:bg-orange-100'
                }`}
              >
                {n.actor ? (
                  <UserAvatar avatarUrl={n.actor.avatar_url} username={n.actor.username} size="sm" />
                ) : (
                  <div className="h-8 w-8 rounded-full bg-resona-saffron flex items-center justify-center text-white text-xs">
                    {n.type === 'trending' ? '🔥' : '•'}
                  </div>
                )}
                <div className="flex-1 min-w-0">
                  <p className="text-sm text-gray-800">{n.message}</p>
                  <p className="text-xs text-gray-400">{timeAgo(n.created_at)}</p>
                </div>
                {!n.read && (
                  <div className="h-2 w-2 rounded-full bg-resona-saffron shrink-0" />
                )}
              </Link>
            )
          })}
        </div>
      )}
    </div>
  )
}
