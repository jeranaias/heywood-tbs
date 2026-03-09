import type { CounselingSession } from '../../lib/types'

interface PrintCounselingProps {
  session: CounselingSession
}

const typeLabels: Record<string, string> = {
  'initial': 'Initial Counseling',
  'progress': 'Progress Review',
  'event-driven': 'Event-Driven Counseling',
  'end-of-phase': 'End of Phase Counseling',
}

export function PrintCounseling({ session }: PrintCounselingProps) {
  return (
    <div className="print-counseling p-8 max-w-3xl mx-auto font-serif text-sm leading-relaxed">
      <style>{`
        @media print {
          .no-print { display: none !important; }
          .print-counseling { padding: 0; max-width: 100%; }
        }
      `}</style>

      <div className="text-center mb-6 border-b-2 border-black pb-4">
        <div className="text-xs font-bold tracking-widest">UNITED STATES MARINE CORPS</div>
        <div className="text-xs">THE BASIC SCHOOL</div>
        <div className="text-xs">MARINE CORPS BASE QUANTICO, VIRGINIA 22134</div>
        <div className="text-lg font-bold mt-3">{typeLabels[session.type] || session.type}</div>
      </div>

      <div className="grid grid-cols-2 gap-4 mb-6 text-xs">
        <div><strong>Student:</strong> {session.studentName}</div>
        <div><strong>Date:</strong> {session.date || '—'}</div>
        <div><strong>Counselor:</strong> {session.counselorName} ({session.counselorRole})</div>
        <div><strong>Status:</strong> {session.status}</div>
      </div>

      {session.outline && (
        <div className="mb-6">
          <h3 className="font-bold text-xs uppercase tracking-wide border-b border-slate-400 pb-1 mb-2">Outline</h3>
          <div className="whitespace-pre-wrap text-xs">{session.outline}</div>
        </div>
      )}

      {session.notes && (
        <div className="mb-6">
          <h3 className="font-bold text-xs uppercase tracking-wide border-b border-slate-400 pb-1 mb-2">Counselor Notes</h3>
          <div className="whitespace-pre-wrap text-xs">{session.notes}</div>
        </div>
      )}

      {session.followUps && session.followUps.length > 0 && (
        <div className="mb-6">
          <h3 className="font-bold text-xs uppercase tracking-wide border-b border-slate-400 pb-1 mb-2">Follow-Up Actions</h3>
          <table className="w-full text-xs border-collapse">
            <thead>
              <tr className="border-b border-slate-300">
                <th className="text-left py-1 pr-4">Action</th>
                <th className="text-left py-1 pr-4">Due Date</th>
                <th className="text-left py-1">Status</th>
              </tr>
            </thead>
            <tbody>
              {session.followUps.map((fu, i) => (
                <tr key={i} className="border-b border-slate-200">
                  <td className="py-1 pr-4">{fu.description}</td>
                  <td className="py-1 pr-4">{fu.dueDate || '—'}</td>
                  <td className="py-1">{fu.status}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <div className="mt-12 grid grid-cols-2 gap-12 text-xs">
        <div>
          <div className="border-b border-black mb-1 h-8" />
          <div>Student Signature / Date</div>
        </div>
        <div>
          <div className="border-b border-black mb-1 h-8" />
          <div>Counselor Signature / Date</div>
        </div>
      </div>

      <div className="mt-4 text-[10px] text-slate-400 text-center italic">
        This counseling outline was generated as a draft. All content has been reviewed and approved by the counselor.
      </div>
    </div>
  )
}
