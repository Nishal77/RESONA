import axios from 'axios'
import { useAuthStore } from '../store/auth.store'

const client = axios.create({
  baseURL: import.meta.env.VITE_API_URL ?? 'http://localhost:8080',
  withCredentials: true, // send httpOnly refresh_token cookie
})

client.interceptors.request.use((config) => {
  const token = useAuthStore.getState().accessToken
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

let isRefreshing = false
let queue: Array<(token: string) => void> = []

client.interceptors.response.use(
  (res) => res,
  async (error) => {
    const original = error.config
    if (error.response?.status === 401 && !original._retry) {
      original._retry = true

      if (isRefreshing) {
        return new Promise((resolve) => {
          queue.push((token) => {
            original.headers.Authorization = `Bearer ${token}`
            resolve(client(original))
          })
        })
      }

      isRefreshing = true
      try {
        const { data } = await axios.post(
          `${import.meta.env.VITE_API_URL ?? 'http://localhost:8080'}/api/v1/auth/refresh`,
          {},
          { withCredentials: true },
        )
        const newToken: string = data.data.access_token
        useAuthStore.getState().setAccessToken(newToken)
        queue.forEach((cb) => cb(newToken))
        queue = []
        original.headers.Authorization = `Bearer ${newToken}`
        return client(original)
      } catch {
        useAuthStore.getState().logout()
        window.location.href = '/'
        return Promise.reject(error)
      } finally {
        isRefreshing = false
      }
    }
    return Promise.reject(error)
  },
)

export default client
