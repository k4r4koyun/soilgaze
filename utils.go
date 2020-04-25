package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
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

// FILES =====================================================

func loadConfig(path string) (*APIKeys, error) {
	apiKeys := new(APIKeys)
	var configFile *os.File
	var err error

	if path == "" {
		configFile, err = os.Open("config.yaml")
	} else {
		configFile, err = os.Open(path)
	}

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
func loadEnvironment() (*APIKeys, error) {
	apiKeys := new(APIKeys)

	apiKeys.Shodan = os.Getenv("SG_SHODAN")
	apiKeys.BinaryEdge = os.Getenv("SG_BINARYEDGE")
	apiKeys.Censys = os.Getenv("SG_CENSYS")
	apiKeys.ZoomEyeU = os.Getenv("SG_ZOOMEYE_U")
	apiKeys.ZoomEyeP = os.Getenv("SG_ZOOMEYE_P")
	apiKeys.Onyphe = os.Getenv("SG_ONYPHE")
	apiKeys.Spyse = os.Getenv("SG_SPYSE")

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
