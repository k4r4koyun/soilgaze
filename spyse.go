package main

import "log"

// Spyse struct that holds the API key
type Spyse struct {
	apiKey string
}

func (s Spyse) check(allHosts *[]HostStruct) {
	log.Println("================== SPYSE ==================")

	log.Println("Spyse is not implemented yet...")

	if s.apiKey == "" {
		log.Println("Spyse: API key value is empty, will skip this resource!")
		return
	}

}
