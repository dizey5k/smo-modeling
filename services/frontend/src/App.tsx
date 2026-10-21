import React, { useState } from 'react'
import { useSimulationStore } from './store/useSimulationStore'
import { useSimulationQueries } from './hooks/useSimulationQueries'
import { Header } from './components/Header'
import { MetricsOverview } from './components/MetricsOverview'
import { ControlPanel } from './components/ControlPanel'
import { BufferRing } from './components/BufferRing'
import { DeviceGrid } from './components/DeviceGrid'
import { EventCalendar } from './components/EventCalendar'
import { AutoTable } from './components/AutoTable'
import { ResearchTab } from './components/ResearchTab'

export const App: React.FC = () => {
  const { mode, lastState, autoResult, setMode } = useSimulationStore()
  const { initMutation, stepMutation, autoMutation, isLoading } =
    useSimulationQueries()
  const [activeTab, setActiveTab] = useState<'SIMULATION' | 'RESEARCH'>(
    'SIMULATION',
  )
  const [bufferMode, setBufferMode] = useState<'RING' | 'ZONED'>('RING')

  const handleToggleBufferMode = (newMode: 'RING' | 'ZONED') => {
    setBufferMode(newMode)
    initMutation.mutate({ buffer_mode: newMode })
  }

  const handleStepBatch = async (count: number) => {
    for (let i = 0; i < count; i++) {
      await stepMutation.mutateAsync()
    }
  }

  return (
    <div className="min-h-screen bg-[#0B0F14] p-4 text-slate-100 sm:p-8">
      <div className="mx-auto max-w-6xl space-y-5">
        <Header />

        <div className="flex border-b border-slate-800 gap-6">
          <button
            onClick={() => setActiveTab('SIMULATION')}
            className={`pb-3 text-sm font-semibold transition ${
              activeTab === 'SIMULATION'
                ? 'border-b-2 border-[#40C9A2] text-[#40C9A2]'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            Стенд симуляции (Этапы 1 и 2)
          </button>
          <button
            onClick={() => setActiveTab('RESEARCH')}
            className={`pb-3 text-sm font-semibold transition ${
              activeTab === 'RESEARCH'
                ? 'border-b-2 border-[#2F9C95] text-[#A3F7B5]'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            Синтез и оптимизация ВС (Этапы 3 и 4)
          </button>
        </div>

        {activeTab === 'SIMULATION' ? (
          <>
            <ControlPanel
              mode={mode}
              loading={isLoading}
              bufferMode={bufferMode}
              lastAction={lastState?.last_action}
              onSetMode={setMode}
              onToggleBufferMode={handleToggleBufferMode}
              onStep={() => stepMutation.mutate()}
              onStepBatch={handleStepBatch}
              onAuto={() => autoMutation.mutate({ buffer_mode: bufferMode })}
              onReset={() => initMutation.mutate({ buffer_mode: bufferMode })}
            />

            {mode === 'STEP' && lastState && (
              <>
                <MetricsOverview
                  time={lastState.current_time}
                  generated={lastState.total_gen}
                  dropped={lastState.total_drop}
                  done={lastState.total_done}
                />

                <div className="grid grid-cols-1 gap-5 lg:grid-cols-2">
                  <BufferRing
                    slots={lastState.buffer}
                    pointerUB={lastState.pointer_ub}
                    bufferMode={bufferMode}
                  />
                  <DeviceGrid
                    devices={lastState.devices}
                    pointerUP={lastState.pointer_up}
                  />
                </div>

                <EventCalendar events={lastState.calendar} />
              </>
            )}

            {mode === 'AUTO' && autoResult && <AutoTable result={autoResult} />}
          </>
        ) : (
          <ResearchTab />
        )}
      </div>
    </div>
  )
}

export default App
