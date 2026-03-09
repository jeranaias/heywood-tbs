import { useState } from 'react'
import { X, Send } from 'lucide-react'
import { api } from '../../lib/api'

interface Props {
  replyTo?: { id: string; subject: string; from: string } | null
  onClose: () => void
  onSent: () => void
}

export function ComposeMailModal({ replyTo, onClose, onSent }: Props) {
  const [to, setTo] = useState(replyTo ? replyTo.from : '')
  const [subject, setSubject] = useState(replyTo ? `RE: ${replyTo.subject}` : '')
  const [body, setBody] = useState('')
  const [sending, setSending] = useState(false)
  const [error, setError] = useState('')

  async function handleSend() {
    if (!replyTo && (!to.trim() || !subject.trim())) {
      setError('To and subject are required')
      return
    }
    if (!body.trim()) {
      setError('Message body is required')
      return
    }

    setSending(true)
    setError('')
    try {
      if (replyTo) {
        await api.replyToMail(replyTo.id, body)
      } else {
        const recipients = to.split(',').map(s => s.trim()).filter(Boolean)
        await api.sendMail(recipients, subject, body)
      }
      onSent()
    } catch {
      setError('Failed to send. Check connection and try again.')
    } finally {
      setSending(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/30 flex items-center justify-center z-50" onClick={onClose}>
      <div className="bg-white rounded-xl border border-slate-200 shadow-lg w-full max-w-lg mx-4" onClick={e => e.stopPropagation()}>
        <div className="flex items-center justify-between px-4 py-3 border-b border-slate-200">
          <h3 className="text-sm font-semibold text-slate-800">{replyTo ? 'Reply' : 'Compose Email'}</h3>
          <button onClick={onClose} className="p-1 hover:bg-slate-100 rounded">
            <X className="w-4 h-4 text-slate-500" />
          </button>
        </div>

        <div className="p-4 space-y-3">
          {!replyTo && (
            <>
              <div>
                <label className="block text-xs font-medium text-slate-600 mb-1">To</label>
                <input
                  type="text"
                  value={to}
                  onChange={e => setTo(e.target.value)}
                  placeholder="email@example.com (comma-separated)"
                  className="w-full px-3 py-2 text-sm border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-200"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-slate-600 mb-1">Subject</label>
                <input
                  type="text"
                  value={subject}
                  onChange={e => setSubject(e.target.value)}
                  className="w-full px-3 py-2 text-sm border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-200"
                />
              </div>
            </>
          )}

          <div>
            <label className="block text-xs font-medium text-slate-600 mb-1">
              {replyTo ? `Reply to ${replyTo.from}` : 'Message'}
            </label>
            <textarea
              value={body}
              onChange={e => setBody(e.target.value)}
              rows={6}
              className="w-full px-3 py-2 text-sm border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-200 resize-none"
              placeholder="Type your message..."
            />
          </div>

          {error && <p className="text-xs text-red-600">{error}</p>}
        </div>

        <div className="flex justify-end gap-2 px-4 py-3 border-t border-slate-200">
          <button onClick={onClose} className="px-3 py-1.5 text-sm text-slate-600 hover:bg-slate-100 rounded-lg">
            Cancel
          </button>
          <button
            onClick={handleSend}
            disabled={sending}
            className="flex items-center gap-1 px-4 py-1.5 text-sm font-medium bg-[var(--color-navy)] text-white rounded-lg hover:bg-[var(--color-navy-light)] disabled:opacity-50"
          >
            <Send className="w-3.5 h-3.5" />
            {sending ? 'Sending...' : 'Send'}
          </button>
        </div>
      </div>
    </div>
  )
}
