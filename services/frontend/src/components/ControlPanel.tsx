import React from 'react'
import { Play, RotateCcw, FastForward } from 'lucide-react'

interface Props {
  mode: 'STEP' | 'AUTO'
  loading: boolean
  lastAction?: string
  onSetMode: (m: 'STEP' | 'AUTO') => void
  onStep: () => void
  onAuto: () => void
  onReset: () => void
}

export const ControlPanel: React.FC<Props> = ({
  mode,
  loading,
  lastAction,
  onSetMode,
  onStep,
  onAuto,
  onReset,
}) => {
  return (
    <div className="rounded-2xl border border-slate-800/90 bg-[#121820] p-4 shadow-sm">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex flex-wrap items-center gap-2">
          <button
            onClick={() => onSetMode('STEP')}
            className={`rounded-xl px-4 py-2 text-xs font-semibold transition-all ${
              mode === 'STEP'
                ? 'bg-[#2F9C95] text-white shadow'
                : 'bg-[#18212C] text-slate-300 hover:bg-[#202B39] hover:text-white'
            }`}
          >
            Пошаговый режим
          </button>

          <button
            onClick={onAuto}
            disabled={loading}
            className={`flex items-center gap-1.5 rounded-xl px-4 py-2 text-xs font-semibold transition-all ${
              mode === 'AUTO'
                ? 'bg-[#40C9A2] text-slate-950 shadow'
                : 'bg-[#18212C] text-slate-300 hover:bg-[#202B39] hover:text-white'
            }`}
          >
            <FastForward className="h-3.5 w-3.5" />
            {loading ? 'Расчёт...' : 'Автоматический (ОР1)'}
          </button>

          <button
            onClick={onReset}
            className="flex items-center gap-1.5 rounded-xl border border-[#664147]/70 bg-[#664147]/25 px-3.5 py-2 text-xs font-medium text-[#f3cbd0] transition hover:bg-[#664147]/45 active:scale-95"
          >
            <RotateCcw className="h-3 w-3" />
            Сброс
          </button>
        </div>

        {mode === 'STEP' && (
          <button
            onClick={onStep}
            className="flex items-center justify-center gap-2 rounded-xl bg-[#40C9A2] px-6 py-2.5 font-bold text-slate-950 transition hover:bg-[#A3F7B5] active:scale-95 shadow-md shadow-[#40C9A2]/20"
          >
            <Play className="h-4 w-4 fill-current" />
            <span className="text-xs uppercase tracking-wider">
              Следующий шаг
            </span>
          </button>
        )}
      </div>

      {mode === 'STEP' && lastAction && (
        <div className="mt-3.5 border-t border-slate-800/80 pt-3 text-xs">
          <span className="font-semibold text-[#2F9C95]">Событие:</span>{' '}
          <span className="font-mono text-slate-200">{lastAction}</span>
        </div>
      )}
    </div>
  )
}
