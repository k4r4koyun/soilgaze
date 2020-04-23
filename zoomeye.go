package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

// Zoomeye struct that holds the API key
type Zoomeye struct {
	username string
	password string
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

func (z Zoomeye) acquireJWT() string {
	data := url.Values{}
	data.Set("username", z.username)
	data.Set("password", z.password)

	response, err := z.postRequest("https://api.zoomeye.org/user/login", data.Encode(), "")

	if err != nil {
		log.Fatal(err)
	}

	return gjson.Get(response, "access_token").String()
}

func (z Zoomeye) check(allHosts *[]HostStruct) {
	log.Println("================== ZOOMEYE ==================")

	log.Println("Zoomeye is not implemented yet...")
	return

	if z.username == "" || z.password == "" {
		log.Println("Zoomeye: One or more credential values are empty, will skip this resource!")
		return
	}

	zoomeyeJWT := z.acquireJWT()

	response, err := z.getRequest("https://api.zoomeye.org/resources-info", zoomeyeJWT)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(response)
}
