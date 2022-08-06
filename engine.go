package sentinel

// Main engine interface. It must provide functions like search, download file from backend
type dlEngine interface {
	// SearchDataset(datasetName string)
	Download(ProductID string, dst string) error
}
