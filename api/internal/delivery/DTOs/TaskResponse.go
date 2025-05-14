package DTOs

import (
	"HITS_ToDoList_Tests/internal/domain/enums"
	"github.com/google/uuid"
	"time"
)

type TaskResponse struct {
	ID          uuid.UUID      `binding:"required" json:"id"`
	CreatedAt   time.Time      `binding:"required" json:"createdAt"`
	ChangedAt   *time.Time     `json:"changedAt"`
	Name        string         `binding:"required" json:"name"`
	Description *string        `json:"description"`
	Deadline    *time.Time     `json:"deadline"`
	Status      enums.Status   `binding:"required" json:"status"`
	Priority    enums.Priority `binding:"required" json:"priority"`
}
