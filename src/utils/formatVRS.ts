export function formatVRS(score: number): string {
  return score.toFixed(3)
}

export function vrsColor(score: number): 'green' | 'amber' | 'gray' {
  if (score >= 0.8) return 'green'
  if (score >= 0.4) return 'amber'
  return 'gray'
}

export function vrsBgClass(score: number): string {
  const c = vrsColor(score)
  return {
    green: 'bg-green-100 text-green-700 border-green-300',
    amber: 'bg-amber-100 text-amber-700 border-amber-300',
    gray:  'bg-gray-100 text-gray-500 border-gray-300',
  }[c]
}
