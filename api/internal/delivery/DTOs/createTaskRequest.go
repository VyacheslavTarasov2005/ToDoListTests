package DTOs

import (
	"HITS_ToDoList_Tests/internal/domain/enums"
	"time"
)

type CreateTaskRequest struct {
	Name        *string         `json:"name" binding:"required"`
	Description *string         `json:"description"`
	Deadline    *time.Time      `json:"deadline"`
	Priority    *enums.Priority `json:"priority"`
}
