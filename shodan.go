package main

import (
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

// Shodan struct that holds the API key
type Shodan struct {
	apiKey string
}

func (s Shodan) check(allHosts *[]HostStruct) {
	if s.apiKey == "" {
		log.Println("Shodan: API key value is empty, will skip this resource!")
		return
	}

	for index := range *allHosts {
		response, err := sendGETRequest("https://api.shodan.io/shodan/host/" + (*allHosts)[index].IPAddress + "?key=" + s.apiKey)

		if err != nil {
			log.Fatal("An error happened while checking Shodan results.")
		} else {
			value := gjson.Get(response, "ports")
			log.Println("Shodan: Open ports for " + (*allHosts)[index].IPAddress + ": " + value.String())

			// println(value.Array()[0].String())

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
