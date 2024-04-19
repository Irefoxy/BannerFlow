package handlers

import (
	"BannerFlow/internal/handlers/mocks"
	"bytes"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"testing"
)

type HandlersTest struct {
	TestSuite
}

func TestHandlers(t *testing.T) {
	suite.Run(t, new(HandlersTest))
}

func (s *HandlersTest) SetupTest() {
	ctrl := s.CTRL()
	s.service = mocks.NewMockService(ctrl)
	s.handlers.srv = s.service

	s.InitRouter()
	s.router.GET("/user_banner", s.handlers.handleUserGetBanner)
	adminGroup := s.router.Group("/banner")
	adminGroup.GET("", s.handlers.handleListBanners)
	adminGroup.POST("", s.handlers.handleCreateBanner)
	adminGroup.DELETE("/:id", s.handlers.handleDeleteBanner)
	adminGroup.PATCH("/:id", s.handlers.handleUpdateBanner)
	adminGroup.GET("/versions/:id", s.handlers.handleListBannerHistory)
	adminGroup.PUT("/versions/:id/activate", s.handlers.handleSelectBannerVersion)
	adminGroup.DELETE("/del", s.handlers.handleDeleteBannerByTagOrFeature)
	s.StartSrv()
}

func (s *HandlersTest) TestCreateBannerWrongRequestBody() {
	const method = "POST"
	uri := "/banner"
	tests := []struct {
		name           string
		requestBody    io.Reader
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "wrong feature",
			requestBody:    bytes.NewReader(marshalBody(setBannerRequestFields(map[string]any{"title": "some"}, nil, false, []int{1}), &s.Suite)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in request body: Key: 'BannerRequest.FeatureId' Error:Field validation for 'FeatureId' failed on the 'required' tag"), &s.Suite),
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			req := s.prepareReq(uri, method, test.requestBody, "")
			resp := s.doReq(req)
			s.compareResponse(resp, &StatusBodyPair{
				status: test.expectedStatus,
				body:   test.expectedBody,
			})
		})
	}
}
