import { useState, useEffect } from 'react'
import { Shield, AlertCircle, Clock, CheckCircle } from 'lucide-react'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell } from 'recharts'
import { api } from '../lib/api'
import type { QualStats, QualRecord, Instructor } from '../lib/types'
import { KPICard } from '../components/dashboard/KPICard'

export function InstructorQuals() {
  const [stats, setStats] = useState<QualStats | null>(null)
  const [records, setRecords] = useState<QualRecord[]>([])
  const [instructors, setInstructors] = useState<Instructor[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      try {
        const [s, r, i] = await Promise.all([
          api.getQualStats(),
          api.getQualRecords(),
          api.getInstructors(),
        ])
        setStats(s)
        setRecords(r as unknown as QualRecord[])
        setInstructors((i as { instructors: Instructor[] }).instructors || [])
      } catch (err) {
        console.error('Failed to load qual data:', err)
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  if (loading) {
    return <div className="animate-pulse h-96 bg-slate-200 rounded-lg" />
  }

  // Workload chart data
  const workloadData = instructors
    .filter(i => i.studentsAssigned > 0)
    .map(i => ({
      name: `${i.rank} ${i.lastName}`,
      students: i.studentsAssigned,
      role: i.role,
    }))
    .sort((a, b) => b.students - a.students)

  // Group records by status for the matrix
  const statusOrder: Record<string, number> = { 'Expired': 0, 'Critical': 1, 'Warning': 2, 'Caution': 3, 'Current': 4 }

  function statusColor(status: string): string {
    if (status.includes('Expired')) return 'bg-red-600 text-white'
    if (status.includes('Critical')) return 'bg-red-100 text-red-800'
    if (status.includes('Warning')) return 'bg-orange-100 text-orange-800'
    if (status.includes('Caution')) return 'bg-yellow-100 text-yellow-800'
    return 'bg-green-100 text-green-800'
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <div className="p-2 rounded-lg bg-blue-100">
          <Shield className="w-5 h-5 text-[var(--color-navy)]" />
        </div>
        <div>
          <h2 className="text-xl font-bold text-slate-900">Instructor Qualifications</h2>
          <p className="text-sm text-slate-500">{records.length} qualification records across {instructors.length} instructors</p>
        </div>
      </div>

      {/* KPI Cards */}
      {stats && (
        <div className="grid grid-cols-2 lg:grid-cols-5 gap-4">
          <KPICard title="Expired" value={stats.expiredCount} icon={AlertCircle} variant="danger" />
          <KPICard title="Critical (30d)" value={stats.expiring30} icon={Clock} variant="danger" />
          <KPICard title="Warning (60d)" value={stats.expiring60} icon={Clock} variant="warning" />
          <KPICard title="Caution (90d)" value={stats.expiring90} icon={Clock} variant="warning" />
          <KPICard title="Current" value={stats.currentCount} icon={CheckCircle} variant="success" />
        </div>
      )}

      {/* Coverage Gaps */}
      {stats && stats.coverageGaps && stats.coverageGaps.length > 0 && (
        <div className="bg-red-50 rounded-lg border border-red-200 p-4">
          <h3 className="text-sm font-semibold text-red-800 mb-3">Coverage Gaps (Below Minimum Required)</h3>
          <div className="space-y-2">
            {stats.coverageGaps.map(gap => (
              <div key={gap.qualCode} className="flex items-center justify-between bg-white rounded-lg p-3 border border-red-100">
                <div>
                  <div className="text-sm font-medium text-slate-800">{gap.qualName}</div>
                  <div className="text-xs text-slate-500">{gap.qualCode}</div>
                </div>
                <div className="text-right">
                  <div className="text-sm font-bold text-red-600">{gap.qualifiedCount} / {gap.requiredCount}</div>
                  <div className="text-xs text-red-500">Gap: {gap.gap}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Qualification Records Table */}
      <div className="bg-white rounded-lg border border-slate-200">
        <div className="px-4 py-3 border-b border-slate-200">
          <h3 className="text-sm font-semibold text-slate-700">All Qualification Records</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-200 bg-slate-50">
                <th className="px-3 py-2.5 text-left font-medium text-slate-600">Instructor</th>
                <th className="px-3 py-2.5 text-left font-medium text-slate-600">Qualification</th>
                <th className="px-3 py-2.5 text-left font-medium text-slate-600">Date Earned</th>
                <th className="px-3 py-2.5 text-left font-medium text-slate-600">Expires</th>
                <th className="px-3 py-2.5 text-left font-medium text-slate-600">Days Left</th>
                <th className="px-3 py-2.5 text-left font-medium text-slate-600">Status</th>
              </tr>
            </thead>
            <tbody>
              {records
                .sort((a, b) => (statusOrder[a.expirationStatus.split(' ')[0]] ?? 5) - (statusOrder[b.expirationStatus.split(' ')[0]] ?? 5))
                .map(r => (
                  <tr key={r.id} className="border-b border-slate-100 hover:bg-slate-50">
                    <td className="px-3 py-2 text-slate-700">{r.instructorName}</td>
                    <td className="px-3 py-2 text-slate-600">{r.qualName}</td>
                    <td className="px-3 py-2 text-slate-500">{r.dateEarned}</td>
                    <td className="px-3 py-2 text-slate-500">{r.expirationDate}</td>
                    <td className="px-3 py-2 text-slate-500">{r.daysUntilExpiration}</td>
                    <td className="px-3 py-2">
                      <span className={`inline-block px-2 py-0.5 rounded text-xs font-medium ${statusColor(r.expirationStatus)}`}>
                        {r.expirationStatus}
                      </span>
                    </td>
                  </tr>
                ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Workload Distribution */}
      {workloadData.length > 0 && (
        <div className="bg-white rounded-lg border border-slate-200 p-4">
          <h3 className="text-sm font-semibold text-slate-700 mb-4">SPC Workload Distribution</h3>
          <ResponsiveContainer width="100%" height={200}>
            <BarChart data={workloadData} layout="vertical">
              <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
              <XAxis type="number" tick={{ fontSize: 12 }} />
              <YAxis dataKey="name" type="category" width={120} tick={{ fontSize: 11 }} />
              <Tooltip />
              <Bar dataKey="students" radius={[0, 4, 4, 0]}>
                {workloadData.map((_, idx) => (
                  <Cell key={idx} fill={idx % 2 === 0 ? '#1a365d' : '#2a4a7f'} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>
      )}
    </div>
  )
}
