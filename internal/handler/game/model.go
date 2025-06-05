package handler

type GetGamesByIDsRequest struct {
	IDs []int64 `json:"ids" binding:"required,min=1"`
}
