import { createContext, useContext, useState, useCallback, useEffect } from 'react'
import type { AuthInfo, Role } from '../lib/types'
import { api } from '../lib/api'

interface AuthContextType {
  auth: AuthInfo
  switchRole: (role: Role, studentId?: string) => Promise<void>
  loading: boolean
}

export const AuthContext = createContext<AuthContextType>({
  auth: { role: 'staff', company: '', name: 'TBS Staff' },
  switchRole: async () => {},
  loading: false,
})

export function useAuth() {
  return useContext(AuthContext)
}

export function useAuthProvider() {
  const [auth, setAuth] = useState<AuthInfo>({
    role: 'staff',
    company: '',
    name: 'TBS Staff',
  })
  const [loading, setLoading] = useState(true)

  // Hydrate auth state from server cookie on mount
  useEffect(() => {
    api.getAuthMe()
      .then(info => setAuth(info))
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  const switchRole = useCallback(async (role: Role, studentId?: string) => {
    setLoading(true)
    try {
      const company = (role === 'spc') ? 'Alpha' : ''
      const sid = role === 'student' ? (studentId || 'STU-001') : ''
      const info = await api.switchRole(role, company, sid)
      setAuth(info)
    } catch (err) {
      console.error('Failed to switch role:', err)
    } finally {
      setLoading(false)
    }
  }, [])

  return { auth, switchRole, loading }
}
