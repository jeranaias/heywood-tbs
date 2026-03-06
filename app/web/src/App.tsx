import { Routes, Route, Navigate } from 'react-router-dom'
import { AuthContext } from './hooks/useAuth'
import { useAuthProvider } from './hooks/useAuth'
import { ChatProvider } from './hooks/ChatContext'
import { AppShell } from './components/layout/AppShell'
import { Dashboard } from './pages/Dashboard'
import { StudentsPage } from './pages/StudentsPage'
import { StudentDetailPage } from './pages/StudentDetailPage'
import { AtRisk } from './pages/AtRisk'
import { ChatPage } from './pages/ChatPage'
import { InstructorQuals } from './pages/InstructorQuals'
import { Schedule } from './pages/Schedule'
import { MyRecord } from './pages/MyRecord'

function App() {
  const authProvider = useAuthProvider()

  return (
    <AuthContext.Provider value={authProvider}>
      <ChatProvider>
        <Routes>
          <Route element={<AppShell />}>
            <Route path="/" element={<Dashboard />} />
            <Route path="/chat" element={<ChatPage />} />
            <Route path="/students" element={<StudentsPage />} />
            <Route path="/students/:id" element={<StudentDetailPage />} />
            <Route path="/at-risk" element={<AtRisk />} />
            <Route path="/instructor-quals" element={<InstructorQuals />} />
            <Route path="/schedule" element={<Schedule />} />
            <Route path="/my-record" element={<MyRecord />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Route>
        </Routes>
      </ChatProvider>
    </AuthContext.Provider>
  )
}

export default App
