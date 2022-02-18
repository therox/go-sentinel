package sentinel

// Main engine interface. It must provide functions like search, download file from backend
type SentinelEngine interface {
	SearchDataset(datasetName string)
	Download(datasetName string)
}
