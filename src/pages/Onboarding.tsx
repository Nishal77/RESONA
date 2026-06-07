import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import { SUPPORTED_LANGUAGES } from '../constants/languages'
import { usersApi } from '../api/users.api'
import { communitiesApi } from '../api/communities.api'
import { useAuthStore } from '../store/auth.store'
import type { Language } from '../types/user.types'

const STARTER_COMMUNITIES = [
  { slug: 'kannada-memes',     language: 'kannada',   name: 'Kannada Memes' },
  { slug: 'mangalore-vibes',   language: 'kannada',   name: 'Mangalore Vibes' },
  { slug: 'bharat-builders',   language: 'kannada',   name: 'Bharat Builders' },
  { slug: 'tamil-poetry',      language: 'tamil',     name: 'Tamil Poetry' },
  { slug: 'telugu-trends',     language: 'telugu',    name: 'Telugu Trends' },
  { slug: 'malayalam-humour',  language: 'malayalam', name: 'Malayalam Humour' },
  { slug: 'hindi-shayari',     language: 'hindi',     name: 'Hindi Shayari' },
  { slug: 'coastal-karnataka', language: 'kannada',   name: 'Coastal Karnataka' },
]

export default function Onboarding() {
  const navigate = useNavigate()
  const { user: _user, setUser } = useAuthStore()
  const [step, setStep] = useState(1)
  const [language, setLanguage] = useState<Language>('kannada')
  const [state, setState] = useState('')
  const [city, setCity] = useState('')
  const [joinedSlugs, setJoinedSlugs] = useState<string[]>([])
  const [loading, setLoading] = useState(false)

  const suggestedCommunities = STARTER_COMMUNITIES.filter(
    c => c.language === language || c.slug === 'bharat-builders',
  )

  const toggleCommunity = (slug: string) => {
    setJoinedSlugs(prev =>
      prev.includes(slug) ? prev.filter(s => s !== slug) : [...prev, slug],
    )
  }

  const finish = async () => {
    setLoading(true)
    try {
      const updated = await usersApi.updateMe({
        primary_language: language,
        state: state || undefined,
        city: city || undefined,
        onboarding_completed: true,
      } as any)
      setUser(updated.data.data)

      // Join selected communities
      for (const slug of joinedSlugs) {
        try {
          const comm = await communitiesApi.getCommunity(slug)
          await communitiesApi.join(comm.data.data.id)
        } catch { /* community may not exist yet — non-fatal */ }
      }

      toast.success('Welcome to Resona!')
      navigate('/home')
    } catch {
      toast.error('Setup failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-resona-cream to-white flex items-center justify-center px-4">
      <div className="w-full max-w-lg bg-white rounded-2xl shadow-lg border border-gray-100 p-6">
        {/* Progress */}
        <div className="flex items-center gap-2 mb-6">
          {[1, 2, 3].map(s => (
            <div
              key={s}
              className={`flex-1 h-1.5 rounded-full transition-colors ${s <= step ? 'bg-resona-saffron' : 'bg-gray-200'}`}
            />
          ))}
        </div>

        {step === 1 && (
          <div>
            <h2 className="text-xl font-bold text-gray-900 mb-1">Your Language</h2>
            <p className="text-sm text-gray-500 mb-4">Pick the language you create in most</p>
            <div className="grid grid-cols-2 gap-3">
              {SUPPORTED_LANGUAGES.filter(l => l.code !== 'english').map(lang => (
                <button
                  key={lang.code}
                  onClick={() => setLanguage(lang.code)}
                  className={`rounded-xl border-2 p-4 text-left transition-all ${
                    language === lang.code
                      ? 'border-resona-saffron bg-orange-50'
                      : 'border-gray-200 hover:border-gray-300'
                  }`}
                >
                  <div className="text-xl mb-1">{lang.flag}</div>
                  <div className="font-medium text-sm">{lang.nativeLabel}</div>
                  <div className="text-xs text-gray-500">{lang.label}</div>
                </button>
              ))}
            </div>
            <button
              onClick={() => setStep(2)}
              className="mt-6 w-full bg-resona-saffron text-white rounded-lg py-2.5 text-sm font-medium hover:bg-orange-500"
            >
              Continue →
            </button>
          </div>
        )}

        {step === 2 && (
          <div>
            <h2 className="text-xl font-bold text-gray-900 mb-1">Where Are You From?</h2>
            <p className="text-sm text-gray-500 mb-4">Optional — helps with local discovery</p>
            <input
              value={state}
              onChange={e => setState(e.target.value)}
              placeholder="State (e.g. Karnataka)"
              className="w-full border border-gray-200 rounded-lg px-3 py-2.5 text-sm mb-3 focus:outline-none focus:ring-2 focus:ring-resona-saffron"
            />
            <input
              value={city}
              onChange={e => setCity(e.target.value)}
              placeholder="City (e.g. Mangalore)"
              className="w-full border border-gray-200 rounded-lg px-3 py-2.5 text-sm mb-6 focus:outline-none focus:ring-2 focus:ring-resona-saffron"
            />
            <div className="flex gap-3">
              <button onClick={() => setStep(1)} className="flex-1 border border-gray-200 rounded-lg py-2.5 text-sm">← Back</button>
              <button onClick={() => setStep(3)} className="flex-1 bg-resona-saffron text-white rounded-lg py-2.5 text-sm font-medium">Continue →</button>
            </div>
          </div>
        )}

        {step === 3 && (
          <div>
            <h2 className="text-xl font-bold text-gray-900 mb-1">Join Communities</h2>
            <p className="text-sm text-gray-500 mb-4">Pick at least 3 to get started</p>
            <div className="space-y-2 mb-6">
              {suggestedCommunities.map(c => (
                <button
                  key={c.slug}
                  onClick={() => toggleCommunity(c.slug)}
                  className={`w-full flex items-center justify-between rounded-xl border-2 px-4 py-3 transition-all ${
                    joinedSlugs.includes(c.slug) ? 'border-resona-saffron bg-orange-50' : 'border-gray-200'
                  }`}
                >
                  <span className="text-sm font-medium">{c.name}</span>
                  {joinedSlugs.includes(c.slug) && <span className="text-resona-saffron text-xs font-medium">Joined ✓</span>}
                </button>
              ))}
            </div>
            <div className="flex gap-3">
              <button onClick={() => setStep(2)} className="flex-1 border border-gray-200 rounded-lg py-2.5 text-sm">← Back</button>
              <button
                onClick={finish}
                disabled={loading}
                className="flex-1 bg-resona-saffron text-white rounded-lg py-2.5 text-sm font-medium disabled:opacity-50"
              >
                {loading ? 'Setting up…' : 'Enter Resona →'}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
