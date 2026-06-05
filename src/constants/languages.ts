import type { Language } from '../types/user.types'

export interface LanguageOption {
  code: Language
  label: string
  nativeLabel: string
  flag: string
}

export const SUPPORTED_LANGUAGES: LanguageOption[] = [
  { code: 'kannada',   label: 'Kannada',   nativeLabel: 'ಕನ್ನಡ',      flag: '🇮🇳' },
  { code: 'tamil',     label: 'Tamil',     nativeLabel: 'தமிழ்',       flag: '🇮🇳' },
  { code: 'telugu',    label: 'Telugu',    nativeLabel: 'తెలుగు',      flag: '🇮🇳' },
  { code: 'malayalam', label: 'Malayalam', nativeLabel: 'മലയാളം',    flag: '🇮🇳' },
  { code: 'hindi',     label: 'Hindi',     nativeLabel: 'हिन्दी',       flag: '🇮🇳' },
  { code: 'english',   label: 'English',   nativeLabel: 'English',     flag: '🌐' },
]

export const LANGUAGE_MAP: Record<Language, string> = {
  kannada:   'ಕನ್ನಡ',
  tamil:     'தமிழ்',
  telugu:    'తెలుగు',
  malayalam: 'മലയാളം',
  hindi:     'हिन्दी',
  english:   'English',
}

export const LANGUAGE_LOCALITY_SCORES: Record<Language, number> = {
  kannada: 1.0, tamil: 1.0, telugu: 1.0,
  malayalam: 1.0, hindi: 1.0, english: 0.3,
}
