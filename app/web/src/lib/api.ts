import type { Student, StudentStats, QualStats, Instructor, TrainingEvent, QualRecord, AuthInfo, ChatMessage } from './types'

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

export const api = {
  // Students
  getStudents: (params?: Record<string, string>) =>
    get<{ students: Student[]; total: number; filtered: number }>('/students', params),
  getStudent: (id: string) =>
    get<Student>(`/students/${id}`),
  getStudentStats: (params?: Record<string, string>) =>
    get<StudentStats>('/students/stats', params),
  getAtRiskStudents: () =>
    get<{ students: Student[]; total: number }>('/students/at-risk'),

  // Instructors
  getInstructors: (params?: Record<string, string>) =>
    get<{ instructors: Instructor[]; total: number }>('/instructors', params),
  getInstructor: (id: string) =>
    get<Instructor>(`/instructors/${id}`),

  // Qualifications
  getQualifications: () =>
    get<Qualification[]>('/qualifications'),
  getQualRecords: () =>
    get<QualRecord[]>('/qual-records'),
  getQualStats: () =>
    get<QualStats>('/qual-records/stats'),

  // Schedule
  getSchedule: (params?: Record<string, string>) =>
    get<{ events: TrainingEvent[]; total: number }>('/schedule', params),

  // Chat
  chat: (message: string, history: ChatMessage[]) =>
    post<{ response: string }>('/chat', { message, history }),

  // Auth
  getAuthMe: () =>
    get<AuthInfo>('/auth/me'),
  switchRole: (role: string, company: string, studentId?: string) =>
    post<AuthInfo>('/auth/switch', { role, company, studentId: studentId || '' }),
}
