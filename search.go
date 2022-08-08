package sentinel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (sc *SentinelClient) Query(params SearchParameters) (QueryResponse, error) {
	// fmt.Printf("%+v\n", params)

	// searchURL := sc.searchURL
	urlParams := ""

	paramList := make([]string, 0)

	if len(params.Platforms) > 0 {
		innerParamList := make([]string, len(params.Platforms))
		for i := range params.Platforms {
			innerParamList[i] = fmt.Sprintf("platformname:'%s'", params.Platforms[i])
		}
		paramList = append(paramList, fmt.Sprintf("(%s)", strings.Join(innerParamList, " OR ")))
	}

	if len(params.TileIDs) > 0 {
		innerParamList := make([]string, len(params.TileIDs))
		for i := range params.TileIDs {
			innerParamList[i] = fmt.Sprintf("tileid:%s", params.TileIDs[i])
		}
		paramList = append(paramList, fmt.Sprintf("(%s)", strings.Join(innerParamList, " OR ")))
	}

	if len(params.Filenames) > 0 {
		innerParamList := make([]string, len(params.Filenames))
		for i := range params.Filenames {
			innerParamList[i] = fmt.Sprintf("filename:%s", params.Filenames[i])
		}
		paramList = append(paramList, fmt.Sprintf("(%s)", strings.Join(innerParamList, " OR ")))
	}

	if len(params.ProductTypes) > 0 {
		innerParamList := make([]string, len(params.ProductTypes))
		for i := range params.ProductTypes {
			innerParamList[i] = fmt.Sprintf("producttype:%s", params.ProductTypes[i])
		}
		paramList = append(paramList, fmt.Sprintf("(%s)", strings.Join(innerParamList, " OR ")))
	}

	if params.EndDate != nil {
		// [2014-01-01T00:00:00.000Z TO NOW]]
		paramList = append(paramList, fmt.Sprintf("beginposition:[%s TO %s]", params.BeginDate.Format("2006-01-02T15:04:05.000Z"), params.EndDate.Format("2006-01-02T15:04:05.000Z")))
	} else {
		paramList = append(paramList, fmt.Sprintf("beginposition:[%s TO NOW]", params.BeginDate.Format("2006-01-02T15:04:05.000Z")))
	}

	if params.Footprint != "" {
		areaRelation := AreaRelationIntersects
		if params.AreaRelation != "" {
			isFound := false
			for _, ar := range []AreaRelation{AreaRelationIntersects, AreaRelationContains, AreaRelationIsWithin} {
				if strings.EqualFold(strings.ToLower(string(params.AreaRelation)), strings.ToLower(string(ar))) {
					isFound = true
					areaRelation = ar
					break
				}
			}
			if !isFound {
				return QueryResponse{}, fmt.Errorf("incorrect AOI relation provided: %s", params.AreaRelation)
			}
		}
		paramList = append(paramList, fmt.Sprintf("footprint:\"%s(%s)\"", areaRelation, params.Footprint))
	}

	if len(params.CloudCoverPercentage) > 0 {
		innerParamList := make([]string, len(params.CloudCoverPercentage))
		for i := range params.CloudCoverPercentage {
			innerParamList[i] = fmt.Sprintf("cloudcoverpercentage:%s", params.CloudCoverPercentage[i])
		}
		paramList = append(paramList, fmt.Sprintf("(%s)", strings.Join(innerParamList, " OR ")))
	}

	//  Union of params
	urlParams += strings.Join(paramList, " AND ")
	// fmt.Printf("> %+v\n", urlParams)

	urlParams = url.QueryEscape(urlParams)
	urlParams += fmt.Sprintf("&format=json&rows=%d", sc.rows)

	return sc.doQuery(fmt.Sprintf("%s%s", sc.searchURL, urlParams))
}

func (sc *SentinelClient) doQuery(queryURL string) (QueryResponse, error) {
	fmt.Printf("%s\n", queryURL)
	var qr QueryResponse
	// ======= requesting first data page =====
	req, err := http.NewRequest(http.MethodGet, queryURL, nil)
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
	qr, err = processQueryResponse(bs)
	if err != nil {
		return qr, err
	}
	fmt.Printf("Found %d results\n", qr.Feed.TotalResults)

	offset := sc.rows
	for {
		if len(qr.Feed.Entries) < qr.Feed.TotalResults {

			nextURL := queryURL + fmt.Sprintf("&start=%d", offset)

			req, err := http.NewRequest(http.MethodGet, nextURL, nil)
			if err != nil {
				return qr, err
			}
			req.SetBasicAuth(sc.user, sc.password)
			req.Header.Add("Content-Type", "application/json")
			resp, err := sc.httpClient.Do(req)
			if err != nil {
				return qr, err
			}
			bs, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				resp.Body.Close()
				return qr, err
			}
			fmt.Println(resp.StatusCode)
			if resp.StatusCode != 200 {
				fmt.Printf("Response error: %s. Repeating.", bs)
				resp.Body.Close()
				continue
			}

			tempQR, err := processQueryResponse(bs)
			if err != nil {
				resp.Body.Close()
				return qr, err
			}
			qr.Feed.Entries = append(qr.Feed.Entries, tempQR.Feed.Entries...)
			resp.Body.Close()
			offset += sc.rows
		} else {
			break
		}
	}
	// Repeat until we get TotalResults items

	return qr, nil
}

func processQueryResponse(bs []byte) (QueryResponse, error) {
	// fmt.Printf("%s", bs)
	var res QueryResponse
	err := json.Unmarshal(bs, &res)
	if err != nil {
		fmt.Println(string(bs))
		return res, err
	}

	res.Feed.TotalResults, _ = strconv.Atoi(res.Feed.TotalResultsStr)
	res.Feed.StartIndex, _ = strconv.Atoi(res.Feed.StartIndexStr)
	res.Feed.ItemsPerPage, _ = strconv.Atoi(res.Feed.ItemsPerPageStr)

	// fmt.Printf("%+v\n", res.Feed.Entries[0])
	// тут мы анмаршаллим Str, Date, Int, Double

	res.Feed.Entries, err = unpackQueryEntryResponse(res.Feed.EntriesRaw)
	if err != nil {
		return res, err
	}
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

func unpackQueryEntryResponse(bs []byte) (res []QueryEntryResponse, err error) {
	if len(bs) == 0 {
		return res, nil
	}
	switch bs[0] {
	case '{':
		var tempRes QueryEntryResponse
		err = json.Unmarshal(bs, &tempRes)
		return []QueryEntryResponse{tempRes}, err
	case '[':
		err = json.Unmarshal(bs, &res)
		return
	}
	return
}
