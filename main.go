package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ManEelyMan/OktaEnumerator/okta"
)

const Domain = "nydig.okta.com"

func main() {

	var oktaDomain string
	var oktaToken string
	if !getAndValidateArgs(&oktaDomain, &oktaToken) {
		usage()
		return
	}

	// Get our groups
	oktaGroups, err := getGroups(oktaDomain, oktaToken)
	fmt.Printf("%v, %v\n", oktaGroups, err)

	// TODO: fill in the groups' users (and apps?)
}

func getAndValidateArgs(oktaDomain *string, oktaToken *string) bool {

	numArgs := len(os.Args)

	if numArgs == 2 {
		*oktaDomain = os.Args[0]
		*oktaToken = os.Args[1]
	} else if numArgs == 0 {
		fmt.Println("No arguments given, looking for parameters from the environment...")
		*oktaDomain = getEnv("OKTA_DOMAIN", "") // Only use the fallback for dev. DO NOT CHECK IN SENSITIVE DATA!
		if len(*oktaDomain) == 0 {
			fmt.Println("No value given for OKTA_DOMAIN environment variable.")
			return false
		}
		*oktaToken = getEnv("OKTA_TOKEN", "") // Only use the fallback for dev. DO NOT CHECK IN SENSITIVE DATA!
		if len(*oktaToken) == 0 {
			fmt.Println("No value given for OKTA_TOKEN environment variable.")
			return false
		}
	} else {
		fmt.Println("Incorrect number of arguments received.")
		return false
	}
	return true
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		fmt.Printf("Found environment variable '%v'\n", key)
		return val
	} else {
		fmt.Printf("Environment Variable '%v' not found. Using fallback value of '%v'.\n", key, fallback)
		return fallback
	}
}

func usage() {
	// TODO:
}

func getGroups(oktaDomain string, oktaToken string) (*okta.OktaGroups, error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("SSWS %v", oktaToken)

	url := fmt.Sprintf("https://%v/api/v1/groups?limit=200", oktaDomain) // Docs said group limit was 200. That... should be enough. Otherwise we'll have to do some lame paging stuff. :( )

	result, err := httpRequest("GET", url, headers, make([]byte, 0))
	if err != nil {
		return nil, err
	}

	var groups okta.OktaGroups
	err = json.Unmarshal([]byte(result), &groups)

	if err != nil {
		return nil, err
	}

	return &groups, nil
}

func httpRequest(method string, url string, headers map[string]string, data []byte) (string, error) {
	client := http.DefaultClient

	rq, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	// Add any headers
	for k, v := range headers {
		rq.Header.Add(k, v)
	}

	rsp, err := client.Do(rq)
	if err != nil {
		return "", err
	}

	b, _ := io.ReadAll(rsp.Body)
	return string(b), nil
}
