import client from './client'

export const uploadApi = {
  uploadImage: (file: File) => {
    const form = new FormData()
    form.append('file', file)
    return client.post<{ success: boolean; data: { url: string } }>('/api/v1/upload/image', form, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  uploadVideo: (file: File) => {
    const form = new FormData()
    form.append('file', file)
    return client.post<{ success: boolean; data: { url: string } }>('/api/v1/upload/video', form, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
}
