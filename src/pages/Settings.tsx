import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import { usersApi } from '../api/users.api'
import { authApi } from '../api/auth.api'
import { useAuthStore } from '../store/auth.store'
import { SUPPORTED_LANGUAGES } from '../constants/languages'
import type { Language } from '../types/user.types'

interface SettingsForm {
  full_name: string
  bio: string
  state: string
  city: string
  primary_language: Language
}

export default function Settings() {
  const navigate = useNavigate()
  const { user, setUser, logout } = useAuthStore()
  const { register, handleSubmit, reset, formState: { isSubmitting } } = useForm<SettingsForm>()

  useEffect(() => {
    if (user) {
      reset({
        full_name: user.full_name ?? '',
        bio: user.bio ?? '',
        state: user.state ?? '',
        city: user.city ?? '',
        primary_language: user.primary_language,
      })
    }
  }, [user, reset])

  const onSave = async (data: SettingsForm) => {
    try {
      const res = await usersApi.updateMe({
        full_name: data.full_name || undefined,
        bio: data.bio || undefined,
        state: data.state || undefined,
        city: data.city || undefined,
        primary_language: data.primary_language,
      } as any)
      setUser(res.data.data)
      toast.success('Profile updated')
    } catch {
      toast.error('Update failed')
    }
  }

  const handleDeleteAccount = async () => {
    const confirmed = window.prompt('Type DELETE to confirm account deletion:')
    if (confirmed !== 'DELETE') return
    try {
      await authApi.logout()
      logout()
      navigate('/')
      toast.success('Account deleted')
    } catch {
      toast.error('Deletion failed')
    }
  }

  if (!user) return null

  return (
    <div className="max-w-lg mx-auto">
      <h1 className="text-xl font-bold text-gray-900 mb-6">Settings</h1>

      <form onSubmit={handleSubmit(onSave)} className="space-y-5">
        {/* Profile */}
        <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-4 space-y-4">
          <h2 className="font-semibold text-gray-800">Profile</h2>
          <div>
            <label className="text-xs text-gray-500 block mb-1">Full Name</label>
            <input
              {...register('full_name')}
              className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-resona-saffron"
            />
          </div>
          <div>
            <label className="text-xs text-gray-500 block mb-1">Bio (max 160 chars)</label>
            <textarea
              {...register('bio')}
              maxLength={160}
              rows={3}
              className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-resona-saffron"
            />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="text-xs text-gray-500 block mb-1">State</label>
              <input {...register('state')} className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-resona-saffron" />
            </div>
            <div>
              <label className="text-xs text-gray-500 block mb-1">City</label>
              <input {...register('city')} className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-resona-saffron" />
            </div>
          </div>
        </div>

        {/* Language */}
        <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-4">
          <h2 className="font-semibold text-gray-800 mb-3">Language</h2>
          <label className="text-xs text-gray-500 block mb-1">Primary Language</label>
          <select
            {...register('primary_language')}
            className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-resona-saffron"
          >
            {SUPPORTED_LANGUAGES.map(l => (
              <option key={l.code} value={l.code}>{l.label} — {l.nativeLabel}</option>
            ))}
          </select>
        </div>

        <button
          type="submit"
          disabled={isSubmitting}
          className="w-full bg-resona-saffron text-white rounded-xl py-3 text-sm font-medium hover:bg-orange-500 transition-colors disabled:opacity-50"
        >
          {isSubmitting ? 'Saving…' : 'Save Changes'}
        </button>
      </form>

      {/* Danger zone */}
      <div className="mt-6 bg-white rounded-xl border border-red-100 shadow-sm p-4">
        <h2 className="font-semibold text-red-600 mb-2">Danger Zone</h2>
        <p className="text-xs text-gray-500 mb-3">
          Deleting your account is permanent. All posts, comments, and engagement data will be removed.
        </p>
        <button
          onClick={handleDeleteAccount}
          className="w-full border border-red-200 text-red-500 rounded-lg py-2 text-sm hover:bg-red-50 transition-colors"
        >
          Delete Account
        </button>
      </div>
    </div>
  )
}
