import { useState } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import type { ChatMessage } from '../../lib/types'
import { User, Bot, Copy, Check } from 'lucide-react'

interface ChatMessageBubbleProps {
  message: ChatMessage
}

export function ChatMessageBubble({ message }: ChatMessageBubbleProps) {
  const [copied, setCopied] = useState(false)
  const isUser = message.role === 'user'

  function handleCopy() {
    navigator.clipboard.writeText(message.content || '').then(() => {
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    })
  }

  return (
    <div className={`flex gap-3 ${isUser ? 'flex-row-reverse' : ''} group`}>
      <div className={`flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center ${
        isUser ? 'bg-[var(--color-navy)] text-white' : 'bg-slate-100 text-slate-600'
      }`}>
        {isUser ? <User className="w-4 h-4" /> : <Bot className="w-4 h-4" />}
      </div>
      <div className={`relative max-w-[80%] rounded-lg px-4 py-3 ${
        isUser
          ? 'bg-[var(--color-navy)] text-white'
          : 'bg-white border border-slate-200 text-slate-800'
      }`}>
        {isUser ? (
          <p className="text-sm">{message.content}</p>
        ) : message.streaming && !message.content ? (
          <div className="flex items-center gap-1.5 py-1">
            <span className="w-1.5 h-1.5 bg-slate-400 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
            <span className="w-1.5 h-1.5 bg-slate-400 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
            <span className="w-1.5 h-1.5 bg-slate-400 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
          </div>
        ) : (
          <div className={`chat-markdown text-sm leading-relaxed ${message.streaming ? 'streaming-cursor' : ''}`}>
            <ReactMarkdown remarkPlugins={[remarkGfm]}>
              {message.content || ''}
            </ReactMarkdown>
          </div>
        )}

        {/* Copy button for assistant messages */}
        {!isUser && message.content && !message.streaming && (
          <button
            onClick={handleCopy}
            className="absolute -bottom-3 right-2 p-1 rounded bg-white border border-slate-200 text-slate-400 hover:text-slate-600 opacity-0 group-hover:opacity-100 transition-opacity shadow-sm"
            title="Copy to clipboard"
          >
            {copied ? <Check className="w-3.5 h-3.5 text-green-500" /> : <Copy className="w-3.5 h-3.5" />}
          </button>
        )}
      </div>
    </div>
  )
}
