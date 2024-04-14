package tests

import (
	"net/http"
	"strings"
)

func (t *E2ETest) TestBannersAdd() {
	Tests := []Test{
		{
			name: "addOK",
			Req: PrepareRequest(http.MethodPost, "http://"+t.address+"/banner", "application/json", t.adminToken,
				strings.NewReader("{\n  \"tag_ids\": [\n    0\n  ],\n  \"feature_id\": 0,\n  \"content\": {\n    \"title\": \"some_title\",\n    \"text\": \"some_text\",\n    \"url\": \"some_url\"\n  },\n  \"is_active\": true\n}")),
			Resp: &Response{
				StatusCode: http.StatusCreated,
				Body:       MarshalStruct(NewBannerIdResponse(1)),
			},
		},
		{
			name: "addConflict",
			Req: PrepareRequest(http.MethodPost, "http://"+t.address+"/banner", "application/json", t.adminToken,
				strings.NewReader("{\n  \"tag_ids\": [\n    0\n  ],\n  \"feature_id\": 0,\n  \"content\": {\n    \"title\": \"some_title\",\n    \"text\": \"some_text\",\n    \"url\": \"some_url\"\n  },\n  \"is_active\": true\n}")),
			Resp: &Response{
				StatusCode: http.StatusBadRequest,
				Body:       MarshalStruct(NewBannerErrorResponse("request is invalid: banner already exists")),
			},
		},
	}

	for _, test := range Tests {
		t.Run(test.name, func() {
			t.doRequest(test.Req, test.Resp)
		})
	}
}
