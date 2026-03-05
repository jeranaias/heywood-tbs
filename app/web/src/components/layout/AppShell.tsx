import { useState } from 'react'
import { Outlet } from 'react-router-dom'
import { Sidebar } from './Sidebar'
import { Header } from './Header'
import { ChatBar } from './ChatBar'
import { useAuth } from '../../hooks/useAuth'
import { useChat } from '../../hooks/useChat'
import { ChatHistory } from '../chat/ChatHistory'

export function AppShell() {
  const { auth } = useAuth()
  const chat = useChat()
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [chatExpanded, setChatExpanded] = useState(false)

  return (
    <div className="flex h-screen bg-slate-50">
      {/* Mobile sidebar overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <div className={`
        fixed inset-y-0 left-0 z-50 w-64 transform transition-transform lg:relative lg:translate-x-0
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
      `}>
        <Sidebar role={auth.role} onClose={() => setSidebarOpen(false)} />
      </div>

      {/* Main content area */}
      <div className="flex flex-1 flex-col min-w-0">
        <Header
          onMenuClick={() => setSidebarOpen(true)}
        />

        {/* Content + Chat area */}
        <div className="flex flex-1 overflow-hidden">
          {/* Page content */}
          <main className={`flex-1 overflow-y-auto p-6 ${chatExpanded ? 'hidden lg:block lg:w-1/2' : 'w-full'}`}>
            <Outlet />
          </main>

          {/* Expanded chat panel */}
          {chatExpanded && (
            <div className="flex-1 flex flex-col border-l border-slate-200 bg-white lg:max-w-[50%]">
              <div className="flex items-center justify-between px-4 py-3 border-b border-slate-200">
                <h3 className="font-semibold text-slate-800">Chat with Heywood</h3>
                <button
                  onClick={() => setChatExpanded(false)}
                  className="text-slate-400 hover:text-slate-600 text-xl leading-none"
                >
                  ×
                </button>
              </div>
              <ChatHistory messages={chat.messages} loading={chat.loading} />
            </div>
          )}
        </div>

        {/* Persistent chat bar */}
        <ChatBar
          onSend={chat.sendMessage}
          loading={chat.loading}
          onExpand={() => setChatExpanded(true)}
          expanded={chatExpanded}
          hasMessages={chat.messages.length > 0}
        />
      </div>
    </div>
  )
}
