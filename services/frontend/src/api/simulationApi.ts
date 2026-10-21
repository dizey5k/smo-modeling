import { apiClient } from './client'
import type {
  SimulationState,
  AutoResult,
  ExperimentResult,
  ExperimentParams,
} from '../types/smo'

export const simulationApi = {
  init: async (config?: unknown): Promise<SimulationState> => {
    const { data } = await apiClient.post<SimulationState>(
      '/init',
      config ?? {},
    )
    return data
  },

  step: async (): Promise<SimulationState> => {
    const { data } = await apiClient.post<SimulationState>('/step')
    return data
  },

  auto: async (config?: unknown): Promise<AutoResult> => {
    const { data } = await apiClient.post<AutoResult>('/auto', config ?? {})
    return data
  },

  experiment: async (params?: ExperimentParams): Promise<ExperimentResult> => {
    const { data } = await apiClient.post<ExperimentResult>(
      '/experiment',
      { experiment: params ?? {} },
      { timeout: 60000 },
    )
    return data
  },
}
