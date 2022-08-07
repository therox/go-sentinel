package sentinel

// Main engine interface. It must provide functions like search, download file from backend
type dlEngine interface {
	// SearchDataset(datasetName string)
	Download(productID string, dst string) error
	IsOnline(productID string) (bool, error)
}
