import { useEffect, useRef, useCallback, useState } from 'react'

export interface SSEEvent {
  type: string
  data: Record<string, unknown>
}

type SSEHandler = (event: SSEEvent) => void

/**
 * Opens an EventSource to the SSE stream endpoint and dispatches events to handlers.
 * Automatically reconnects on disconnect with exponential backoff.
 */
export function useSSE(onEvent: SSEHandler) {
  const [connected, setConnected] = useState(false)
  const eventSourceRef = useRef<EventSource | null>(null)
  const handlerRef = useRef(onEvent)
  handlerRef.current = onEvent

  const connect = useCallback(() => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close()
    }

    const es = new EventSource('/api/v1/events/stream')
    eventSourceRef.current = es

    es.onopen = () => {
      setConnected(true)
    }

    es.onmessage = (event) => {
      try {
        const parsed = JSON.parse(event.data) as SSEEvent
        handlerRef.current(parsed)
      } catch {
        // Skip malformed events
      }
    }

    es.onerror = () => {
      setConnected(false)
      es.close()
      // Reconnect after 5 seconds
      setTimeout(connect, 5000)
    }
  }, [])

  useEffect(() => {
    connect()
    return () => {
      if (eventSourceRef.current) {
        eventSourceRef.current.close()
        eventSourceRef.current = null
      }
    }
  }, [connect])

  return { connected }
}
