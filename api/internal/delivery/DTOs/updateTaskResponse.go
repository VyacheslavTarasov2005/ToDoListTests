package DTOs

import "HITS_ToDoList_Tests/internal/domain/enums"

type UpdateTaskResponse struct {
	Status enums.Status `json:"status" binding:"required"`
}
