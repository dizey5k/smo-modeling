import React, { useEffect } from 'react'
import { useSimulationStore } from './store/useSimulationStore'
import { useSimulationQueries } from './hooks/useSimulationQueries'
import { MetricsOverview } from './components/MetricsOverview'
import { ControlPanel } from './components/ControlPanel'
import { BufferRing } from './components/BufferRing'
import { DeviceGrid } from './components/DeviceGrid'
import { EventCalendar } from './components/EventCalendar'
import { AutoTable } from './components/AutoTable'
import { Header } from './components/Header'

export const App: React.FC = () => {
  const { mode, lastState, autoResult, setMode } = useSimulationStore()
  const { initMutation, stepMutation, autoMutation, isLoading } =
    useSimulationQueries()

  useEffect(() => {
    initMutation.mutate()
  }, [])

  const error = initMutation.error || stepMutation.error || autoMutation.error

  return (
    <div className="min-h-screen bg-[#0B0F14] p-4 text-slate-100 sm:p-8">
      <div className="mx-auto max-w-6xl space-y-5">
        <Header />

        {error && (
          <div className="rounded-lg border border-rose-500/30 bg-rose-950/40 p-3 text-xs text-rose-300">
            {error.message} (проверьте доступность engine-service на :8080)
          </div>
        )}

        <ControlPanel
          mode={mode}
          loading={isLoading}
          lastAction={lastState?.last_action}
          onSetMode={setMode}
          onStep={() => stepMutation.mutate()}
          onAuto={() => autoMutation.mutate()}
          onReset={() => initMutation.mutate()}
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
      </div>
    </div>
  )
}

export default App
