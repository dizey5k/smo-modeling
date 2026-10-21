import React, { useState } from 'react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts'
import { simulationApi } from '../api/simulationApi'
import type { ExperimentResult, ExperimentParams } from '../types/smo'
import { Sparkles, CheckCircle2, Sliders, RotateCcw } from 'lucide-react'

const DEFAULT_PARAMS: ExperimentParams = {
  target_n: 1000,
  device_options: [2, 3, 4],
  buffer_options: [4, 8, 12, 16],
  modes: ['RING', 'ZONED'],
}

export const ResearchTab: React.FC = () => {
  const [data, setData] = useState<ExperimentResult | null>(null)
  const [loading, setLoading] = useState(false)
  const [showSettings, setShowSettings] = useState(false)

  // Пользовательские параметры с дефолтами
  const [params, setParams] = useState<ExperimentParams>(DEFAULT_PARAMS)

  const runExperiment = async () => {
    setLoading(true)
    try {
      const res = await simulationApi.experiment(params)
      setData(res)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const handleResetDefaults = () => {
    setParams(DEFAULT_PARAMS)
  }

  // Динамические буферы и приборы для графиков из результатов
  const bufferList = Array.from(
    new Set(data?.runs.map((r) => r.buffer_cap) || [4, 8, 12, 16]),
  ).sort((a, b) => a - b)
  const devList = Array.from(
    new Set(data?.runs.map((r) => r.devices_count) || [2, 3, 4]),
  ).sort((a, b) => a - b)

  const chartData = bufferList.map((buf) => {
    const row: Record<string, number> = { buffer: buf }
    devList.forEach((m) => {
      const run = data?.runs.find(
        (r) =>
          r.buffer_cap === buf &&
          r.devices_count === m &&
          r.buffer_mode === 'RING',
      )
      if (run) {
        row[`p_otk_m${m}`] = Number((run.p_otk * 100).toFixed(2))
        row[`k_isp_m${m}`] = Number((run.k_isp * 100).toFixed(2))
      }
    })
    return row
  })

  const colors = ['#e3a8af', '#40C9A2', '#A3F7B5', '#2F9C95', '#FCD34D']

  return (
    <div className="space-y-6">
      {/* Главная плашка управления */}
      <div className="rounded-2xl border border-slate-800 bg-[#121820] p-5 shadow-sm">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h2 className="text-lg font-bold text-[#E5F9E0]">
              Параметрический синтез и оптимизация ВС (Этапы 3 и 4)
            </h2>
            <p className="text-xs text-slate-400 mt-1">
              Критерии: P_отк &le; 10%, K_исп &ge; 80%, &rho; &asymp; 1. Выборка
              N = <strong className="text-[#40C9A2]">{params.target_n}</strong>{' '}
              заявок на конфигурацию.
            </p>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={() => setShowSettings(!showSettings)}
              className={`flex items-center gap-1.5 rounded-xl border px-3.5 py-2.5 text-xs font-semibold transition ${
                showSettings
                  ? 'border-[#2F9C95] bg-[#2F9C95]/20 text-[#A3F7B5]'
                  : 'border-slate-800 bg-[#0B0F14] text-slate-300 hover:bg-slate-800'
              }`}
            >
              <Sliders className="h-3.5 w-3.5" />
              Параметры сетки
            </button>

            <button
              onClick={runExperiment}
              disabled={loading}
              className="flex items-center gap-2 rounded-xl bg-[#40C9A2] px-5 py-2.5 text-xs font-bold text-slate-950 transition hover:bg-[#A3F7B5] disabled:opacity-50 shadow-md shadow-[#40C9A2]/20 active:scale-95"
            >
              <Sparkles className="h-4 w-4" />
              {loading ? 'Расчёт...' : 'Запустить исследование'}
            </button>
          </div>
        </div>

        {/* Панель настройки параметров эксперимента */}
        {showSettings && (
          <div className="mt-4 border-t border-slate-800/80 pt-4 grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs animate-in fade-in duration-150">
            {/* Объем N */}
            <div className="rounded-xl border border-slate-800/80 bg-[#0B0F14] p-3">
              <label className="text-slate-400 block mb-1.5 font-medium">
                Объём выборки (N заявок на точку):
              </label>
              <div className="flex items-center gap-2">
                <input
                  type="number"
                  min="200"
                  max="10000"
                  step="100"
                  value={params.target_n}
                  onChange={(e) =>
                    setParams({
                      ...params,
                      target_n: Math.max(100, Number(e.target.value)),
                    })
                  }
                  className="w-full rounded-lg border border-slate-700 bg-[#121820] px-3 py-1.5 font-mono text-[#E5F9E0] focus:border-[#40C9A2] focus:outline-none"
                />
                <div className="flex gap-1">
                  {[500, 1000, 3000].map((val) => (
                    <button
                      key={val}
                      onClick={() => setParams({ ...params, target_n: val })}
                      className={`px-2 py-1 rounded text-[10px] font-mono border ${
                        params.target_n === val
                          ? 'border-[#40C9A2] bg-[#40C9A2]/20 text-[#40C9A2]'
                          : 'border-slate-800 bg-[#18212C] text-slate-400 hover:text-slate-200'
                      }`}
                    >
                      {val}
                    </button>
                  ))}
                </div>
              </div>
            </div>

            {/* Приборы M */}
            <div className="rounded-xl border border-slate-800/80 bg-[#0B0F14] p-3">
              <label className="text-slate-400 block mb-1.5 font-medium">
                Число приборов NPU (через запятую):
              </label>
              <input
                type="text"
                value={params.device_options?.join(', ')}
                onChange={(e) => {
                  const arr = e.target.value
                    .split(',')
                    .map((x) => parseInt(x.trim(), 10))
                    .filter((x) => !isNaN(x) && x > 0)
                  setParams({ ...params, device_options: arr })
                }}
                placeholder="2, 3, 4"
                className="w-full rounded-lg border border-slate-700 bg-[#121820] px-3 py-1.5 font-mono text-[#E5F9E0] focus:border-[#40C9A2] focus:outline-none"
              />
            </div>

            {/* Емкость буфера */}
            <div className="rounded-xl border border-slate-800/80 bg-[#0B0F14] p-3 flex flex-col justify-between">
              <div>
                <label className="text-slate-400 block mb-1.5 font-medium">
                  Емкости буфера (слоты):
                </label>
                <input
                  type="text"
                  value={params.buffer_options?.join(', ')}
                  onChange={(e) => {
                    const arr = e.target.value
                      .split(',')
                      .map((x) => parseInt(x.trim(), 10))
                      .filter((x) => !isNaN(x) && x > 0)
                    setParams({ ...params, buffer_options: arr })
                  }}
                  placeholder="4, 8, 12, 16"
                  className="w-full rounded-lg border border-slate-700 bg-[#121820] px-3 py-1.5 font-mono text-[#E5F9E0] focus:border-[#40C9A2] focus:outline-none"
                />
              </div>

              <div className="mt-2 flex justify-end">
                <button
                  onClick={handleResetDefaults}
                  className="flex items-center gap-1 text-[11px] text-[#e3a8af] hover:underline"
                >
                  <RotateCcw className="h-3 w-3" />
                  Сбросить к дефолтам (1000, [2..4], [4..16])
                </button>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Карточка оптимума */}
      {data?.optimal_config && (
        <div className="rounded-2xl border border-[#40C9A2]/40 bg-[#40C9A2]/10 p-5 shadow-sm">
          <div className="flex items-center gap-2 text-sm font-bold text-[#A3F7B5]">
            <CheckCircle2 className="h-5 w-5" />
            Рекомендованная оптимальная конфигурация архитектуры ВС:
          </div>
          <div className="mt-3 grid grid-cols-2 gap-4 sm:grid-cols-5 font-mono text-xs">
            <div className="rounded-lg bg-[#0B0F14] p-3 border border-slate-800">
              <span className="text-slate-400 block text-[10px]">
                Приборы (NPU)
              </span>
              <strong className="text-[#E5F9E0] text-base">
                {data.optimal_config.devices_count} ядра
              </strong>
            </div>
            <div className="rounded-lg bg-[#0B0F14] p-3 border border-slate-800">
              <span className="text-slate-400 block text-[10px]">
                Буфер (ОЗУ)
              </span>
              <strong className="text-[#E5F9E0] text-base">
                {data.optimal_config.buffer_cap} слотов
              </strong>
            </div>
            <div className="rounded-lg bg-[#0B0F14] p-3 border border-slate-800">
              <span className="text-slate-400 block text-[10px]">
                Вероятность P_отк
              </span>
              <strong className="text-[#A3F7B5] text-base">
                {(data.optimal_config.p_otk * 100).toFixed(2)}%
              </strong>
            </div>
            <div className="rounded-lg bg-[#0B0F14] p-3 border border-slate-800">
              <span className="text-slate-400 block text-[10px]">
                Загрузка K_исп
              </span>
              <strong className="text-[#40C9A2] text-base">
                {(data.optimal_config.k_isp * 100).toFixed(2)}%
              </strong>
            </div>
            <div className="rounded-lg bg-[#0B0F14] p-3 border border-slate-800">
              <span className="text-slate-400 block text-[10px]">
                Стоимость решения
              </span>
              <strong className="text-[#E5F9E0] text-base">
                {data.optimal_config.cost.toLocaleString()} руб.
              </strong>
            </div>
          </div>
        </div>
      )}

      {/* Графики и Таблица */}
      {data && (
        <>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="rounded-2xl border border-slate-800 bg-[#121820] p-5 shadow-sm">
              <h3 className="text-sm font-semibold text-[#E5F9E0] mb-4">
                График: Вероятность отказа P_отк (%) от емкости буфера
              </h3>
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#1B2927" />
                    <XAxis dataKey="buffer" stroke="#94a3b8" />
                    <YAxis stroke="#94a3b8" unit="%" />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: '#0B0F14',
                        borderColor: '#1B2927',
                      }}
                    />
                    <Legend />
                    {devList.map((m, idx) => (
                      <Line
                        key={m}
                        type="monotone"
                        dataKey={`p_otk_m${m}`}
                        name={`${m} NPU`}
                        stroke={colors[idx % colors.length]}
                        strokeWidth={2}
                      />
                    ))}
                  </LineChart>
                </ResponsiveContainer>
              </div>
            </div>

            <div className="rounded-2xl border border-slate-800 bg-[#121820] p-5 shadow-sm">
              <h3 className="text-sm font-semibold text-[#E5F9E0] mb-4">
                График: Коэффициент использования K_исп (%) от емкости буфера
              </h3>
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#1B2927" />
                    <XAxis dataKey="buffer" stroke="#94a3b8" />
                    <YAxis stroke="#94a3b8" unit="%" />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: '#0B0F14',
                        borderColor: '#1B2927',
                      }}
                    />
                    <Legend />
                    {devList.map((m, idx) => (
                      <Line
                        key={m}
                        type="monotone"
                        dataKey={`k_isp_m${m}`}
                        name={`${m} NPU`}
                        stroke={colors[idx % colors.length]}
                        strokeWidth={2}
                      />
                    ))}
                  </LineChart>
                </ResponsiveContainer>
              </div>
            </div>
          </div>

          <div className="rounded-2xl border border-slate-800 bg-[#121820] p-5 shadow-sm">
            <h3 className="text-sm font-semibold text-[#E5F9E0] mb-3">
              Сводная матрица конфигураций (Результаты перебора)
            </h3>
            <div className="overflow-x-auto">
              <table className="w-full text-left font-mono text-xs">
                <thead>
                  <tr className="border-b border-slate-800 text-slate-400">
                    <th className="py-2">Режим</th>
                    <th className="py-2">Приборы (M)</th>
                    <th className="py-2">Буфер (N)</th>
                    <th className="py-2">Загрузка (&rho;)</th>
                    <th className="py-2">P_отк (%)</th>
                    <th className="py-2">K_исп (%)</th>
                    <th className="py-2">T_преб</th>
                    <th className="py-2">Стоимость</th>
                    <th className="py-2">Статус ТЗ</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-800/60">
                  {data.runs.map((r, i) => (
                    <tr
                      key={i}
                      className={`hover:bg-slate-800/30 transition-colors ${
                        r.is_optimal ? 'bg-[#40C9A2]/10 font-bold' : ''
                      }`}
                    >
                      <td className="py-2 text-[#40C9A2]">{r.buffer_mode}</td>
                      <td className="py-2">{r.devices_count}</td>
                      <td className="py-2">{r.buffer_cap}</td>
                      <td className="py-2">{r.load_rho.toFixed(2)}</td>
                      <td
                        className={`py-2 ${r.p_otk > 0.1 ? 'text-[#e3a8af]' : 'text-[#A3F7B5]'}`}
                      >
                        {(r.p_otk * 100).toFixed(2)}%
                      </td>
                      <td className="py-2 text-[#40C9A2]">
                        {(r.k_isp * 100).toFixed(2)}%
                      </td>
                      <td className="py-2">{r.t_stay.toFixed(4)}</td>
                      <td className="py-2">{r.cost.toLocaleString()} р.</td>
                      <td className="py-2">
                        {r.is_optimal ? (
                          <span className="rounded bg-emerald-950 border border-emerald-500/40 text-[#A3F7B5] px-2 py-0.5 text-[10px]">
                            Оптимально
                          </span>
                        ) : (
                          <span className="text-slate-500 text-[10px]">
                            Вне нормы
                          </span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </>
      )}
    </div>
  )
}
