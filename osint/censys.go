package osint

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// Censys struct that holds the API key
type Censys struct {
	APIKey string
}

type censysQueryBody struct {
	Query   string   `json:"query"`
	Page    int      `json:"page"`
	Fields  []string `json:"fields"`
	Flatten bool     `json:"flatten"`
}

// WEB FUNCTIONS =============================================

func (c Censys) basicAuth() string {
	auth := c.APIKey
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
func (c Censys) redirectPolicy(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+c.basicAuth())
	return nil
}

func (c Censys) getRequest(address string) (string, error) {
	client := &http.Client{
		CheckRedirect: c.redirectPolicy,
	}

	req, err := http.NewRequest("GET", address, nil)

	if err != nil {
		fmt.Printf("An error occured while preparing Censys request")
		return "", err
	}

	req.Header.Add("Authorization", "Basic "+c.basicAuth())

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("An error occured while sending Censys request")
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("An error occured while reading Censys response")
		return "", err
	}

	return string(body), nil
}

func (c Censys) postRequest(address string, payload string) (string, error) {
	client := &http.Client{
		CheckRedirect: c.redirectPolicy,
	}

	req, err := http.NewRequest("POST", address, strings.NewReader(payload))

	if err != nil {
		fmt.Printf("An error occured while preparing Censys request")
		return "", err
	}

	req.Header.Add("Authorization", "Basic "+c.basicAuth())

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("An error occured while sending Censys request")
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("An error occured while reading Censys response")
		return "", err
	}

	return string(body), nil
}

// MAIN FUNCTIONS =============================================

func (c Censys) queryQuota() (int, int, error) {
	response, err := c.getRequest("https://censys.io/api/v1/account")

	if err != nil {
		return 0, 0, err
	}

	usedQueries := int(gjson.Get(response, "quota").Get("used").Int())
	allowance := int(gjson.Get(response, "quota").Get("allowance").Int())

	return usedQueries, allowance, nil
}

func (c Censys) searchPorts(allHosts *[]HostStruct, filterString string) error {
	queryBody := censysQueryBody{
		filterString,
		1, // Pages start from 1
		[]string{"ip", "ports"},
		true,
	}

	queryString, err := json.Marshal(&queryBody)
	if err != nil {
		return err
	}

	response, err := c.postRequest("https://censys.io/api/v1/search/ipv4", string(queryString))

	if err != nil {
		return err
	}

	results := gjson.Get(response, "results")
	for _, result := range results.Array() {
		resultIP := result.Get("ip").String()

		for index, hostStruct := range *allHosts {
			if hostStruct.IPAddress == resultIP {
				var portArray []int
				for _, port := range result.Get("ports").Array() {
					portInt, _ := strconv.Atoi(port.String())
					portArray = append(portArray, portInt)
				}

				sort.Ints(portArray)
				(*allHosts)[index].OSINTResponse.Censys.OpenPorts = portArray
			}
		}
	}

	return nil
}

// Check is the interface generic method
func (c Censys) Check(allHosts *[]HostStruct) {
	log.Println("================== CENSYS ==================")

	if c.APIKey == "" {
		log.Println("Censys: API key value is empty, will skip this resource!")
		return
	}

	usedQueries, allowance, err := c.queryQuota()

	if err != nil {
		log.Println("Could not query quota, will skip this resource.")
		return
	}

	log.Println("Censys used queries: " + strconv.Itoa(usedQueries))
	log.Println("Censys allowance: " + strconv.Itoa(allowance))

	limitCalculation := (allowance - usedQueries) - len(*allHosts)

	if limitCalculation < 0 {
		log.Println("Censys: Remaining query allowance is not enough for the host list, skipping this OSINT resource...")
		return
	}

	queryString := ""
	batchCount := 25

	for i := 0; i < len(*allHosts); i += batchCount {
		queryString = ""
		end := i + batchCount

		if end > len(*allHosts) {
			end = len(*allHosts)
		}

		currentBatch := (*allHosts)[i:end]

		for _, hostStruct := range currentBatch {
			if queryString == "" {
				queryString = "ip:" + hostStruct.IPAddress
			} else {
				queryString += " OR ip:" + hostStruct.IPAddress
			}
		}

		err := c.searchPorts(allHosts, queryString)
		if err != nil {
			log.Println("Could not query a batch of hosts.")
		}

		time.Sleep(3 * time.Second)
	}
}
