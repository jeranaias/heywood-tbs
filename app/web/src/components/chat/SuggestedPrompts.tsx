import { useState, useEffect } from 'react'
import { Sparkles } from 'lucide-react'
import { api } from '../../lib/api'

interface SuggestedPromptsProps {
  onSelect: (prompt: string) => void
}

export function SuggestedPrompts({ onSelect }: SuggestedPromptsProps) {
  const [prompts, setPrompts] = useState<string[]>([])

  useEffect(() => {
    api.getSuggestedPrompts()
      .then(res => setPrompts(res.prompts))
      .catch(() => {})
  }, [])

  if (prompts.length === 0) return null

  return (
    <div className="flex flex-wrap gap-2 mt-4">
      {prompts.map((prompt) => (
        <button
          key={prompt}
          onClick={() => onSelect(prompt)}
          className="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm text-slate-600 bg-slate-50 border border-slate-200 rounded-full hover:bg-slate-100 hover:text-slate-800 hover:border-slate-300 transition-colors"
        >
          <Sparkles className="w-3.5 h-3.5 text-amber-500" />
          {prompt}
        </button>
      ))}
    </div>
  )
}
