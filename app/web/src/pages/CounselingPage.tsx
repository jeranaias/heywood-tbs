import { useState, useEffect } from 'react'
import { ClipboardList, Plus, Eye, Edit2, Printer } from 'lucide-react'
import { api } from '../lib/api'
import type { CounselingSession } from '../lib/types'
import { CounselingForm } from '../components/counseling/CounselingForm'
import { PrintCounseling } from '../components/counseling/PrintCounseling'

const typeLabels: Record<string, string> = {
  'initial': 'Initial',
  'progress': 'Progress',
  'event-driven': 'Event-Driven',
  'end-of-phase': 'End of Phase',
}

const statusColors: Record<string, string> = {
  'draft': 'bg-yellow-100 text-yellow-800',
  'conducted': 'bg-blue-100 text-blue-800',
  'completed': 'bg-green-100 text-green-800',
}

export function CounselingPage() {
  const [sessions, setSessions] = useState<CounselingSession[]>([])
  const [loading, setLoading] = useState(true)
  const [filterStudent, setFilterStudent] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [editSession, setEditSession] = useState<CounselingSession | null>(null)
  const [printSession, setPrintSession] = useState<CounselingSession | null>(null)

  function loadSessions() {
    const params: Record<string, string> = {}
    if (filterStudent) params.studentId = filterStudent
    api.getCounselings(params)
      .then(res => setSessions(res.sessions || []))
      .catch(() => {})
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadSessions()
  }, [filterStudent])

  function handleSaved() {
    setShowForm(false)
    setEditSession(null)
    loadSessions()
  }

  if (printSession) {
    return (
      <div>
        <div className="no-print mb-4 flex gap-2">
          <button
            onClick={() => window.print()}
            className="flex items-center gap-1 px-3 py-2 text-sm font-medium bg-[var(--color-navy)] text-white rounded-lg hover:bg-[var(--color-navy-light)]"
          >
            <Printer className="w-4 h-4" /> Print
          </button>
          <button
            onClick={() => setPrintSession(null)}
            className="px-3 py-2 text-sm text-slate-600 hover:bg-slate-100 rounded-lg"
          >
            Back
          </button>
        </div>
        <PrintCounseling session={printSession} />
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-slate-900 flex items-center gap-2">
          <ClipboardList className="w-5 h-5" /> Counseling Sessions
        </h1>
        <button
          onClick={() => { setEditSession(null); setShowForm(true) }}
          className="flex items-center gap-1 px-3 py-2 text-sm font-medium bg-[var(--color-navy)] text-white rounded-lg hover:bg-[var(--color-navy-light)]"
        >
          <Plus className="w-4 h-4" /> New Counseling
        </button>
      </div>

      {/* Filter */}
      <div className="flex gap-2">
        <input
          type="text"
          value={filterStudent}
          onChange={e => setFilterStudent(e.target.value)}
          placeholder="Filter by Student ID..."
          className="px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20"
        />
      </div>

      {/* Sessions list */}
      {loading ? (
        <div className="animate-pulse space-y-3">
          {[1, 2, 3].map(i => (
            <div key={i} className="h-20 bg-slate-200 rounded-lg" />
          ))}
        </div>
      ) : sessions.length === 0 ? (
        <div className="text-center py-12 text-slate-500">
          No counseling sessions yet. Click "New Counseling" to start.
        </div>
      ) : (
        <div className="space-y-3">
          {sessions.map(session => (
            <div key={session.id} className="bg-white rounded-lg border border-slate-200 p-4 hover:border-slate-300 transition-colors">
              <div className="flex items-start justify-between">
                <div>
                  <div className="flex items-center gap-2 mb-1">
                    <span className="font-semibold text-slate-900">{session.studentName}</span>
                    <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${statusColors[session.status] || 'bg-slate-100 text-slate-600'}`}>
                      {session.status}
                    </span>
                    <span className="px-2 py-0.5 bg-slate-100 text-slate-600 rounded-full text-xs">
                      {typeLabels[session.type] || session.type}
                    </span>
                  </div>
                  <div className="text-xs text-slate-500">
                    {session.date || 'No date'} | Counselor: {session.counselorName}
                    {session.followUps && session.followUps.length > 0 && (
                      <span className="ml-2 text-blue-600">
                        {session.followUps.filter(f => f.status === 'pending').length} pending follow-ups
                      </span>
                    )}
                  </div>
                  {session.notes && (
                    <p className="text-sm text-slate-600 mt-1 line-clamp-2">{session.notes}</p>
                  )}
                </div>
                <div className="flex gap-1">
                  <button
                    onClick={() => setPrintSession(session)}
                    className="p-1.5 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded"
                    title="Print"
                  >
                    <Printer className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => { setEditSession(session); setShowForm(true) }}
                    className="p-1.5 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded"
                    title="Edit"
                  >
                    <Edit2 className="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Form modal */}
      {showForm && (
        <CounselingForm
          session={editSession}
          onSave={handleSaved}
          onCancel={() => { setShowForm(false); setEditSession(null) }}
        />
      )}
    </div>
  )
}
