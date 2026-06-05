interface Props {
  avatarUrl?: string | null
  username: string
  size?: 'xs' | 'sm' | 'md' | 'lg'
}

const sizeMap = {
  xs: 'h-6 w-6 text-xs',
  sm: 'h-8 w-8 text-sm',
  md: 'h-10 w-10 text-base',
  lg: 'h-16 w-16 text-xl',
}

export function UserAvatar({ avatarUrl, username, size = 'md' }: Props) {
  const initials = username.slice(0, 2).toUpperCase()
  return avatarUrl ? (
    <img
      src={avatarUrl}
      alt={username}
      className={`${sizeMap[size]} rounded-full object-cover bg-gray-200`}
    />
  ) : (
    <div
      className={`${sizeMap[size]} rounded-full bg-resona-saffron text-white flex items-center justify-center font-semibold`}
    >
      {initials}
    </div>
  )
}
