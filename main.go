package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
)

// OSINT bla bla
type OSINT interface {
	check(*[]HostStruct)
}

// APIKeys : A struct that holds API key values
type APIKeys struct {
	Shodan     string `yaml:"api_shodan"`
	BinaryEdge string `yaml:"api_binaryedge"`
	Censys     string `yaml:"api_censys"`
	ZoomEyeU   string `yaml:"api_zoomeye_u"`
	ZoomEyeP   string `yaml:"api_zoomeye_p"`
	Onyphe     string `yaml:"api_onyphe"`
	Spyse      string `yaml:"api_spyse"`
}

// HostStruct holds temporary values before converting to JSON
type HostStruct struct {
	IPAddress     string
	Hostname      []string
	OSINTResponse OSINTResponses
}

// OSINTResponses holds responses from OSINT resources
type OSINTResponses struct {
	Shodan     OSINTInfo
	Binaryedge OSINTInfo
	Censys     OSINTInfo
}

// OSINTInfo is a generic struct for OSINT information
type OSINTInfo struct {
	OpenPorts []int
}

func main() {
	var err error

	var allHosts []HostStruct

	var hostFile string
	var osintList string
	var outFile string

	flag.StringVar(&hostFile, "host-file", "", "File location to read target hosts")
	flag.StringVar(&osintList, "osint-list", "", "OSINT resources to gather information from. Example: --osint-list=shodan,binaryedge")
	flag.StringVar(&outFile, "out-file", "", "File to write the results in JSON format. If not given, results will only be printed to console.")
	flag.Parse()

	if hostFile == "" {
		log.Fatal("A list of hosts should be provided!")
	}

	var apiKeys *APIKeys
	apiKeys, err = loadConfig()

	if err != nil {
		log.Fatal(err)
	}

	hostsFileContents, err := readLines(hostFile)
	if err != nil {
		log.Fatal("An error occured while reading hosts file.")
	}

	prepareHostStruct(hostsFileContents, &allHosts)

	var shodan OSINT = Shodan{apiKeys.Shodan}
	var binaryedge OSINT = Binaryedge{apiKeys.BinaryEdge}
	var censys OSINT = Censys{apiKeys.Censys}
	var zoomeye OSINT = Zoomeye{apiKeys.ZoomEyeU, apiKeys.ZoomEyeP}
	var onyphe OSINT = Onyphe{apiKeys.Onyphe}
	var spyse OSINT = Spyse{apiKeys.Spyse}

	if osintList == "" {
		shodan.check(&allHosts)
		binaryedge.check(&allHosts)
		censys.check(&allHosts)
		// zoomeye.check(&allHosts)
		// onyphe.check(&allHosts)
		// spyse.check(&allHosts)
	} else {
		resources := strings.Split(osintList, ",")

		for _, name := range resources {
			if name == "shodan" {
				shodan.check(&allHosts)
			} else if name == "binaryedge" {
				binaryedge.check(&allHosts)
			} else if name == "censys" {
				censys.check(&allHosts)
			} else if name == "zoomeye" {
				zoomeye.check(&allHosts)
			} else if name == "onyphe" {
				onyphe.check(&allHosts)
			} else if name == "spyse" {
				spyse.check(&allHosts)
			}
		}
	}

	finalResult, err := json.Marshal(&allHosts)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(finalResult))

	if outFile != "" {
		if err := writeStringToFile(string(finalResult), outFile); err != nil {
			log.Fatalf("writeStringToFile: %s", err)
		}
	}
}
