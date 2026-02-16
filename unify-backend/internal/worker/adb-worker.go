package worker

import (
	"fmt"
	"time"

	// "unify-backend/cmd"
	"unify-backend/internal/adb"
	"unify-backend/internal/database"
	"unify-backend/internal/http/sse"
	"unify-backend/internal/job"
	"unify-backend/internal/services"
	"unify-backend/models"

	"github.com/google/uuid"
)

func StartADBWorkerPool(manager *Manager, count int, jobs <-chan job.ADBJob) {
	sseManager := manager.GetSSE()
	for i := 1; i <= count; i++ {
		go adbWorker(i, jobs, sseManager)
	}
}

type UpdateDataResAdbProps struct {
	Status    adb.AdbStatus `json:"status"`
	StartedAt time.Time     `json:"startedAt"`
	FinishAt  time.Time     `json:"finishAt"`
	Result    string        `json:"result"`
	ID        uuid.UUID     `json:"id"`
}

func updateDataResADB(j job.ADBJob, data UpdateDataResAdbProps) error {
	return database.DB.Model(&models.AdbResult{}).
		Where("id = ?", j.ID).
		Updates(map[string]interface{}{
			"status":      data.Status,
			"start_time":  data.StartedAt,
			"finish_time": data.FinishAt,
			"result":      data.Result,
		}).Error
}

const ServiceRunningADB string = "running-adb"

func adbWorker(id int, jobs <-chan job.ADBJob, sseManager *sse.SSEManager) {
	for j := range jobs {

		services.LogInfo(
			ServiceRunningADB,
			fmt.Sprintf("[ADB Worker %d] Processing Job %s", id, j.ID),
		)

		startedAt := time.Now().UTC()
		status, output := adb.RunJob(j, adb.DefaultConfigPath)

		finishAt := time.Now().UTC()

		err := updateDataResADB(j, UpdateDataResAdbProps{
			Status:    status,
			StartedAt: startedAt,
			FinishAt:  finishAt,
			Result:    output,
		})

		if err != nil {
			services.LogError(ServiceRunningADB, "failed update adb result")
		}

		services.LogInfo(
			ServiceRunningADB,
			fmt.Sprintf("[ADB Worker %d] Finished Job %s\n", id, j.ID),
		)

		res := sse.ServicesEvent{
			Type: ServiceRunningADB,
			Data: UpdateDataResAdbProps{
				Status:    status,
				StartedAt: startedAt,
				FinishAt:  finishAt,
				Result:    output,
				ID:        j.ID,
			},
		}

		if sseManager != nil {
			sseManager.Broadcast(sse.SSEChannelServices, res)
		}
	}
}
