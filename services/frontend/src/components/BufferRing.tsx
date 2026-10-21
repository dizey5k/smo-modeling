import React from 'react'
import type { BufferSlot } from '../types/smo'

interface Props {
  slots: BufferSlot[]
  pointerUB: number
}

export const BufferRing: React.FC<Props> = ({ slots, pointerUB }) => {
  return (
    <div className="rounded-2xl border border-slate-800/90 bg-[#121820] p-5 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <h3 className="text-sm font-semibold text-[#E5F9E0]">
          Кольцевой буфер задач (Д1ОЗ1, Д1ОО1)
        </h3>
        <span className="rounded-lg border border-slate-800 bg-[#0B0F14] px-2.5 py-1 font-mono text-xs text-slate-400">
          УБ = <strong className="text-[#40C9A2]">{pointerUB}</strong>
        </span>
      </div>

      <div className="grid grid-cols-4 gap-2.5 sm:grid-cols-8">
        {slots.map((slot) => {
          const isPointer = slot.index === pointerUB
          return (
            <div
              key={slot.index}
              className={`relative flex flex-col items-center justify-center rounded-xl border p-3 transition-all duration-150 ${
                slot.is_occupied
                  ? 'border-[#40C9A2]/60 bg-[#40C9A2]/15 text-[#E5F9E0]'
                  : 'border-slate-800/80 bg-[#0B0F14]/70 text-slate-600'
              } ${isPointer ? 'ring-2 ring-[#40C9A2] border-transparent shadow-lg shadow-[#40C9A2]/10' : ''}`}
            >
              {isPointer && (
                <span className="pointer-badge absolute -top-2.5 rounded-full bg-[#40C9A2] px-1.5 py-0.2 font-mono text-[9px] font-black text-slate-950 shadow">
                  УБ
                </span>
              )}
              <span className="font-mono text-[10px] text-slate-500">
                #{slot.index}
              </span>
              {slot.is_occupied && slot.task ? (
                <div className="mt-1 text-center font-mono">
                  <div className="text-xs font-bold text-[#A3F7B5]">
                    №{slot.task.id}
                  </div>
                  <div className="text-[10px] text-slate-400">
                    И{slot.task.source_id}
                  </div>
                </div>
              ) : (
                <span className="mt-2 text-[11px] font-bold text-slate-800">
                  —
                </span>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}
