package osint

import "log"

// Spyse struct that holds the API key
type Spyse struct {
	APIKey string
}

// Check is the interface generic method
func (s Spyse) Check(allHosts *[]HostStruct) {
	log.Println("================== SPYSE ==================")

	log.Println("Spyse is not implemented yet...")

	if s.APIKey == "" {
		log.Println("Spyse: API key value is empty, will skip this resource!")
		return
	}

}
