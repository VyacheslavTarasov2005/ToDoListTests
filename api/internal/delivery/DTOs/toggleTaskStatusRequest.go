package DTOs

type ToggleTaskStatusRequest struct {
	IsDone bool `json:"isDone" binding:"required"`
}
