package osint

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

// Binaryedge struct that holds the API key
type Binaryedge struct {
	APIKey string
}

// WEB FUNCTIONS =============================================

func (b Binaryedge) redirectPolicy(req *http.Request, via []*http.Request) error {
	req.Header.Add("X-Key", b.APIKey)
	return nil
}

func (b Binaryedge) getRequest(address string) (string, error) {
	client := &http.Client{
		CheckRedirect: b.redirectPolicy,
	}

	req, err := http.NewRequest("GET", address, nil)

	if err != nil {
		fmt.Printf("An error occured while preparing Binaryedge request")
		return "", err
	}

	req.Header.Add("X-Key", b.APIKey)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("An error occured while sending Binaryedge request")
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("An error occured while reading Binaryedge response")
		return "", err
	}

	return string(body), nil
}

// MAIN FUNCTIONS =============================================

func (b Binaryedge) queryQuota() (int, int, error) {
	response, err := b.getRequest("https://api.binaryedge.io/v2/user/subscription")

	if err != nil {
		return 0, 0, err
	}

	queriesLeft := int(gjson.Get(response, "requests_left").Int())
	apiPlan := int(gjson.Get(response, "requests_plan").Int())

	return queriesLeft, apiPlan, nil
}

func (b Binaryedge) searchPorts(allHosts *[]HostStruct, queryString string) error {
	response, err := b.getRequest("https://api.binaryedge.io/v2/query/search?only_ips=1&query=" + queryString)

	if err != nil {
		return err
	}

	results := gjson.Get(response, "events")
	for _, result := range results.Array() {
		resultIP := result.Get("ip").String()

		for index, hostStruct := range *allHosts {
			if hostStruct.IPAddress == resultIP {
				currentPorts := (*allHosts)[index].OSINTResponse.Binaryedge.OpenPorts
				currentPorts = append(currentPorts, int(result.Get("port").Int()))

				sort.Ints(currentPorts)
				(*allHosts)[index].OSINTResponse.Binaryedge.OpenPorts = currentPorts
			}
		}
	}

	return nil
}

// Check is the interface generic method
func (b Binaryedge) Check(allHosts *[]HostStruct) {
	log.Println("================== BINARYEDGE ==================")

	if b.APIKey == "" {
		log.Println("Binaryedge: API key value is empty, will skip this resource!")
		return
	}

	queriesLeft, apiPlan, err := b.queryQuota()

	if err != nil {
		log.Println("Could not query quota, will skip this resource.")
		return
	}

	log.Println("Binaryedge left queries: " + strconv.Itoa(queriesLeft))
	log.Println("Binaryedge allowance: " + strconv.Itoa(apiPlan))

	hostSize := float64(len(*allHosts)) / 5
	limitCalculation := queriesLeft - int(math.Ceil(hostSize))

	if limitCalculation < 0 {
		log.Println("Binaryedge: Remaining query allowance is not enough for the host list, skipping this OSINT resource...")
		return
	}
	log.Println("Number of queries left are sufficient, beginning operation.")

	queryString := ""
	batchCount := 5

	for i := 0; i < len(*allHosts); i += batchCount {
		queryString = ""
		end := i + batchCount

		if end > len(*allHosts) {
			end = len(*allHosts)
		}

		currentBatch := (*allHosts)[i:end]

		for _, hostStruct := range currentBatch {
			if queryString == "" {
				queryString = "(ip:" + hostStruct.IPAddress
			} else {
				queryString += "%20OR%20ip:" + hostStruct.IPAddress
			}
		}

		queryString += ")"
		queryString += "%20AND%20protocol:tcp"
		queryString += "%20AND%20type:service-simple"

		err := b.searchPorts(allHosts, queryString)
		if err != nil {
			log.Println("Could not query a batch of hosts.")
		}

		time.Sleep(2 * time.Second)
	}
}
