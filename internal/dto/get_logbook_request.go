package dto

type GetLogbookRequest struct {
	Start *int64 `json:"start"`
	End   *int64 `json:"end"`
}
