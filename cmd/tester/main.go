package main

import (
	"BannerFlow/internal/tests"
	"BannerFlow/pkg/api"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
)

const (
	flagName = "addr"
	envName  = "SERVER_ADDR"
)

func main() {
	addr := getServerAddress()
	if addr == "" {
		panic("no server address")
	}
	client := &http.Client{}

	userToken, err := getUserToken(addr, client)
	if err != nil {
		panic(err)
	}
	adminToken, err := getAdminToken(addr, client)
	if err != nil {
		panic(err)
	}

	fmt.Println("userToken:", userToken)
	fmt.Println("adminToken:", adminToken)
	tests.AddBannersTests(addr, adminToken)

	for _, t := range tests.GetTests() {
		r, err := client.Do(t.Req)
		if err != nil {
			fmt.Println("Error sending request to server:", err)
			return
		}
		defer r.Body.Close()
		if r.StatusCode != t.Resp.StatusCode {
			fmt.Println("Bad response from server:", r.StatusCode)
		}
		responseBody, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		if slices.Equal(responseBody, t.Resp.Body) {
			fmt.Println("Bad response from server:", string(responseBody))
		}
	}
}

func getUserToken(address string, client *http.Client) (string, error) {
	return getToken(address, "/get_token/", client)
}

func getAdminToken(address string, client *http.Client) (string, error) {
	return getToken(address, "/get_token/admin", client)
}

func getToken(address, uri string, client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", "http://"+address+uri, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	r, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server:", err)
		return "", err
	}
	defer r.Body.Close()
	reader := json.NewDecoder(r.Body)
	token := &api.TokenResponse{}
	err = reader.Decode(token)
	if err != nil {
		fmt.Println("Error decoding token:", err)
	}
	return token.Token, nil
}

func getServerAddress() string {
	var address string
	flag.StringVar(&address, flagName, "", "server address")
	flag.Parse()
	if address == "" {
		address = os.Getenv(envName)
	}
	return address
}
