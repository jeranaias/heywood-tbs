import { useState, useEffect, useCallback, createContext, useContext } from 'react'
import { X, AlertTriangle, CheckCircle, Info } from 'lucide-react'

interface ToastItem {
  id: number
  message: string
  severity: 'info' | 'success' | 'warning' | 'error'
  dismissAt: number
}

interface ToastContextValue {
  addToast: (message: string, severity?: ToastItem['severity']) => void
}

const ToastContext = createContext<ToastContextValue>({ addToast: () => {} })

export function useToast() {
  return useContext(ToastContext)
}

let nextId = 1

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<ToastItem[]>([])

  const addToast = useCallback((message: string, severity: ToastItem['severity'] = 'info') => {
    const id = nextId++
    setToasts(prev => [...prev, { id, message, severity, dismissAt: Date.now() + 5000 }])
  }, [])

  const removeToast = useCallback((id: number) => {
    setToasts(prev => prev.filter(t => t.id !== id))
  }, [])

  // Auto-dismiss expired toasts
  useEffect(() => {
    if (toasts.length === 0) return
    const timer = setInterval(() => {
      const now = Date.now()
      setToasts(prev => prev.filter(t => t.dismissAt > now))
    }, 500)
    return () => clearInterval(timer)
  }, [toasts.length])

  const severityConfig = {
    info: { bg: 'bg-blue-50 border-blue-200', text: 'text-blue-800', Icon: Info },
    success: { bg: 'bg-green-50 border-green-200', text: 'text-green-800', Icon: CheckCircle },
    warning: { bg: 'bg-amber-50 border-amber-200', text: 'text-amber-800', Icon: AlertTriangle },
    error: { bg: 'bg-red-50 border-red-200', text: 'text-red-800', Icon: AlertTriangle },
  }

  return (
    <ToastContext.Provider value={{ addToast }}>
      {children}
      {/* Toast container — fixed top-right */}
      <div className="fixed top-4 right-4 z-[60] flex flex-col gap-2 max-w-sm">
        {toasts.map(toast => {
          const config = severityConfig[toast.severity]
          return (
            <div
              key={toast.id}
              className={`flex items-start gap-2 px-4 py-3 rounded-lg border shadow-md animate-slide-in ${config.bg}`}
            >
              <config.Icon className={`w-4 h-4 mt-0.5 flex-shrink-0 ${config.text}`} />
              <p className={`text-sm flex-1 ${config.text}`}>{toast.message}</p>
              <button
                onClick={() => removeToast(toast.id)}
                className={`flex-shrink-0 p-0.5 rounded hover:bg-black/5 ${config.text}`}
              >
                <X className="w-3.5 h-3.5" />
              </button>
            </div>
          )
        })}
      </div>
    </ToastContext.Provider>
  )
}
