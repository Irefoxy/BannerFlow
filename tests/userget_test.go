package tests

import (
	"net/http"
)

func (t *E2ETest) TestUserGet() {
	Tests := []Test{
		{
			name: "wrong content",
			Req: PrepareRequest(http.MethodGet, "http://"+t.address+"/user_banner?feature_id=100&tag_id=100", "application/xml", t.userToken,
				nil),
			Resp: &Response{
				StatusCode: http.StatusBadRequest,
				Body:       MarshalStruct(NewBannerIdResponse(1)),
			},
		},
		{
			name: "wrong token",
			Req: PrepareRequest(http.MethodGet, "http://"+t.address+"/user_banner?feature_id=100&tag_id=100", "application/json", "asdaasda",
				nil),
			Resp: &Response{
				StatusCode: http.StatusUnauthorized,
				Body:       MarshalStruct(NewBannerIdResponse(1)),
			},
		},
	}

	for _, test := range Tests {
		t.Run(test.name, func() {
			t.doRequest(test.Req, test.Resp)
		})
	}
}
