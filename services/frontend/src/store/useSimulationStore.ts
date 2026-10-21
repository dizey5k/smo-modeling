import { create } from 'zustand'
import { devtools } from 'zustand/middleware'
import type { SimulationState, AutoResult } from '../types/smo'

interface SimulationStore {
  mode: 'STEP' | 'AUTO'
  isPlaying: boolean
  stepSpeed: number // мс между шагами при авто-воспроизведении
  history: string[] // лог последних действий
  lastState: SimulationState | null
  autoResult: AutoResult | null

  setMode: (mode: 'STEP' | 'AUTO') => void
  setIsPlaying: (playing: boolean) => void
  setStepSpeed: (speed: number) => void
  setLastState: (state: SimulationState) => void
  setAutoResult: (result: AutoResult) => void
  resetStore: () => void
}

export const useSimulationStore = create<SimulationStore>()(
  devtools((set) => ({
    mode: 'STEP',
    isPlaying: false,
    stepSpeed: 500,
    history: [],
    lastState: null,
    autoResult: null,

    setMode: (mode) => set({ mode }),
    setIsPlaying: (isPlaying) => set({ isPlaying }),
    setStepSpeed: (stepSpeed) => set({ stepSpeed }),
    setLastState: (state) =>
      set((prev) => ({
        lastState: state,
        history: state.last_action
          ? [state.last_action, ...prev.history.slice(0, 19)]
          : prev.history,
      })),
    setAutoResult: (autoResult) => set({ autoResult, mode: 'AUTO' }),
    resetStore: () =>
      set({
        isPlaying: false,
        history: [],
        autoResult: null,
      }),
  })),
)
