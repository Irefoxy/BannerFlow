package tests

import (
	"BannerFlow/pkg/api"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

const (
	envName = "SERVER_ADDR"
)

type E2ETest struct {
	suite.Suite
	userToken  string
	adminToken string
	address    string
	client     *http.Client
}

func (t *E2ETest) SetupTest() {
	t.getServerAddress()
	t.Assert().NotEmpty(t.address, "Server address should not be empty")
	t.client = &http.Client{}
	t.getUserToken()
	t.getAdminToken()
	t.setDefaultBanners()
}

func TestE2E(t *testing.T) {
	suite.Run(t, new(E2ETest))
}

func (t *E2ETest) doRequest(req *http.Request, expected *Response) {
	r, err := t.client.Do(req)
	if t.NoError(err, "Error sending request to server") {
		t.Equal(expected.StatusCode, r.StatusCode, "Bad response from server")
		responseBody, err := io.ReadAll(r.Body)
		t.NoError(err, "Error reading response body")
		t.Equal(expected.Body, responseBody, "Bad response body from server")
		r.Body.Close()
	}
}

func (t *E2ETest) getUserToken() {
	token, err := getToken(t.address, "/get_token/", t.client)
	t.Assert().NoError(err, "Fetching user token failed")
	t.Assert().NotEmpty(token, "User token should not be empty")
	t.userToken = token
}

func (t *E2ETest) getAdminToken() {
	token, err := getToken(t.address, "/get_token/admin", t.client)
	t.Assert().NoError(err, "Fetching admin token failed")
	t.Assert().NotEmpty(token, "Admin token should not be empty")
	t.adminToken = token
}

func getToken(address, uri string, client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", "http://"+address+uri, nil)
	if err != nil {
		return "", err
	}
	r, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	reader := json.NewDecoder(r.Body)
	token := &api.TokenResponse{}
	err = reader.Decode(token)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

func (t *E2ETest) getServerAddress() {
	t.address = os.Getenv(envName)
}

func (t *E2ETest) setDefaultBanners() {
	Tests := []Test{
		{
			Req: PrepareRequest(http.MethodPost, "http://"+t.address+"/banner", "application/json", t.adminToken,
				strings.NewReader("{\n  \"tag_ids\": [\n    100\n  ],\n  \"feature_id\": 100,\n  \"content\": {\n    \"title\": \"some_title\",\n    \"text\": \"some_text\",\n    \"url\": \"some_url\"\n  },\n  \"is_active\": true\n}")),
			Resp: &Response{
				StatusCode: http.StatusCreated,
				Body:       MarshalStruct(NewBannerIdResponse(1)),
			},
		},
	}

	for _, test := range Tests {
		t.doRequest(test.Req, test.Resp)
	}
}
