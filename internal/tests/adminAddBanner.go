package tests

import (
	"net/http"
	"strings"
)

func AddBannersTests(address, adminToken string) {
	bannerTests := []Test{
		{
			Req: PrepareRequest(http.MethodPost, "http://"+address+"/banner", "application/json", adminToken,
				strings.NewReader("{\n  \"tag_ids\": [\n    0\n  ],\n  \"feature_id\": 0,\n  \"content\": {\n    \"title\": \"some_title\",\n    \"text\": \"some_text\",\n    \"url\": \"some_url\"\n  },\n  \"is_active\": true\n}")),
			Resp: &Response{
				StatusCode: http.StatusCreated,
				Body:       MarshalStruct(NewBannerIdResponse(1)),
			},
		},
	}

	tests = append(tests, bannerTests...)
}
