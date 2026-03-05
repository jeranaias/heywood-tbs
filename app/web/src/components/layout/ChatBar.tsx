import { useState, useRef, useEffect } from 'react'
import { Send, MessageSquare, Loader2 } from 'lucide-react'

interface ChatBarProps {
  onSend: (message: string) => void
  loading: boolean
  onExpand: () => void
  expanded: boolean
  hasMessages: boolean
}

export function ChatBar({ onSend, loading, onExpand, expanded, hasMessages }: ChatBarProps) {
  const [input, setInput] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    function handleKeydown(e: KeyboardEvent) {
      // Cmd/Ctrl+K to focus chat
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault()
        inputRef.current?.focus()
      }
    }
    document.addEventListener('keydown', handleKeydown)
    return () => document.removeEventListener('keydown', handleKeydown)
  }, [])

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (input.trim() && !loading) {
      onSend(input.trim())
      setInput('')
      if (!expanded) onExpand()
    }
  }

  return (
    <div className="border-t border-slate-200 bg-white px-4 py-3">
      <form onSubmit={handleSubmit} className="flex items-center gap-3 max-w-4xl mx-auto">
        <button
          type="button"
          onClick={onExpand}
          className={`flex-shrink-0 p-2 rounded-lg transition-colors ${
            hasMessages
              ? 'text-[var(--color-navy)] bg-blue-50 hover:bg-blue-100'
              : 'text-slate-400 hover:text-slate-600 hover:bg-slate-100'
          }`}
          title="Open chat panel"
        >
          <MessageSquare className="w-5 h-5" />
        </button>
        <div className="flex-1 relative">
          <input
            ref={inputRef}
            type="text"
            value={input}
            onChange={e => setInput(e.target.value)}
            placeholder="Ask Heywood anything... (Ctrl+K)"
            disabled={loading}
            className="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-lg text-sm placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20 focus:border-[var(--color-navy)] disabled:opacity-50"
          />
        </div>
        <button
          type="submit"
          disabled={!input.trim() || loading}
          className="flex-shrink-0 p-2.5 bg-[var(--color-navy)] text-white rounded-lg hover:bg-[var(--color-navy-light)] disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {loading ? (
            <Loader2 className="w-5 h-5 animate-spin" />
          ) : (
            <Send className="w-5 h-5" />
          )}
        </button>
      </form>
    </div>
  )
}
