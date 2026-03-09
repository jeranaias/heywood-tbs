import { useEffect, useRef } from 'react'
import type { ChatMessage } from '../../lib/types'
import { ChatMessageBubble } from './ChatMessage'
import { SuggestedPrompts } from './SuggestedPrompts'
import { Loader2 } from 'lucide-react'

interface ChatHistoryProps {
  messages: ChatMessage[]
  loading: boolean
  onSuggestedPrompt?: (prompt: string) => void
}

export function ChatHistory({ messages, loading, onSuggestedPrompt }: ChatHistoryProps) {
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages, loading])

  if (messages.length === 0 && !loading) {
    return (
      <div className="flex-1 flex items-center justify-center p-8 text-center">
        <div>
          <div className="w-16 h-16 rounded-full bg-slate-100 flex items-center justify-center mx-auto mb-4">
            <span className="text-2xl">🎖️</span>
          </div>
          <h3 className="text-lg font-semibold text-slate-800 mb-2">Heywood is ready</h3>
          <p className="text-sm text-slate-500 max-w-sm">
            Ask about student performance, prepare counseling outlines, analyze AARs,
            or generate training scenarios.
          </p>
          {onSuggestedPrompt && (
            <SuggestedPrompts onSelect={onSuggestedPrompt} />
          )}
        </div>
      </div>
    )
  }

  return (
    <div className="flex-1 overflow-y-auto p-4 space-y-4 chat-scroll">
      {messages.map((msg, i) => (
        <ChatMessageBubble key={i} message={msg} />
      ))}
      {loading && (
        <div className="flex items-center gap-2 text-slate-500 text-sm">
          <Loader2 className="w-4 h-4 animate-spin" />
          Heywood is thinking...
        </div>
      )}
      <div ref={bottomRef} />
    </div>
  )
}
