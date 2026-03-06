import { NavLink } from 'react-router-dom'
import {
  MessageSquare,
  LayoutDashboard,
  Users,
  AlertTriangle,
  Shield,
  Calendar,
  User,
} from 'lucide-react'

interface SidebarProps {
  role: string
  onClose: () => void
}

export function Sidebar({ role, onClose }: SidebarProps) {
  const links = [
    { to: '/chat', label: 'Chat', icon: MessageSquare, roles: ['xo', 'staff', 'spc', 'student'] },
    { to: '/', label: 'Dashboard', icon: LayoutDashboard, roles: ['xo', 'staff', 'spc'] },
    { to: '/students', label: 'Students', icon: Users, roles: ['xo', 'staff', 'spc'] },
    { to: '/at-risk', label: 'At-Risk', icon: AlertTriangle, roles: ['xo', 'staff', 'spc'] },
    { to: '/instructor-quals', label: 'Instructor Quals', icon: Shield, roles: ['xo', 'staff'] },
    { to: '/schedule', label: 'Schedule', icon: Calendar, roles: ['xo', 'staff', 'spc'] },
    { to: '/my-record', label: 'My Record', icon: User, roles: ['student'] },
  ]

  return (
    <div className="flex h-full flex-col bg-[var(--color-sidebar)] text-white">
      {/* Logo/Brand */}
      <div className="flex items-center gap-3 px-6 py-5 border-b border-slate-700">
        <div className="w-8 h-8 rounded-lg bg-[var(--color-scarlet)] flex items-center justify-center font-bold text-sm">
          H
        </div>
        <div>
          <h1 className="text-lg font-bold tracking-tight">Heywood</h1>
          <p className="text-xs text-slate-400">TBS Digital Staff Officer</p>
        </div>
        <button
          onClick={onClose}
          className="ml-auto text-slate-400 hover:text-white lg:hidden"
        >
          ×
        </button>
      </div>

      {/* Navigation */}
      <nav className="flex-1 px-3 py-4 space-y-1">
        {links
          .filter(link => link.roles.includes(role))
          .map(link => (
            <NavLink
              key={link.to}
              to={link.to}
              onClick={onClose}
              className={({ isActive }) =>
                `flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-colors ${
                  isActive
                    ? 'bg-[var(--color-sidebar-active)] text-white font-medium'
                    : 'text-slate-300 hover:bg-[var(--color-sidebar-hover)] hover:text-white'
                }`
              }
            >
              <link.icon className="w-5 h-5 flex-shrink-0" />
              {link.label}
            </NavLink>
          ))}
      </nav>

      {/* Role indicator */}
      <div className="px-4 py-3 border-t border-slate-700">
        <div className="text-xs text-slate-500 uppercase tracking-wider">Current Role</div>
        <div className="text-sm text-slate-300 mt-1 capitalize">{role}</div>
      </div>
    </div>
  )
}
