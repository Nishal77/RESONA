import { LANGUAGE_MAP } from '../../constants/languages'
import type { Language } from '../../types/user.types'

interface Props {
  language: Language | string
  size?: 'xs' | 'sm'
}

const colorMap: Record<string, string> = {
  kannada:   'bg-yellow-100 text-yellow-700',
  tamil:     'bg-red-100 text-red-700',
  telugu:    'bg-blue-100 text-blue-700',
  malayalam: 'bg-green-100 text-green-700',
  hindi:     'bg-orange-100 text-orange-700',
  english:   'bg-gray-100 text-gray-600',
}

export function LanguageBadge({ language, size = 'sm' }: Props) {
  const label = LANGUAGE_MAP[language as Language] ?? language
  const color = colorMap[language] ?? 'bg-gray-100 text-gray-600'
  const sizeClass = size === 'xs' ? 'text-[10px] px-1.5 py-0.5' : 'text-xs px-2 py-0.5'

  return (
    <span className={`inline-block rounded-full font-medium ${sizeClass} ${color}`}>
      {label}
    </span>
  )
}
