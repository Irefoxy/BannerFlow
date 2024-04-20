package handlers

import (
	"BannerFlow/internal/domain/models"
	"BannerFlow/internal/handlers/mocks"
	"bytes"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"strconv"
	"testing"
)

type HandlersTest struct {
	TestSuite
	service *mocks.MockService
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
	const (
		method = "POST"
		uri    = "/banner"
	)

	tests := []struct {
		name           string
		requestBody    io.Reader
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "wrong feature - required",
			requestBody:    bytes.NewReader(marshalBody(setBannerRequestFields(map[string]any{"title": "some"}, nil, false, []int{1}), &s.Suite)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in request body: Key: 'BannerRequest.FeatureId' Error:Field validation for 'FeatureId' failed on the 'required' tag"), &s.Suite),
		},
		{
			name:           "wrong feature - gte 0",
			requestBody:    bytes.NewReader(marshalBody(setBannerRequestFields(map[string]any{"title": "some"}, -1, false, []int{1}), &s.Suite)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in request body: Key: 'BannerRequest.FeatureId' Error:Field validation for 'FeatureId' failed on the 'gte' tag"), &s.Suite),
		},
		{
			name:           "wrong tagIds - required",
			requestBody:    bytes.NewReader(marshalBody(setBannerRequestFields(map[string]any{"title": "some"}, 1, false, nil), &s.Suite)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in request body: Key: 'BannerRequest.TagIds' Error:Field validation for 'TagIds' failed on the 'required' tag"), &s.Suite),
		},
		{
			name:           "wrong tagIds - gte 1",
			requestBody:    bytes.NewReader(marshalBody(setBannerRequestFields(map[string]any{"title": "some"}, 1, false, []int{}), &s.Suite)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in request body: Key: 'BannerRequest.TagIds' Error:Field validation for 'TagIds' failed on the 'gte' tag"), &s.Suite),
		},
		{
			name:           "wrong isActive - required",
			requestBody:    bytes.NewReader(marshalBody(setBannerRequestFields(map[string]any{"title": "some"}, 1, nil, []int{1}), &s.Suite)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in request body: Key: 'BannerRequest.IsActive' Error:Field validation for 'IsActive' failed on the 'required' tag"), &s.Suite),
		},
		{
			name:           "wrong content - required",
			requestBody:    bytes.NewReader(marshalBody(setBannerRequestFields(nil, 1, true, []int{1}), &s.Suite)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in request body: Key: 'BannerRequest.Content' Error:Field validation for 'Content' failed on the 'required' tag"), &s.Suite),
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

func (s *HandlersTest) TestCreateBannerService() {
	const (
		method = "POST"
		uri    = "/banner"
	)
	expectedArg := &models.Banner{
		BaseBanner: models.BaseBanner{
			UserBanner: models.UserBanner{Content: map[string]any{"title": "some"}},
			FeatureId:  1,
			TagIds:     []int{1},
		},
		IsActive: false,
	}
	requestBody :=
		marshalBody(setBannerRequestFields(expectedArg.Content, expectedArg.FeatureId, expectedArg.IsActive, expectedArg.TagIds), &s.Suite)

	tests := []struct {
		name           string
		mockedErr      error
		mockedId       int
		expectedStatus int
		expectedBody   []byte
	}{
		{"ServiceError", TestError("any"), 0, http.StatusInternalServerError, marshalBody(NewBannerErrorResponse("something went wrong"), &s.Suite)},
		{"ServiceOK", nil, 123, http.StatusCreated, marshalBody(NewBannerIdResponse(123), &s.Suite)},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.service.EXPECT().CreateBanner(gomock.Any(), expectedArg).Return(test.mockedId, test.mockedErr)
			req := s.prepareReq(uri, method, bytes.NewReader(requestBody), "")
			resp := s.doReq(req)
			s.compareResponse(resp, &StatusBodyPair{
				status: test.expectedStatus,
				body:   test.expectedBody,
			})
		})
	}
}

func (s *HandlersTest) TestUserBannerWrongParams() {
	const (
		method     = "GET"
		defaultUri = "/user_banner"
	)
	tests := []struct {
		name           string
		uri            string
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "wrong feature id - required",
			uri:            defaultUri + "?tag_id=1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: Key: 'UserBannerParams.FeatureId' Error:Field validation for 'FeatureId' failed on the 'required' tag"), &s.Suite),
		},
		{
			name:           "wrong feature id - gte 0",
			uri:            defaultUri + "?tag_id=1&feature_id=-1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: Key: 'UserBannerParams.FeatureId' Error:Field validation for 'FeatureId' failed on the 'gte' tag"), &s.Suite),
		},
		{
			name:           "wrong tag id - required",
			uri:            defaultUri + "?feature_id=1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: Key: 'UserBannerParams.TagId' Error:Field validation for 'TagId' failed on the 'required' tag"), &s.Suite),
		},
		{
			name:           "wrong tag id - gte 0",
			uri:            defaultUri + "?feature_id=1&tag_id=-1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: Key: 'UserBannerParams.TagId' Error:Field validation for 'TagId' failed on the 'gte' tag"), &s.Suite),
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			req := s.prepareReq(test.uri, method, nil, "")
			resp := s.doReq(req)
			s.compareResponse(resp, &StatusBodyPair{
				status: test.expectedStatus,
				body:   test.expectedBody,
			})
		})
	}
}

func (s *HandlersTest) TestUserBannerService() {
	expectedArg := &models.BannerIdentOptions{
		FeatureId: 1,
		TagId:     1,
	}
	const method = "GET"
	defaultUri := "/user_banner?tag_id=" + strconv.Itoa(expectedArg.TagId) + "&feature_id=" + strconv.Itoa(expectedArg.FeatureId)
	defaultUserBanner := &models.UserBanner{Content: map[string]any{"title": "some"}}

	tests := []struct {
		name           string
		uri            string
		useFlag        bool
		mockedErr      error
		mockedBanner   *models.UserBanner
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "use_last_revision=true with service error",
			uri:            defaultUri + "&use_last_revision=true",
			useFlag:        true,
			mockedErr:      TestError("any"),
			mockedBanner:   nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   marshalBody(NewBannerErrorResponse("something went wrong"), &s.Suite),
		},
		{
			name:           "use_last_revision=false with service OK",
			uri:            defaultUri,
			useFlag:        false,
			mockedErr:      nil,
			mockedBanner:   defaultUserBanner,
			expectedStatus: http.StatusOK,
			expectedBody:   marshalBody(defaultUserBanner.Content, &s.Suite),
		},
	}

	for _, test := range tests {
		s.service.EXPECT().UserGetBanners(gomock.Any(), &models.BannerUserOptions{
			BannerIdentOptions: *expectedArg,
			UseLastRevision:    test.useFlag,
		}).Return(test.mockedBanner, test.mockedErr)
		s.Run(test.name, func() {
			req := s.prepareReq(test.uri, method, nil, "")
			resp := s.doReq(req)
			s.compareResponse(resp, &StatusBodyPair{
				status: test.expectedStatus,
				body:   test.expectedBody,
			})
		})
	}
}

func (s *HandlersTest) TestDeleteBannerWrongParam() {
	const (
		method = "DELETE"
		uri    = "/banner/"
	)

	expectedStatus := http.StatusNotFound
	expectedBody := []byte("404 page not found")

	req := s.prepareReq(uri, method, nil, "")
	resp := s.doReq(req)
	s.compareResponse(resp, &StatusBodyPair{
		status: expectedStatus,
		body:   expectedBody,
	})
}

func (s *HandlersTest) TestDeleteBannerService() {
	expectedId := 123
	const method = "DELETE"
	defaultUri := "/banner/" + strconv.Itoa(expectedId)

	tests := []struct {
		name           string
		mockedErr      error
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "deleteBanner with service error",
			mockedErr:      TestError("any"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   marshalBody(NewBannerErrorResponse("something went wrong"), &s.Suite),
		},
		{
			name:           "deleteBanner with service OK",
			mockedErr:      nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   []byte{},
		},
	}

	for _, test := range tests {
		s.service.EXPECT().DeleteBanner(gomock.Any(), expectedId).Return(test.mockedErr)
		s.Run(test.name, func() {
			req := s.prepareReq(defaultUri, method, nil, "")
			resp := s.doReq(req)
			s.compareResponse(resp, &StatusBodyPair{
				status: test.expectedStatus,
				body:   test.expectedBody,
			})
		})
	}
}
