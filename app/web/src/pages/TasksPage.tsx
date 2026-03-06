import { useState, useEffect, useCallback } from 'react'
import { ClipboardList, Mail, Clock, CheckCircle2, AlertCircle, ChevronRight } from 'lucide-react'
import { api } from '../lib/api'
import type { Task, Message } from '../lib/types'

type TabType = 'tasks' | 'messages'
type FilterType = 'all' | 'pending' | 'in_progress' | 'completed'

export function TasksPage() {
  const [tab, setTab] = useState<TabType>('tasks')
  const [tasks, setTasks] = useState<Task[]>([])
  const [messages, setMessages] = useState<Message[]>([])
  const [filter, setFilter] = useState<FilterType>('all')
  const [loading, setLoading] = useState(true)
  const [selectedTask, setSelectedTask] = useState<Task | null>(null)

  const loadData = useCallback(async () => {
    try {
      setLoading(true)
      const [t, m] = await Promise.all([
        api.getTasks({ all: 'true' }),
        api.getMessages({ all: 'true' }),
      ])
      setTasks(t || [])
      setMessages(m || [])
    } catch (err) {
      console.error('Failed to load tasks/messages:', err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { loadData() }, [loadData])

  const filteredTasks = filter === 'all'
    ? tasks
    : tasks.filter(t => t.status === filter)

  const handleStatusUpdate = async (taskId: string, status: string) => {
    try {
      await api.updateTask(taskId, { status: status as Task['status'] })
      await loadData()
      if (selectedTask?.id === taskId) {
        setSelectedTask(prev => prev ? { ...prev, status: status as Task['status'] } : null)
      }
    } catch (err) {
      console.error('Failed to update task:', err)
    }
  }

  const handleMarkRead = async (msgId: string) => {
    try {
      await api.markMessageRead(msgId)
      setMessages(prev => prev.map(m => m.id === msgId ? { ...m, read: true } : m))
    } catch (err) {
      console.error('Failed to mark message read:', err)
    }
  }

  const priorityColor = (p: string) => {
    switch (p) {
      case 'high': return 'text-red-600 bg-red-50 border-red-200'
      case 'medium': return 'text-amber-600 bg-amber-50 border-amber-200'
      case 'low': return 'text-green-600 bg-green-50 border-green-200'
      default: return 'text-slate-600 bg-slate-50 border-slate-200'
    }
  }

  const statusIcon = (s: string) => {
    switch (s) {
      case 'pending': return <Clock className="w-4 h-4 text-amber-500" />
      case 'in_progress': return <AlertCircle className="w-4 h-4 text-blue-500" />
      case 'completed': return <CheckCircle2 className="w-4 h-4 text-green-500" />
      default: return <Clock className="w-4 h-4 text-slate-400" />
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-[var(--color-navy)]" />
      </div>
    )
  }

  return (
    <div className="max-w-5xl mx-auto space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-slate-900">Task Inbox</h1>
        <p className="text-sm text-slate-500 mt-1">Tasks and messages from Heywood</p>
      </div>

      {/* Tabs */}
      <div className="flex gap-1 border-b border-slate-200">
        <button
          onClick={() => setTab('tasks')}
          className={`flex items-center gap-2 px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
            tab === 'tasks'
              ? 'border-[var(--color-navy)] text-[var(--color-navy)]'
              : 'border-transparent text-slate-500 hover:text-slate-700'
          }`}
        >
          <ClipboardList className="w-4 h-4" />
          Tasks
          {tasks.filter(t => t.status !== 'completed').length > 0 && (
            <span className="px-1.5 py-0.5 text-xs bg-[var(--color-navy)] text-white rounded-full">
              {tasks.filter(t => t.status !== 'completed').length}
            </span>
          )}
        </button>
        <button
          onClick={() => setTab('messages')}
          className={`flex items-center gap-2 px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
            tab === 'messages'
              ? 'border-[var(--color-navy)] text-[var(--color-navy)]'
              : 'border-transparent text-slate-500 hover:text-slate-700'
          }`}
        >
          <Mail className="w-4 h-4" />
          Messages
          {messages.filter(m => !m.read).length > 0 && (
            <span className="px-1.5 py-0.5 text-xs bg-[var(--color-scarlet)] text-white rounded-full">
              {messages.filter(m => !m.read).length}
            </span>
          )}
        </button>
      </div>

      {tab === 'tasks' && (
        <div className="space-y-4">
          {/* Filter pills */}
          <div className="flex gap-2">
            {(['all', 'pending', 'in_progress', 'completed'] as const).map(f => (
              <button
                key={f}
                onClick={() => setFilter(f)}
                className={`px-3 py-1.5 text-xs font-medium rounded-full transition-colors ${
                  filter === f
                    ? 'bg-[var(--color-navy)] text-white'
                    : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                }`}
              >
                {f === 'all' ? 'All' : f === 'in_progress' ? 'In Progress' : f.charAt(0).toUpperCase() + f.slice(1)}
              </button>
            ))}
          </div>

          {/* Task list */}
          {filteredTasks.length === 0 ? (
            <div className="text-center py-12 text-slate-500">
              <ClipboardList className="w-12 h-12 mx-auto mb-3 text-slate-300" />
              <p>No tasks yet. Heywood will create tasks when you give orders.</p>
            </div>
          ) : (
            <div className="space-y-2">
              {filteredTasks.map(task => (
                <div
                  key={task.id}
                  onClick={() => setSelectedTask(selectedTask?.id === task.id ? null : task)}
                  className={`bg-white rounded-lg border p-4 cursor-pointer transition-all hover:shadow-md ${
                    selectedTask?.id === task.id ? 'ring-2 ring-[var(--color-navy)] border-transparent' : 'border-slate-200'
                  }`}
                >
                  <div className="flex items-start gap-3">
                    {statusIcon(task.status)}
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 flex-wrap">
                        <h3 className="text-sm font-semibold text-slate-900">{task.title}</h3>
                        <span className={`px-2 py-0.5 text-xs font-medium rounded-full border ${priorityColor(task.priority)}`}>
                          {task.priority}
                        </span>
                        <span className="text-xs text-slate-400">{task.id}</span>
                      </div>
                      <p className="text-xs text-slate-500 mt-1">
                        Assigned to: {task.assignedTo} | Created by: {task.createdBy}
                        {task.dueDate && ` | Due: ${task.dueDate}`}
                      </p>
                    </div>
                    <ChevronRight className={`w-4 h-4 text-slate-400 transition-transform ${selectedTask?.id === task.id ? 'rotate-90' : ''}`} />
                  </div>

                  {selectedTask?.id === task.id && (
                    <div className="mt-3 pt-3 border-t border-slate-100">
                      {task.description && (
                        <p className="text-sm text-slate-700 mb-3">{task.description}</p>
                      )}
                      {task.relatedId && (
                        <p className="text-xs text-slate-500 mb-3">Related: {task.relatedId}</p>
                      )}
                      <div className="flex gap-2">
                        {task.status !== 'in_progress' && task.status !== 'completed' && (
                          <button
                            onClick={(e) => { e.stopPropagation(); handleStatusUpdate(task.id, 'in_progress') }}
                            className="px-3 py-1.5 text-xs font-medium bg-blue-50 text-blue-700 rounded-lg hover:bg-blue-100"
                          >
                            Start Work
                          </button>
                        )}
                        {task.status !== 'completed' && (
                          <button
                            onClick={(e) => { e.stopPropagation(); handleStatusUpdate(task.id, 'completed') }}
                            className="px-3 py-1.5 text-xs font-medium bg-green-50 text-green-700 rounded-lg hover:bg-green-100"
                          >
                            Mark Complete
                          </button>
                        )}
                        {task.status === 'completed' && (
                          <button
                            onClick={(e) => { e.stopPropagation(); handleStatusUpdate(task.id, 'pending') }}
                            className="px-3 py-1.5 text-xs font-medium bg-slate-50 text-slate-700 rounded-lg hover:bg-slate-100"
                          >
                            Reopen
                          </button>
                        )}
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {tab === 'messages' && (
        <div className="space-y-2">
          {messages.length === 0 ? (
            <div className="text-center py-12 text-slate-500">
              <Mail className="w-12 h-12 mx-auto mb-3 text-slate-300" />
              <p>No messages yet.</p>
            </div>
          ) : (
            messages.map(msg => (
              <div
                key={msg.id}
                onClick={() => !msg.read && handleMarkRead(msg.id)}
                className={`bg-white rounded-lg border p-4 cursor-pointer transition-all hover:shadow-md ${
                  msg.read ? 'border-slate-200' : 'border-blue-200 bg-blue-50/30'
                }`}
              >
                <div className="flex items-start gap-3">
                  <Mail className={`w-4 h-4 mt-0.5 ${msg.read ? 'text-slate-400' : 'text-blue-500'}`} />
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <h3 className={`text-sm ${msg.read ? 'text-slate-700' : 'font-semibold text-slate-900'}`}>
                        {msg.subject}
                      </h3>
                      {!msg.read && (
                        <span className="w-2 h-2 rounded-full bg-blue-500" />
                      )}
                    </div>
                    <p className="text-xs text-slate-500 mt-0.5">From: {msg.from}</p>
                    <p className="text-sm text-slate-700 mt-2">{msg.body}</p>
                    <p className="text-xs text-slate-400 mt-2">
                      {new Date(msg.createdAt).toLocaleString()}
                    </p>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      )}
    </div>
  )
}
