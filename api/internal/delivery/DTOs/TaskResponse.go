package DTOs

import (
	"HITS_ToDoList_Tests/internal/domain/enums"
	"github.com/google/uuid"
	"time"
)

type TaskResponse struct {
	ID          uuid.UUID `binding:"required"`
	CreatedAt   time.Time `binding:"required"`
	ChangedAt   *time.Time
	Name        string `binding:"required"`
	Description *string
	Deadline    *time.Time
	Status      enums.Status   `binding:"required"`
	Priority    enums.Priority `binding:"required"`
}
