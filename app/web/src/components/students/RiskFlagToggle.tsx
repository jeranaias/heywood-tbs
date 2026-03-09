import { useState } from 'react'
import { AlertTriangle, X, Plus } from 'lucide-react'
import { api } from '../../lib/api'
import type { Student } from '../../lib/types'

interface RiskFlagToggleProps {
  student: Student
  onUpdate: (student: Student) => void
}

export function RiskFlagToggle({ student, onUpdate }: RiskFlagToggleProps) {
  const [newFlag, setNewFlag] = useState('')
  const [saving, setSaving] = useState(false)

  async function toggleAtRisk() {
    setSaving(true)
    try {
      const updated = await api.updateStudent(student.id, { atRisk: !student.atRisk })
      onUpdate(updated)
    } catch {
      // ignore
    } finally {
      setSaving(false)
    }
  }

  async function removeFlag(flag: string) {
    const flags = student.riskFlags.filter(f => f !== flag)
    setSaving(true)
    try {
      const updated = await api.updateStudent(student.id, { riskFlags: flags })
      onUpdate(updated)
    } catch {
      // ignore
    } finally {
      setSaving(false)
    }
  }

  async function addFlag() {
    if (!newFlag.trim()) return
    const flags = [...student.riskFlags, newFlag.trim()]
    setSaving(true)
    try {
      const updated = await api.updateStudent(student.id, { riskFlags: flags })
      onUpdate(updated)
      setNewFlag('')
    } catch {
      // ignore
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className={`rounded-lg border p-4 ${student.atRisk ? 'bg-red-50 border-red-200' : 'bg-white border-slate-200'}`}>
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-semibold text-slate-700 flex items-center gap-2">
          <AlertTriangle className={`w-4 h-4 ${student.atRisk ? 'text-red-500' : 'text-slate-400'}`} />
          Risk Status
        </h3>
        <button
          onClick={toggleAtRisk}
          disabled={saving}
          className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
            student.atRisk ? 'bg-red-500' : 'bg-slate-300'
          }`}
        >
          <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
            student.atRisk ? 'translate-x-6' : 'translate-x-1'
          }`} />
        </button>
      </div>

      {/* Risk flags */}
      <div className="flex flex-wrap gap-1.5 mb-3">
        {student.riskFlags.map(flag => (
          <span key={flag} className="inline-flex items-center gap-1 px-2 py-0.5 bg-red-100 text-red-800 text-xs rounded-full font-medium">
            {flag}
            <button onClick={() => removeFlag(flag)} className="hover:text-red-600" disabled={saving}>
              <X className="w-3 h-3" />
            </button>
          </span>
        ))}
      </div>

      {/* Add new flag */}
      <div className="flex items-center gap-2">
        <input
          type="text"
          value={newFlag}
          onChange={e => setNewFlag(e.target.value)}
          onKeyDown={e => e.key === 'Enter' && addFlag()}
          placeholder="Add flag..."
          className="flex-1 px-2 py-1 border border-slate-200 rounded text-xs focus:outline-none focus:ring-2 focus:ring-[var(--color-navy)]/20"
        />
        <button
          onClick={addFlag}
          disabled={!newFlag.trim() || saving}
          className="p-1 text-slate-400 hover:text-slate-600 disabled:opacity-50"
        >
          <Plus className="w-4 h-4" />
        </button>
      </div>
    </div>
  )
}
