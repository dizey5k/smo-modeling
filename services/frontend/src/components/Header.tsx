import React from 'react'

export const Header: React.FC = () => {
  return (
    <header className="flex flex-col gap-3 rounded-2xl border border-slate-800/80 bg-[#121820] p-5 shadow-sm sm:flex-row sm:items-center sm:justify-between">
      <div>
        <div className="flex items-center gap-3">
          <h1 className="text-xl font-bold tracking-tight text-[#E5F9E0]">
            Edge AI Gateway • Симуляция СМО
          </h1>
          <span className="rounded-md border border-[#40C9A2]/40 bg-[#40C9A2]/10 px-2.5 py-0.5 font-mono text-xs font-semibold text-[#40C9A2]">
            Вариант 18
          </span>
        </div>
        <p className="mt-1 font-mono text-xs text-slate-400">
          Федоров Олег Александрович • 5130904/40106
        </p>
      </div>

      <div className="rounded-xl border border-slate-800 bg-[#0B0F14] px-3.5 py-2 font-mono text-[11px] text-[#A3F7B5]">
        ИБ | ИЗ2 | ПЗ1 | Д1ОЗ1 | Д1ОО1 | Д2П2 | Д2Б4 | ОР1 | ОД3
      </div>
    </header>
  )
}
