// Dual-mode API layer
// - Normal mode: fetches from Go backend (/api/v1/*)
// - Static mode (GitHub Pages): uses bundled JSON data + mock AI

import type {
  Student, StudentStats, QualStats,
  Instructor, TrainingEvent, Qualification, QualRecord,
  AuthInfo, ChatMessage, Task, Message, Notification,
  AppSettings, SystemInfo, CalendarEvent, MailSummary,
} from './types'

const STATIC = import.meta.env.MODE === 'static'

// ---- Backend API (normal mode) ----

const API = '/api/v1'

async function get<T>(path: string, params?: Record<string, string>): Promise<T> {
  const url = params
    ? `${API}${path}?${new URLSearchParams(params)}`
    : `${API}${path}`
  const res = await fetch(url, { credentials: 'include' })
  if (!res.ok) throw new Error(`API ${res.status}: ${res.statusText}`)
  return res.json()
}

async function post<T>(path: string, body?: unknown): Promise<T> {
  const res = await fetch(`${API}${path}`, {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: body ? JSON.stringify(body) : undefined,
  })
  if (!res.ok) throw new Error(`API ${res.status}: ${res.statusText}`)
  return res.json()
}

async function patch<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${API}${path}`, {
    method: 'PATCH',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) throw new Error(`API ${res.status}: ${res.statusText}`)
  return res.json()
}

async function put<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${API}${path}`, {
    method: 'PUT',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) throw new Error(`API ${res.status}: ${res.statusText}`)
  return res.json()
}

// ---- Static data imports (tree-shaken when not in static mode) ----

async function getStaticModule() {
  return import('./staticData')
}

async function getMockChatModule() {
  return import('./mockChat')
}

// ---- Unified API ----

export const api = {
  // Students
  getStudents: async (params?: Record<string, string>) => {
    if (STATIC) {
      const sd = await getStaticModule()
      return sd.listStudents(params)
    }
    return get<{ students: Student[]; total: number; filtered: number }>('/students', params)
  },

  getStudent: async (id: string) => {
    if (STATIC) {
      const sd = await getStaticModule()
      const s = sd.getStudent(id)
      if (!s) throw new Error('Student not found')
      return s
    }
    return get<Student>(`/students/${id}`)
  },

  getStudentStats: async (params?: Record<string, string>) => {
    if (STATIC) {
      const sd = await getStaticModule()
      return sd.getStudentStats(params)
    }
    return get<StudentStats>('/students/stats', params)
  },

  getAtRiskStudents: async () => {
    if (STATIC) {
      const sd = await getStaticModule()
      return sd.getAtRiskStudents()
    }
    return get<{ students: Student[]; total: number }>('/students/at-risk')
  },

  // Instructors
  getInstructors: async (params?: Record<string, string>) => {
    if (STATIC) {
      const sd = await getStaticModule()
      return sd.listInstructors(params)
    }
    return get<{ instructors: Instructor[]; total: number }>('/instructors', params)
  },

  getInstructor: async (id: string) => {
    if (STATIC) {
      const sd = await getStaticModule()
      const i = sd.getInstructor(id)
      if (!i) throw new Error('Instructor not found')
      return i
    }
    return get<Instructor>(`/instructors/${id}`)
  },

  // Qualifications
  getQualifications: async () => {
    if (STATIC) {
      const sd = await getStaticModule()
      return sd.listQualifications()
    }
    return get<Qualification[]>('/qualifications')
  },

  getQualRecords: async () => {
    if (STATIC) {
      const sd = await getStaticModule()
      return sd.listQualRecords()
    }
    return get<QualRecord[]>('/qual-records')
  },

  getQualStats: async () => {
    if (STATIC) {
      const sd = await getStaticModule()
      return sd.getQualStats()
    }
    return get<QualStats>('/qual-records/stats')
  },

  // Schedule
  getSchedule: async (params?: Record<string, string>) => {
    if (STATIC) {
      const sd = await getStaticModule()
      return sd.listSchedule(params)
    }
    return get<{ events: TrainingEvent[]; total: number }>('/schedule', params)
  },

  // Chat (non-streaming)
  chat: async (message: string, history: ChatMessage[]) => {
    if (STATIC) {
      const sd = await getStaticModule()
      const mc = await getMockChatModule()
      const auth = sd.getAuth()
      await new Promise(r => setTimeout(r, 300 + Math.random() * 700))
      return { response: mc.generateMockResponse(message, history, auth.role) }
    }
    return post<{ response: string }>('/chat', { message, history })
  },

  // Chat (SSE streaming)
  chatStream: async (
    message: string,
    history: ChatMessage[],
    onChunk: (chunk: string) => void,
    onDone: () => void,
    signal?: AbortSignal,
  ) => {
    if (STATIC) {
      // Simulate streaming for static mode
      const sd = await getStaticModule()
      const mc = await getMockChatModule()
      const auth = sd.getAuth()
      const full = mc.generateMockResponse(message, history, auth.role)
      // Simulate token-by-token output
      const words = full.split(/(\s+)/)
      for (let i = 0; i < words.length; i++) {
        if (signal?.aborted) return
        onChunk(words[i])
        await new Promise(r => setTimeout(r, 15 + Math.random() * 25))
      }
      onDone()
      return
    }

    const res = await fetch(`${API}/chat`, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ message, history, stream: true }),
      signal,
    })
    if (!res.ok) throw new Error(`API ${res.status}: ${res.statusText}`)
    if (!res.body) throw new Error('No response body')

    const reader = res.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''

      for (const line of lines) {
        const trimmed = line.trim()
        if (!trimmed || !trimmed.startsWith('data: ')) continue
        const data = trimmed.slice(6)
        if (data === '[DONE]') {
          onDone()
          return
        }
        try {
          const parsed = JSON.parse(data) as { content?: string; error?: string }
          if (parsed.content) onChunk(parsed.content)
        } catch {
          // skip malformed chunks
        }
      }
    }
    onDone()
  },

  // Auth
  getAuthMe: async () => {
    if (STATIC) {
      const sd = await getStaticModule()
      return sd.getAuth()
    }
    return get<AuthInfo>('/auth/me')
  },

  switchRole: async (role: string, company: string, studentId?: string) => {
    if (STATIC) {
      const sd = await getStaticModule()
      return sd.switchAuth(role, company, studentId)
    }
    return post<AuthInfo>('/auth/switch', { role, company, studentId: studentId || '' })
  },

  // Tasks
  getTasks: async (params?: Record<string, string>) => {
    return get<Task[]>('/tasks', params)
  },

  getTask: async (id: string) => {
    return get<Task>(`/tasks/${id}`)
  },

  updateTask: async (id: string, updates: Partial<Pick<Task, 'status' | 'priority' | 'assignedTo'>>) => {
    return patch<Task>(`/tasks/${id}`, updates)
  },

  // Messages
  getMessages: async (params?: Record<string, string>) => {
    return get<Message[]>('/messages', params)
  },

  markMessageRead: async (id: string) => {
    return post<{ status: string }>(`/messages/${id}/read`)
  },

  // Notifications
  getNotifications: async (params?: Record<string, string>) => {
    return get<Notification[]>('/notifications', params)
  },

  getNotificationCount: async () => {
    return get<{ count: number }>('/notifications/count')
  },

  markNotificationRead: async (id: string) => {
    return post<{ status: string }>(`/notifications/${id}/read`)
  },

  // Calendar
  getCalendarEvents: async (params?: Record<string, string>) => {
    return get<{ events: CalendarEvent[]; start: string; end: string }>('/calendar/events', params)
  },

  getCalendarToday: async () => {
    return get<{ events: CalendarEvent[]; date: string }>('/calendar/today')
  },

  getMailSummary: async () => {
    return get<{ messages: MailSummary[]; unreadCount: number }>('/mail/summary')
  },

  getMailUnreadCount: async () => {
    return get<{ count: number }>('/mail/unread-count')
  },

  // Settings (XO/Staff only)
  getSettings: async () => {
    return get<AppSettings>('/settings')
  },

  updateSettings: async (settings: AppSettings) => {
    return put<{ status: string; note: string }>('/settings', settings)
  },

  testConnection: async (params: { type: string; connectionString?: string; tenantId?: string; clientId?: string; clientSecret?: string; siteUrl?: string }) => {
    return post<{ status: string; message: string }>('/settings/test-connection', params)
  },

  uploadFile: async (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    const res = await fetch(`${API}/settings/upload`, {
      method: 'POST',
      credentials: 'include',
      body: formData,
    })
    if (!res.ok) throw new Error(`Upload failed: ${res.statusText}`)
    return res.json() as Promise<{ status: string; filename: string; path: string; size: number; type: string }>
  },

  getSystemInfo: async () => {
    return get<SystemInfo>('/settings/system-info')
  },
}
