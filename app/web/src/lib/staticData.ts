// Static data layer for GitHub Pages deployment
// Imports JSON data at build time via Vite — no backend needed

import type {
  Student, StudentStats, QualStats, CoverageGap,
  Instructor, TrainingEvent, Qualification, QualRecord,
  AuthInfo,
} from './types'

// @ts-ignore — resolved by Vite @data alias
import studentsRaw from '@data/students.json'
// @ts-ignore
import instructorsRaw from '@data/instructors.json'
// @ts-ignore
import scheduleRaw from '@data/schedule.json'
// @ts-ignore
import qualificationsRaw from '@data/qualifications.json'
// @ts-ignore
import qualRecordsRaw from '@data/qual-records.json'

const students = studentsRaw as Student[]
const instructors = instructorsRaw as Instructor[]
const schedule = scheduleRaw as TrainingEvent[]
const qualifications = qualificationsRaw as Qualification[]
const qualRecords = qualRecordsRaw as QualRecord[]

// --------------- Auth (localStorage-backed) ---------------

let currentAuth: AuthInfo = {
  role: 'staff',
  company: '',
  name: 'TBS Staff',
}

const saved = localStorage.getItem('heywood-auth')
if (saved) {
  try { currentAuth = JSON.parse(saved) } catch { /* ignore */ }
}

export function getAuth(): AuthInfo {
  return { ...currentAuth }
}

export function switchAuth(role: string, company: string, studentId?: string): AuthInfo {
  const names: Record<string, string> = {
    xo: 'Executive Officer',
    staff: 'TBS Staff',
    spc: `SPC (${company} Company)`,
    student: `Student ${studentId || 'STU-001'}`,
  }
  currentAuth = {
    role,
    company,
    studentId: studentId || undefined,
    name: names[role] || 'Unknown',
  }
  localStorage.setItem('heywood-auth', JSON.stringify(currentAuth))
  return { ...currentAuth }
}

// --------------- Students ---------------

export function listStudents(params?: Record<string, string>): { students: Student[]; total: number; filtered: number } {
  let result = [...students]
  const auth = getAuth()

  // Role-based filtering
  if (auth.role === 'spc' && auth.company) {
    result = result.filter(s => s.company === auth.company)
  } else if (auth.role === 'student' && auth.studentId) {
    result = result.filter(s => s.id === auth.studentId)
  }

  // Query params
  if (params?.company) result = result.filter(s => s.company === params.company)
  if (params?.phase) result = result.filter(s => s.phase === params.phase)
  if (params?.atRisk === 'true') result = result.filter(s => s.atRisk)

  if (params?.search) {
    const q = params.search.toLowerCase()
    result = result.filter(s =>
      s.firstName.toLowerCase().includes(q) ||
      s.lastName.toLowerCase().includes(q) ||
      s.id.toLowerCase().includes(q)
    )
  }

  return { students: result, total: students.length, filtered: result.length }
}

export function getStudent(id: string): Student | undefined {
  return students.find(s => s.id === id)
}

export function getStudentStats(params?: Record<string, string>): StudentStats {
  let filtered = [...students]
  const auth = getAuth()

  if (auth.role === 'spc' && auth.company) {
    filtered = filtered.filter(s => s.company === auth.company)
  }
  if (params?.company) filtered = filtered.filter(s => s.company === params.company)

  const active = filtered.filter(s => s.status === 'Active')
  const atRisk = active.filter(s => s.atRisk)
  const avgComposite = active.length > 0
    ? active.reduce((sum, s) => sum + s.overallComposite, 0) / active.length
    : 0

  const byPhase: Record<string, number> = {}
  const byStandingThird: Record<string, number> = {}

  for (const s of active) {
    byPhase[s.phase] = (byPhase[s.phase] || 0) + 1
    byStandingThird[s.classStandingThird] = (byStandingThird[s.classStandingThird] || 0) + 1
  }

  return {
    activeStudents: active.length,
    avgComposite,
    atRiskCount: atRisk.length,
    atRiskPercent: active.length > 0 ? (atRisk.length / active.length) * 100 : 0,
    byPhase,
    byStandingThird,
  }
}

export function getAtRiskStudents(): { students: Student[]; total: number } {
  let result = students.filter(s => s.atRisk)
  const auth = getAuth()

  if (auth.role === 'spc' && auth.company) {
    result = result.filter(s => s.company === auth.company)
  }

  return { students: result, total: result.length }
}

// --------------- Instructors ---------------

export function listInstructors(params?: Record<string, string>): { instructors: Instructor[]; total: number } {
  let result = [...instructors]
  if (params?.company) result = result.filter(i => i.company === params.company)
  return { instructors: result, total: result.length }
}

export function getInstructor(id: string): Instructor | undefined {
  return instructors.find(i => i.id === id)
}

// --------------- Qualifications ---------------

export function listQualifications(): Qualification[] {
  return [...qualifications]
}

export function listQualRecords(): QualRecord[] {
  return [...qualRecords]
}

export function getQualStats(): QualStats {
  const expired = qualRecords.filter(r => r.expirationStatus === 'Expired').length
  const expiring30 = qualRecords.filter(r => r.daysUntilExpiration > 0 && r.daysUntilExpiration <= 30).length
  const expiring60 = qualRecords.filter(r => r.daysUntilExpiration > 30 && r.daysUntilExpiration <= 60).length
  const expiring90 = qualRecords.filter(r => r.daysUntilExpiration > 60 && r.daysUntilExpiration <= 90).length
  const current = qualRecords.filter(r => r.expirationStatus === 'Current').length

  const coverageGaps: CoverageGap[] = []
  for (const q of qualifications) {
    if (q.minimumPerEvent <= 0) continue
    const qualified = qualRecords.filter(
      r => r.qualCode === q.code && r.expirationStatus === 'Current'
    ).length
    if (qualified < q.minimumPerEvent) {
      coverageGaps.push({
        qualCode: q.code,
        qualName: q.name,
        qualifiedCount: qualified,
        requiredCount: q.minimumPerEvent,
        gap: q.minimumPerEvent - qualified,
      })
    }
  }

  return {
    totalRecords: qualRecords.length,
    expiredCount: expired,
    expiring30,
    expiring60,
    expiring90,
    currentCount: current,
    coverageGaps,
  }
}

// --------------- Schedule ---------------

export function listSchedule(params?: Record<string, string>): { events: TrainingEvent[]; total: number } {
  let result = [...schedule]
  if (params?.phase) result = result.filter(e => e.phase === params.phase)
  return { events: result, total: result.length }
}

// --------------- Export raw data for chat context ---------------

export function getAllStudents(): Student[] { return students }
export function getAllInstructors(): Instructor[] { return instructors }
