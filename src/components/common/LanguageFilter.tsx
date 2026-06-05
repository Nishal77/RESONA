import { SUPPORTED_LANGUAGES } from '../../constants/languages'
import type { Language } from '../../types/user.types'

interface Props {
  value: Language | ''
  onChange: (lang: Language | '') => void
}

export function LanguageFilter({ value, onChange }: Props) {
  return (
    <div className="flex items-center gap-1 overflow-x-auto pb-1 scrollbar-hide">
      <button
        onClick={() => onChange('')}
        className={`shrink-0 rounded-full px-3 py-1 text-sm font-medium transition-colors ${
          value === '' ? 'bg-resona-saffron text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
        }`}
      >
        All
      </button>
      {SUPPORTED_LANGUAGES.filter(l => l.code !== 'english').map(lang => (
        <button
          key={lang.code}
          onClick={() => onChange(lang.code)}
          className={`shrink-0 rounded-full px-3 py-1 text-sm font-medium transition-colors ${
            value === lang.code ? 'bg-resona-saffron text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
          }`}
        >
          {lang.nativeLabel}
        </button>
      ))}
    </div>
  )
}
