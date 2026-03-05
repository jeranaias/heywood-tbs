// Dual-mode API layer
// - Normal mode: fetches from Go backend (/api/v1/*)
// - Static mode (GitHub Pages): uses bundled JSON data + mock AI

import type {
  Student, StudentStats, QualStats,
  Instructor, TrainingEvent, Qualification, QualRecord,
  AuthInfo, ChatMessage,
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

async function post<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${API}${path}`, {
    method: 'POST',
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

  // Chat
  chat: async (message: string, history: ChatMessage[]) => {
    if (STATIC) {
      const sd = await getStaticModule()
      const mc = await getMockChatModule()
      const auth = sd.getAuth()
      // Simulate network delay for realism
      await new Promise(r => setTimeout(r, 300 + Math.random() * 700))
      return { response: mc.generateMockResponse(message, history, auth.role) }
    }
    return post<{ response: string }>('/chat', { message, history })
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
}
