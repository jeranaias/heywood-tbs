import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function scoreColor(score: number): string {
  if (score >= 85) return 'text-green-700 bg-green-50'
  if (score >= 75) return 'text-yellow-700 bg-yellow-50'
  return 'text-red-700 bg-red-50'
}

export function scoreBadge(score: number): string {
  if (score >= 85) return 'bg-green-100 text-green-800'
  if (score >= 75) return 'bg-yellow-100 text-yellow-800'
  return 'bg-red-100 text-red-800'
}

export function scoreLabel(score: number): string {
  if (score >= 85) return 'Strong'
  if (score >= 75) return 'Satisfactory'
  if (score > 0) return 'Below Standard'
  return 'N/A'
}

export function trendIcon(trend: string): string {
  switch (trend) {
    case 'up': return '↑'
    case 'down': return '↓'
    default: return '→'
  }
}

export function trendColor(trend: string): string {
  switch (trend) {
    case 'up': return 'text-green-600'
    case 'down': return 'text-red-600'
    default: return 'text-gray-500'
  }
}

export function formatScore(score: number): string {
  if (score === 0) return '—'
  return score.toFixed(1)
}
