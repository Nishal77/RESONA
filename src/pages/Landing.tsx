import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import toast from 'react-hot-toast'
import { GoogleLogin, type CredentialResponse } from '@react-oauth/google'
import { authApi } from '../api/auth.api'
import { useAuthStore } from '../store/auth.store'
import { SUPPORTED_LANGUAGES } from '../constants/languages'

type Tab = 'login' | 'register'
interface LoginForm    { email: string; password: string }
interface RegisterForm { username: string; email: string; password: string; full_name: string }

export default function Landing() {
  const [tab, setTab] = useState<Tab>('login')
  const navigate = useNavigate()
  const { setAuth, user } = useAuthStore()

  const loginForm    = useForm<LoginForm>()
  const registerForm = useForm<RegisterForm>()

  // Already logged in → skip landing
  useEffect(() => {
    if (user) {
      navigate(user.onboarding_completed ? '/home' : '/onboarding', { replace: true })
    }
  }, [user, navigate])

  const redirectAfterAuth = (u: { onboarding_completed: boolean }) => {
    navigate(u.onboarding_completed ? '/home' : '/onboarding', { replace: true })
  }

  const onLogin = async (data: LoginForm) => {
    try {
      const res = await authApi.login(data)
      const { access_token, user: u } = res.data.data
      setAuth(u, access_token)
      toast.success(`Welcome back, ${u.username}!`)
      redirectAfterAuth(u)
    } catch (err: any) {
      toast.error(err.response?.data?.error ?? 'Login failed')
    }
  }

  const onRegister = async (data: RegisterForm) => {
    try {
      const res = await authApi.register(data)
      const { access_token, user: u } = res.data.data
      setAuth(u, access_token)
      toast.success('Welcome to Resona!')
      redirectAfterAuth(u)
    } catch (err: any) {
      toast.error(err.response?.data?.error ?? 'Registration failed')
    }
  }

  // Google One Tap — credential is an id_token JWT, sent directly to backend
  const onGoogleSuccess = async (credentialResponse: CredentialResponse) => {
    if (!credentialResponse.credential) {
      toast.error('Google login failed — no credential returned')
      return
    }
    try {
      const res = await authApi.googleAuth(credentialResponse.credential)
      const { access_token, user: u } = res.data.data
      setAuth(u, access_token)
      toast.success(`Welcome, ${u.full_name ?? u.username}!`)
      redirectAfterAuth(u)
    } catch (err: any) {
      toast.error(err.response?.data?.error ?? 'Google login failed')
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-resona-cream via-white to-orange-50">
      <div className="mx-auto max-w-5xl px-4 pt-12 pb-8 grid md:grid-cols-2 gap-12 items-center">

        {/* Left — Hero */}
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

          <div className="flex flex-wrap gap-2 mb-8">
            {SUPPORTED_LANGUAGES.filter(l => l.code !== 'english').map(lang => (
              <span key={lang.code} className="bg-white rounded-full px-3 py-1 text-sm font-medium shadow-sm border border-gray-100">
                {lang.nativeLabel}
              </span>
            ))}
          </div>

          <div className="bg-white rounded-xl p-4 border border-orange-100 shadow-sm">
            <div className="text-xs text-gray-500 mb-1 font-mono">VERNACULAR RESONANCE SCORE</div>
            <div className="font-mono text-sm text-gray-700">
              VRS = (Engagement × <span className="text-resona-saffron font-bold">Language Score</span> × Share Velocity) ÷ Time
            </div>
            <p className="text-xs text-gray-500 mt-1">Pure Kannada gets 1.0x. English gets 0.3x. Your language matters.</p>
          </div>
        </div>

        {/* Right — Auth */}
        <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-6">

          {/* Tab switcher */}
          <div className="flex rounded-lg bg-gray-100 p-1 mb-5">
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

          {/* Google Login — rendered by Google SDK, most reliable */}
          <div className="flex justify-center mb-4">
            <GoogleLogin
              onSuccess={onGoogleSuccess}
              onError={() => toast.error('Google login failed')}
              theme="outline"
              size="large"
              width="340"
              text="continue_with"
              shape="rectangular"
            />
          </div>

          {/* Divider */}
          <div className="flex items-center gap-3 mb-4">
            <div className="flex-1 h-px bg-gray-200" />
            <span className="text-xs text-gray-400">or</span>
            <div className="flex-1 h-px bg-gray-200" />
          </div>

          {/* Login form */}
          {tab === 'login' ? (
            <form onSubmit={loginForm.handleSubmit(onLogin)} className="space-y-4">
              <input
                {...loginForm.register('email', { required: true })}
                type="email"
                placeholder="Email"
                className="w-full border border-gray-200 rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-resona-saffron"
              />
              <input
                {...loginForm.register('password', { required: true })}
                type="password"
                placeholder="Password"
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
                type="email"
                placeholder="Email"
                className="w-full border border-gray-200 rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-resona-saffron"
              />
              <input
                {...registerForm.register('password', { required: true, minLength: 8 })}
                type="password"
                placeholder="Password (min 8 chars)"
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
