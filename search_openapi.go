package sentinel

type platformNameType string

const (
	Sentinel1          platformNameType = "Sentinel-1"
	Sentinel2          platformNameType = "Sentinel-2"
	Sentinel3          platformNameType = "Sentinel-3"
	Sentinel5Precursor platformNameType = "Sentinel-5 Precursor"
)

type OpenAPISearchParams struct {
	PlatformName platformNameType
}
