import React from 'react'
import { Play, RotateCcw, FastForward, Layers } from 'lucide-react'

interface Props {
  mode: 'STEP' | 'AUTO'
  loading: boolean
  bufferMode: 'RING' | 'ZONED'
  lastAction?: string
  onSetMode: (m: 'STEP' | 'AUTO') => void
  onToggleBufferMode: (m: 'RING' | 'ZONED') => void
  onStep: () => void
  onStepBatch: (count: number) => void
  onAuto: () => void
  onReset: () => void
}

export const ControlPanel: React.FC<Props> = ({
  mode,
  loading,
  bufferMode,
  lastAction,
  onSetMode,
  onToggleBufferMode,
  onStep,
  onStepBatch,
  onAuto,
  onReset,
}) => {
  return (
    <div className="rounded-2xl border border-slate-800/90 bg-[#121820] p-4 shadow-sm">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        {/* Переключение режимов интерфейса */}
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

          {/* Переключатель типа буфера (Кольцевой / Зонный) */}
          <div className="ml-2 flex items-center rounded-xl border border-slate-800 bg-[#0B0F14] p-1">
            <button
              onClick={() => onToggleBufferMode('RING')}
              className={`rounded-lg px-2.5 py-1 text-[11px] font-semibold transition ${
                bufferMode === 'RING'
                  ? 'bg-[#40C9A2] text-slate-950'
                  : 'text-slate-400 hover:text-slate-200'
              }`}
            >
              Кольцевой буфер
            </button>
            <button
              onClick={() => onToggleBufferMode('ZONED')}
              className={`rounded-lg px-2.5 py-1 text-[11px] font-semibold transition ${
                bufferMode === 'ZONED'
                  ? 'bg-[#2F9C95] text-white'
                  : 'text-slate-400 hover:text-slate-200'
              }`}
            >
              Зонный буфер
            </button>
          </div>

          <button
            onClick={onReset}
            className="flex items-center gap-1.5 rounded-xl border border-[#664147]/70 bg-[#664147]/25 px-3 py-1.5 text-xs font-medium text-[#f3cbd0] transition hover:bg-[#664147]/45 active:scale-95"
          >
            <RotateCcw className="h-3 w-3" />
            Сброс
          </button>
        </div>

        {/* Кнопки шагов */}
        {mode === 'STEP' && (
          <div className="flex items-center gap-2">
            <button
              onClick={() => onStepBatch(10)}
              className="flex items-center gap-1.5 rounded-xl border border-slate-700 bg-[#18212C] px-3.5 py-2 text-xs font-semibold text-slate-200 transition hover:bg-[#202B39] active:scale-95"
            >
              <Layers className="h-3.5 w-3.5 text-[#40C9A2]" />
              +10 шагов
            </button>

            <button
              onClick={onStep}
              className="flex items-center justify-center gap-2 rounded-xl bg-[#40C9A2] px-5 py-2 font-bold text-slate-950 transition hover:bg-[#A3F7B5] active:scale-95 shadow-md shadow-[#40C9A2]/20"
            >
              <Play className="h-4 w-4 fill-current" />
              <span className="text-xs uppercase tracking-wider">
                Следующий шаг
              </span>
            </button>
          </div>
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
