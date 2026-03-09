import { useState, useEffect } from 'react'
import { X, Sparkles, Loader2 } from 'lucide-react'
import { api } from '../../lib/api'
import type { Student, CounselingSession } from '../../lib/types'

interface CounselingFormProps {
  session?: CounselingSession | null
  onSave: () => void
  onCancel: () => void
}

const counselingTypes = [
  { value: 'initial', label: 'Initial Counseling' },
  { value: 'progress', label: 'Progress Review' },
  { value: 'event-driven', label: 'Event-Driven' },
  { value: 'end-of-phase', label: 'End of Phase' },
]

const statusOptions = [
  { value: 'draft', label: 'Draft' },
  { value: 'conducted', label: 'Conducted' },
  { value: 'completed', label: 'Completed' },
]

export function CounselingForm({ session, onSave, onCancel }: CounselingFormProps) {
  const [students, setStudents] = useState<Student[]>([])
  const [studentId, setStudentId] = useState(session?.studentId || '')
  const [type, setType] = useState(session?.type || 'initial')
  const [date, setDate] = useState(session?.date || new Date().toISOString().slice(0, 10))
  const [outline, setOutline] = useState(session?.outline || '')
  const [notes, setNotes] = useState(session?.notes || '')
  const [status, setStatus] = useState(session?.status || 'draft')
  const [followUps, setFollowUps] = useState(session?.followUps || [])
  const [saving, setSaving] = useState(false)
  const [generating, setGenerating] = useState(false)

  useEffect(() => {
    api.getStudents().then(res => setStudents(res.students || [])).catch(() => {})
  }, [])

  async function handleGenerate() {
    if (!studentId) return
    setGenerating(true)
    try {
      const res = await api.generateCounselingOutline(studentId, type)
      setOutline(res.outline)
    } catch {
      // ignore
    } finally {
      setGenerating(false)
    }
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!studentId) return
    setSaving(true)
    try {
      if (session?.id) {
        await api.updateCounseling(session.id, { notes, outline, status, type, date, followUps })
      } else {
        await api.createCounseling({ studentId, type, date, outline, notes, status, followUps })
      }
      onSave()
    } catch {
      // ignore
    } finally {
      setSaving(false)
    }
  }

  function addFollowUp() {
    setFollowUps([...followUps, { description: '', dueDate: '', status: 'pending' }])
  }

  function updateFollowUp(idx: number, field: string, value: string) {
    const updated = [...followUps]
    updated[idx] = { ...updated[idx], [field]: value }
    setFollowUps(updated)
  }

  function removeFollowUp(idx: number) {
    setFollowUps(followUps.filter((_, i) => i !== idx))
  }

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-xl shadow-2xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between p-4 border-b border-slate-200">
          <h2 className="text-lg font-bold text-slate-900">
            {session?.id ? 'Edit Counseling Session' : 'New Counseling Session'}
          </h2>
          <button onClick={onCancel} className="p-1 hover:bg-slate-100 rounded">
            <X className="w-5 h-5 text-slate-500" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-4 space-y-4">
          <div className="grid grid-cols-2 gap-4">
            {/* Student */}
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">Student</label>
              <select
                value={studentId}
                onChange={e => setStudentId(e.target.value)}
                disabled={!!session?.id}
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20"
              >
                <option value="">Select student...</option>
                {students.map(s => (
                  <option key={s.id} value={s.id}>
                    {s.rank} {s.lastName}, {s.firstName} ({s.company})
                  </option>
                ))}
              </select>
            </div>

            {/* Type */}
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">Type</label>
              <select
                value={type}
                onChange={e => setType(e.target.value)}
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20"
              >
                {counselingTypes.map(ct => (
                  <option key={ct.value} value={ct.value}>{ct.label}</option>
                ))}
              </select>
            </div>

            {/* Date */}
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">Date</label>
              <input
                type="date"
                value={date}
                onChange={e => setDate(e.target.value)}
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20"
              />
            </div>

            {/* Status */}
            <div>
              <label className="block text-xs font-medium text-slate-600 mb-1">Status</label>
              <select
                value={status}
                onChange={e => setStatus(e.target.value)}
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20"
              >
                {statusOptions.map(so => (
                  <option key={so.value} value={so.value}>{so.label}</option>
                ))}
              </select>
            </div>
          </div>

          {/* AI Outline */}
          <div>
            <div className="flex items-center justify-between mb-1">
              <label className="text-xs font-medium text-slate-600">Outline</label>
              <button
                type="button"
                onClick={handleGenerate}
                disabled={!studentId || generating}
                className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-[var(--color-navy)] hover:bg-slate-100 rounded disabled:opacity-50"
              >
                {generating ? <Loader2 className="w-3 h-3 animate-spin" /> : <Sparkles className="w-3 h-3" />}
                {generating ? 'Generating...' : 'Generate with AI'}
              </button>
            </div>
            <textarea
              value={outline}
              onChange={e => setOutline(e.target.value)}
              rows={6}
              placeholder="AI-generated outline will appear here, or write your own..."
              className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20 resize-none font-mono"
            />
          </div>

          {/* Notes */}
          <div>
            <label className="block text-xs font-medium text-slate-600 mb-1">Counselor Notes</label>
            <textarea
              value={notes}
              onChange={e => setNotes(e.target.value)}
              rows={4}
              placeholder="Your notes from the counseling session..."
              className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20 resize-none"
            />
          </div>

          {/* Follow-ups */}
          <div>
            <div className="flex items-center justify-between mb-2">
              <label className="text-xs font-medium text-slate-600">Follow-Up Actions</label>
              <button
                type="button"
                onClick={addFollowUp}
                className="text-xs text-[var(--color-navy)] hover:underline"
              >
                + Add Follow-Up
              </button>
            </div>
            {followUps.map((fu, idx) => (
              <div key={idx} className="flex items-center gap-2 mb-2">
                <input
                  type="text"
                  value={fu.description}
                  onChange={e => updateFollowUp(idx, 'description', e.target.value)}
                  placeholder="Action description..."
                  className="flex-1 px-2 py-1.5 border border-slate-200 rounded text-xs focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20"
                />
                <input
                  type="date"
                  value={fu.dueDate}
                  onChange={e => updateFollowUp(idx, 'dueDate', e.target.value)}
                  className="px-2 py-1.5 border border-slate-200 rounded text-xs focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20"
                />
                <select
                  value={fu.status}
                  onChange={e => updateFollowUp(idx, 'status', e.target.value)}
                  className="px-2 py-1.5 border border-slate-200 rounded text-xs"
                >
                  <option value="pending">Pending</option>
                  <option value="completed">Done</option>
                </select>
                <button type="button" onClick={() => removeFollowUp(idx)} className="text-slate-400 hover:text-red-500">
                  <X className="w-4 h-4" />
                </button>
              </div>
            ))}
          </div>

          {/* Actions */}
          <div className="flex justify-end gap-2 pt-2 border-t border-slate-200">
            <button
              type="button"
              onClick={onCancel}
              className="px-4 py-2 text-sm text-slate-600 hover:bg-slate-100 rounded-lg"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={!studentId || saving}
              className="px-4 py-2 text-sm font-medium bg-[var(--color-navy)] text-white rounded-lg hover:bg-[var(--color-navy-light)] disabled:opacity-50"
            >
              {saving ? 'Saving...' : session?.id ? 'Update' : 'Create'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
