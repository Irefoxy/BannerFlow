package tests

import (
	"net/http"
	"strings"
)

func (t *E2ETest) TestUserGet() {
	Tests := []Test{
		{
			name: "OK",
			Req: PrepareRequest(http.MethodGet, "http://"+t.address+"/user_banner?feature_id=100&tag_id=100", "application/json", t.userToken,
				nil),
			Resp: &Response{
				StatusCode: http.StatusOK,
				Body:       MarshalStruct(map[string]any{"text": "some_text", "title": "some_title", "url": "some_url"}),
			},
		},
		{
			name: "wrong token",
			Req: PrepareRequest(http.MethodGet, "http://"+t.address+"/user_banner?feature_id=200&tag_id=200", "application/json", "asdaasda",
				nil),
			Resp: &Response{
				StatusCode: http.StatusUnauthorized,
				Body:       nil,
			},
		},
		{
			name: "NotExist",
			Req: PrepareRequest(http.MethodGet, "http://"+t.address+"/user_banner?feature_id=100&tag_id=200", "application/json", t.userToken,
				nil),
			Resp: &Response{
				StatusCode: http.StatusNotFound,
				Body:       nil,
			},
		},
		{
			name: "admin updating",
			Req: PrepareRequest(http.MethodPatch, "http://"+t.address+"/banner/2", "application/json", t.adminToken,
				strings.NewReader("{\n  \"tag_ids\": [\n    200\n  ],\n  \"feature_id\": 200,\n  \"content\": {\n    \"title\": \"some_title\",\n    \"text\": \"some_text\",\n    \"url\": \"some_url\"\n  },\n  \"is_active\": false\n}")),
			Resp: &Response{
				StatusCode: http.StatusOK,
				Body:       nil,
			},
		},
		{
			name: "NotVisible",
			Req: PrepareRequest(http.MethodGet, "http://"+t.address+"/user_banner?feature_id=200&tag_id=200&use_last_revision=true", "application/json", t.userToken,
				nil),
			Resp: &Response{
				StatusCode: http.StatusNotFound,
				Body:       nil,
			},
		},
	}

	for _, test := range Tests {
		t.Run(test.name, func() {
			t.doRequest(test.Req, test.Resp)
		})
	}
}
