export interface APIResponse<T> {
  success: boolean
  data: T
  message?: string
  error?: string
  status_code?: number
}

export interface PaginationMeta {
  page: number
  limit: number
  total: number
  has_more: boolean
}

export interface PaginatedResponse<T> {
  success: boolean
  data: T[]
  meta: PaginationMeta
}
