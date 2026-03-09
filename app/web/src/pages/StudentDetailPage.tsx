import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { ArrowLeft, AlertTriangle } from 'lucide-react'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ReferenceLine, ResponsiveContainer, Cell } from 'recharts'
import { api } from '../lib/api'
import type { Student } from '../lib/types'
import { formatScore, scoreBadge, scoreLabel, trendIcon, trendColor } from '../lib/utils'
import { useChat } from '../hooks/useChat'
import { useAuth } from '../hooks/useAuth'
import { RiskFlagToggle } from '../components/students/RiskFlagToggle'
import { StudentNotesPanel } from '../components/students/StudentNotesPanel'

export function StudentDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [student, setStudent] = useState<Student | null>(null)
  const [loading, setLoading] = useState(true)
  const chat = useChat()
  const { auth } = useAuth()
  const canEdit = auth.role === 'xo' || auth.role === 'staff' || auth.role === 'spc'

  useEffect(() => {
    if (!id) return
    api.getStudent(id)
      .then(setStudent)
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [id])

  if (loading) {
    return <div className="animate-pulse h-96 bg-slate-200 rounded-lg" />
  }

  if (!student) {
    return <div className="text-center py-12 text-slate-500">Student not found</div>
  }

  const examData = [
    { name: 'Exam 1', score: student.exam1 },
    { name: 'Exam 2', score: student.exam2 },
    { name: 'Exam 3', score: student.exam3 },
    { name: 'Exam 4', score: student.exam4 },
    { name: 'Quiz Avg', score: student.quizAvg },
  ].filter(d => d.score > 0)

  const pillarData = [
    { name: 'Academic\n(32%)', score: student.academicComposite, weight: 32 },
    { name: 'Mil Skills\n(32%)', score: student.milSkillsComposite, weight: 32 },
    { name: 'Leadership\n(36%)', score: student.leadershipComposite, weight: 36 },
  ]

  const milSkills = [
    { label: 'PFT', value: student.pftScore > 0 ? `${student.pftScore}` : '—', pass: student.pftScore >= 235 },
    { label: 'CFT', value: student.cftScore > 0 ? `${student.cftScore}` : '—', pass: student.cftScore >= 235 },
    { label: 'Rifle', value: student.rifleQual, pass: student.rifleQual !== 'Unqualified' && student.rifleQual !== '' },
    { label: 'Pistol', value: student.pistolQual, pass: student.pistolQual !== 'Unqualified' && student.pistolQual !== '' },
    { label: 'Land Nav (Day)', value: student.landNavDay, pass: student.landNavDay === 'Pass' },
    { label: 'Land Nav (Night)', value: student.landNavNight, pass: student.landNavNight === 'Pass' },
    { label: 'Land Nav Written', value: student.landNavWritten > 0 ? formatScore(student.landNavWritten) : '—', pass: student.landNavWritten >= 70 },
    { label: 'Obstacle Course', value: student.obstacleCourse, pass: student.obstacleCourse === 'Pass' },
    { label: 'Endurance Course', value: student.enduranceCourse, pass: student.enduranceCourse === 'Pass' },
  ]

  function barColor(score: number): string {
    if (score >= 85) return '#22c55e'
    if (score >= 75) return '#eab308'
    return '#ef4444'
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <Link to="/students" className="inline-flex items-center gap-1 text-sm text-slate-500 hover:text-slate-700 mb-2">
          <ArrowLeft className="w-4 h-4" /> Back to Roster
        </Link>
        <div className="flex items-start justify-between">
          <div>
            <h2 className="text-xl font-bold text-slate-900">
              {student.rank} {student.lastName}, {student.firstName}
            </h2>
            <p className="text-sm text-slate-500 mt-1">
              {student.company} Company | {student.platoon} Platoon | SPC: {student.spc}
            </p>
            <div className="flex items-center gap-3 mt-2">
              <span className="text-sm text-slate-600">{student.phase}</span>
              <span className={`inline-flex items-center gap-1 text-sm font-medium ${trendColor(student.trend)}`}>
                {trendIcon(student.trend)} {student.trend}
              </span>
              {student.atRisk && (
                <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-red-100 text-red-800 text-xs font-medium">
                  <AlertTriangle className="w-3 h-3" /> At Risk
                </span>
              )}
            </div>
          </div>
          <div className="text-right">
            <div className="text-3xl font-bold text-slate-900">{formatScore(student.overallComposite)}</div>
            <div className="text-sm text-slate-500">Overall Composite</div>
            {student.companyRank > 0 && (
              <div className="text-xs text-slate-400 mt-1">Rank #{student.companyRank} / {student.classStandingThird}</div>
            )}
          </div>
        </div>
      </div>

      {/* Three Pillar Gauges */}
      <div className="grid grid-cols-3 gap-4">
        {pillarData.map(p => (
          <div key={p.name} className="bg-white rounded-lg border border-slate-200 p-4 text-center">
            <div className={`text-3xl font-bold ${p.score >= 85 ? 'text-green-600' : p.score >= 75 ? 'text-yellow-600' : 'text-red-600'}`}>
              {formatScore(p.score)}
            </div>
            <div className="text-sm font-medium text-slate-700 mt-1">{p.name.replace('\n', ' ')}</div>
            <div className="text-xs text-slate-500">{scoreLabel(p.score)}</div>
            <div className="mt-2 h-2 bg-slate-100 rounded-full overflow-hidden">
              <div
                className={`h-full rounded-full ${p.score >= 85 ? 'bg-green-500' : p.score >= 75 ? 'bg-yellow-500' : 'bg-red-500'}`}
                style={{ width: `${Math.min(100, p.score)}%` }}
              />
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Exam Scores Bar Chart */}
        {examData.length > 0 && (
          <div className="bg-white rounded-lg border border-slate-200 p-4">
            <h3 className="text-sm font-semibold text-slate-700 mb-4">Academic Exam Scores</h3>
            <ResponsiveContainer width="100%" height={200}>
              <BarChart data={examData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
                <XAxis dataKey="name" tick={{ fontSize: 12 }} />
                <YAxis domain={[0, 100]} tick={{ fontSize: 12 }} />
                <Tooltip />
                <ReferenceLine y={75} stroke="#ef4444" strokeDasharray="5 5" label={{ value: 'Min 75', fill: '#ef4444', fontSize: 10 }} />
                <Bar dataKey="score" radius={[4, 4, 0, 0]}>
                  {examData.map((entry, idx) => (
                    <Cell key={idx} fill={barColor(entry.score)} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>
        )}

        {/* Military Skills Table */}
        <div className="bg-white rounded-lg border border-slate-200 p-4">
          <h3 className="text-sm font-semibold text-slate-700 mb-4">Military Skills</h3>
          <div className="space-y-2">
            {milSkills.map(skill => (
              <div key={skill.label} className="flex items-center justify-between py-1.5 border-b border-slate-100 last:border-0">
                <span className="text-sm text-slate-600">{skill.label}</span>
                <span className={`text-sm font-medium ${
                  skill.value === '—' || skill.value === 'Not Yet Tested'
                    ? 'text-slate-400'
                    : skill.pass ? 'text-green-600' : 'text-red-600'
                }`}>
                  {skill.value === 'Pass' ? '✓ Pass' :
                   skill.value === 'Fail' ? '✗ Fail' :
                   skill.value === 'Not Yet Tested' ? '○ Pending' :
                   skill.value}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Risk Flags — editable for staff/SPC/XO, read-only for students */}
      {canEdit ? (
        <RiskFlagToggle student={student} onUpdate={setStudent} />
      ) : student.riskFlags.length > 0 && (
        <div className="bg-red-50 rounded-lg border border-red-200 p-4">
          <h3 className="text-sm font-semibold text-red-800 mb-2">Risk Flags</h3>
          <div className="flex flex-wrap gap-2">
            {student.riskFlags.map(flag => (
              <span key={flag} className="px-2 py-1 bg-red-100 text-red-800 text-xs rounded-full font-medium">
                {flag}
              </span>
            ))}
          </div>
        </div>
      )}

      {/* Staff Notes — visible to staff/SPC/XO only */}
      {canEdit && id && (
        <StudentNotesPanel studentId={id} />
      )}

      {/* Legacy notes field */}
      {student.notes && (
        <div className="bg-yellow-50 rounded-lg border border-yellow-200 p-4">
          <h3 className="text-sm font-semibold text-yellow-800 mb-1">Notes</h3>
          <p className="text-sm text-yellow-700">{student.notes}</p>
        </div>
      )}
    </div>
  )
}
