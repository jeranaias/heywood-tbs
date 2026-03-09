import { createContext, useContext, type ReactNode } from 'react'
import { useChat } from './useChat'
import type { ChatMessage } from '../lib/types'

interface ChatContextValue {
  messages: ChatMessage[]
  loading: boolean
  sendMessage: (content: string) => Promise<void>
  clearMessages: () => void
  loadMessages: (msgs: ChatMessage[]) => void
}

const ChatContext = createContext<ChatContextValue | null>(null)

export function ChatProvider({ children }: { children: ReactNode }) {
  const chat = useChat()
  return <ChatContext.Provider value={chat}>{children}</ChatContext.Provider>
}

export function useChatContext(): ChatContextValue {
  const ctx = useContext(ChatContext)
  if (!ctx) throw new Error('useChatContext must be used within ChatProvider')
  return ctx
}
