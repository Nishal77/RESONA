import { vrsBgClass, formatVRS } from '../../utils/formatVRS'

interface Props {
  score: number
  size?: 'sm' | 'md' | 'lg'
}

export function VRSBadge({ score, size = 'sm' }: Props) {
  const sizeClass = {
    sm: 'text-xs px-1.5 py-0.5',
    md: 'text-sm px-2 py-1',
    lg: 'text-base px-3 py-1.5 font-semibold',
  }[size]

  return (
    <span
      className={`inline-flex items-center gap-1 rounded-full border font-mono ${sizeClass} ${vrsBgClass(score)}`}
      title="Vernacular Resonance Score"
    >
      <span className="opacity-60 text-[10px]">VRS</span>
      {formatVRS(score)}
    </span>
  )
}
