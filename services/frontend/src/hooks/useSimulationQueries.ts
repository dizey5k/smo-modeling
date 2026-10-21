import { useMutation, useQueryClient } from '@tanstack/react-query'
import { simulationApi } from '../api/simulationApi'
import { useSimulationStore } from '../store/useSimulationStore'

export const SIMULATION_QUERY_KEY = ['simulationState']

export const useSimulationQueries = () => {
  const queryClient = useQueryClient()
  const { setLastState, setAutoResult, resetStore } = useSimulationStore()

  const initMutation = useMutation({
    mutationFn: simulationApi.init,
    onSuccess: (data) => {
      resetStore()
      setLastState(data)
      queryClient.setQueryData(SIMULATION_QUERY_KEY, data)
    },
  })

  const stepMutation = useMutation({
    mutationFn: simulationApi.step,
    onSuccess: (data) => {
      setLastState(data)
      queryClient.setQueryData(SIMULATION_QUERY_KEY, data)
    },
  })

  const autoMutation = useMutation({
    mutationFn: simulationApi.auto,
    onSuccess: (data) => {
      setAutoResult(data)
    },
  })

  return {
    initMutation,
    stepMutation,
    autoMutation,
    isLoading:
      initMutation.isPending ||
      stepMutation.isPending ||
      autoMutation.isPending,
  }
}
