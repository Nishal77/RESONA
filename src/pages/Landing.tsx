import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import toast from 'react-hot-toast'
import { useGoogleLogin } from '@react-oauth/google'
import { authApi } from '../api/auth.api'
import { useAuthStore } from '../store/auth.store'
import { SUPPORTED_LANGUAGES } from '../constants/languages'

type Tab = 'login' | 'register'

interface LoginForm { email: string; password: string }
interface RegisterForm { username: string; email: string; password: string; full_name: string }

export default function Landing() {
  const [tab, setTab] = useState<Tab>('login')
  const [googleLoading, setGoogleLoading] = useState(false)
  const navigate = useNavigate()
  const { setAuth } = useAuthStore()

  const loginForm = useForm<LoginForm>()
  const registerForm = useForm<RegisterForm>()

  const onLogin = async (data: LoginForm) => {
    try {
      const res = await authApi.login(data)
      const { access_token, user } = res.data.data
      setAuth(user, access_token)
      toast.success(`Welcome back, ${user.username}!`)
      navigate(user.onboarding_completed ? '/home' : '/onboarding')
    } catch (err: any) {
      toast.error(err.response?.data?.error ?? 'Login failed')
    }
  }

  const onRegister = async (data: RegisterForm) => {
    try {
      const res = await authApi.register(data)
      const { access_token, user } = res.data.data
      setAuth(user, access_token)
      toast.success('Account created!')
      navigate('/onboarding')
    } catch (err: any) {
      toast.error(err.response?.data?.error ?? 'Registration failed')
    }
  }

  // Google OAuth — gets credential JWT from Google, sends to backend
  const googleLogin = useGoogleLogin({
    onSuccess: async (tokenResponse) => {
      setGoogleLoading(true)
      try {
        // tokenResponse.access_token is the OAuth access token
        // Backend verifyGoogleToken uses tokeninfo endpoint — needs id_token (credential)
        // useGoogleLogin flow=implicit gives access_token, not id_token
        // So fetch userinfo and send access_token — backend needs to handle this
        const res = await authApi.googleAuth(tokenResponse.access_token)
        const { access_token, user } = res.data.data
        setAuth(user, access_token)
        toast.success(`Welcome, ${user.full_name ?? user.username}!`)
        navigate(user.onboarding_completed ? '/home' : '/onboarding')
      } catch (err: any) {
        toast.error(err.response?.data?.error ?? 'Google login failed')
      } finally {
        setGoogleLoading(false)
      }
    },
    onError: () => {
      toast.error('Google login cancelled')
    },
    flow: 'implicit',
  })

  return (
    <div className="min-h-screen bg-gradient-to-br from-resona-cream via-white to-orange-50">
      {/* Hero */}
      <div className="mx-auto max-w-5xl px-4 pt-12 pb-8 grid md:grid-cols-2 gap-12 items-center">
        <div>
          <div className="flex items-center gap-2 mb-4">
            <span className="text-5xl font-bold text-resona-saffron">ভ</span>
            <span className="text-3xl font-bold text-gray-900">Resona</span>
          </div>
          <h1 className="text-4xl font-bold text-gray-900 leading-tight mb-4">
            Bharat's Voice,<br />
            <span className="text-resona-saffron">In Your Language</span>
          </h1>
          <p className="text-gray-600 text-lg mb-6">
            The first vernacular-first social platform where regional language creators get discovered — without speaking English.
          </p>

          {/* Language showcase */}
          <div className="flex flex-wrap gap-2 mb-8">
            {SUPPORTED_LANGUAGES.filter(l => l.code !== 'english').map(lang => (
              <span key={lang.code} className="bg-white rounded-full px-3 py-1 text-sm font-medium shadow-sm border border-gray-100">
                {lang.nativeLabel}
              </span>
            ))}
          </div>

          {/* VRS explanation */}
          <div className="bg-white rounded-xl p-4 border border-orange-100 shadow-sm">
            <div className="text-xs text-gray-500 mb-1 font-mono">VERNACULAR RESONANCE SCORE</div>
            <div className="font-mono text-sm text-gray-700">
              VRS = (Engagement × <span className="text-resona-saffron font-bold">Language Score</span> × Share Velocity) ÷ Time
            </div>
            <p className="text-xs text-gray-500 mt-1">Pure Kannada gets 1.0x. English gets 0.3x. Your language matters.</p>
          </div>
        </div>

        {/* Auth form */}
        <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-6">
          {/* Tab switcher */}
          <div className="flex rounded-lg bg-gray-100 p-1 mb-6">
            <button
              onClick={() => setTab('login')}
              className={`flex-1 rounded-md py-2 text-sm font-medium transition-colors ${tab === 'login' ? 'bg-white shadow text-gray-900' : 'text-gray-500'}`}
            >
              Login
            </button>
            <button
              onClick={() => setTab('register')}
              className={`flex-1 rounded-md py-2 text-sm font-medium transition-colors ${tab === 'register' ? 'bg-white shadow text-gray-900' : 'text-gray-500'}`}
            >
              Join Resona
            </button>
          </div>

          {/* Google login button */}
          <button
            onClick={() => googleLogin()}
            disabled={googleLoading}
            className="w-full flex items-center justify-center gap-3 border border-gray-200 rounded-lg py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors mb-4 disabled:opacity-50"
          >
            {/* Google SVG icon */}
            <svg className="w-5 h-5" viewBox="0 0 24 24">
              <path
                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                fill="#4285F4"
              />
              <path
                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                fill="#34A853"
              />
              <path
                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                fill="#FBBC05"
              />
              <path
                d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                fill="#EA4335"
              />
            </svg>
            {googleLoading ? 'Signing in…' : 'Continue with Google'}
          </button>

          {/* Divider */}
          <div className="flex items-center gap-3 mb-4">
            <div className="flex-1 h-px bg-gray-200" />
            <span className="text-xs text-gray-400">or</span>
            <div className="flex-1 h-px bg-gray-200" />
          </div>

          {/* Email/password forms */}
          {tab === 'login' ? (
            <form onSubmit={loginForm.handleSubmit(onLogin)} className="space-y-4">
              <input
                {...loginForm.register('email', { required: true })}
                type="email" placeholder="Email"
                className="w-full border border-gray-200 rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-resona-saffron"
              />
              <input
                {...loginForm.register('password', { required: true })}
                type="password" placeholder="Password"
                className="w-full border border-gray-200 rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-resona-saffron"
              />
              <button
                type="submit"
                disabled={loginForm.formState.isSubmitting}
                className="w-full bg-resona-saffron text-white rounded-lg py-2.5 text-sm font-medium hover:bg-orange-500 transition-colors disabled:opacity-50"
              >
                {loginForm.formState.isSubmitting ? 'Logging in…' : 'Login'}
              </button>
            </form>
          ) : (
            <form onSubmit={registerForm.handleSubmit(onRegister)} className="space-y-4">
              <input
                {...registerForm.register('full_name', { required: true })}
                placeholder="Full Name"
                className="w-full border border-gray-200 rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-resona-saffron"
              />
              <input
                {...registerForm.register('username', { required: true, minLength: 3 })}
                placeholder="Username"
                className="w-full border border-gray-200 rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-resona-saffron"
              />
              <input
                {...registerForm.register('email', { required: true })}
                type="email" placeholder="Email"
                className="w-full border border-gray-200 rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-resona-saffron"
              />
              <input
                {...registerForm.register('password', { required: true, minLength: 8 })}
                type="password" placeholder="Password (min 8 chars)"
                className="w-full border border-gray-200 rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-resona-saffron"
              />
              <button
                type="submit"
                disabled={registerForm.formState.isSubmitting}
                className="w-full bg-resona-saffron text-white rounded-lg py-2.5 text-sm font-medium hover:bg-orange-500 transition-colors disabled:opacity-50"
              >
                {registerForm.formState.isSubmitting ? 'Creating account…' : 'Create Account'}
              </button>
            </form>
          )}
        </div>
      </div>
    </div>
  )
}
