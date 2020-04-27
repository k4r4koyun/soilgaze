package osint

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

// Shodan struct that holds the API key
type Shodan struct {
	APIKey string
}

func (s Shodan) sendGETRequest(address string) (string, error) {
	resp, err := http.Get(address)

	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", errors.New("Non-2XX/3XX HTTP code received. ")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Check is the interface generic method
func (s Shodan) Check(allHosts *[]HostStruct) {
	log.Println("================== SHODAN ==================")

	if s.APIKey == "" {
		log.Println("Shodan: API key value is empty, will skip this resource!")
		return
	}

	for index := range *allHosts {
		response, err := s.sendGETRequest("https://api.shodan.io/shodan/host/" + (*allHosts)[index].IPAddress + "?key=" + s.APIKey)

		if err != nil {
			log.Println(fmt.Errorf("An error occured while checking results: %v", err))
		} else {
			value := gjson.Get(response, "ports")
			log.Println("Open ports for " + (*allHosts)[index].IPAddress + ": " + value.String())

			var portArray []int
			for _, port := range value.Array() {
				portInt, _ := strconv.Atoi(port.String())
				portArray = append(portArray, portInt)
			}

			sort.Ints(portArray)
			(*allHosts)[index].OSINTResponse.Shodan.OpenPorts = portArray
		}

		time.Sleep(2 * time.Second)
	}
}
