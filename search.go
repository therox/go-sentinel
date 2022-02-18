package sentinel

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (sc *SentinelClient) Query(params SearchParameters) error {
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
		return err
	}
	req.SetBasicAuth(sc.user, sc.password)
	resp, err := sc.httpClient.Do(req)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(bs))

	return nil
}
