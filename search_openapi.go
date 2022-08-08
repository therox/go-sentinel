package sentinel

import (
	"encoding/json"
	"time"
)

type Platform string

type AreaRelation string

const (
	AreaRelationIntersects AreaRelation = "Intersects"
	AreaRelationContains   AreaRelation = "Contains"
	AreaRelationIsWithin   AreaRelation = "IsWithin"
)

const (
	PlanformSentinel1          Platform = "Sentinel-1"
	PlanformSentinel2          Platform = "Sentinel-2"
	PlanformSentinel3          Platform = "Sentinel-3"
	PlanformSentinel5Precursor Platform = "Sentinel-5 Precursor"
)

type SearchParameters struct {
	Platforms            []Platform
	Footprint            string
	AreaRelation         AreaRelation
	TileIDs              []string   // tileid:37UCU
	BeginDate            time.Time  // Ingestion date from [2014-01-01T00:00:00.000Z TO NOW]
	EndDate              *time.Time // Ingestion date to, NOW if not set
	ProductTypes         []string
	Filenames            []string
	CloudCoverPercentage []string // [0 TO 100]
}

type TypedCommonData struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type QueryEntryResponse struct {
	Title string `json:"title"`
	Link  []struct {
		Rel  string `json:"rel"`
		HREF string `json:"href"`
	} `json:"link"`
	ID                          string `json:"id"`
	Summary                     string `json:"summary"`
	OnDemandStr                 string `json:"ondemand"`
	OnDemand                    bool
	Date                        json.RawMessage `json:"date"`
	Int                         json.RawMessage `json:"int"`
	Double                      json.RawMessage `json:"double"`
	Str                         json.RawMessage `json:"str"`
	DataTakeSensingStart        time.Time
	GenerationDate              time.Time
	BeginPosition               time.Time
	EndPosition                 time.Time
	IngestionDate               time.Time
	OrbitNumber                 int
	RelativeOrbitNumber         int
	CloudCoverPercentage        float64
	SensorOperationalMode       string
	GMLFootprint                string
	Footprint                   string
	Level1CPDIdentifier         string
	TileId                      string
	HVOrderTileid               string
	Format                      string
	ProcessingBaseline          string
	PlatformName                string
	FileName                    string
	InstrumentName              string
	InstrumentShortName         string
	Size                        string
	S2DataTakeID                string
	ProductType                 string
	PlatformIdentifier          string
	OrbitDirection              string
	PlatformSerialIdentifier    string
	ProcessingLevel             string
	DataStripIdentifier         string
	GranuleIdentifier           string
	Identifier                  string
	UUID                        string
	IlluminationAzimuthAngle    float64
	IlluminationZenithAngle     float64
	VegetationPercentage        float64
	NotVegetatedPercentage      float64
	WaterPercentage             float64
	UnclassifiedPercentage      float64
	MediumProbaCloudsPercentage float64
	HighProbaCloudsPercentage   float64
	SnowIcePercentage           float64
}

type QueryResponse struct {
	Feed struct {
		OpenSearch string    `json:"xmlns:opensearch"`
		NS         string    `json:"xmlns"`
		Title      string    `json:"title"`
		Subtitle   string    `json:"subtitle"`
		Updated    time.Time `json:"updated"`
		Author     struct {
			Name string `json:"name"`
		} `json:"athor"`
		ID              string `json:"id"`
		TotalResultsStr string `json:"opensearch:totalResults"`
		TotalResults    int
		StartIndexStr   string `json:"opensearch:startIndex"`
		StartIndex      int
		ItemsPerPageStr string `json:"opensearch:itemsPerPage"`
		ItemsPerPage    int
		Query           struct {
			Role        string `json:"role"`
			SearchTerms string `json:"searchTerms"`
			StartPage   string `json:"startPage"`
		} `json:"opensearch:Query"`
		Link []struct {
			Rel  string `json:"rel"`
			Type string `json:"type"`
			HREF string `json:"href"`
		}
		Entries    []QueryEntryResponse
		EntriesRaw json.RawMessage `json:"entry"`
	} `json:"feed"`
}

func (qer *QueryEntryResponse) GetID() string {
	return qer.ID
}
