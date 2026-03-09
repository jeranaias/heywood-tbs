import { useState, useRef } from 'react'
import { Send, Loader2, Trash2 } from 'lucide-react'
import { ChatHistory } from '../components/chat/ChatHistory'
import { ChatSessionSidebar } from '../components/chat/ChatSessionSidebar'
import { useChatContext } from '../hooks/ChatContext'
import type { ChatMessage } from '../lib/types'

export function ChatPage() {
  const { messages, loading, sendMessage, clearMessages, loadMessages } = useChatContext()
  const [input, setInput] = useState('')
  const [sessionId, setSessionId] = useState<string | undefined>()
  const inputRef = useRef<HTMLInputElement>(null)

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (input.trim() && !loading) {
      sendMessage(input.trim())
      setInput('')
    }
  }

  function handleSuggestedPrompt(prompt: string) {
    if (!loading) {
      sendMessage(prompt)
    }
  }

  function handleLoadSession(id: string, msgs: Array<{ role: string; content: string }>) {
    setSessionId(id)
    const chatMsgs: ChatMessage[] = msgs.map(m => ({
      role: m.role as 'user' | 'assistant',
      content: m.content,
    }))
    loadMessages(chatMsgs)
  }

  return (
    <div className="flex gap-4 h-[calc(100vh-8rem)]">
      {/* Chat session sidebar */}
      <ChatSessionSidebar
        onLoadSession={handleLoadSession}
        currentSessionId={sessionId}
      />

      {/* Main chat area */}
      <div className="flex-1 flex flex-col min-w-0">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h2 className="text-xl font-bold text-slate-900">Chat with Heywood</h2>
            <p className="text-sm text-slate-500">Your digital staff officer for TBS</p>
          </div>
          {messages.length > 0 && (
            <button
              onClick={() => { clearMessages(); setSessionId(undefined) }}
              className="flex items-center gap-1.5 px-3 py-1.5 text-sm text-slate-500 hover:text-slate-700 border border-slate-200 rounded-lg hover:bg-slate-50"
            >
              <Trash2 className="w-4 h-4" /> Clear
            </button>
          )}
        </div>

        <div className="flex-1 bg-white rounded-lg border border-slate-200 flex flex-col overflow-hidden">
          <ChatHistory
            messages={messages}
            loading={loading}
            onSuggestedPrompt={handleSuggestedPrompt}
          />

          <div className="border-t border-slate-200 p-4">
            <form onSubmit={handleSubmit} className="flex items-center gap-3">
              <input
                ref={inputRef}
                type="text"
                value={input}
                onChange={e => setInput(e.target.value)}
                placeholder="Ask about students, counseling, AARs, scenarios..."
                disabled={loading}
                className="flex-1 px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-lg text-sm placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20 focus:border-[var(--color-navy)] disabled:opacity-50"
              />
              <button
                type="submit"
                disabled={!input.trim() || loading}
                className="flex-shrink-0 p-2.5 bg-[var(--color-navy)] text-white rounded-lg hover:bg-[var(--color-navy-light)] disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                {loading ? <Loader2 className="w-5 h-5 animate-spin" /> : <Send className="w-5 h-5" />}
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  )
}
