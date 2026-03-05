import type { LucideIcon } from 'lucide-react'
import { cn } from '../../lib/utils'

interface KPICardProps {
  title: string
  value: string | number
  subtitle?: string
  icon: LucideIcon
  variant?: 'default' | 'danger' | 'warning' | 'success'
}

const variants = {
  default: 'bg-white border-slate-200',
  danger: 'bg-red-50 border-red-200',
  warning: 'bg-yellow-50 border-yellow-200',
  success: 'bg-green-50 border-green-200',
}

const iconVariants = {
  default: 'text-[var(--color-navy)] bg-blue-50',
  danger: 'text-red-600 bg-red-100',
  warning: 'text-yellow-600 bg-yellow-100',
  success: 'text-green-600 bg-green-100',
}

export function KPICard({ title, value, subtitle, icon: Icon, variant = 'default' }: KPICardProps) {
  return (
    <div className={cn('rounded-lg border p-4', variants[variant])}>
      <div className="flex items-start justify-between">
        <div>
          <p className="text-xs font-medium text-slate-500 uppercase tracking-wider">{title}</p>
          <p className="text-2xl font-bold text-slate-900 mt-1">{value}</p>
          {subtitle && <p className="text-xs text-slate-500 mt-1">{subtitle}</p>}
        </div>
        <div className={cn('p-2 rounded-lg', iconVariants[variant])}>
          <Icon className="w-5 h-5" />
        </div>
      </div>
    </div>
  )
}
