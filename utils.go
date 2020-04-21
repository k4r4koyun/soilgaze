package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

// DATA =====================================================

func isIpv4Net(host string) bool {
	return net.ParseIP(host) != nil
}

func extractIP(host string) string {
	ip, err := net.LookupIP(host)

	if err != nil {
		log.Println("Could not resolve: " + host)
		return ""
	}

	return ip[0].String()
}

func prepareHostStruct(hostFile []string, allHosts *[]HostStruct) {
	for _, host := range hostFile {
		if isIpv4Net(host) {
			shouldAdd := true

			for _, hostStruct := range *allHosts {
				if hostStruct.IPAddress == host {
					shouldAdd = false
				}
			}

			if shouldAdd {
				var hostStruct HostStruct
				hostStruct.IPAddress = host
				hostStruct.Hostname = []string{}

				*allHosts = append(*allHosts, hostStruct)
			} else {
				log.Println("Skipping duplicate IP on the host list: " + host)
			}
		} else {
			ipAddress := extractIP(host)

			if host == "" {
				log.Println("Skipping unresolved host")
				return
			}

			shouldAdd := true
			for _, hostStruct := range *allHosts {
				if hostStruct.IPAddress == ipAddress {
					log.Println("Domain resolved to the same IP: " + host)
					hostStruct.Hostname = append(hostStruct.Hostname, host)
					shouldAdd = false
				}
			}

			if shouldAdd {
				var hostStruct HostStruct
				hostStruct.IPAddress = ipAddress
				hostStruct.Hostname = []string{host}

				*allHosts = append(*allHosts, hostStruct)
			}
		}
	}
}

// NETWORK =====================================================

func sendGETRequest(address string) (string, error) {
	resp, err := http.Get(address)

	if err != nil {
		fmt.Printf("An error occured while sending request")
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("An error occured while reading body")
		return "", err
	}

	return string(body), nil
}

// FILES =====================================================

func loadConfig() (*APIKeys, error) {
	apiKeys := new(APIKeys)

	configFile, err := os.Open("config.yaml")
	if err != nil {
		return apiKeys, errors.New("Could not read the config file. ")
	}
	defer configFile.Close()

	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(apiKeys)
	if err != nil {
		return apiKeys, errors.New("Could not parse the config file. ")
	}

	return apiKeys, nil

}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
func writeStringToFile(line string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintln(w, line)

	return w.Flush()
}
