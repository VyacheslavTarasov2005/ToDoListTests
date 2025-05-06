package DTOs

import (
	"HITS_ToDoList_Tests/internal/domain/enums"
	"time"
)

type CreateTaskRequest struct {
	Name        *string `binding:"required"`
	Description *string
	Deadline    *time.Time
	Priority    *enums.Priority `binding:"omitempty,oneof=Low Medium High Critical" msg:"Incorrect Priority"`
}
