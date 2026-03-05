import { useState, useRef, useEffect } from 'react'
import { Menu, ChevronDown, User, Shield, Users } from 'lucide-react'
import { useAuth } from '../../hooks/useAuth'
import type { Role } from '../../lib/types'

interface HeaderProps {
  onMenuClick: () => void
}

const roles: { value: Role; label: string; desc: string; icon: typeof User }[] = [
  { value: 'staff', label: 'Staff', desc: 'Full TBS-wide access', icon: Shield },
  { value: 'spc', label: 'SPC (Alpha Co)', desc: 'Company-filtered view', icon: Users },
  { value: 'student', label: 'Student', desc: 'Individual record only', icon: User },
]

export function Header({ onMenuClick }: HeaderProps) {
  const { auth, switchRole } = useAuth()
  const [dropdownOpen, setDropdownOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setDropdownOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [])

  const currentRole = roles.find(r => r.value === auth.role) || roles[0]

  return (
    <header className="flex items-center gap-4 border-b border-slate-200 bg-white px-4 py-3 lg:px-6">
      <button
        onClick={onMenuClick}
        className="lg:hidden text-slate-500 hover:text-slate-700"
      >
        <Menu className="w-6 h-6" />
      </button>

      <div className="flex-1" />

      {/* Role Switcher */}
      <div className="relative" ref={dropdownRef}>
        <button
          onClick={() => setDropdownOpen(!dropdownOpen)}
          className="flex items-center gap-2 px-3 py-2 rounded-lg border border-slate-200 hover:bg-slate-50 transition-colors"
        >
          <currentRole.icon className="w-4 h-4 text-slate-500" />
          <span className="text-sm font-medium text-slate-700">{auth.name}</span>
          <ChevronDown className="w-4 h-4 text-slate-400" />
        </button>

        {dropdownOpen && (
          <div className="absolute right-0 mt-2 w-64 bg-white rounded-lg border border-slate-200 shadow-lg py-1 z-50">
            <div className="px-3 py-2 text-xs text-slate-500 uppercase tracking-wider border-b border-slate-100">
              Switch Role (Demo)
            </div>
            {roles.map(role => (
              <button
                key={role.value}
                onClick={() => {
                  switchRole(role.value)
                  setDropdownOpen(false)
                }}
                className={`w-full flex items-center gap-3 px-3 py-2.5 text-left hover:bg-slate-50 transition-colors ${
                  auth.role === role.value ? 'bg-slate-50' : ''
                }`}
              >
                <role.icon className={`w-5 h-5 ${auth.role === role.value ? 'text-[var(--color-navy)]' : 'text-slate-400'}`} />
                <div>
                  <div className={`text-sm ${auth.role === role.value ? 'font-semibold text-[var(--color-navy)]' : 'text-slate-700'}`}>
                    {role.label}
                  </div>
                  <div className="text-xs text-slate-500">{role.desc}</div>
                </div>
                {auth.role === role.value && (
                  <div className="ml-auto w-2 h-2 rounded-full bg-[var(--color-navy)]" />
                )}
              </button>
            ))}
          </div>
        )}
      </div>
    </header>
  )
}
