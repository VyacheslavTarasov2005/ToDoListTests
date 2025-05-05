package DTOs

import (
	"HITS_ToDoList_Tests/internal/domain/enums"
	"time"
)

type UpdateTaskRequest struct {
	Name        *string         `json:"name"`
	Description *string         `json:"description"`
	Deadline    *time.Time      `json:"deadline"`
	Priority    *enums.Priority `json:"priority"`
}
