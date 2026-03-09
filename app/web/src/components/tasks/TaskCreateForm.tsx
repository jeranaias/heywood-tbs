import { useState } from 'react'
import { X } from 'lucide-react'
import { api } from '../../lib/api'

interface TaskCreateFormProps {
  onCreated: () => void
  onClose: () => void
}

export function TaskCreateForm({ onCreated, onClose }: TaskCreateFormProps) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [priority, setPriority] = useState('medium')
  const [assignedTo, setAssignedTo] = useState('staff')
  const [dueDate, setDueDate] = useState('')
  const [relatedId, setRelatedId] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!title.trim()) {
      setError('Title is required')
      return
    }
    setSubmitting(true)
    setError('')
    try {
      await api.createTask({
        title: title.trim(),
        description: description.trim(),
        priority,
        assignedTo,
        dueDate: dueDate || undefined,
        relatedId: relatedId.trim() || undefined,
      })
      onCreated()
    } catch (err) {
      setError('Failed to create task')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-lg mx-4 p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-bold text-slate-900">New Task</h2>
          <button onClick={onClose} className="p-1 text-slate-400 hover:text-slate-600">
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">Title *</label>
            <input
              type="text"
              value={title}
              onChange={e => setTitle(e.target.value)}
              placeholder="e.g., Review PT scores for at-risk students"
              className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20 focus:border-[var(--color-navy)]"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">Description</label>
            <textarea
              value={description}
              onChange={e => setDescription(e.target.value)}
              rows={3}
              placeholder="Details about what needs to be done..."
              className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20 focus:border-[var(--color-navy)] resize-none"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">Priority</label>
              <select
                value={priority}
                onChange={e => setPriority(e.target.value)}
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20 focus:border-[var(--color-navy)]"
              >
                <option value="high">High</option>
                <option value="medium">Medium</option>
                <option value="low">Low</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">Assign To</label>
              <select
                value={assignedTo}
                onChange={e => setAssignedTo(e.target.value)}
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20 focus:border-[var(--color-navy)]"
              >
                <option value="xo">XO</option>
                <option value="staff">Staff</option>
                <option value="spc">SPC</option>
              </select>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">Due Date</label>
              <input
                type="date"
                value={dueDate}
                onChange={e => setDueDate(e.target.value)}
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20 focus:border-[var(--color-navy)]"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">Related ID</label>
              <input
                type="text"
                value={relatedId}
                onChange={e => setRelatedId(e.target.value)}
                placeholder="e.g., STU-042"
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20 focus:border-[var(--color-navy)]"
              />
            </div>
          </div>

          {error && (
            <p className="text-sm text-red-600">{error}</p>
          )}

          <div className="flex justify-end gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm text-slate-600 hover:text-slate-800"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="px-4 py-2 text-sm font-medium bg-[var(--color-navy)] text-white rounded-lg hover:bg-[var(--color-navy-light)] disabled:opacity-50"
            >
              {submitting ? 'Creating...' : 'Create Task'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
