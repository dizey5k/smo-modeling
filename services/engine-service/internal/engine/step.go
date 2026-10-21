package engine

import (
	"fmt"

	"smo-gateway/internal/model"
)

func (se *SimulationEngine) Step() model.SimulationState {
	if len(se.Calendar) == 0 {
		return se.GetState("Очередь событий пуста", true)
	}

	event := se.Calendar[0]
	se.Calendar = se.Calendar[1:]
	se.CurrentTime = event.Time

	var desc string
	switch event.Type {
	case model.EventTaskGenerated:
		desc = se.handleTaskGenerated(event.SourceID)
	case model.EventDeviceFreed:
		desc = se.handleDeviceFreed(event.DeviceID)
	}

	return se.GetState(desc, false)
}

func (se *SimulationEngine) handleTaskGenerated(srcID int) string {
	se.TaskSeq++
	se.TotalGen++
	se.SourceGenCount[srcID]++

	task := &model.InferenceTask{
		ID:             se.TaskSeq,
		SourceID:       srcID,
		Priority:       srcID,
		GenerationTime: se.CurrentTime,
		Status:         model.StatusCreated,
	}
	se.Tasks[task.ID] = task

	// Планируем следующую генерацию источника по ИЗ2 через микросервис generator-service
	for _, s := range se.Config.Sources {
		if s.ID == srcID {
			nextInterval := se.GenClient.GetUniform(s.RangeA, s.RangeB)
			se.ScheduleEvent(model.Event{
				Time:     se.CurrentTime + nextInterval,
				Type:     model.EventTaskGenerated,
				SourceID: srcID,
			})
			break
		}
	}

	// Постановка в буфер (Д1ОЗ1 / Д1ОО1)
	dropped, slot := se.Buffer.PutTask(task, se.CurrentTime)
	var desc string
	if dropped != nil {
		se.TotalDrop++
		se.SourceDropCount[dropped.SourceID]++
		desc = fmt.Sprintf("Генерация И%d (#%d). Буфер полон: выбит #%d под УБ. Новый в слоте %d.",
			srcID, task.ID, dropped.ID, slot)
	} else {
		desc = fmt.Sprintf("Генерация И%d (#%d). Записан в слот %d (УБ=%d).",
			srcID, task.ID, slot, se.Buffer.GetPointer())
	}

	if dispatchMsg := se.tryDispatch(); dispatchMsg != "" {
		desc += dispatchMsg
	}
	return desc
}

func (se *SimulationEngine) handleDeviceFreed(devID int) string {
	dev := &se.Pool.Devices[devID]
	if dev.CurrentTask != nil {
		t := dev.CurrentTask
		t.ServiceEndTime = se.CurrentTime
		t.Status = model.StatusCompleted
		se.TotalDone++
		se.SourceDoneCount[t.SourceID]++

		wait := t.ServiceStartTime - t.BufferEntryTime
		serv := t.ServiceEndTime - t.ServiceStartTime
		stay := t.ServiceEndTime - t.GenerationTime

		se.WaitDurations[t.SourceID] = append(se.WaitDurations[t.SourceID], wait)
		se.ServDurations[t.SourceID] = append(se.ServDurations[t.SourceID], serv)
		se.StayDurations[t.SourceID] = append(se.StayDurations[t.SourceID], stay)

		dev.BusyDuration += serv
		dev.CurrentTask = nil
	}
	dev.IsBusy = false

	desc := fmt.Sprintf("П%d освободился.", devID)
	if dispatchMsg := se.tryDispatch(); dispatchMsg != "" {
		desc += dispatchMsg
	} else {
		desc += " Буфер пуст, П простаивает."
	}
	return desc
}

func (se *SimulationEngine) tryDispatch() string {
	if se.Buffer.IsEmpty() {
		return ""
	}
	freeDev := se.Pool.SelectDeviceRoundRobin()
	if freeDev == nil {
		return ""
	}
	task := se.Buffer.ExtractTaskByPriority()
	if task == nil {
		return ""
	}

	freeDev.IsBusy = true
	freeDev.CurrentTask = task
	task.ServiceStartTime = se.CurrentTime
	task.Status = model.StatusInService

	// Время инференса по закону ПЗ1 генерируется микросервисом generator-service
	servTime := se.GenClient.GetExponential(se.Pool.MuRate)
	se.ScheduleEvent(model.Event{
		Time:     se.CurrentTime + servTime,
		Type:     model.EventDeviceFreed,
		DeviceID: freeDev.ID,
		Task:     task,
	})

	return fmt.Sprintf(" ДВ передал задачу #%d (И%d) на П%d по кольцу УП (длительность %.3f).",
		task.ID, task.SourceID, freeDev.ID, servTime)
}
