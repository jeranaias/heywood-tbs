import { useState, useRef } from 'react'
import { MessageSquare, X, Send, Loader2, Trash2 } from 'lucide-react'
import { useLocation } from 'react-router-dom'
import { ChatHistory } from './ChatHistory'
import { useChatContext } from '../../hooks/ChatContext'

export function ChatWidget() {
  const [open, setOpen] = useState(false)
  const { messages, loading, sendMessage, clearMessages } = useChatContext()
  const [input, setInput] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)
  const location = useLocation()

  // Don't show the FAB on the dedicated chat page
  if (location.pathname === '/chat') return null

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (input.trim() && !loading) {
      sendMessage(input.trim())
      setInput('')
    }
  }

  return (
    <>
      {/* Slide-out panel */}
      {open && (
        <>
          <div
            className="fixed inset-0 z-[60] bg-black/20 lg:bg-transparent"
            onClick={() => setOpen(false)}
          />
          <div className="fixed right-0 top-0 bottom-0 z-[70] w-full max-w-md bg-white shadow-2xl border-l border-slate-200 flex flex-col animate-slide-in">
            {/* Header */}
            <div className="flex items-center justify-between px-4 py-3 border-b border-slate-200 bg-[var(--color-navy)] text-white">
              <div className="flex items-center gap-2">
                <MessageSquare className="w-5 h-5" />
                <span className="font-semibold text-sm">Ask Heywood</span>
              </div>
              <div className="flex items-center gap-1">
                {messages.length > 0 && (
                  <button
                    onClick={clearMessages}
                    className="p-1.5 hover:bg-white/10 rounded transition-colors"
                    title="Clear chat"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                )}
                <button
                  onClick={() => setOpen(false)}
                  className="p-1.5 hover:bg-white/10 rounded transition-colors"
                >
                  <X className="w-4 h-4" />
                </button>
              </div>
            </div>

            {/* Chat area */}
            <div className="flex-1 overflow-hidden flex flex-col">
              <ChatHistory messages={messages} loading={loading} />
            </div>

            {/* Input */}
            <div className="border-t border-slate-200 p-3">
              <form onSubmit={handleSubmit} className="flex items-center gap-2">
                <input
                  ref={inputRef}
                  type="text"
                  value={input}
                  onChange={e => setInput(e.target.value)}
                  placeholder="Ask Heywood about this page..."
                  disabled={loading}
                  className="flex-1 px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20 focus:border-[var(--color-navy)] disabled:opacity-50"
                  autoFocus
                />
                <button
                  type="submit"
                  disabled={!input.trim() || loading}
                  className="flex-shrink-0 p-2 bg-[var(--color-navy)] text-white rounded-lg hover:bg-[var(--color-navy-light)] disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : <Send className="w-4 h-4" />}
                </button>
              </form>
            </div>
          </div>
        </>
      )}

      {/* FAB */}
      {!open && (
        <button
          onClick={() => {
            setOpen(true)
            setTimeout(() => inputRef.current?.focus(), 100)
          }}
          className="fixed bottom-6 right-6 z-[55] w-14 h-14 bg-[var(--color-navy)] text-white rounded-full shadow-lg hover:shadow-xl hover:scale-105 transition-all flex items-center justify-center group"
          title="Ask Heywood"
        >
          <MessageSquare className="w-6 h-6" />
          {messages.length > 0 && (
            <span className="absolute -top-1 -right-1 w-5 h-5 bg-red-500 text-white text-xs font-bold rounded-full flex items-center justify-center">
              {messages.filter(m => m.role === 'assistant').length}
            </span>
          )}
        </button>
      )}
    </>
  )
}
