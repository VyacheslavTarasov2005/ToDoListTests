package DTOs

import (
	"HITS_ToDoList_Tests/internal/domain/enums"
	"github.com/google/uuid"
	"time"
)

type TaskResponse struct {
	ID          uuid.UUID      `json:"id" binding:"required"`
	CreatedAt   time.Time      `json:"createdAt" binding:"required"`
	ChangedAt   *time.Time     `json:"changedAt"`
	Name        string         `json:"name" binding:"required"`
	Description *string        `json:"description"`
	Deadline    *time.Time     `json:"deadline"`
	Status      enums.Status   `json:"status" binding:"required"`
	Priority    enums.Priority `json:"priority" binding:"required"`
}
