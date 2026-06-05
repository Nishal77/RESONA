import { Outlet } from 'react-router-dom'
import { Navbar } from './Navbar'
import { BottomNav } from './BottomNav'

export function AppLayout() {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <main className="mx-auto max-w-5xl px-4 py-4 pb-20 sm:pb-6">
        <Outlet />
      </main>
      <BottomNav />
    </div>
  )
}
