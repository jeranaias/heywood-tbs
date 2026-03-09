import { useState, useEffect } from 'react'
import { BarChart3 } from 'lucide-react'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell } from 'recharts'
import { ExportButton } from '../components/common/ExportButton'

interface CompanySummary {
  company: string
  studentCount: number
  avgAcademic: number
  avgMilSkills: number
  avgLeadership: number
  avgOverall: number
  atRiskCount: number
  atRiskPct: number
}

export function ReportsPage() {
  const [companies, setCompanies] = useState<CompanySummary[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/v1/reports/company-summary', { credentials: 'include' })
      .then(r => r.json())
      .then(data => {
        setCompanies(data.companies || [])
        setTotal(data.total || 0)
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  function barColor(value: number) {
    if (value >= 85) return '#22c55e'
    if (value >= 75) return '#eab308'
    return '#ef4444'
  }

  if (loading) {
    return <div className="animate-pulse h-96 bg-slate-200 rounded-lg" />
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-slate-900 flex items-center gap-2">
          <BarChart3 className="w-5 h-5" /> Company Performance Reports
        </h1>
        <div className="flex gap-2">
          <ExportButton url="/api/v1/export/students" filename="students.csv" label="Export Roster" />
          <ExportButton url="/api/v1/export/at-risk" filename="at-risk.csv" label="Export At-Risk" />
          <ExportButton url="/api/v1/export/qual-records" filename="qual-records.csv" label="Export Quals" />
          <ExportButton url="/api/v1/export/counselings" filename="counselings.csv" label="Export Counselings" />
        </div>
      </div>

      <div className="text-sm text-slate-500">Total students: {total}</div>

      {/* Overall Composite by Company */}
      {companies.length > 0 && (
        <div className="bg-white rounded-lg border border-slate-200 p-4">
          <h2 className="text-sm font-semibold text-slate-700 mb-4">Average Overall Composite by Company</h2>
          <ResponsiveContainer width="100%" height={250}>
            <BarChart data={companies}>
              <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
              <XAxis dataKey="company" tick={{ fontSize: 12 }} />
              <YAxis domain={[0, 100]} tick={{ fontSize: 12 }} />
              <Tooltip />
              <Bar dataKey="avgOverall" name="Avg Overall" radius={[4, 4, 0, 0]}>
                {companies.map((entry, idx) => (
                  <Cell key={idx} fill={barColor(entry.avgOverall)} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>
      )}

      {/* Pillar breakdown */}
      {companies.length > 0 && (
        <div className="bg-white rounded-lg border border-slate-200 p-4">
          <h2 className="text-sm font-semibold text-slate-700 mb-4">Pillar Averages by Company</h2>
          <ResponsiveContainer width="100%" height={250}>
            <BarChart data={companies}>
              <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
              <XAxis dataKey="company" tick={{ fontSize: 12 }} />
              <YAxis domain={[0, 100]} tick={{ fontSize: 12 }} />
              <Tooltip />
              <Bar dataKey="avgAcademic" name="Academic" fill="#3b82f6" radius={[2, 2, 0, 0]} />
              <Bar dataKey="avgMilSkills" name="Mil Skills" fill="#10b981" radius={[2, 2, 0, 0]} />
              <Bar dataKey="avgLeadership" name="Leadership" fill="#f59e0b" radius={[2, 2, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>
      )}

      {/* Company table */}
      <div className="bg-white rounded-lg border border-slate-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="bg-slate-50 border-b border-slate-200">
              <th className="text-left px-4 py-2 font-medium text-slate-600">Company</th>
              <th className="text-right px-4 py-2 font-medium text-slate-600">Students</th>
              <th className="text-right px-4 py-2 font-medium text-slate-600">Academic</th>
              <th className="text-right px-4 py-2 font-medium text-slate-600">Mil Skills</th>
              <th className="text-right px-4 py-2 font-medium text-slate-600">Leadership</th>
              <th className="text-right px-4 py-2 font-medium text-slate-600">Overall</th>
              <th className="text-right px-4 py-2 font-medium text-slate-600">At-Risk</th>
            </tr>
          </thead>
          <tbody>
            {companies.map(co => (
              <tr key={co.company} className="border-b border-slate-100 hover:bg-slate-50">
                <td className="px-4 py-2 font-medium text-slate-900">{co.company}</td>
                <td className="px-4 py-2 text-right text-slate-600">{co.studentCount}</td>
                <td className="px-4 py-2 text-right" style={{ color: barColor(co.avgAcademic) }}>{co.avgAcademic.toFixed(1)}</td>
                <td className="px-4 py-2 text-right" style={{ color: barColor(co.avgMilSkills) }}>{co.avgMilSkills.toFixed(1)}</td>
                <td className="px-4 py-2 text-right" style={{ color: barColor(co.avgLeadership) }}>{co.avgLeadership.toFixed(1)}</td>
                <td className="px-4 py-2 text-right font-semibold" style={{ color: barColor(co.avgOverall) }}>{co.avgOverall.toFixed(1)}</td>
                <td className="px-4 py-2 text-right">
                  <span className={co.atRiskPct > 15 ? 'text-red-600 font-medium' : 'text-slate-600'}>
                    {co.atRiskCount} ({co.atRiskPct.toFixed(0)}%)
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
