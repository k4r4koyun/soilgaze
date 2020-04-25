package main

import (
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

// Zoomeye struct that holds the API key
type Zoomeye struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (z Zoomeye) getRequest(address string, JWT string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", address, nil)

	if err != nil {
		fmt.Printf("An error occured while preparing Zoomeye request")
		return "", err
	}

	if JWT != "" {
		req.Header.Add("Authorization", "JWT "+JWT)
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("An error occured while sending Zoomeye request")
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("An error occured while reading Zoomeye response")
		return "", err
	}

	return string(body), nil
}

func (z Zoomeye) postRequest(address string, payload string, JWT string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", address, strings.NewReader(payload))

	if err != nil {
		fmt.Printf("An error occured while preparing Zoomeye request")
		return "", err
	}

	if JWT != "" {
		req.Header.Add("Authorization", "JWT "+JWT)
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("An error occured while sending Zoomeye request")
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("An error occured while reading Zoomeye response")
		return "", err
	}

	return string(body), nil
}

func (z Zoomeye) acquireJWT() (string, error) {
	credentials, err := json.Marshal(&z)
	if err != nil {
		return "", err
	}

	response, err := z.postRequest("https://api.zoomeye.org/user/login", string(credentials), "")
	if err != nil {
		return "", err
	}

	return gjson.Get(response, "access_token").String(), nil
}

func (z Zoomeye) check(allHosts *[]HostStruct) {
	log.Println("================== ZOOMEYE ==================")

	if z.Username == "" || z.Password == "" {
		log.Println("Zoomeye: One or more credential values are empty, will skip this resource!")
		return
	}

	zoomeyeJWT, err := z.acquireJWT()
	if err != nil {
		log.Println("Could not fetch a Zoomeye JWT, not continuing.")
		return
	}

	time.Sleep(2 * time.Second)

	response, err := z.getRequest("https://api.zoomeye.org/resources-info", zoomeyeJWT)
	if err != nil {
		log.Println("Could not fetch remaining queries, not continuing.")
		return
	}

	time.Sleep(2 * time.Second)

	remainingCredits := int(gjson.Get(response, "resources").Get("search").Int())
	log.Println("Zoomeye remaining queries: " + strconv.Itoa(remainingCredits))

	if remainingCredits-len(*allHosts) < 0 {
		log.Println("Zoomeye: Remaining query allowance is not enough for the host list, skipping this OSINT resource...")
		return
	}

	for index := range *allHosts {
		response, err = z.getRequest("https://api.zoomeye.org/host/search?query=ip:"+(*allHosts)[index].IPAddress, zoomeyeJWT)

		if err != nil {
			log.Println("An error happened while checking single host.")
		} else {
			matches := gjson.Get(response, "matches")

			var portArray []int
			for _, match := range matches.Array() {

				portInt := int(match.Get("portinfo").Get("port").Int())
				portArray = append(portArray, portInt)
			}

			sort.Ints(portArray)
			(*allHosts)[index].OSINTResponse.Zoomeye.OpenPorts = portArray
		}

		time.Sleep(2 * time.Second)
	}
}
