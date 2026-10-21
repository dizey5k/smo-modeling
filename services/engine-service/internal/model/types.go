package model

type TaskStatus string

const (
	StatusCreated   TaskStatus = "CREATED"
	StatusInBuffer  TaskStatus = "IN_BUFFER"
	StatusInService TaskStatus = "IN_SERVICE"
	StatusCompleted TaskStatus = "COMPLETED"
	StatusDropped   TaskStatus = "DROPPED"
)

type EventType string

const (
	EventTaskGenerated EventType = "TASK_GENERATED"
	EventDeviceFreed   EventType = "DEVICE_FREED"
)

// InferenceTask представляет заявку на обработку кадра (18 вариант)
type InferenceTask struct {
	ID               int        `json:"id"`
	SourceID         int        `json:"source_id"` // 1, 2 или 3
	Priority         int        `json:"priority"`  // 1 - наивысший, 2 - средний, 3 - базовый
	GenerationTime   float64    `json:"generation_time"`
	BufferEntryTime  float64    `json:"buffer_entry_time"`
	ServiceStartTime float64    `json:"service_start_time"`
	ServiceEndTime   float64    `json:"service_end_time"`
	Status           TaskStatus `json:"status"`
}

// Event - особое событие в календаре
type Event struct {
	Time     float64        `json:"time"`
	Type     EventType      `json:"type"`
	SourceID int            `json:"source_id,omitempty"` // Для генерации
	DeviceID int            `json:"device_id,omitempty"` // Для освобождения прибора
	Task     *InferenceTask `json:"task,omitempty"`
}

// BufferSlot - ячейка кольцевой буферной памяти
type BufferSlot struct {
	Index      int            `json:"index"`
	IsOccupied bool           `json:"is_occupied"`
	Task       *InferenceTask `json:"task,omitempty"`
}

// Device - аппаратный вычислитель (NPU / GPU core)
type Device struct {
	ID           int            `json:"id"`
	IsBusy       bool           `json:"is_busy"`
	CurrentTask  *InferenceTask `json:"current_task,omitempty"`
	BusyDuration float64        `json:"busy_duration"`
}

// SourceConfig - параметры источника
type SourceConfig struct {
	ID       int     `json:"id"`
	Priority int     `json:"priority"`
	RangeA   float64 `json:"range_a"` // Равномерный закон ИЗ2: интервал [A, B]
	RangeB   float64 `json:"range_b"`
}

// SimulationConfig - начальная конфигурация стенда
type SimulationConfig struct {
	SourcesCount int            `json:"sources_count"` // По умолчанию 3
	BufferCap    int            `json:"buffer_cap"`    // Емкость N
	DevicesCount int            `json:"devices_count"` // Число вычислителей M
	MuRate       float64        `json:"mu_rate"`       // Интенсивность экспоненциального закона ПЗ1
	Sources      []SourceConfig `json:"sources"`
	BufferMode   string         `json:"buffer_mode"` // "RING" или "ZONED" (для бонуса)
}

// SourceStat - выходная статистика источника по ОР1
type SourceStat struct {
	SourceID    int     `json:"source_id"`
	Generated   int     `json:"generated"`
	Dropped     int     `json:"dropped"`
	Completed   int     `json:"completed"`
	ProbRefusal float64 `json:"p_otk"`  // P_отк
	TStay       float64 `json:"t_stay"` // T_преб
	TWait       float64 `json:"t_wait"` // T_бп
	TService    float64 `json:"t_serv"` // T_обсл
	DispWait    float64 `json:"d_wait"` // Д_бп
	DispService float64 `json:"d_serv"` // Д_обсл
}

// DeviceStat - выходная статистика прибора по ОР1
type DeviceStat struct {
	DeviceID   int     `json:"device_id"`
	BusyTime   float64 `json:"busy_time"`
	UsageCoeff float64 `json:"k_isp"` // K_исп
}

// SimulationState - срез текущего состояния системы (для пошагового режима)
type SimulationState struct {
	CurrentTime float64      `json:"current_time"`
	PointerUB   int          `json:"pointer_ub"` // Указатель буфера
	PointerUP   int          `json:"pointer_up"` // Указатель приборов
	Buffer      []BufferSlot `json:"buffer"`
	Devices     []Device     `json:"devices"`
	Events      []Event      `json:"calendar"`
	LastAction  string       `json:"last_action"`
	TotalGen    int          `json:"total_gen"`
	TotalDrop   int          `json:"total_drop"`
	TotalDone   int          `json:"total_done"`
	IsFinished  bool         `json:"is_finished"`
}
