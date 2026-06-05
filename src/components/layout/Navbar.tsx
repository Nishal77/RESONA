import { Link, useNavigate } from 'react-router-dom'
import { Bell, Search, PlusCircle, LogOut, Settings } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import { useAuthStore } from '../../store/auth.store'
import { notificationsApi } from '../../api/notifications.api'
import { authApi } from '../../api/auth.api'
import { UserAvatar } from '../user/UserAvatar'
import toast from 'react-hot-toast'

export function Navbar() {
  const navigate = useNavigate()
  const { user, logout } = useAuthStore()

  const { data: unreadData } = useQuery({
    queryKey: ['notifications', 'unread'],
    queryFn: () => notificationsApi.unreadCount(),
    enabled: !!user,
    refetchInterval: 30_000,
  })
  const unreadCount = unreadData?.data?.data?.count ?? 0

  const handleLogout = async () => {
    try { await authApi.logout() } catch { /* ignore */ }
    logout()
    navigate('/')
    toast.success('Logged out')
  }

  return (
    <header className="sticky top-0 z-40 bg-white border-b border-gray-100 shadow-sm">
      <div className="mx-auto max-w-5xl flex items-center justify-between px-4 h-14">
        {/* Logo */}
        <Link to={user ? '/home' : '/'} className="flex items-center gap-1.5">
          <span className="text-xl font-bold text-resona-saffron">ভ</span>
          <span className="font-bold text-gray-900 hidden sm:block">Resona</span>
        </Link>

        {/* Search */}
        <Link
          to="/explore"
          className="flex items-center gap-2 bg-gray-100 rounded-full px-3 py-1.5 text-sm text-gray-500 hover:bg-gray-200 transition-colors w-40 sm:w-64"
        >
          <Search className="h-4 w-4 shrink-0" />
          <span className="truncate">Search in Bharat…</span>
        </Link>

        {/* Right actions */}
        <div className="flex items-center gap-2">
          {user ? (
            <>
              <Link
                to="/create"
                className="hidden sm:flex items-center gap-1.5 bg-resona-saffron text-white rounded-full px-3 py-1.5 text-sm font-medium hover:bg-orange-500 transition-colors"
              >
                <PlusCircle className="h-4 w-4" />
                Create
              </Link>

              <Link to="/notifications" className="relative p-2 text-gray-600 hover:text-gray-900">
                <Bell className="h-5 w-5" />
                {unreadCount > 0 && (
                  <span className="absolute top-1 right-1 h-4 w-4 bg-red-500 text-white text-[10px] rounded-full flex items-center justify-center font-bold">
                    {unreadCount > 9 ? '9+' : unreadCount}
                  </span>
                )}
              </Link>

              <div className="relative group">
                <button className="flex items-center">
                  <UserAvatar avatarUrl={user.avatar_url} username={user.username} size="sm" />
                </button>
                <div className="absolute right-0 top-full mt-1 w-44 bg-white rounded-xl shadow-lg border border-gray-100 py-1 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-50">
                  <Link to={`/${user.username}`} className="flex items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:bg-gray-50">
                    <UserAvatar avatarUrl={user.avatar_url} username={user.username} size="xs" />
                    {user.username}
                  </Link>
                  <Link to="/settings" className="flex items-center gap-2 px-3 py-2 text-sm text-gray-700 hover:bg-gray-50">
                    <Settings className="h-4 w-4" /> Settings
                  </Link>
                  <button onClick={handleLogout} className="flex items-center gap-2 w-full px-3 py-2 text-sm text-red-500 hover:bg-red-50">
                    <LogOut className="h-4 w-4" /> Logout
                  </button>
                </div>
              </div>
            </>
          ) : (
            <Link
              to="/"
              className="bg-resona-saffron text-white rounded-full px-4 py-1.5 text-sm font-medium hover:bg-orange-500 transition-colors"
            >
              Join Resona
            </Link>
          )}
        </div>
      </div>
    </header>
  )
}
