package sentinel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func (sc *SentinelClient) Query(params SearchParameters) (QueryResponse, error) {
	var qr QueryResponse
	fmt.Printf("%+v\n", params)

	searchURL := sc.searchURL
	if len(params.Platforms) > 0 {
		searchURL += "("
		for i := range params.Platforms {
			if i > 0 {
				searchURL += " OR "
			}
			searchURL += fmt.Sprintf("platformname:'%s'", params.Platforms[i])
		}
		searchURL += ")"
	}
	searchURL += "&format=json&rows=" + strconv.Itoa(sc.rows)
	fmt.Printf("%+v\n", searchURL)

	// ======= requesting data =====
	req, err := http.NewRequest(http.MethodGet, searchURL, nil)
	if err != nil {
		return qr, err
	}
	req.SetBasicAuth(sc.user, sc.password)
	resp, err := sc.httpClient.Do(req)
	if err != nil {
		return qr, err
	}
	req.Header.Add("Content-Type", "application/json")
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return qr, err
	}
	return processQueryResponse(bs)
}

func processQueryResponse(bs []byte) (QueryResponse, error) {
	// fmt.Printf("%s", bs)
	var res QueryResponse
	err := json.Unmarshal(bs, &res)
	if err != nil {
		return res, err
	}
	// fmt.Printf("%+v\n", res.Feed.Entries[0])
	// тут мы анмаршаллим Str, Date, Int, Double
	for i := range res.Feed.Entries {
		strList, err := unpackTypedCommonData(res.Feed.Entries[i].Str)
		if err != nil {
			return res, err
		}
		for j := range strList {
			switch strList[j].Name {
			case "sensoroperationalmode":
				res.Feed.Entries[i].SensorOperationalMode = strList[j].Content

			case "gmlfootprint":
				res.Feed.Entries[i].GMLFootprint = strList[j].Content

			case "footprint":
				res.Feed.Entries[i].Footprint = strList[j].Content

			case "tileid":
				res.Feed.Entries[i].TileId = strList[j].Content

			case "hv_order_tileid":
				res.Feed.Entries[i].HVOrderTileid = strList[j].Content

			case "format":
				res.Feed.Entries[i].Format = strList[j].Content

			case "processingbaseline":
				res.Feed.Entries[i].ProcessingBaseline = strList[j].Content

			case "platformname":
				res.Feed.Entries[i].PlatformName = strList[j].Content

			case "filename":
				res.Feed.Entries[i].FileName = strList[j].Content

			case "instrumentname":
				res.Feed.Entries[i].InstrumentName = strList[j].Content

			case "instrumentshortname":
				res.Feed.Entries[i].InstrumentShortName = strList[j].Content

			case "size":
				res.Feed.Entries[i].Size = strList[j].Content

			case "s2datatakeid":
				res.Feed.Entries[i].S2DataTakeID = strList[j].Content

			case "producttype":
				res.Feed.Entries[i].ProductType = strList[j].Content

			case "platformidentifier":
				res.Feed.Entries[i].PlatformIdentifier = strList[j].Content

			case "level1cpdiidentifier":
				res.Feed.Entries[i].Level1CPDIdentifier = strList[j].Content

			case "orbitdirection":
				res.Feed.Entries[i].OrbitDirection = strList[j].Content

			case "platformserialidentifier":
				res.Feed.Entries[i].PlatformSerialIdentifier = strList[j].Content

			case "processinglevel":
				res.Feed.Entries[i].ProcessingLevel = strList[j].Content

			case "datastripidentifier":
				res.Feed.Entries[i].DataStripIdentifier = strList[j].Content

			case "granuleidentifier":
				res.Feed.Entries[i].GranuleIdentifier = strList[j].Content

			case "identifier":
				res.Feed.Entries[i].Identifier = strList[j].Content

			case "uuid":
				res.Feed.Entries[i].UUID = strList[j].Content

			}
		}

		intList, err := unpackTypedCommonData(res.Feed.Entries[i].Int)
		if err != nil {
			return res, err
		}
		for j := range intList {
			switch intList[j].Name {
			case "orbitnumber":
				res.Feed.Entries[i].OrbitNumber, _ = strconv.Atoi(intList[j].Content)

			case "relativeorbitnumber":
				res.Feed.Entries[i].RelativeOrbitNumber, _ = strconv.Atoi(intList[j].Content)

			}
		}

		doubleList, err := unpackTypedCommonData(res.Feed.Entries[i].Double)
		if err != nil {
			return res, err
		}
		for j := range doubleList {
			switch doubleList[j].Name {
			case "cloudcoverpercentage":
				res.Feed.Entries[i].CloudCoverPercentage, _ = strconv.ParseFloat(doubleList[j].Content, 64)

			case "illuminationazimuthangle":
				res.Feed.Entries[i].IlluminationAzimuthAngle, _ = strconv.ParseFloat(doubleList[j].Content, 64)
			case "illuminationzenithangle":
				res.Feed.Entries[i].IlluminationZenithAngle, _ = strconv.ParseFloat(doubleList[j].Content, 64)
			case "vegetationpercentage":
				res.Feed.Entries[i].VegetationPercentage, _ = strconv.ParseFloat(doubleList[j].Content, 64)
			case "notvegetatedpercentage":
				res.Feed.Entries[i].NotVegetatedPercentage, _ = strconv.ParseFloat(doubleList[j].Content, 64)
			case "waterpercentage":
				res.Feed.Entries[i].WaterPercentage, _ = strconv.ParseFloat(doubleList[j].Content, 64)
			case "unclassifiedpercentage":
				res.Feed.Entries[i].UnclassifiedPercentage, _ = strconv.ParseFloat(doubleList[j].Content, 64)
			case "mediumprobacloudspercentage":
				res.Feed.Entries[i].MediumProbaCloudsPercentage, _ = strconv.ParseFloat(doubleList[j].Content, 64)
			case "highprobacloudspercentage":
				res.Feed.Entries[i].HighProbaCloudsPercentage, _ = strconv.ParseFloat(doubleList[j].Content, 64)
			case "snowicepercentage":
				res.Feed.Entries[i].SnowIcePercentage, _ = strconv.ParseFloat(doubleList[j].Content, 64)

			}
		}

		dateList, err := unpackTypedCommonData(res.Feed.Entries[i].Date)
		if err != nil {
			return res, err
		}
		for j := range dateList {
			switch dateList[j].Name {
			case "datatakesensingstart":
				res.Feed.Entries[i].DataTakeSensingStart, err = time.Parse(time.RFC3339, dateList[j].Content)
			case "generationdate":
				res.Feed.Entries[i].GenerationDate, err = time.Parse(time.RFC3339, dateList[j].Content)
			case "beginposition":
				res.Feed.Entries[i].BeginPosition, err = time.Parse(time.RFC3339, dateList[j].Content)
			case "endposition":
				res.Feed.Entries[i].EndPosition, err = time.Parse(time.RFC3339, dateList[j].Content)
			case "ingestiondate":
				res.Feed.Entries[i].IngestionDate, err = time.Parse(time.RFC3339, dateList[j].Content)
			}
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	fmt.Println(res.Feed.Entries[1].TileId, res.Feed.Entries[0].FileName)
	fmt.Println(string(res.Feed.Entries[0].Str))
	return res, nil
}

func unpackTypedCommonData(bs []byte) (res []TypedCommonData, err error) {
	if len(bs) == 0 {
		return res, nil
	}
	switch bs[0] {
	case '{':
		var tempRes TypedCommonData
		err = json.Unmarshal(bs, &tempRes)
		return []TypedCommonData{tempRes}, err
	case '[':
		err = json.Unmarshal(bs, &res)
		return
	}
	return
}
