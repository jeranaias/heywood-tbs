import { useCallback, useState } from 'react'
import { useSSE } from './useSSE'
import { useToast } from '../components/common/Toast'
import type { SSEEvent } from './useSSE'

/**
 * Combines SSE events with toast notifications and badge count tracking.
 * Returns the live notification count from SSE events.
 */
export function useNotifications() {
  const { addToast } = useToast()
  const [sseNotifCount, setSseNotifCount] = useState(0)

  const handleEvent = useCallback((event: SSEEvent) => {
    switch (event.type) {
      case 'at-risk-alert': {
        const msg = (event.data.message as string) || 'At-risk status changed'
        addToast(msg, 'warning')
        setSseNotifCount(prev => prev + 1)
        break
      }
      case 'task': {
        const action = event.data.action as string
        const task = event.data.task as { title?: string } | undefined
        const title = task?.title || 'Task'
        if (action === 'created') {
          addToast(`New task: ${title}`, 'info')
        } else if (action === 'updated') {
          addToast(`Task updated: ${title}`, 'info')
        }
        setSseNotifCount(prev => prev + 1)
        break
      }
      case 'notification': {
        const msg = (event.data.message as string) || 'New notification'
        addToast(msg, 'info')
        setSseNotifCount(prev => prev + 1)
        break
      }
      // 'connected' events are ignored
    }
  }, [addToast])

  const { connected } = useSSE(handleEvent)

  return { connected, sseNotifCount }
}
