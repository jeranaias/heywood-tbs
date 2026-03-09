import { useState, useEffect } from 'react'
import { Calendar, Filter, Plus, Edit2, Trash2 } from 'lucide-react'
import { api } from '../lib/api'
import type { TrainingEvent } from '../lib/types'
import { useAuth } from '../hooks/useAuth'
import { EventForm } from '../components/schedule/EventForm'

export function Schedule() {
  const [events, setEvents] = useState<TrainingEvent[]>([])
  const [loading, setLoading] = useState(true)
  const [phaseFilter, setPhaseFilter] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [editEvent, setEditEvent] = useState<TrainingEvent | null>(null)
  const { auth } = useAuth()
  const canManage = auth.role === 'xo' || auth.role === 'staff'

  function loadEvents() {
    api.getSchedule(phaseFilter ? { phase: phaseFilter } : undefined)
      .then(res => setEvents(res.events || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadEvents()
  }, [phaseFilter])

  async function handleDelete(id: string) {
    try {
      await api.deleteTrainingEvent(id)
      loadEvents()
    } catch { /* ignore */ }
  }

  function handleSaved() {
    setShowForm(false)
    setEditEvent(null)
    loadEvents()
  }

  const phases = [...new Set(events.map(e => e.phase))].sort()

  const categoryColors: Record<string, string> = {
    'Academic': 'bg-blue-100 text-blue-800 border-blue-200',
    'Field': 'bg-green-100 text-green-800 border-green-200',
    'Range': 'bg-orange-100 text-orange-800 border-orange-200',
    'Physical Training': 'bg-purple-100 text-purple-800 border-purple-200',
    'Admin': 'bg-slate-100 text-slate-800 border-slate-200',
    'Evaluation': 'bg-red-100 text-red-800 border-red-200',
  }

  function getCategoryStyle(category: string): string {
    return categoryColors[category] || 'bg-slate-100 text-slate-800 border-slate-200'
  }

  if (loading) {
    return <div className="animate-pulse h-96 bg-slate-200 rounded-lg" />
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-blue-100">
            <Calendar className="w-5 h-5 text-[var(--color-navy)]" />
          </div>
          <div>
            <h2 className="text-xl font-bold text-slate-900">Training Schedule</h2>
            <p className="text-sm text-slate-500">{events.length} events</p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Filter className="w-4 h-4 text-slate-400" />
          <select
            value={phaseFilter}
            onChange={e => setPhaseFilter(e.target.value)}
            className="text-sm border border-slate-200 rounded-lg px-3 py-1.5 bg-white"
          >
            <option value="">All Phases</option>
            {phases.map(p => (
              <option key={p} value={p}>{p}</option>
            ))}
          </select>
          {canManage && (
            <button
              onClick={() => { setEditEvent(null); setShowForm(true) }}
              className="flex items-center gap-1 px-3 py-1.5 text-sm font-medium bg-[var(--color-navy)] text-white rounded-lg hover:bg-[var(--color-navy-light)]"
            >
              <Plus className="w-4 h-4" /> New Event
            </button>
          )}
        </div>
      </div>

      {/* Event Cards */}
      <div className="space-y-3">
        {events.map(evt => (
          <div key={evt.id} className="bg-white rounded-lg border border-slate-200 p-4 hover:border-slate-300 transition-colors">
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <div className="flex items-center gap-2 mb-1">
                  <h3 className="text-sm font-semibold text-slate-800">{evt.title}</h3>
                  {evt.isGraded && (
                    <span className="px-1.5 py-0.5 bg-red-100 text-red-700 text-xs font-medium rounded">GRADED</span>
                  )}
                </div>
                <div className="flex items-center gap-3 text-xs text-slate-500">
                  <span>{evt.code}</span>
                  <span>•</span>
                  <span>{evt.startDate}</span>
                  <span>•</span>
                  <span>{evt.startTime} - {evt.endTime}</span>
                  <span>•</span>
                  <span>{evt.durationHours}h</span>
                </div>
                <div className="text-xs text-slate-500 mt-1">{evt.location}</div>
                {evt.leadInstructor && (
                  <div className="text-xs text-slate-400 mt-1">Lead: {evt.leadInstructor}</div>
                )}
              </div>
              <div className="flex flex-col items-end gap-1.5 ml-4">
                <span className={`px-2 py-0.5 text-xs font-medium rounded border ${getCategoryStyle(evt.category)}`}>
                  {evt.category}
                </span>
                {evt.gradePillar && evt.gradePillar !== 'Not Graded' && (
                  <span className="text-xs text-slate-500">{evt.gradePillar}</span>
                )}
                <span className={`text-xs ${
                  evt.status === 'Complete' ? 'text-green-600' :
                  evt.status === 'Scheduled' ? 'text-blue-600' :
                  'text-slate-500'
                }`}>{evt.status}</span>
                {canManage && (
                  <div className="flex gap-1 mt-1">
                    <button
                      onClick={() => { setEditEvent(evt); setShowForm(true) }}
                      className="p-1 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded"
                      title="Edit"
                    >
                      <Edit2 className="w-3.5 h-3.5" />
                    </button>
                    <button
                      onClick={() => handleDelete(evt.id)}
                      className="p-1 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded"
                      title="Delete"
                    >
                      <Trash2 className="w-3.5 h-3.5" />
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>

      {showForm && (
        <EventForm
          event={editEvent}
          onSave={handleSaved}
          onCancel={() => { setShowForm(false); setEditEvent(null) }}
        />
      )}
    </div>
  )
}
