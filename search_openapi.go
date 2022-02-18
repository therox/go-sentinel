package sentinel

type Platform string

const (
	PlanformSentinel1          Platform = "Sentinel-1"
	PlanformSentinel2          Platform = "Sentinel-2"
	PlanformSentinel3          Platform = "Sentinel-3"
	PlanformSentinel5Precursor Platform = "Sentinel-5 Precursor"
)

type SearchParameters struct {
	Platforms []Platform
}
