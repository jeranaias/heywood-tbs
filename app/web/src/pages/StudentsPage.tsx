import { useState, useEffect } from 'react'
import { api } from '../lib/api'
import type { Student } from '../lib/types'
import { StudentRoster } from '../components/students/StudentRoster'
import { Search, Users } from 'lucide-react'

export function StudentsPage() {
  const [students, setStudents] = useState<Student[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [phaseFilter, setPhaseFilter] = useState('')
  const [atRiskFilter, setAtRiskFilter] = useState(false)

  useEffect(() => {
    async function load() {
      try {
        const params: Record<string, string> = {}
        if (search) params.search = search
        if (phaseFilter) params.phase = phaseFilter
        if (atRiskFilter) params.atRisk = 'true'
        const res = await api.getStudents(params)
        setStudents(res.students || [])
      } catch (err) {
        console.error('Failed to load students:', err)
      } finally {
        setLoading(false)
      }
    }
    setLoading(true)
    load()
  }, [search, phaseFilter, atRiskFilter])

  if (loading && students.length === 0) {
    return (
      <div className="animate-pulse space-y-4">
        <div className="h-10 bg-slate-200 rounded-lg w-1/3" />
        <div className="h-96 bg-slate-200 rounded-lg" />
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-bold text-slate-900">Students</h2>
          <p className="text-sm text-slate-500 mt-0.5">
            {students.length} student{students.length !== 1 ? 's' : ''} shown
          </p>
        </div>
        <div className="flex items-center gap-2 text-sm text-slate-500">
          <Users className="w-4 h-4" />
          Alpha Company — Class 1-26
        </div>
      </div>

      {/* Filters */}
      <div className="flex flex-wrap items-center gap-3">
        <div className="relative flex-1 min-w-[200px] max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
          <input
            type="text"
            value={search}
            onChange={e => setSearch(e.target.value)}
            placeholder="Search by name or ID..."
            className="w-full pl-9 pr-4 py-2 bg-white border border-slate-200 rounded-lg text-sm placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20 focus:border-[var(--color-navy)]"
          />
        </div>
        <select
          value={phaseFilter}
          onChange={e => setPhaseFilter(e.target.value)}
          className="px-3 py-2 bg-white border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20"
        >
          <option value="">All Phases</option>
          <option value="Phase 1 - Foundations">Phase 1</option>
          <option value="Phase 2 - Warfighting">Phase 2</option>
          <option value="Phase 3 - Leadership">Phase 3</option>
        </select>
        <label className="flex items-center gap-2 text-sm text-slate-600 cursor-pointer">
          <input
            type="checkbox"
            checked={atRiskFilter}
            onChange={e => setAtRiskFilter(e.target.checked)}
            className="rounded border-slate-300"
          />
          At-Risk Only
        </label>
      </div>

      {/* Roster */}
      <StudentRoster students={students} hideSearch />
    </div>
  )
}
