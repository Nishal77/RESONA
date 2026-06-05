import { Suspense, lazy } from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Toaster } from 'react-hot-toast'
import { AppLayout } from './components/layout/AppLayout'
import { Spinner } from './components/common/Spinner'
import { useAuthStore } from './store/auth.store'

const Landing      = lazy(() => import('./pages/Landing'))
const Onboarding   = lazy(() => import('./pages/Onboarding'))
const Home         = lazy(() => import('./pages/Home'))
const Explore      = lazy(() => import('./pages/Explore'))
const PostDetail   = lazy(() => import('./pages/PostDetail'))
const CreatePost   = lazy(() => import('./pages/CreatePost'))
const Community    = lazy(() => import('./pages/Community'))
const Profile      = lazy(() => import('./pages/Profile'))
const Notifications = lazy(() => import('./pages/Notifications'))
const Settings     = lazy(() => import('./pages/Settings'))

const qc = new QueryClient({
  defaultOptions: {
    queries: { staleTime: 30_000, retry: 1 },
  },
})

function RequireAuth({ children }: { children: React.ReactNode }) {
  const { user } = useAuthStore()
  if (!user) return <Navigate to="/" replace />
  return <>{children}</>
}

function LoadingFallback() {
  return (
    <div className="flex items-center justify-center min-h-screen">
      <Spinner size="lg" />
    </div>
  )
}

export default function App() {
  return (
    <QueryClientProvider client={qc}>
      <BrowserRouter>
        <Suspense fallback={<LoadingFallback />}>
          <Routes>
            {/* Public */}
            <Route path="/" element={<Landing />} />

            {/* Auth-required standalone */}
            <Route path="/onboarding" element={<RequireAuth><Onboarding /></RequireAuth>} />

            {/* App layout routes */}
            <Route element={<AppLayout />}>
              <Route path="/home" element={<RequireAuth><Home /></RequireAuth>} />
              <Route path="/explore" element={<Explore />} />
              <Route path="/posts/:id" element={<PostDetail />} />
              <Route path="/create" element={<RequireAuth><CreatePost /></RequireAuth>} />
              <Route path="/c/:slug" element={<Community />} />
              <Route path="/notifications" element={<RequireAuth><Notifications /></RequireAuth>} />
              <Route path="/settings" element={<RequireAuth><Settings /></RequireAuth>} />
              {/* Profile must be last — catches /:username */}
              <Route path="/:username" element={<Profile />} />
            </Route>

            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </Suspense>
      </BrowserRouter>

      <Toaster
        position="bottom-center"
        toastOptions={{
          duration: 3000,
          style: { borderRadius: '12px', fontSize: '14px' },
        }}
      />
    </QueryClientProvider>
  )
}
