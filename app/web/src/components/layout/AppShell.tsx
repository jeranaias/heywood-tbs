import { useState } from 'react'
import { Outlet } from 'react-router-dom'
import { Sidebar } from './Sidebar'
import { Header } from './Header'
import { ChatWidget } from '../chat/ChatWidget'
import { useAuth } from '../../hooks/useAuth'
import { ToastProvider } from '../common/Toast'
import { useNotifications } from '../../hooks/useNotifications'

function AppShellInner() {
  const { auth } = useAuth()
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const { sseNotifCount } = useNotifications()

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
        <Header onMenuClick={() => setSidebarOpen(true)} sseNotifCount={sseNotifCount} />
        <main className="flex-1 overflow-y-auto p-6">
          <Outlet />
        </main>
      </div>

      {/* Floating Heywood chat — available on every page except /chat */}
      <ChatWidget />
    </div>
  )
}

export function AppShell() {
  return (
    <ToastProvider>
      <AppShellInner />
    </ToastProvider>
  )
}
