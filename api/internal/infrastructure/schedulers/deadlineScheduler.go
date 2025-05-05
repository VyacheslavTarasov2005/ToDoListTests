package schedulers

import (
	"HITS_ToDoList_Tests/internal/application/interfaces"
	"time"
)

func StartTasksDeadlineScheduling(service interfaces.TasksService, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			service.UpdateTaskStatuses()
		}
	}()
}
