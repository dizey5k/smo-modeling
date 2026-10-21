export interface InferenceTask {
  id: number
  source_id: number
  priority: number
  generation_time: number
  buffer_entry_time: number
  service_start_time: number
  service_end_time: number
  status: 'CREATED' | 'IN_BUFFER' | 'IN_SERVICE' | 'COMPLETED' | 'DROPPED'
}

export interface BufferSlot {
  index: number
  is_occupied: boolean
  task?: InferenceTask
}

export interface Device {
  id: number
  is_busy: boolean
  current_task?: InferenceTask
  busy_duration: number
}

export interface EventItem {
  time: number
  type: 'TASK_GENERATED' | 'DEVICE_FREED'
  source_id?: number
  device_id?: number
  task?: InferenceTask
}

export interface SimulationState {
  current_time: number
  pointer_ub: number
  pointer_up: number
  buffer: BufferSlot[]
  devices: Device[]
  calendar: EventItem[]
  last_action: string
  total_gen: number
  total_drop: number
  total_done: number
  is_finished: boolean
}

export interface SourceStat {
  source_id: number
  generated: number
  dropped: number
  completed: number
  p_otk: number
  t_stay: number
  t_wait: number
  t_serv: number
  d_wait: number
  d_serv: number
}

export interface DeviceStat {
  device_id: number
  busy_time: number
  k_isp: number
}

export interface AutoResult {
  target_n: number
  total_time: number
  source_stats: SourceStat[]
  device_stats: DeviceStat[]
}
