import React from 'react'

interface Props {
  time: number
  generated: number
  dropped: number
  done: number
}

export const MetricsOverview: React.FC<Props> = ({
  time,
  generated,
  dropped,
  done,
}) => {
  return (
    <div className="grid grid-cols-2 gap-3.5 sm:grid-cols-4">
      <div className="rounded-xl border border-slate-800/90 bg-[#121820] p-4">
        <div className="text-[11px] font-medium uppercase tracking-wider text-slate-400">
          Модельное время (T)
        </div>
        <div className="mt-1.5 font-mono text-2xl font-bold text-[#40C9A2]">
          {time.toFixed(4)}
        </div>
      </div>

      <div className="rounded-xl border border-slate-800/90 bg-[#121820] p-4">
        <div className="text-[11px] font-medium uppercase tracking-wider text-slate-400">
          Сгенерировано
        </div>
        <div className="mt-1.5 font-mono text-2xl font-bold text-[#E5F9E0]">
          {generated}
        </div>
      </div>

      <div className="rounded-xl border border-[#664147]/70 bg-[#664147]/15 p-4">
        <div className="text-[11px] font-medium uppercase tracking-wider text-[#e3a8af]">
          Сброшено (Отказы)
        </div>
        <div className="mt-1.5 font-mono text-2xl font-bold text-[#f5c2c7]">
          {dropped}
        </div>
      </div>

      <div className="rounded-xl border border-slate-800/90 bg-[#121820] p-4">
        <div className="text-[11px] font-medium uppercase tracking-wider text-slate-400">
          Обработано
        </div>
        <div className="mt-1.5 font-mono text-2xl font-bold text-[#A3F7B5]">
          {done}
        </div>
      </div>
    </div>
  )
}
