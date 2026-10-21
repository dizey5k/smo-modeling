import React from 'react'
import type { EventItem } from '../types/smo'

export const EventCalendar: React.FC<{ events: EventItem[] }> = ({
  events,
}) => {
  return (
    <div className="rounded-2xl border border-slate-800/90 bg-[#121820] p-5 shadow-sm">
      <h3 className="mb-3 text-sm font-semibold text-[#E5F9E0]">
        Календарь событий (Особые состояния)
      </h3>
      <div className="max-h-56 overflow-y-auto pr-1">
        <table className="w-full text-left font-mono text-xs">
          <thead>
            <tr className="border-b border-slate-800 text-slate-400">
              <th className="pb-2.5">Время T</th>
              <th className="pb-2.5">Тип</th>
              <th className="pb-2.5">Субъект</th>
              <th className="pb-2.5">Задача</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800/60">
            {events.map((ev, i) => (
              <tr key={i} className="hover:bg-slate-800/30 transition-colors">
                <td className="py-2 text-[#40C9A2] font-semibold">
                  {ev.time.toFixed(4)}
                </td>
                <td className="py-2">
                  <span
                    className={`rounded px-2 py-0.5 text-[10px] font-medium border ${
                      ev.type === 'TASK_GENERATED'
                        ? 'border-[#2F9C95]/40 bg-[#2F9C95]/15 text-[#E5F9E0]'
                        : 'border-[#40C9A2]/40 bg-[#40C9A2]/15 text-[#A3F7B5]'
                    }`}
                  >
                    {ev.type === 'TASK_GENERATED'
                      ? 'Генерация'
                      : 'Освобождение'}
                  </span>
                </td>
                <td className="py-2 text-slate-300">
                  {ev.source_id ? `И${ev.source_id}` : `П${ev.device_id}`}
                </td>
                <td className="py-2 text-slate-400">
                  {ev.task ? `№${ev.task.id} (И${ev.task.source_id})` : '—'}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
