import { useState, useEffect } from 'react'
import { Link, Navigate } from 'react-router-dom'
import { Users, TrendingUp, AlertTriangle, Shield, Calendar, MessageSquare } from 'lucide-react'
import { api } from '../lib/api'
import type { Student, StudentStats, QualStats } from '../lib/types'
import { KPICard } from '../components/dashboard/KPICard'
import { useAuth } from '../hooks/useAuth'

export function Dashboard() {
  const { auth } = useAuth()

  // Students should only see their own record, not the staff dashboard
  if (auth.role === 'student') return <Navigate to="/my-record" replace />
  const [stats, setStats] = useState<StudentStats | null>(null)
  const [qualStats, setQualStats] = useState<QualStats | null>(null)
  const [atRisk, setAtRisk] = useState<Student[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      try {
        const [statsRes, atRiskRes, qualRes] = await Promise.all([
          api.getStudentStats(),
          api.getAtRiskStudents(),
          api.getQualStats().catch(() => null),
        ])
        setStats(statsRes)
        setAtRisk(atRiskRes.students || [])
        if (qualRes) setQualStats(qualRes)
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
        <div className="h-48 bg-slate-200 rounded-lg" />
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
            title="Qual Alerts"
            value={qualStats ? qualStats.expiredCount + qualStats.expiring30 : '—'}
            subtitle="Expired + critical"
            icon={Shield}
            variant={qualStats && (qualStats.expiredCount + qualStats.expiring30) > 0 ? 'danger' : 'success'}
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
                <div className="text-xs text-slate-500">{phase.replace(/ - .*/, '')}</div>
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* At-Risk Students */}
        <div className="bg-white rounded-lg border border-slate-200 p-4">
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-sm font-semibold text-slate-700">At-Risk Students</h3>
            <Link to="/at-risk" className="text-xs text-[var(--color-navy)] hover:underline">
              View all {atRisk.length} →
            </Link>
          </div>
          {atRisk.length === 0 ? (
            <p className="text-sm text-slate-500">No students currently at-risk.</p>
          ) : (
            <div className="space-y-2">
              {atRisk.slice(0, 5).map(s => (
                <Link
                  key={s.id}
                  to={`/students/${s.id}`}
                  className="flex items-center justify-between p-2.5 rounded-lg hover:bg-slate-50 transition-colors"
                >
                  <div>
                    <span className="text-sm font-medium text-slate-800">{s.id}</span>
                    <span className="text-sm text-slate-500 ml-2">{s.rank} {s.lastName}</span>
                  </div>
                  <div className="flex items-center gap-3 text-xs">
                    <span className="text-slate-500">Overall: <strong className={s.overallComposite < 78 ? 'text-red-600' : 'text-amber-600'}>{s.overallComposite.toFixed(1)}</strong></span>
                    <span className={`px-1.5 py-0.5 rounded ${s.trend === 'down' ? 'bg-red-100 text-red-700' : 'bg-slate-100 text-slate-600'}`}>
                      {s.trend === 'down' ? '↓' : s.trend === 'up' ? '↑' : '→'}
                    </span>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </div>

        {/* Quick Actions */}
        <div className="bg-white rounded-lg border border-slate-200 p-4">
          <h3 className="text-sm font-semibold text-slate-700 mb-3">Quick Actions</h3>
          <div className="space-y-2">
            <Link to="/chat" className="flex items-center gap-3 p-3 rounded-lg hover:bg-slate-50 border border-slate-100 transition-colors">
              <MessageSquare className="w-5 h-5 text-[var(--color-navy)]" />
              <div>
                <div className="text-sm font-medium text-slate-800">Chat with Heywood</div>
                <div className="text-xs text-slate-500">Morning brief, student analysis, counseling prep</div>
              </div>
            </Link>
            <Link to="/students" className="flex items-center gap-3 p-3 rounded-lg hover:bg-slate-50 border border-slate-100 transition-colors">
              <Users className="w-5 h-5 text-[var(--color-navy)]" />
              <div>
                <div className="text-sm font-medium text-slate-800">Student Roster</div>
                <div className="text-xs text-slate-500">Search, filter, and review all students</div>
              </div>
            </Link>
            <Link to="/schedule" className="flex items-center gap-3 p-3 rounded-lg hover:bg-slate-50 border border-slate-100 transition-colors">
              <Calendar className="w-5 h-5 text-[var(--color-navy)]" />
              <div>
                <div className="text-sm font-medium text-slate-800">Training Schedule</div>
                <div className="text-xs text-slate-500">Upcoming events and instructor assignments</div>
              </div>
            </Link>
          </div>
        </div>
      </div>

      {/* Qual Coverage Gaps */}
      {qualStats && qualStats.coverageGaps.length > 0 && (
        <div className="bg-white rounded-lg border border-red-200 p-4">
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-sm font-semibold text-red-700">Qualification Coverage Gaps</h3>
            <Link to="/instructor-quals" className="text-xs text-[var(--color-navy)] hover:underline">
              Full details →
            </Link>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
            {qualStats.coverageGaps.slice(0, 6).map(g => (
              <div key={g.qualCode} className="flex items-center justify-between p-2.5 bg-red-50 rounded-lg">
                <span className="text-sm text-slate-700 truncate mr-2">{g.qualName}</span>
                <span className="text-xs font-bold text-red-700 whitespace-nowrap">-{g.gap}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
