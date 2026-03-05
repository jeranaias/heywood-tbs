import { useState, useEffect } from 'react'
import { Users, TrendingUp, AlertTriangle, Percent } from 'lucide-react'
import { api } from '../lib/api'
import type { Student, StudentStats } from '../lib/types'
import { KPICard } from '../components/dashboard/KPICard'
import { StudentRoster } from '../components/students/StudentRoster'

export function Dashboard() {
  const [stats, setStats] = useState<StudentStats | null>(null)
  const [students, setStudents] = useState<Student[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      try {
        const [statsRes, studentsRes] = await Promise.all([
          api.getStudentStats(),
          api.getStudents(),
        ])
        setStats(statsRes)
        setStudents(studentsRes.students || [])
      } catch (err) {
        console.error('Failed to load dashboard:', err)
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  if (loading) {
    return (
      <div className="animate-pulse space-y-6">
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
          {[1,2,3,4].map(i => <div key={i} className="h-24 bg-slate-200 rounded-lg" />)}
        </div>
        <div className="h-96 bg-slate-200 rounded-lg" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-bold text-slate-900">Dashboard</h2>
        <p className="text-sm text-slate-500 mt-1">TBS Alpha Company — Class 1-26</p>
      </div>

      {/* KPI Cards */}
      {stats && (
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
          <KPICard
            title="Active Students"
            value={stats.activeStudents}
            icon={Users}
          />
          <KPICard
            title="Avg Composite"
            value={stats.avgComposite.toFixed(1)}
            icon={TrendingUp}
            variant={stats.avgComposite >= 80 ? 'success' : 'warning'}
          />
          <KPICard
            title="At-Risk"
            value={stats.atRiskCount}
            subtitle={`${stats.atRiskPercent.toFixed(1)}% of class`}
            icon={AlertTriangle}
            variant={stats.atRiskCount > 0 ? 'danger' : 'success'}
          />
          <KPICard
            title="At-Risk Rate"
            value={`${stats.atRiskPercent.toFixed(1)}%`}
            icon={Percent}
            variant={stats.atRiskPercent > 20 ? 'danger' : stats.atRiskPercent > 10 ? 'warning' : 'success'}
          />
        </div>
      )}

      {/* Phase Distribution */}
      {stats && Object.keys(stats.byPhase).length > 0 && (
        <div className="bg-white rounded-lg border border-slate-200 p-4">
          <h3 className="text-sm font-semibold text-slate-700 mb-3">Students by Phase</h3>
          <div className="flex gap-3">
            {Object.entries(stats.byPhase).map(([phase, count]) => (
              <div key={phase} className="flex-1 bg-slate-50 rounded-lg p-3 text-center">
                <div className="text-lg font-bold text-slate-900">{count}</div>
                <div className="text-xs text-slate-500">{phase.replace(' - ', '\n').split('\n')[0]}</div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Student Roster */}
      <div>
        <h3 className="text-lg font-semibold text-slate-800 mb-3">Student Roster</h3>
        <StudentRoster students={students} />
      </div>
    </div>
  )
}
