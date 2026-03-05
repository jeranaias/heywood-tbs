import { useState, useEffect } from 'react'
import { api } from '../lib/api'
import type { Student } from '../lib/types'
import { useAuth } from '../hooks/useAuth'
import { formatScore, scoreLabel } from '../lib/utils'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ReferenceLine, ResponsiveContainer, Cell } from 'recharts'

export function MyRecord() {
  const { auth } = useAuth()
  const [student, setStudent] = useState<Student | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (auth.studentId) {
      api.getStudent(auth.studentId)
        .then(setStudent)
        .catch(console.error)
        .finally(() => setLoading(false))
    } else {
      setLoading(false)
    }
  }, [auth.studentId])

  if (loading) {
    return <div className="animate-pulse h-96 bg-slate-200 rounded-lg" />
  }

  if (!student) {
    return (
      <div className="text-center py-12">
        <p className="text-slate-500">No student record found. Make sure you're logged in as a student.</p>
      </div>
    )
  }

  const examData = [
    { name: 'Exam 1', score: student.exam1 },
    { name: 'Exam 2', score: student.exam2 },
    { name: 'Exam 3', score: student.exam3 },
    { name: 'Exam 4', score: student.exam4 },
    { name: 'Quiz Avg', score: student.quizAvg },
  ].filter(d => d.score > 0)

  function barColor(score: number): string {
    if (score >= 85) return '#22c55e'
    if (score >= 75) return '#eab308'
    return '#ef4444'
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-bold text-slate-900">My Record</h2>
        <p className="text-sm text-slate-500">{student.rank} {student.lastName}, {student.firstName} — {student.phase}</p>
      </div>

      {/* Overall Composite */}
      <div className="bg-white rounded-lg border border-slate-200 p-6 text-center">
        <div className="text-4xl font-bold text-slate-900">{formatScore(student.overallComposite)}</div>
        <div className="text-sm text-slate-500 mt-1">Overall Composite — {scoreLabel(student.overallComposite)}</div>
      </div>

      {/* Three Pillars */}
      <div className="grid grid-cols-3 gap-4">
        {[
          { label: 'Academic (32%)', score: student.academicComposite },
          { label: 'Mil Skills (32%)', score: student.milSkillsComposite },
          { label: 'Leadership (36%)', score: student.leadershipComposite },
        ].map(p => (
          <div key={p.label} className="bg-white rounded-lg border border-slate-200 p-4 text-center">
            <div className={`text-2xl font-bold ${p.score >= 85 ? 'text-green-600' : p.score >= 75 ? 'text-yellow-600' : 'text-red-600'}`}>
              {formatScore(p.score)}
            </div>
            <div className="text-xs text-slate-500 mt-1">{p.label}</div>
            <div className="mt-2 h-2 bg-slate-100 rounded-full overflow-hidden">
              <div
                className={`h-full rounded-full ${p.score >= 85 ? 'bg-green-500' : p.score >= 75 ? 'bg-yellow-500' : 'bg-red-500'}`}
                style={{ width: `${Math.min(100, p.score)}%` }}
              />
            </div>
          </div>
        ))}
      </div>

      {/* Exam Chart */}
      {examData.length > 0 && (
        <div className="bg-white rounded-lg border border-slate-200 p-4">
          <h3 className="text-sm font-semibold text-slate-700 mb-4">My Exam Scores</h3>
          <ResponsiveContainer width="100%" height={200}>
            <BarChart data={examData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
              <XAxis dataKey="name" tick={{ fontSize: 12 }} />
              <YAxis domain={[0, 100]} tick={{ fontSize: 12 }} />
              <Tooltip />
              <ReferenceLine y={75} stroke="#ef4444" strokeDasharray="5 5" />
              <Bar dataKey="score" radius={[4, 4, 0, 0]}>
                {examData.map((entry, idx) => (
                  <Cell key={idx} fill={barColor(entry.score)} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>
      )}
    </div>
  )
}
