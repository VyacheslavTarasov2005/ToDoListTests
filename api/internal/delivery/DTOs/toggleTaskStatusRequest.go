package DTOs

type ToggleTaskStatusRequest struct {
	IsDone *bool `binding:"required"`
}
