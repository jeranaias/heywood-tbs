export interface Student {
  id: string
  edipi: string
  rank: string
  lastName: string
  firstName: string
  company: string
  platoon: string
  spc: string
  classNumber: string
  classStartDate: string
  phase: string
  exam1: number
  exam2: number
  exam3: number
  exam4: number
  quizAvg: number
  academicComposite: number
  pftScore: number
  cftScore: number
  rifleQual: string
  pistolQual: string
  landNavDay: string
  landNavNight: string
  landNavWritten: number
  obstacleCourse: string
  enduranceCourse: string
  milSkillsComposite: number
  leadershipWeek12: number
  leadershipWeek22: number
  peerEvalWeek12: number
  peerEvalWeek22: number
  leadershipComposite: number
  overallComposite: number
  classStandingThird: string
  companyRank: number
  trend: string
  atRisk: boolean
  riskFlags: string[]
  status: string
  notes: string
}

export interface Instructor {
  id: string
  edipi: string
  lastName: string
  firstName: string
  rank: string
  role: string
  company: string
  platoon: string
  classNumber: string
  dateAssigned: string
  prd: string
  studentsAssigned: number
  eventsThisWeek: number
  eventsThisMonth: number
  counselingsOverdue: number
  status: string
  phone: string
  email: string
  notes: string
}

export interface TrainingEvent {
  id: string
  title: string
  code: string
  phase: string
  category: string
  gradePillar: string
  isGraded: boolean
  startDate: string
  endDate: string
  startTime: string
  endTime: string
  durationHours: number
  location: string
  company: string
  classNumber: string
  leadInstructor: string
  supportInstructors: string
  instructorCountRequired: number
  prerequisiteEvents: string
  specialEquipment: string
  status: string
  weatherContingency: string
  notes: string
}

export interface Qualification {
  id: string
  code: string
  name: string
  category: string
  issuingAuthority: string
  validityMonths: number
  renewalProcess: string
  requiredForEvents: string
  minimumPerEvent: number
  orderReference: string
  status: string
  notes: string
}

export interface QualRecord {
  id: string
  instructorEdipi: string
  instructorName: string
  qualCode: string
  qualName: string
  dateEarned: string
  expirationDate: string
  daysUntilExpiration: number
  expirationStatus: string
  certificateNumber: string
  issuedBy: string
  renewalStatus: string
  renewalDate: string
  notes: string
}

export interface StudentStats {
  activeStudents: number
  avgComposite: number
  atRiskCount: number
  atRiskPercent: number
  byPhase: Record<string, number>
  byStandingThird: Record<string, number>
}

export interface QualStats {
  totalRecords: number
  expiredCount: number
  expiring30: number
  expiring60: number
  expiring90: number
  currentCount: number
  coverageGaps: CoverageGap[]
}

export interface CoverageGap {
  qualCode: string
  qualName: string
  qualifiedCount: number
  requiredCount: number
  gap: number
}

export interface ChatMessage {
  role: 'user' | 'assistant'
  content: string
  streaming?: boolean
}

export interface AuthInfo {
  role: string
  company: string
  studentId?: string
  name: string
}

export type Role = 'xo' | 'staff' | 'spc' | 'student'

export interface Task {
  id: string
  title: string
  description: string
  assignedTo: string
  createdBy: string
  priority: 'high' | 'medium' | 'low'
  status: 'pending' | 'in_progress' | 'completed'
  dueDate: string
  relatedId: string
  createdAt: string
  updatedAt: string
}

export interface Message {
  id: string
  from: string
  to: string
  subject: string
  body: string
  read: boolean
  relatedId: string
  createdAt: string
}

export interface Notification {
  id: string
  userRole: string
  type: 'task' | 'message' | 'alert'
  title: string
  body: string
  read: boolean
  actionUrl: string
  createdAt: string
}

export interface AppSettings {
  dataSource: {
    type: string
    jsonDir: string
    excelPath: string
    sharepoint: {
      tenantId: string
      clientId: string
      clientSecret: string
      siteUrl: string
      cloud: string
    }
    database: {
      type: string
      connectionString: string
    }
  }
  ai: {
    provider: string
    model: string
    searxngUrl: string
  }
  outlook: {
    enabled: boolean
    tenantId: string
    clientId: string
    clientSecret: string
    cloud: string
    masterCalendarId: string
    syncIntervalMinutes: number
  }
  auth: {
    mode: string
  }
}

export interface CalendarEvent {
  id: string
  title: string
  start: string
  end: string
  location: string
  description: string
  source: string
  calendarId?: string
  isAllDay: boolean
  organizer?: string
  category?: string
  company?: string
}

export interface MailSummary {
  id: string
  subject: string
  from: string
  preview: string
  received: string
  isRead: boolean
  hasAttachments: boolean
}

export interface SystemInfo {
  version: string
  authMode: string
  dataSource: string
  studentCount: number
  ai: {
    status: string
    keyHint: string
    model: string
  }
  searxng: {
    url: string
    status: string
  }
  outlook: {
    enabled: boolean
  }
}
