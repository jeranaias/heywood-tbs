import { useState, useEffect } from 'react'
import { MessageSquare, Trash2, Download, ChevronLeft, ChevronRight, Loader2 } from 'lucide-react'
import { api } from '../../lib/api'

interface ChatSession {
  id: string
  title: string
  userRole: string
  createdAt: string
  updatedAt: string
}

interface ChatSessionSidebarProps {
  onLoadSession: (sessionId: string, messages: Array<{ role: string; content: string }>) => void
  currentSessionId?: string
}

export function ChatSessionSidebar({ onLoadSession, currentSessionId }: ChatSessionSidebarProps) {
  const [sessions, setSessions] = useState<ChatSession[]>([])
  const [collapsed, setCollapsed] = useState(true)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!collapsed) {
      loadSessions()
    }
  }, [collapsed])

  async function loadSessions() {
    setLoading(true)
    try {
      const res = await api.getChatSessions()
      setSessions(res.sessions || [])
    } catch {
      // Chat history may not be available
    } finally {
      setLoading(false)
    }
  }

  async function handleLoad(sessionId: string) {
    try {
      const res = await api.getChatSessionMessages(sessionId)
      onLoadSession(sessionId, res.messages || [])
    } catch {
      // ignore
    }
  }

  async function handleDelete(sessionId: string) {
    try {
      await api.deleteChatSession(sessionId)
      setSessions(prev => prev.filter(s => s.id !== sessionId))
    } catch {
      // ignore
    }
  }

  function handleExport(sessionId: string) {
    window.open(`/api/v1/chat/sessions/${sessionId}/export`, '_blank')
  }

  if (collapsed) {
    return (
      <button
        onClick={() => setCollapsed(false)}
        className="flex items-center gap-1 px-2 py-1.5 text-sm text-slate-500 hover:text-slate-700 border border-slate-200 rounded-lg hover:bg-slate-50 transition-colors"
        title="Chat history"
      >
        <ChevronLeft className="w-4 h-4" />
        <MessageSquare className="w-4 h-4" />
      </button>
    )
  }

  return (
    <div className="w-64 flex-shrink-0 bg-white border border-slate-200 rounded-lg flex flex-col overflow-hidden">
      <div className="flex items-center justify-between p-3 border-b border-slate-200">
        <h3 className="text-sm font-semibold text-slate-700">History</h3>
        <button
          onClick={() => setCollapsed(true)}
          className="p-1 text-slate-400 hover:text-slate-600"
        >
          <ChevronRight className="w-4 h-4" />
        </button>
      </div>

      <div className="flex-1 overflow-y-auto">
        {loading ? (
          <div className="flex items-center justify-center p-4">
            <Loader2 className="w-4 h-4 animate-spin text-slate-400" />
          </div>
        ) : sessions.length === 0 ? (
          <p className="p-3 text-sm text-slate-400">No previous chats</p>
        ) : (
          <div className="divide-y divide-slate-100">
            {sessions.map(session => (
              <div
                key={session.id}
                className={`p-3 hover:bg-slate-50 cursor-pointer group ${
                  session.id === currentSessionId ? 'bg-slate-50 border-l-2 border-l-[var(--color-navy)]' : ''
                }`}
              >
                <button
                  onClick={() => handleLoad(session.id)}
                  className="w-full text-left"
                >
                  <p className="text-sm font-medium text-slate-700 truncate">
                    {session.title || 'Untitled chat'}
                  </p>
                  <p className="text-xs text-slate-400 mt-0.5">
                    {new Date(session.createdAt).toLocaleDateString()}
                  </p>
                </button>
                <div className="flex gap-1 mt-1 opacity-0 group-hover:opacity-100 transition-opacity">
                  <button
                    onClick={(e) => { e.stopPropagation(); handleExport(session.id) }}
                    className="p-1 text-slate-400 hover:text-slate-600"
                    title="Export"
                  >
                    <Download className="w-3.5 h-3.5" />
                  </button>
                  <button
                    onClick={(e) => { e.stopPropagation(); handleDelete(session.id) }}
                    className="p-1 text-slate-400 hover:text-red-500"
                    title="Delete"
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
