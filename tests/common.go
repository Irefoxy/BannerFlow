package tests

import (
	"BannerFlow/pkg/api"
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	StatusCode int
	Body       []byte
}

type Test struct {
	name string
	Req  *http.Request
	Resp *Response
}

func PrepareRequest(method, target, contentType, token string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, target, body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("token", token)
	return req
}

func MarshalStruct(body any) []byte {
	b, _ := json.Marshal(&body)
	return b
}

func NewBannerIdResponse(id int) api.BannerIdResponse {
	return api.BannerIdResponse{
		BannerId: &id,
	}
}

func NewBannerErrorResponse(msg string) api.BannerErrorResponse {
	return api.BannerErrorResponse{
		Error: &msg,
	}
}
