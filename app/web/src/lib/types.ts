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
}

export interface AuthInfo {
  role: string
  company: string
  studentId?: string
  name: string
}

export type Role = 'staff' | 'spc' | 'student'
