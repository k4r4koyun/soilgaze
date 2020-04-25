package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

// Onyphe struct that holds the API key
type Onyphe struct {
	apiKey string
}

// WEB FUNCTIONS =============================================

func (o Onyphe) redirectPolicy(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "apikey "+o.apiKey)
	return nil
}

func (o Onyphe) getRequest(address string) (string, error) {
	client := &http.Client{
		CheckRedirect: o.redirectPolicy,
	}

	req, err := http.NewRequest("GET", address, nil)

	if err != nil {
		fmt.Printf("An error occured while preparing Onyphe request")
		return "", err
	}

	req.Header.Add("Authorization", "apikey "+o.apiKey)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("An error occured while sending Onyphe request")
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("An error occured while reading Onyphe response")
		return "", err
	}

	return string(body), nil
}

// MAIN FUNCTIONS =============================================

// Onyphe returns multiple entries for the same port on different dates
func (o Onyphe) unique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (o Onyphe) queryQuota() (int, error) {
	response, err := o.getRequest("https://www.onyphe.io/api/v2/user")

	if err != nil {
		return 0, err
	}

	firstResult := gjson.Get(response, "results.0")
	credits := int(firstResult.Get("credits").Int())

	return credits, nil
}

func (o Onyphe) check(allHosts *[]HostStruct) {
	log.Println("================== ONYPHE ==================")

	if o.apiKey == "" {
		log.Println("API key value is empty, will skip this resource!")
		return
	}

	credits, err := o.queryQuota()
	if err != nil {
		log.Println("Could not query quota, will skip this resource.")
		return
	}

	log.Println("Remaining credits: " + strconv.Itoa(credits))

	limitCalculation := credits - len(*allHosts)

	if limitCalculation < 0 {
		log.Println("Remaining query allowance is not enough for the host list, skipping this OSINT resource...")
		return
	}

	for index := range *allHosts {
		response, err := o.getRequest("https://www.onyphe.io/api/v2/simple/synscan/" + (*allHosts)[index].IPAddress)

		if err != nil {
			log.Println("An error happened while checking a single host.")
		} else {
			var portArray []int
			for _, singleHost := range gjson.Get(response, "results").Array() {
				portArray = append(portArray, int(singleHost.Get("port").Int()))
			}

			sort.Ints(portArray)
			(*allHosts)[index].OSINTResponse.Onyphe.OpenPorts = o.unique(portArray)
		}

		time.Sleep(2 * time.Second)
	}
}
