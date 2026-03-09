import { useState, useEffect } from 'react'
import { api } from '../lib/api'
import type { Student } from '../lib/types'
import { StudentRoster } from '../components/students/StudentRoster'
import { ExportButton } from '../components/common/ExportButton'
import { AlertTriangle } from 'lucide-react'

export function AtRisk() {
  const [students, setStudents] = useState<Student[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getAtRiskStudents()
      .then(res => setStudents(res.students || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) {
    return <div className="animate-pulse h-96 bg-slate-200 rounded-lg" />
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <div className="p-2 rounded-lg bg-red-100">
          <AlertTriangle className="w-5 h-5 text-red-600" />
        </div>
        <div>
          <h2 className="text-xl font-bold text-slate-900">At-Risk Students</h2>
          <p className="text-sm text-slate-500">{students.length} students flagged — any pillar below 75 or negative trend</p>
        </div>
      </div>

      <div className="flex justify-end">
        <ExportButton url="/api/v1/export/at-risk" filename="at-risk-students.csv" />
      </div>

      <StudentRoster students={students} />
    </div>
  )
}
