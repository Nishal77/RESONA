import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { ImagePlus, Video, X } from 'lucide-react'
import toast from 'react-hot-toast'
import { postsApi } from '../api/posts.api'
import { uploadApi } from '../api/upload.api'
import { LanguageBadge } from '../components/common/LanguageBadge'
import { Spinner } from '../components/common/Spinner'
import { SUPPORTED_LANGUAGES } from '../constants/languages'
import type { Language } from '../types/user.types'

interface FormData {
  content_text: string
  manual_language?: Language
}

export default function CreatePost() {
  const navigate = useNavigate()
  const { register, handleSubmit, watch } = useForm<FormData>()
  const [detectedLang, setDetectedLang] = useState<string | null>(null)
  const [confidence, setConfidence] = useState<number>(0)
  const [mediaUrl, setMediaUrl] = useState<string | null>(null)
  const [mediaType, setMediaType] = useState<'image' | 'video' | null>(null)
  const [uploading, setUploading] = useState(false)
  const [submitting, setSubmitting] = useState(false)

  const content = watch('content_text')

  // Debounced language detection preview (client-side hint only)
  useEffect(() => {
    if (!content || content.length < 10) { setDetectedLang(null); return }
    const timer = setTimeout(() => {
      // Basic heuristic for preview — actual detection happens server-side
      const kannadaRange = /[ಀ-೿]/
      const tamilRange = /[஀-௿]/
      const teluguRange = /[ఀ-౿]/
      const malayalamRange = /[ഀ-ൿ]/
      const devanagari = /[ऀ-ॿ]/

      if (kannadaRange.test(content)) setDetectedLang('kannada')
      else if (tamilRange.test(content)) setDetectedLang('tamil')
      else if (teluguRange.test(content)) setDetectedLang('telugu')
      else if (malayalamRange.test(content)) setDetectedLang('malayalam')
      else if (devanagari.test(content)) setDetectedLang('hindi')
      else setDetectedLang('english')
    }, 800)
    return () => clearTimeout(timer)
  }, [content])

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    setUploading(true)
    try {
      const isVideo = file.type.startsWith('video/')
      const res = isVideo ? await uploadApi.uploadVideo(file) : await uploadApi.uploadImage(file)
      setMediaUrl(res.data.data.url)
      setMediaType(isVideo ? 'video' : 'image')
      toast.success('Media uploaded')
    } catch {
      toast.error('Upload failed')
    } finally {
      setUploading(false)
    }
  }

  const onSubmit = async (data: FormData) => {
    if (!data.content_text?.trim() && !mediaUrl) {
      toast.error('Add some content or media')
      return
    }
    setSubmitting(true)
    try {
      const res = await postsApi.createPost({
        content_text: data.content_text || undefined,
        media_url: mediaUrl || undefined,
        media_type: mediaType || undefined,
        manual_language: data.manual_language || undefined,
      })
      toast.success('Posted!')
      navigate(`/posts/${res.data.data.id}`)
    } catch {
      toast.error('Failed to post')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="max-w-2xl mx-auto">
      <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-5">
        <h1 className="text-lg font-bold text-gray-900 mb-4">Create Post</h1>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="relative">
            <textarea
              {...register('content_text')}
              placeholder="ಏನಾದರೂ ಹೇಳಿ… (What's on your mind?)"
              rows={5}
              maxLength={2000}
              className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-resona-saffron"
            />
            <div className="flex items-center justify-between mt-1 px-1">
              <div className="flex items-center gap-1.5">
                {detectedLang && (
                  <>
                    <span className="text-xs text-gray-400">Detected:</span>
                    <LanguageBadge language={detectedLang} size="xs" />
                  </>
                )}
              </div>
              <span className="text-xs text-gray-400">{content?.length ?? 0}/2000</span>
            </div>
          </div>

          {/* Manual language override */}
          {detectedLang === 'english' && (
            <div>
              <label className="text-xs text-gray-500 mb-1 block">
                Override language (if content is regional but detected as English):
              </label>
              <select
                {...register('manual_language')}
                className="border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none w-full"
              >
                <option value="">Auto-detect</option>
                {SUPPORTED_LANGUAGES.map(l => (
                  <option key={l.code} value={l.code}>{l.label} — {l.nativeLabel}</option>
                ))}
              </select>
            </div>
          )}

          {/* Media preview */}
          {mediaUrl && (
            <div className="relative">
              {mediaType === 'image' ? (
                <img src={mediaUrl} alt="" className="w-full rounded-xl max-h-60 object-cover" />
              ) : (
                <video src={mediaUrl} controls className="w-full rounded-xl max-h-60" />
              )}
              <button
                type="button"
                onClick={() => { setMediaUrl(null); setMediaType(null) }}
                className="absolute top-2 right-2 bg-white rounded-full p-1 shadow"
              >
                <X className="h-4 w-4 text-gray-600" />
              </button>
            </div>
          )}

          {/* Media upload */}
          {!mediaUrl && (
            <div className="flex gap-2">
              <label className="flex items-center gap-1.5 cursor-pointer bg-gray-100 hover:bg-gray-200 rounded-lg px-3 py-2 text-sm text-gray-600 transition-colors">
                {uploading ? <Spinner size="sm" /> : <ImagePlus className="h-4 w-4" />}
                Image
                <input type="file" accept="image/jpeg,image/png,image/webp" onChange={handleFileUpload} className="hidden" />
              </label>
              <label className="flex items-center gap-1.5 cursor-pointer bg-gray-100 hover:bg-gray-200 rounded-lg px-3 py-2 text-sm text-gray-600 transition-colors">
                {uploading ? <Spinner size="sm" /> : <Video className="h-4 w-4" />}
                Video
                <input type="file" accept="video/mp4,video/webm" onChange={handleFileUpload} className="hidden" />
              </label>
            </div>
          )}

          <button
            type="submit"
            disabled={submitting}
            className="w-full bg-resona-saffron text-white rounded-xl py-3 font-medium hover:bg-orange-500 transition-colors disabled:opacity-50"
          >
            {submitting ? 'Posting…' : 'Post to Resona'}
          </button>
        </form>
      </div>
    </div>
  )
}
