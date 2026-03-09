import { useState, useEffect } from 'react'
import { MessageSquare, Send } from 'lucide-react'
import { api } from '../../lib/api'

interface StudentNote {
  id: string
  studentId: string
  authorRole: string
  authorName: string
  content: string
  type: string
  createdAt: string
}

interface StudentNotesPanelProps {
  studentId: string
}

const typeLabels: Record<string, string> = {
  note: 'Note',
  observation: 'Observation',
  'counseling-note': 'Counseling Note',
}

export function StudentNotesPanel({ studentId }: StudentNotesPanelProps) {
  const [notes, setNotes] = useState<StudentNote[]>([])
  const [content, setContent] = useState('')
  const [noteType, setNoteType] = useState('note')
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    api.getStudentNotes(studentId)
      .then(res => setNotes(res.notes || []))
      .catch(() => {})
  }, [studentId])

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!content.trim()) return
    setSubmitting(true)
    try {
      const res = await api.createStudentNote(studentId, content.trim(), noteType)
      setNotes(res.notes || [])
      setContent('')
    } catch {
      // ignore
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="bg-white rounded-lg border border-slate-200 p-4">
      <h3 className="text-sm font-semibold text-slate-700 mb-3 flex items-center gap-2">
        <MessageSquare className="w-4 h-4" /> Staff Notes
      </h3>

      {/* Note list */}
      {notes.length > 0 ? (
        <div className="space-y-3 mb-4 max-h-64 overflow-y-auto">
          {notes.map(note => (
            <div key={note.id} className="border-l-2 border-slate-200 pl-3 py-1">
              <div className="flex items-center gap-2 text-xs text-slate-500">
                <span className="font-medium text-slate-600">{note.authorName || note.authorRole}</span>
                <span className="px-1.5 py-0.5 bg-slate-100 rounded text-xs">{typeLabels[note.type] || note.type}</span>
                <span>{new Date(note.createdAt).toLocaleDateString()}</span>
              </div>
              <p className="text-sm text-slate-700 mt-1">{note.content}</p>
            </div>
          ))}
        </div>
      ) : (
        <p className="text-sm text-slate-400 mb-4">No notes yet</p>
      )}

      {/* Add note form */}
      <form onSubmit={handleSubmit} className="space-y-2">
        <textarea
          value={content}
          onChange={e => setContent(e.target.value)}
          placeholder="Add a note..."
          rows={2}
          className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20 focus:border-[var(--color-navy)] resize-none"
        />
        <div className="flex items-center gap-2">
          <select
            value={noteType}
            onChange={e => setNoteType(e.target.value)}
            className="px-2 py-1.5 border border-slate-200 rounded text-xs focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20"
          >
            <option value="note">Note</option>
            <option value="observation">Observation</option>
            <option value="counseling-note">Counseling Note</option>
          </select>
          <button
            type="submit"
            disabled={!content.trim() || submitting}
            className="flex items-center gap-1 px-3 py-1.5 text-xs font-medium bg-[var(--color-navy)] text-white rounded-lg hover:bg-[var(--color-navy-light)] disabled:opacity-50"
          >
            <Send className="w-3 h-3" />
            {submitting ? 'Adding...' : 'Add Note'}
          </button>
        </div>
      </form>
    </div>
  )
}
