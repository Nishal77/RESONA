import { formatDistanceToNow } from 'date-fns'

export function timeAgo(dateStr: string): string {
  return formatDistanceToNow(new Date(dateStr), { addSuffix: true })
}
