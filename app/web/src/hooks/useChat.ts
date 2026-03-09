import { useState, useCallback, useRef } from 'react'
import type { ChatMessage } from '../lib/types'
import { api } from '../lib/api'

export function useChat() {
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [loading, setLoading] = useState(false)
  const abortRef = useRef<AbortController | null>(null)

  const sendMessage = useCallback(async (content: string) => {
    if (!content.trim() || loading) return

    // Cancel any in-flight stream
    abortRef.current?.abort()
    const controller = new AbortController()
    abortRef.current = controller

    const userMsg: ChatMessage = { role: 'user', content }
    setMessages(prev => [...prev, userMsg])
    setLoading(true)

    // Add a placeholder streaming message
    const streamMsg: ChatMessage = { role: 'assistant', content: '', streaming: true }
    setMessages(prev => [...prev, streamMsg])

    try {
      const history = messages.slice(-20)

      await api.chatStream(
        content,
        history,
        // onChunk: append to the last message
        (chunk: string) => {
          setMessages(prev => {
            const updated = [...prev]
            const last = updated[updated.length - 1]
            updated[updated.length - 1] = { ...last, content: last.content + chunk }
            return updated
          })
        },
        // onDone: clear streaming flag
        () => {
          setMessages(prev => {
            const updated = [...prev]
            const last = updated[updated.length - 1]
            updated[updated.length - 1] = { ...last, streaming: false }
            return updated
          })
        },
        controller.signal,
      )
    } catch (err) {
      if ((err as Error).name === 'AbortError') return
      setMessages(prev => {
        const updated = [...prev]
        const last = updated[updated.length - 1]
        if (last.streaming) {
          updated[updated.length - 1] = {
            role: 'assistant',
            content: last.content || 'Sorry, I encountered an error. Please try again.',
            streaming: false,
          }
        }
        return updated
      })
    } finally {
      setLoading(false)
      abortRef.current = null
    }
  }, [messages, loading])

  const clearMessages = useCallback(() => {
    abortRef.current?.abort()
    setMessages([])
  }, [])

  const loadMessages = useCallback((msgs: ChatMessage[]) => {
    abortRef.current?.abort()
    setMessages(msgs)
  }, [])

  return { messages, loading, sendMessage, clearMessages, loadMessages }
}
