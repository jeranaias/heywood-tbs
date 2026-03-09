import { lazy, Suspense } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { ErrorBoundary } from './components/ErrorBoundary'
import { AuthContext } from './hooks/useAuth'
import { useAuthProvider } from './hooks/useAuth'
import { ChatProvider } from './hooks/ChatContext'
import { AppShell } from './components/layout/AppShell'

// Eagerly loaded — always needed or tiny
import { Dashboard } from './pages/Dashboard'
import { StudentsPage } from './pages/StudentsPage'
import { AtRisk } from './pages/AtRisk'
import { MyRecord } from './pages/MyRecord'

// Lazy loaded — visited less frequently
const ChatPage = lazy(() => import('./pages/ChatPage').then(m => ({ default: m.ChatPage })))
const StudentDetailPage = lazy(() => import('./pages/StudentDetailPage').then(m => ({ default: m.StudentDetailPage })))
const InstructorQuals = lazy(() => import('./pages/InstructorQuals').then(m => ({ default: m.InstructorQuals })))
const Schedule = lazy(() => import('./pages/Schedule').then(m => ({ default: m.Schedule })))
const TasksPage = lazy(() => import('./pages/TasksPage').then(m => ({ default: m.TasksPage })))
const CalendarPage = lazy(() => import('./pages/CalendarPage').then(m => ({ default: m.CalendarPage })))
const CounselingPage = lazy(() => import('./pages/CounselingPage').then(m => ({ default: m.CounselingPage })))
const ReportsPage = lazy(() => import('./pages/ReportsPage').then(m => ({ default: m.ReportsPage })))
const SettingsPage = lazy(() => import('./pages/SettingsPage').then(m => ({ default: m.SettingsPage })))

function App() {
  const authProvider = useAuthProvider()

  return (
    <ErrorBoundary>
      <AuthContext.Provider value={authProvider}>
        <ChatProvider>
          <Suspense fallback={<div className="flex items-center justify-center h-64 text-zinc-400">Loading...</div>}>
            <Routes>
              <Route element={<AppShell />}>
                <Route path="/" element={<Dashboard />} />
                <Route path="/chat" element={<ChatPage />} />
                <Route path="/students" element={<StudentsPage />} />
                <Route path="/students/:id" element={<StudentDetailPage />} />
                <Route path="/at-risk" element={<AtRisk />} />
                <Route path="/instructor-quals" element={<InstructorQuals />} />
                <Route path="/schedule" element={<Schedule />} />
                <Route path="/tasks" element={<TasksPage />} />
                <Route path="/counseling" element={<CounselingPage />} />
                <Route path="/reports" element={<ReportsPage />} />
                <Route path="/my-record" element={<MyRecord />} />
                <Route path="/calendar" element={<CalendarPage />} />
                <Route path="/settings" element={<SettingsPage />} />
                <Route path="*" element={<Navigate to="/" replace />} />
              </Route>
            </Routes>
          </Suspense>
        </ChatProvider>
      </AuthContext.Provider>
    </ErrorBoundary>
  )
}

export default App
