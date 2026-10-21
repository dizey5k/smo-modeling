import React from 'react'
import type { AutoResult } from '../types/smo'

export const AutoTable: React.FC<{ result: AutoResult }> = ({ result }) => {
  return (
    <div className="space-y-4">
      <div className="rounded-xl border border-[#1B2927] bg-[#111A19] p-4">
        <div className="mb-3 flex items-center justify-between">
          <h3 className="text-sm font-semibold text-[#40C9A2]">
            Таблица 1. Характеристики источников
          </h3>
          <span className="font-mono text-xs text-[#E5F9E0]/70">
            N = <strong className="text-[#E5F9E0]">{result.target_n}</strong> |
            T = {result.total_time.toFixed(2)}
          </span>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left font-mono text-xs">
            <thead>
              <tr className="border-b border-[#1B2927] text-[#2F9C95]">
                <th className="py-2">Ист.</th>
                <th className="py-2">Сгенер.</th>
                <th className="py-2">Отказов</th>
                <th className="py-2">P_отк</th>
                <th className="py-2">T_преб</th>
                <th className="py-2">T_бп</th>
                <th className="py-2">T_обсл</th>
                <th className="py-2">D_бп</th>
                <th className="py-2">D_обсл</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[#1B2927]/60">
              {result.source_stats.map((s) => (
                <tr
                  key={s.source_id}
                  className="hover:bg-[#162221]/50 transition-colors"
                >
                  <td className="py-2 font-bold text-[#E5F9E0]">
                    И{s.source_id}
                  </td>
                  <td className="py-2">{s.generated}</td>
                  <td className="py-2 text-[#E5A8B0]">{s.dropped}</td>
                  <td className="py-2 text-[#40C9A2] font-semibold">
                    {(s.p_otk * 100).toFixed(2)}%
                  </td>
                  <td className="py-2">{s.t_stay.toFixed(4)}</td>
                  <td className="py-2">{s.t_wait.toFixed(4)}</td>
                  <td className="py-2">{s.t_serv.toFixed(4)}</td>
                  <td className="py-2 text-[#E5F9E0]/50">
                    {s.d_wait.toFixed(4)}
                  </td>
                  <td className="py-2 text-[#E5F9E0]/50">
                    {s.d_serv.toFixed(4)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <div className="rounded-xl border border-[#1B2927] bg-[#111A19] p-4">
        <h3 className="mb-3 text-sm font-semibold text-[#2F9C95]">
          Таблица 2. Характеристики приборов
        </h3>
        <div className="overflow-x-auto">
          <table className="w-full text-left font-mono text-xs">
            <thead>
              <tr className="border-b border-[#1B2927] text-[#2F9C95]">
                <th className="py-2">Прибор</th>
                <th className="py-2">Время занятости</th>
                <th className="py-2">Коэффициент использования K_исп</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[#1B2927]/60">
              {result.device_stats.map((d) => (
                <tr
                  key={d.device_id}
                  className="hover:bg-[#162221]/50 transition-colors"
                >
                  <td className="py-2 font-bold text-[#E5F9E0]">
                    NPU_{d.device_id}
                  </td>
                  <td className="py-2">{d.busy_time.toFixed(2)}</td>
                  <td className="py-2 font-bold text-[#A3F7B5]">
                    {(d.k_isp * 100).toFixed(2)}%
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
