import { useState, useMemo } from 'react'
import { Link } from 'react-router-dom'
import { Search, ChevronUp, ChevronDown } from 'lucide-react'
import type { Student } from '../../lib/types'
import { formatScore, scoreBadge, trendIcon, trendColor } from '../../lib/utils'

interface StudentRosterProps {
  students: Student[]
  compact?: boolean
}

type SortKey = 'lastName' | 'overallComposite' | 'academicComposite' | 'milSkillsComposite' | 'leadershipComposite' | 'companyRank'
type SortDir = 'asc' | 'desc'

export function StudentRoster({ students, compact = false }: StudentRosterProps) {
  const [search, setSearch] = useState('')
  const [sortKey, setSortKey] = useState<SortKey>('companyRank')
  const [sortDir, setSortDir] = useState<SortDir>('asc')
  const [page, setPage] = useState(0)
  const pageSize = compact ? 10 : 20

  const filtered = useMemo(() => {
    let result = students
    if (search) {
      const q = search.toLowerCase()
      result = result.filter(s =>
        s.lastName.toLowerCase().includes(q) ||
        s.firstName.toLowerCase().includes(q) ||
        s.id.toLowerCase().includes(q)
      )
    }
    result.sort((a, b) => {
      const av = a[sortKey] ?? 0
      const bv = b[sortKey] ?? 0
      const cmp = typeof av === 'string' ? (av as string).localeCompare(bv as string) : (av as number) - (bv as number)
      return sortDir === 'asc' ? cmp : -cmp
    })
    return result
  }, [students, search, sortKey, sortDir])

  const pageCount = Math.ceil(filtered.length / pageSize)
  const pageStudents = filtered.slice(page * pageSize, (page + 1) * pageSize)

  function toggleSort(key: SortKey) {
    if (sortKey === key) {
      setSortDir(d => d === 'asc' ? 'desc' : 'asc')
    } else {
      setSortKey(key)
      setSortDir(key === 'lastName' ? 'asc' : 'desc')
    }
    setPage(0)
  }

  function SortIcon({ col }: { col: SortKey }) {
    if (sortKey !== col) return null
    return sortDir === 'asc' ? <ChevronUp className="w-3 h-3" /> : <ChevronDown className="w-3 h-3" />
  }

  return (
    <div>
      <div className="flex items-center gap-3 mb-4">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
          <input
            type="text"
            value={search}
            onChange={e => { setSearch(e.target.value); setPage(0) }}
            placeholder="Search students..."
            className="w-full pl-9 pr-4 py-2 bg-white border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20"
          />
        </div>
        <span className="text-sm text-slate-500">{filtered.length} students</span>
      </div>

      <div className="overflow-x-auto bg-white rounded-lg border border-slate-200">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-slate-200 bg-slate-50">
              <th className="px-3 py-2.5 text-left font-medium text-slate-600">Rank</th>
              <th
                className="px-3 py-2.5 text-left font-medium text-slate-600 cursor-pointer hover:text-slate-900"
                onClick={() => toggleSort('companyRank')}
              >
                <span className="flex items-center gap-1"># <SortIcon col="companyRank" /></span>
              </th>
              <th
                className="px-3 py-2.5 text-left font-medium text-slate-600 cursor-pointer hover:text-slate-900"
                onClick={() => toggleSort('lastName')}
              >
                <span className="flex items-center gap-1">Name <SortIcon col="lastName" /></span>
              </th>
              <th className="px-3 py-2.5 text-left font-medium text-slate-600">Phase</th>
              <th
                className="px-3 py-2.5 text-right font-medium text-slate-600 cursor-pointer hover:text-slate-900"
                onClick={() => toggleSort('academicComposite')}
              >
                <span className="flex items-center justify-end gap-1">Acad <SortIcon col="academicComposite" /></span>
              </th>
              <th
                className="px-3 py-2.5 text-right font-medium text-slate-600 cursor-pointer hover:text-slate-900"
                onClick={() => toggleSort('milSkillsComposite')}
              >
                <span className="flex items-center justify-end gap-1">MilSk <SortIcon col="milSkillsComposite" /></span>
              </th>
              <th
                className="px-3 py-2.5 text-right font-medium text-slate-600 cursor-pointer hover:text-slate-900"
                onClick={() => toggleSort('leadershipComposite')}
              >
                <span className="flex items-center justify-end gap-1">Ldr <SortIcon col="leadershipComposite" /></span>
              </th>
              <th
                className="px-3 py-2.5 text-right font-medium text-slate-600 cursor-pointer hover:text-slate-900"
                onClick={() => toggleSort('overallComposite')}
              >
                <span className="flex items-center justify-end gap-1">Overall <SortIcon col="overallComposite" /></span>
              </th>
              <th className="px-3 py-2.5 text-center font-medium text-slate-600">Trend</th>
              {!compact && <th className="px-3 py-2.5 text-left font-medium text-slate-600">Status</th>}
            </tr>
          </thead>
          <tbody>
            {pageStudents.map(s => (
              <tr key={s.id} className="border-b border-slate-100 hover:bg-slate-50 transition-colors">
                <td className="px-3 py-2 text-slate-600">{s.rank}</td>
                <td className="px-3 py-2 text-slate-500">{s.companyRank || '—'}</td>
                <td className="px-3 py-2">
                  <Link
                    to={`/students/${s.id}`}
                    className="text-[var(--color-navy)] hover:underline font-medium"
                  >
                    {s.lastName}, {s.firstName}
                  </Link>
                </td>
                <td className="px-3 py-2 text-slate-600 text-xs">{s.phase.replace('Phase ', 'Ph ').substring(0, 15)}</td>
                <td className="px-3 py-2 text-right">
                  <span className={`inline-block px-2 py-0.5 rounded text-xs font-medium ${scoreBadge(s.academicComposite)}`}>
                    {formatScore(s.academicComposite)}
                  </span>
                </td>
                <td className="px-3 py-2 text-right">
                  <span className={`inline-block px-2 py-0.5 rounded text-xs font-medium ${scoreBadge(s.milSkillsComposite)}`}>
                    {formatScore(s.milSkillsComposite)}
                  </span>
                </td>
                <td className="px-3 py-2 text-right">
                  <span className={`inline-block px-2 py-0.5 rounded text-xs font-medium ${scoreBadge(s.leadershipComposite)}`}>
                    {formatScore(s.leadershipComposite)}
                  </span>
                </td>
                <td className="px-3 py-2 text-right">
                  <span className={`inline-block px-2 py-0.5 rounded text-xs font-bold ${scoreBadge(s.overallComposite)}`}>
                    {formatScore(s.overallComposite)}
                  </span>
                </td>
                <td className={`px-3 py-2 text-center font-medium ${trendColor(s.trend)}`}>
                  {trendIcon(s.trend)}
                </td>
                {!compact && (
                  <td className="px-3 py-2">
                    {s.atRisk ? (
                      <span className="inline-block px-2 py-0.5 rounded text-xs font-medium bg-red-100 text-red-800">
                        At Risk
                      </span>
                    ) : (
                      <span className="text-xs text-slate-500">{s.status}</span>
                    )}
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      {pageCount > 1 && (
        <div className="flex items-center justify-between mt-3">
          <span className="text-sm text-slate-500">
            Page {page + 1} of {pageCount}
          </span>
          <div className="flex gap-2">
            <button
              onClick={() => setPage(p => Math.max(0, p - 1))}
              disabled={page === 0}
              className="px-3 py-1.5 text-sm border border-slate-200 rounded-lg hover:bg-slate-50 disabled:opacity-50"
            >
              Previous
            </button>
            <button
              onClick={() => setPage(p => Math.min(pageCount - 1, p + 1))}
              disabled={page >= pageCount - 1}
              className="px-3 py-1.5 text-sm border border-slate-200 rounded-lg hover:bg-slate-50 disabled:opacity-50"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
