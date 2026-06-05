import { Link, useLocation } from 'react-router-dom'
import { Home, Compass, PlusCircle, Bell, User } from 'lucide-react'
import { useAuthStore } from '../../store/auth.store'

export function BottomNav() {
  const { pathname } = useLocation()
  const { user } = useAuthStore()

  if (!user) return null

  const items = [
    { to: '/home',          icon: Home,       label: 'Home' },
    { to: '/explore',       icon: Compass,    label: 'Explore' },
    { to: '/create',        icon: PlusCircle, label: 'Post', accent: true },
    { to: '/notifications', icon: Bell,       label: 'Alerts' },
    { to: `/${user.username}`, icon: User,    label: 'Profile' },
  ]

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-40 bg-white border-t border-gray-100 sm:hidden">
      <div className="flex items-center justify-around h-14">
        {items.map(({ to, icon: Icon, label, accent }) => {
          const active = pathname === to
          return (
            <Link
              key={to}
              to={to}
              className={`flex flex-col items-center gap-0.5 px-3 py-1 ${
                accent
                  ? 'text-resona-saffron'
                  : active
                  ? 'text-resona-saffron'
                  : 'text-gray-400'
              }`}
            >
              <Icon className={`h-5 w-5 ${accent && !active ? 'scale-110' : ''}`} />
              <span className="text-[10px]">{label}</span>
            </Link>
          )
        })}
      </div>
    </nav>
  )
}
