package model

type TaskStatus string

const (
	StatusCreated   = "CREATED"
	StatusInBuffer  = "IN_BUFFER"
	StatusInService = "IN_SERVICE"
	StatusCompleted = "COMPLETED"
	StatusDropped   = "DROPPED"
)

type EventType string

const (
	EventTaskGenerated = "TASK_GENERATED"
	EventDeviceFreed   = "DEVICE_FREED"
)

// InferenceTask представляет заявку на обработку кадра (18 вариант)
type InferenceTask struct {
	ID               int     `json:"id"`
	SourceID         int     `json:"source_id"`
	Priority         int     `json:"priority"`
	GenerationTime   float64 `json:"generation_time"`
	BufferEntryTime  float64 `json:"buffer_entry_time"`
	ServiceStartTime float64 `json:"service_start_time"`
	ServiceEndTime   float64 `json:"service_end_time"`
	Status           string  `json:"status"`
}

type Event struct {
	Time     float64        `json:"time"`
	Type     string         `json:"type"`
	SourceID int            `json:"source_id,omitempty"`
	DeviceID int            `json:"device_id,omitempty"`
	Task     *InferenceTask `json:"task,omitempty"`
}

type BufferSlot struct {
	Index      int            `json:"index"`
	IsOccupied bool           `json:"is_occupied"`
	Task       *InferenceTask `json:"task,omitempty"`
}

type Device struct {
	ID           int            `json:"id"`
	IsBusy       bool           `json:"is_busy"`
	CurrentTask  *InferenceTask `json:"current_task,omitempty"`
	BusyDuration float64        `json:"busy_duration"`
}

type SourceConfig struct {
	ID       int     `json:"id"`
	Priority int     `json:"priority"`
	RangeA   float64 `json:"range_a"`
	RangeB   float64 `json:"range_b"`
}

type SimulationConfig struct {
	SourcesCount int            `json:"sources_count"`
	BufferCap    int            `json:"buffer_cap"`
	DevicesCount int            `json:"devices_count"`
	MuRate       float64        `json:"mu_rate"`
	BufferMode   string         `json:"buffer_mode"` // "RING" или "ZONED"
	Sources      []SourceConfig `json:"sources"`
}

// Статистика источника по стандарту ОР1
type SourceStat struct {
	SourceID  int     `json:"source_id"`
	Generated int     `json:"generated"`
	Dropped   int     `json:"dropped"`
	Completed int     `json:"completed"`
	POtk      float64 `json:"p_otk"`
	TStay     float64 `json:"t_stay"`
	TWait     float64 `json:"t_wait"`
	TServ     float64 `json:"t_serv"`
	DWait     float64 `json:"d_wait"`
	DServ     float64 `json:"d_serv"`
}

// Статистика вычислителя NPU по стандарту ОР1
type DeviceStat struct {
	DeviceID int     `json:"device_id"`
	BusyTime float64 `json:"busy_time"`
	KIsp     float64 `json:"k_isp"`
}

type SimulationState struct {
	CurrentTime float64      `json:"current_time"`
	PointerUB   int          `json:"pointer_ub"`
	PointerUP   int          `json:"pointer_up"`
	Buffer      []BufferSlot `json:"buffer"`
	Devices     []Device     `json:"devices"`
	Events      []Event      `json:"calendar"`
	LastAction  string       `json:"last_action"`
	TotalGen    int          `json:"total_gen"`
	TotalDrop   int          `json:"total_drop"`
	TotalDone   int          `json:"total_done"`
	IsFinished  bool         `json:"is_finished"`
}

// ExperimentRequest задает параметры для Grid Search
type ExperimentRequest struct {
	TargetN       int      `json:"target_n"`       // По умолчанию 1000
	DeviceOptions []int    `json:"device_options"` // По умолчанию [2, 3, 4]
	BufferOptions []int    `json:"buffer_options"` // По умолчанию [4, 8, 12, 16]
	Modes         []string `json:"modes"`          // По умолчанию ["RING", "ZONED"]
}
