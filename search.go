package sentinel

import "fmt"

func (sc *SentinelClient) SearchOData(params OpenAPISearchParams) error {
	fmt.Printf("%+v\n", params)

	return nil
}
