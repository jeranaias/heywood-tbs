import { Download } from 'lucide-react'

interface ExportButtonProps {
  url: string
  filename: string
  label?: string
}

export function ExportButton({ url, filename, label = 'Export CSV' }: ExportButtonProps) {
  async function handleExport() {
    try {
      const res = await fetch(url, { credentials: 'include' })
      if (!res.ok) throw new Error(`Export failed: ${res.status}`)
      const blob = await res.blob()
      const a = document.createElement('a')
      a.href = URL.createObjectURL(blob)
      a.download = filename
      a.click()
      URL.revokeObjectURL(a.href)
    } catch {
      // ignore
    }
  }

  return (
    <button
      onClick={handleExport}
      className="flex items-center gap-1 px-3 py-1.5 text-xs font-medium text-slate-600 border border-slate-200 rounded-lg hover:bg-slate-50"
    >
      <Download className="w-3.5 h-3.5" />
      {label}
    </button>
  )
}
