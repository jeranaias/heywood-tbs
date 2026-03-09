import { useState, useEffect, useCallback } from 'react'
import { Calendar as CalendarIcon, ChevronLeft, ChevronRight, Mail, Clock, MapPin, Send, Check, X as XIcon, HelpCircle } from 'lucide-react'
import { api } from '../lib/api'
import type { CalendarEvent, MailSummary } from '../lib/types'
import { useAuth } from '../hooks/useAuth'
import { ComposeMailModal } from '../components/mail/ComposeMailModal'

type ViewMode = 'week' | 'day' | 'agenda'

const SOURCE_COLORS: Record<string, string> = {
  'tbs-schedule': 'bg-blue-100 text-blue-800 border-blue-200',
  'outlook': 'bg-green-100 text-green-800 border-green-200',
  'shared': 'bg-orange-100 text-orange-800 border-orange-200',
}

const CATEGORY_DOTS: Record<string, string> = {
  training: 'bg-blue-500',
  admin: 'bg-slate-400',
  personal: 'bg-green-500',
}

export function CalendarPage() {
  const { auth } = useAuth()
  const [events, setEvents] = useState<CalendarEvent[]>([])
  const [mails, setMails] = useState<MailSummary[]>([])
  const [unreadCount, setUnreadCount] = useState(0)
  const [loading, setLoading] = useState(true)
  const [view, setView] = useState<ViewMode>('agenda')
  const [currentDate, setCurrentDate] = useState(new Date())
  const [selectedEvent, setSelectedEvent] = useState<CalendarEvent | null>(null)
  const [showCompose, setShowCompose] = useState(false)
  const [replyTo, setReplyTo] = useState<{ id: string; subject: string; from: string } | null>(null)

  const formatDateParam = (d: Date) => d.toISOString().split('T')[0]

  const loadData = useCallback(async () => {
    try {
      setLoading(true)
      const start = new Date(currentDate)
      const end = new Date(currentDate)

      if (view === 'day') {
        end.setDate(end.getDate() + 1)
      } else {
        // Start of week
        start.setDate(start.getDate() - start.getDay())
        end.setDate(start.getDate() + 7)
      }

      const [calData, mailData] = await Promise.all([
        api.getCalendarEvents({ start: formatDateParam(start), end: formatDateParam(end) }),
        api.getMailSummary(),
      ])

      setEvents(calData.events || [])
      setMails(mailData.messages || [])
      setUnreadCount(mailData.unreadCount || 0)
    } catch (err) {
      console.error('Failed to load calendar:', err)
    } finally {
      setLoading(false)
    }
  }, [currentDate, view])

  useEffect(() => { loadData() }, [loadData])

  const goToday = () => setCurrentDate(new Date())
  const goPrev = () => {
    const d = new Date(currentDate)
    d.setDate(d.getDate() - (view === 'day' ? 1 : 7))
    setCurrentDate(d)
  }
  const goNext = () => {
    const d = new Date(currentDate)
    d.setDate(d.getDate() + (view === 'day' ? 1 : 7))
    setCurrentDate(d)
  }

  // Group events by date for agenda view
  const eventsByDate = events.reduce<Record<string, CalendarEvent[]>>((acc, event) => {
    const date = event.start.split('T')[0]
    if (!acc[date]) acc[date] = []
    acc[date].push(event)
    return acc
  }, {})

  // Sort dates
  const sortedDates = Object.keys(eventsByDate).sort()

  const formatTime = (iso: string) => {
    const timePart = iso.split('T')[1]
    if (!timePart) return ''
    const [h, m] = timePart.split(':')
    const hour = parseInt(h)
    const ampm = hour >= 12 ? 'PM' : 'AM'
    const h12 = hour === 0 ? 12 : hour > 12 ? hour - 12 : hour
    return `${h12}:${m} ${ampm}`
  }

  const formatDateLabel = (dateStr: string) => {
    const d = new Date(dateStr + 'T12:00:00')
    const today = new Date()
    const tomorrow = new Date(today)
    tomorrow.setDate(tomorrow.getDate() + 1)

    if (dateStr === formatDateParam(today)) return 'Today'
    if (dateStr === formatDateParam(tomorrow)) return 'Tomorrow'

    return d.toLocaleDateString('en-US', { weekday: 'long', month: 'short', day: 'numeric' })
  }

  const roleLabel = auth.role === 'xo' ? 'XO' : auth.role === 'spc' ? 'SPC' : auth.role === 'student' ? 'Student' : 'Staff'

  async function handleRespondEvent(eventId: string, response: 'accept' | 'decline' | 'tentativelyAccept') {
    try {
      await api.respondToEvent(eventId, response)
      loadData()
    } catch { /* ignore */ }
  }

  function handleReply(mail: MailSummary) {
    setReplyTo({ id: mail.id, subject: mail.subject, from: mail.from })
    setShowCompose(true)
  }

  return (
    <div className="max-w-5xl mx-auto space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between flex-wrap gap-3">
        <div className="flex items-center gap-3">
          <CalendarIcon className="w-6 h-6 text-slate-700" />
          <h1 className="text-2xl font-bold text-slate-900">Calendar</h1>
          <span className="text-sm text-slate-500">({roleLabel} View)</span>
        </div>

        <div className="flex items-center gap-2">
          <button onClick={goToday} className="px-3 py-1.5 text-sm border border-slate-300 rounded-lg hover:bg-slate-50">
            Today
          </button>
          <button onClick={goPrev} className="p-1.5 border border-slate-300 rounded-lg hover:bg-slate-50">
            <ChevronLeft className="w-4 h-4" />
          </button>
          <button onClick={goNext} className="p-1.5 border border-slate-300 rounded-lg hover:bg-slate-50">
            <ChevronRight className="w-4 h-4" />
          </button>
          <span className="text-sm font-medium text-slate-700 min-w-[120px] text-center">
            {currentDate.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })}
          </span>
        </div>

        <div className="flex border border-slate-300 rounded-lg overflow-hidden">
          {(['day', 'week', 'agenda'] as ViewMode[]).map(v => (
            <button
              key={v}
              onClick={() => setView(v)}
              className={`px-3 py-1.5 text-sm capitalize ${
                view === v ? 'bg-[var(--color-navy)] text-white' : 'bg-white text-slate-600 hover:bg-slate-50'
              }`}
            >
              {v}
            </button>
          ))}
        </div>
      </div>

      <div className="flex gap-4">
        {/* Main calendar area */}
        <div className="flex-1">
          {loading ? (
            <div className="bg-white rounded-xl border border-slate-200 p-12 text-center text-slate-400">
              Loading events...
            </div>
          ) : events.length === 0 ? (
            <div className="bg-white rounded-xl border border-slate-200 p-12 text-center text-slate-400">
              No events for this period
            </div>
          ) : (
            <div className="space-y-4">
              {sortedDates.map(date => (
                <div key={date} className="bg-white rounded-xl border border-slate-200 overflow-hidden">
                  <div className="px-4 py-2.5 bg-slate-50 border-b border-slate-200">
                    <h3 className="text-sm font-semibold text-slate-700">{formatDateLabel(date)}</h3>
                  </div>
                  <div className="divide-y divide-slate-100">
                    {eventsByDate[date]
                      .sort((a, b) => a.start.localeCompare(b.start))
                      .map(event => (
                        <button
                          key={event.id}
                          onClick={() => setSelectedEvent(selectedEvent?.id === event.id ? null : event)}
                          className="w-full text-left px-4 py-3 hover:bg-slate-50 transition-colors"
                        >
                          <div className="flex items-start gap-3">
                            <div className={`w-2 h-2 rounded-full mt-1.5 flex-shrink-0 ${CATEGORY_DOTS[event.category || ''] || 'bg-slate-300'}`} />
                            <div className="flex-1 min-w-0">
                              <div className="flex items-center gap-2">
                                <span className="text-sm font-medium text-slate-900 truncate">{event.title}</span>
                                <span className={`text-xs px-1.5 py-0.5 rounded border ${SOURCE_COLORS[event.source] || 'bg-slate-100 text-slate-600'}`}>
                                  {event.source === 'tbs-schedule' ? 'TBS' : event.source === 'outlook' ? 'Outlook' : event.source}
                                </span>
                              </div>
                              <div className="flex items-center gap-3 mt-0.5 text-xs text-slate-500">
                                <span className="flex items-center gap-1">
                                  <Clock className="w-3 h-3" />
                                  {formatTime(event.start)} - {formatTime(event.end)}
                                </span>
                                {event.location && (
                                  <span className="flex items-center gap-1">
                                    <MapPin className="w-3 h-3" />
                                    {event.location}
                                  </span>
                                )}
                              </div>
                              {/* Expanded detail */}
                              {selectedEvent?.id === event.id && (
                                <div className="mt-2 space-y-2">
                                  {event.description && (
                                    <p className="text-xs text-slate-600 bg-slate-50 p-2 rounded">{event.description}</p>
                                  )}
                                  {event.source === 'outlook' && (
                                    <div className="flex items-center gap-1.5">
                                      <button
                                        onClick={e => { e.stopPropagation(); handleRespondEvent(event.id, 'accept') }}
                                        className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-green-700 bg-green-50 border border-green-200 rounded hover:bg-green-100"
                                      >
                                        <Check className="w-3 h-3" /> Accept
                                      </button>
                                      <button
                                        onClick={e => { e.stopPropagation(); handleRespondEvent(event.id, 'tentativelyAccept') }}
                                        className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-amber-700 bg-amber-50 border border-amber-200 rounded hover:bg-amber-100"
                                      >
                                        <HelpCircle className="w-3 h-3" /> Tentative
                                      </button>
                                      <button
                                        onClick={e => { e.stopPropagation(); handleRespondEvent(event.id, 'decline') }}
                                        className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-red-700 bg-red-50 border border-red-200 rounded hover:bg-red-100"
                                      >
                                        <XIcon className="w-3 h-3" /> Decline
                                      </button>
                                    </div>
                                  )}
                                </div>
                              )}
                            </div>
                          </div>
                        </button>
                      ))}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Mail sidebar */}
        <div className="w-72 flex-shrink-0 hidden lg:block space-y-4">
          <div className="bg-white rounded-xl border border-slate-200 p-4">
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-2">
                <Mail className="w-4 h-4 text-slate-600" />
                <h3 className="text-sm font-semibold text-slate-700">Outlook Mail</h3>
                {unreadCount > 0 && (
                  <span className="bg-[var(--color-scarlet)] text-white text-xs font-bold px-1.5 py-0.5 rounded-full">
                    {unreadCount}
                  </span>
                )}
              </div>
              <button
                onClick={() => { setReplyTo(null); setShowCompose(true) }}
                className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-[var(--color-navy)] border border-slate-200 rounded hover:bg-slate-50"
              >
                <Send className="w-3 h-3" /> Compose
              </button>
            </div>
            {mails.length === 0 ? (
              <p className="text-xs text-slate-400">No recent messages</p>
            ) : (
              <div className="space-y-2.5">
                {mails.map(mail => (
                  <div key={mail.id} className={`text-xs ${mail.isRead ? 'opacity-60' : ''}`}>
                    <div className="flex items-center gap-1.5">
                      {!mail.isRead && <div className="w-1.5 h-1.5 bg-blue-500 rounded-full flex-shrink-0" />}
                      <span className="font-medium text-slate-700 truncate flex-1">{mail.subject}</span>
                      <button
                        onClick={() => handleReply(mail)}
                        className="flex-shrink-0 p-0.5 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded"
                        title="Reply"
                      >
                        <Send className="w-3 h-3" />
                      </button>
                    </div>
                    <div className="text-slate-500 mt-0.5">{mail.from}</div>
                    <div className="text-slate-400 mt-0.5 line-clamp-2">{mail.preview}</div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Legend */}
          <div className="bg-white rounded-xl border border-slate-200 p-4">
            <h3 className="text-sm font-semibold text-slate-700 mb-2">Legend</h3>
            <div className="space-y-1.5">
              <div className="flex items-center gap-2 text-xs">
                <div className="w-2 h-2 rounded-full bg-blue-500" />
                <span className="text-slate-600">Training</span>
              </div>
              <div className="flex items-center gap-2 text-xs">
                <div className="w-2 h-2 rounded-full bg-slate-400" />
                <span className="text-slate-600">Administrative</span>
              </div>
              <div className="flex items-center gap-2 text-xs">
                <div className="w-2 h-2 rounded-full bg-green-500" />
                <span className="text-slate-600">Personal</span>
              </div>
              <div className="mt-2 pt-2 border-t border-slate-100 space-y-1">
                <div className="flex items-center gap-2 text-xs">
                  <span className="px-1 py-0.5 rounded border bg-blue-100 text-blue-800 border-blue-200">TBS</span>
                  <span className="text-slate-600">TBS Schedule</span>
                </div>
                <div className="flex items-center gap-2 text-xs">
                  <span className="px-1 py-0.5 rounded border bg-green-100 text-green-800 border-green-200">Outlook</span>
                  <span className="text-slate-600">Outlook Calendar</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {showCompose && (
        <ComposeMailModal
          replyTo={replyTo}
          onClose={() => { setShowCompose(false); setReplyTo(null) }}
          onSent={() => { setShowCompose(false); setReplyTo(null); loadData() }}
        />
      )}
    </div>
  )
}
