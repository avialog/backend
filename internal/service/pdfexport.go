
type ExportService interface {
	ExportLogbook(userID string) (string, error)
}
