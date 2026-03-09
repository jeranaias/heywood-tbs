// Dual-mode API layer
// - Normal mode: fetches from Go backend (/api/v1/*)
// - Static mode (GitHub Pages): uses bundled JSON data + mock AI

import type {
  Student, StudentStats, QualStats,
  Instructor, TrainingEvent, Qualification, QualRecord,
  AuthInfo, ChatMessage, Task, Message, Notification,
  AppSettings, SystemInfo, CalendarEvent, MailSummary,
  CounselingSession,
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

  updateStudent: async (id: string, updates: { notes?: string; atRisk?: boolean; riskFlags?: string[] }) => {
    return patch<Student>(`/students/${id}`, updates)
  },

  getStudentNotes: async (id: string) => {
    return get<{ notes: Array<{ id: string; studentId: string; authorRole: string; authorName: string; content: string; type: string; createdAt: string }> }>(`/students/${id}/notes`)
  },

  createStudentNote: async (studentId: string, content: string, type: string = 'note') => {
    return post<{ notes: Array<{ id: string; studentId: string; authorRole: string; authorName: string; content: string; type: string; createdAt: string }> }>(`/students/${studentId}/notes`, { content, type })
  },

  // Counseling
  getCounselings: async (params?: Record<string, string>) => {
    return get<{ sessions: CounselingSession[]; total: number }>('/counselings', params)
  },

  getCounseling: async (id: string) => {
    return get<CounselingSession>(`/counselings/${id}`)
  },

  createCounseling: async (session: Partial<CounselingSession>) => {
    return post<CounselingSession>('/counselings', session)
  },

  updateCounseling: async (id: string, session: Partial<CounselingSession>) => {
    return put<CounselingSession>(`/counselings/${id}`, session)
  },

  generateCounselingOutline: async (studentId: string, type?: string) => {
    return post<{ outline: string }>('/counselings/generate-outline', { studentId, type })
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

  createTrainingEvent: async (event: Partial<TrainingEvent>) => {
    return post<TrainingEvent>('/schedule', event)
  },

  updateTrainingEvent: async (id: string, event: Partial<TrainingEvent>) => {
    return put<TrainingEvent>(`/schedule/${id}`, event)
  },

  deleteTrainingEvent: async (id: string) => {
    const res = await fetch(`${API}/schedule/${id}`, {
      method: 'DELETE',
      credentials: 'include',
    })
    if (!res.ok) throw new Error(`API ${res.status}: ${res.statusText}`)
    return res.json()
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

  // Chat suggestions
  getSuggestedPrompts: async () => {
    if (STATIC) return { prompts: ['How am I doing?', 'My schedule today'] }
    return get<{ prompts: string[] }>('/chat/suggestions')
  },

  // Chat sessions
  getChatSessions: async () => {
    return get<{ sessions: Array<{ id: string; title: string; userRole: string; createdAt: string; updatedAt: string }> }>('/chat/sessions')
  },

  getChatSessionMessages: async (sessionId: string) => {
    return get<{ messages: Array<{ role: string; content: string }> }>(`/chat/sessions/${sessionId}/messages`)
  },

  deleteChatSession: async (sessionId: string) => {
    const res = await fetch(`${API}/chat/sessions/${sessionId}`, {
      method: 'DELETE',
      credentials: 'include',
    })
    if (!res.ok) throw new Error(`API ${res.status}: ${res.statusText}`)
    return res.json()
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

  createTask: async (task: { title: string; description?: string; assignedTo?: string; priority?: string; dueDate?: string; relatedId?: string }) => {
    return post<Task>('/tasks', task)
  },

  deleteTask: async (id: string) => {
    const res = await fetch(`${API}/tasks/${id}`, {
      method: 'DELETE',
      credentials: 'include',
    })
    if (!res.ok) throw new Error(`API ${res.status}: ${res.statusText}`)
    return res.json()
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

  createCalendarEvent: async (event: Partial<CalendarEvent>) => {
    return post<CalendarEvent>('/calendar/events', event)
  },

  getMailSummary: async () => {
    return get<{ messages: MailSummary[]; unreadCount: number }>('/mail/summary')
  },

  getMailUnreadCount: async () => {
    return get<{ count: number }>('/mail/unread-count')
  },

  sendMail: async (to: string[], subject: string, body: string) => {
    return post<{ status: string }>('/mail/send', { to, subject, body })
  },

  replyToMail: async (messageId: string, body: string) => {
    return post<{ status: string }>(`/mail/${messageId}/reply`, { body })
  },

  respondToEvent: async (eventId: string, response: 'accept' | 'decline' | 'tentativelyAccept') => {
    return post<{ status: string }>(`/calendar/events/${eventId}/respond`, { response })
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
