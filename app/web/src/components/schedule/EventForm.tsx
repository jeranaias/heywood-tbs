import { useState } from 'react'
import { X } from 'lucide-react'
import { api } from '../../lib/api'
import type { TrainingEvent } from '../../lib/types'

interface EventFormProps {
  event?: TrainingEvent | null
  onSave: () => void
  onCancel: () => void
}

export function EventForm({ event, onSave, onCancel }: EventFormProps) {
  const [title, setTitle] = useState(event?.title || '')
  const [code, setCode] = useState(event?.code || '')
  const [phase, setPhase] = useState(event?.phase || '')
  const [category, setCategory] = useState(event?.category || '')
  const [startDate, setStartDate] = useState(event?.startDate || '')
  const [endDate, setEndDate] = useState(event?.endDate || '')
  const [startTime, setStartTime] = useState(event?.startTime || '')
  const [endTime, setEndTime] = useState(event?.endTime || '')
  const [location, setLocation] = useState(event?.location || '')
  const [leadInstructor, setLeadInstructor] = useState(event?.leadInstructor || '')
  const [isGraded, setIsGraded] = useState(event?.isGraded || false)
  const [gradePillar, setGradePillar] = useState(event?.gradePillar || '')
  const [saving, setSaving] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!title) return
    setSaving(true)
    try {
      const data = {
        title, code, phase, category, startDate, endDate,
        startTime, endTime, location, leadInstructor, isGraded, gradePillar,
      }
      if (event?.id) {
        await api.updateTrainingEvent(event.id, data)
      } else {
        await api.createTrainingEvent(data)
      }
      onSave()
    } catch {
      // ignore
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-xl shadow-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between p-4 border-b border-slate-200">
          <h2 className="text-lg font-bold text-slate-900">
            {event?.id ? 'Edit Event' : 'New Training Event'}
          </h2>
          <button onClick={onCancel} className="p-1 hover:bg-slate-100 rounded">
            <X className="w-5 h-5 text-slate-500" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-4 space-y-3">
          <div>
            <label className="block text-xs font-medium text-slate-600 mb-1">Title *</label>
            <input type="text" value={title} onChange={e => setTitle(e.target.value)} required
              className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20" />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">Code</label>
              <input type="text" value={code} onChange={e => setCode(e.target.value)}
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20" />
            </div>
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">Phase</label>
              <input type="text" value={phase} onChange={e => setPhase(e.target.value)}
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20" />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">Start Date</label>
              <input type="date" value={startDate} onChange={e => setStartDate(e.target.value)}
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20" />
            </div>
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">End Date</label>
              <input type="date" value={endDate} onChange={e => setEndDate(e.target.value)}
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20" />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">Start Time</label>
              <input type="text" value={startTime} onChange={e => setStartTime(e.target.value)} placeholder="0600"
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20" />
            </div>
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">End Time</label>
              <input type="text" value={endTime} onChange={e => setEndTime(e.target.value)} placeholder="1800"
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20" />
            </div>
          </div>

          <div>
            <label className="block text-xs font-medium text-slate-600 mb-1">Location</label>
            <input type="text" value={location} onChange={e => setLocation(e.target.value)}
              className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20" />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">Lead Instructor</label>
              <input type="text" value={leadInstructor} onChange={e => setLeadInstructor(e.target.value)}
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20" />
            </div>
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">Category</label>
              <select value={category} onChange={e => setCategory(e.target.value)}
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20">
                <option value="">Select...</option>
                <option value="Academic">Academic</option>
                <option value="Field">Field</option>
                <option value="Range">Range</option>
                <option value="Physical Training">Physical Training</option>
                <option value="Admin">Admin</option>
                <option value="Evaluation">Evaluation</option>
              </select>
            </div>
          </div>

          <div className="flex items-center gap-4">
            <label className="flex items-center gap-2 text-sm">
              <input type="checkbox" checked={isGraded} onChange={e => setIsGraded(e.target.checked)}
                className="rounded border-slate-300" />
              Graded Event
            </label>
            {isGraded && (
              <select value={gradePillar} onChange={e => setGradePillar(e.target.value)}
                className="px-2 py-1 border border-slate-200 rounded text-xs">
                <option value="">Pillar...</option>
                <option value="Academic">Academic</option>
                <option value="Military Skills">Military Skills</option>
                <option value="Leadership">Leadership</option>
              </select>
            )}
          </div>

          <div className="flex justify-end gap-2 pt-2 border-t border-slate-200">
            <button type="button" onClick={onCancel}
              className="px-4 py-2 text-sm text-slate-600 hover:bg-slate-100 rounded-lg">Cancel</button>
            <button type="submit" disabled={!title || saving}
              className="px-4 py-2 text-sm font-medium bg-[var(--color-navy)] text-white rounded-lg hover:bg-[var(--color-navy-light)] disabled:opacity-50">
              {saving ? 'Saving...' : event?.id ? 'Update' : 'Create'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
