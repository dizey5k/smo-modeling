import React from 'react'
import type { Device } from '../types/smo'

interface Props {
  devices: Device[]
  pointerUP: number
}

export const DeviceGrid: React.FC<Props> = ({ devices, pointerUP }) => {
  return (
    <div className="rounded-2xl border border-slate-800/90 bg-[#121820] p-5 shadow-sm">
      <div className="mb-4 flex items-center justify-between">
        <h3 className="text-sm font-semibold text-[#E5F9E0]">
          Ускорители инференса NPU (Д2П2: кольцо)
        </h3>
        <span className="rounded-lg border border-slate-800 bg-[#0B0F14] px-2.5 py-1 font-mono text-xs text-slate-400">
          УП = <strong className="text-[#2F9C95]">{pointerUP}</strong>
        </span>
      </div>

      <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
        {devices.map((dev) => {
          const isPointer = dev.id === pointerUP
          return (
            <div
              key={dev.id}
              className={`relative rounded-xl border p-3.5 transition-all duration-150 ${
                dev.is_busy
                  ? 'border-[#2F9C95]/70 bg-[#2F9C95]/15'
                  : 'border-slate-800/80 bg-[#0B0F14]/70'
              } ${isPointer ? 'ring-2 ring-[#2F9C95] border-transparent shadow-lg shadow-[#2F9C95]/10' : ''}`}
            >
              {isPointer && (
                <span className="pointer-badge absolute -top-2.5 right-3 rounded-full bg-[#2F9C95] px-1.5 py-0.2 font-mono text-[9px] font-black text-white shadow">
                  УП
                </span>
              )}
              <div className="flex items-center justify-between">
                <span className="font-mono text-xs font-bold text-slate-200">
                  NPU_{dev.id}
                </span>
                <span
                  className={`rounded-md px-2 py-0.5 text-[10px] font-medium ${
                    dev.is_busy
                      ? 'bg-[#2F9C95]/30 text-[#A3F7B5] border border-[#2F9C95]/40'
                      : 'bg-slate-800 text-slate-400'
                  }`}
                >
                  {dev.is_busy ? 'Инференс' : 'Свободен'}
                </span>
              </div>

              {dev.is_busy && dev.current_task ? (
                <div className="mt-3 rounded-lg border border-slate-800 bg-[#0B0F14] p-2.5 font-mono text-[11px]">
                  <div className="flex justify-between text-slate-400">
                    <span>Задача:</span>
                    <span className="font-bold text-[#A3F7B5]">
                      №{dev.current_task.id}
                    </span>
                  </div>
                  <div className="mt-1 flex justify-between text-slate-400">
                    <span>Источник:</span>
                    <span className="text-[#E5F9E0]">
                      И{dev.current_task.source_id}
                    </span>
                  </div>
                </div>
              ) : (
                <div className="mt-3 py-2 text-center text-xs text-slate-600 font-mono">
                  Ожидание задачи
                </div>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}
