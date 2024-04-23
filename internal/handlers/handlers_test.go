package handlers

import (
	"BannerFlow/internal/domain/models"
	"BannerFlow/internal/handlers/mocks"
	"BannerFlow/pkg/api"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"
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
			name:           "wrong tagIds - not empty",
			requestBody:    bytes.NewReader(marshalBody(setBannerRequestFields(map[string]any{"title": "some"}, 1, false, []int{}), &s.Suite)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in request body: Key: 'BannerRequest.TagIds' Error:Field validation for 'TagIds' failed on the 'gte' tag"), &s.Suite),
		},
		{
			name:           "wrong tagIds - elms gte 0",
			requestBody:    bytes.NewReader(marshalBody(setBannerRequestFields(map[string]any{"title": "some"}, 1, false, []int{-1}), &s.Suite)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in request body: Key: 'BannerRequest.TagIds[0]' Error:Field validation for 'TagIds[0]' failed on the 'gte' tag"), &s.Suite),
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

func (s *HandlersTest) TestUpdateBannerWrongRequest() {
	const (
		method     = "PATCH"
		defaultUri = "/banner/"
		id         = 123
	)
	uri := defaultUri + strconv.Itoa(id)

	tests := []struct {
		name           string
		uri            string
		requestBody    io.Reader
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "wrong id - not exist",
			uri:            defaultUri,
			requestBody:    bytes.NewReader(marshalBody(setBannerUpdateRequestFields(map[string]any{"title": "some"}, 1, true, []int{1}), &s.Suite)),
			expectedStatus: http.StatusNotFound,
			expectedBody:   []byte("404 page not found"),
		},
		{
			name:           "wrong id - not int",
			uri:            defaultUri + "asd",
			requestBody:    bytes.NewReader(marshalBody(setBannerUpdateRequestFields(map[string]any{"title": "some"}, 1, true, []int{1}), &s.Suite)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: strconv.ParseInt: parsing \"asd\": invalid syntax"), &s.Suite)},
		{
			name:           "wrong request - all empty",
			uri:            uri,
			requestBody:    bytes.NewReader(marshalBody(setBannerUpdateRequestFields(nil, nil, nil, nil), &s.Suite)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in request body: all fields are empty"), &s.Suite),
		},
		{
			name:           "wrong request - tags elm < 0",
			uri:            uri,
			requestBody:    bytes.NewReader(marshalBody(setBannerUpdateRequestFields(nil, nil, nil, []int{1, -2}), &s.Suite)),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in request body: Key: 'BannerUpdateRequest.TagIds[1]' Error:Field validation for 'TagIds[1]' failed on the 'gte' tag"), &s.Suite),
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			req := s.prepareReq(test.uri, method, test.requestBody, "")
			resp := s.doReq(req)
			s.compareResponse(resp, &StatusBodyPair{
				status: test.expectedStatus,
				body:   test.expectedBody,
			})
		})
	}
}

func (s *HandlersTest) TestUpdateBannerService() {
	const (
		method     = "PATCH"
		defaultUri = "/banner/"
		id         = 123
	)
	uri := defaultUri + strconv.Itoa(id)
	expectedArgs := &models.UpdateBanner{
		Banner: models.Banner{
			BaseBanner: models.BaseBanner{
				UserBanner: models.UserBanner{
					Content: map[string]any{"title": "some"},
				},
				FeatureId: 1,
				TagIds:    []int{1},
			},
		},
		Flags: models.ContentBit | models.FeatureBit | models.TagBit,
	}
	requestBody :=
		marshalBody(setBannerUpdateRequestFields(expectedArgs.Content, expectedArgs.FeatureId, nil, expectedArgs.TagIds), &s.Suite)

	tests := []struct {
		name           string
		mockedErr      error
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "ServiceError",
			mockedErr:      TestError("any"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   marshalBody(NewBannerErrorResponse("something went wrong"), &s.Suite),
		},
		{
			name:           "ServiceOK",
			mockedErr:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   []byte{},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.service.EXPECT().UpdateBanner(gomock.Any(), id, expectedArgs).Return(test.mockedErr)
			req := s.prepareReq(uri, method, bytes.NewReader(requestBody), "")
			resp := s.doReq(req)
			s.compareResponse(resp, &StatusBodyPair{
				status: test.expectedStatus,
				body:   test.expectedBody,
			})
		})
	}
}

func (s *HandlersTest) TestListBannerWrongParams() {
	const (
		method     = "GET"
		defaultUri = "/banner"
	)

	tests := []struct {
		name           string
		uri            string
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "tag id < 0",
			uri:            defaultUri + "?feature_id=1&tag_id=-1&offset=0&limit=1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: Key: 'ListBannerParams.TagId' Error:Field validation for 'TagId' failed on the 'gte' tag"), &s.Suite),
		},
		{
			name:           "feature id < 0",
			uri:            defaultUri + "?feature_id=-1&tag_ids=1&offset=0&limit=1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: Key: 'ListBannerParams.FeatureId' Error:Field validation for 'FeatureId' failed on the 'gte' tag"), &s.Suite),
		},
		{
			name:           "offset < 0",
			uri:            defaultUri + "?feature_id=1&tag_ids=1&offset=-1&limit=1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: Key: 'ListBannerParams.Offset' Error:Field validation for 'Offset' failed on the 'gte' tag"), &s.Suite),
		},
		{
			name:           "limit < 1",
			uri:            defaultUri + "?feature_id=1&tag_ids=1&offset=0&limit=0",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: Key: 'ListBannerParams.Limit' Error:Field validation for 'Limit' failed on the 'gte' tag"), &s.Suite),
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

func (s *HandlersTest) TestListBannerService() {
	const (
		method     = "GET"
		defaultUri = "/banner"
	)
	defaultResponse := models.BannerExt{
		BannerId: 0,
		Banner: models.Banner{
			BaseBanner: models.BaseBanner{
				UserBanner: models.UserBanner{
					Content: map[string]any{"title": "some"},
				},
				FeatureId: 1,
				TagIds:    []int{1, 2, 3},
			},
			IsActive: true,
		},
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		uri            string
		expectedArgs   *models.BannerListOptions
		mockedErr      error
		mockedResponse []models.BannerExt
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name: "ServiceError",
			uri:  defaultUri + "?feature_id=1&limit=10&offset=6&tag_id=2",
			expectedArgs: &models.BannerListOptions{
				BannerIdentOptions: models.BannerIdentOptions{
					FeatureId: 1,
					TagId:     2,
				},
				Limit:  10,
				Offset: 6,
			},
			mockedErr:      TestError("any"),
			mockedResponse: nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   marshalBody(NewBannerErrorResponse("something went wrong"), &s.Suite),
		},
		{
			name: "ServiceOK",
			uri:  defaultUri,
			expectedArgs: &models.BannerListOptions{
				BannerIdentOptions: models.BannerIdentOptions{
					FeatureId: models.ZeroValue,
					TagId:     models.ZeroValue,
				},
				Limit:  models.ZeroValue,
				Offset: models.ZeroValue,
			},
			mockedErr:      nil,
			mockedResponse: []models.BannerExt{defaultResponse},
			expectedStatus: http.StatusOK,
			expectedBody:   marshalBody([]api.BannerResponse{NewBannerResponse(defaultResponse.BannerId, defaultResponse.Content, defaultResponse.CreatedAt, defaultResponse.UpdatedAt, defaultResponse.IsActive, defaultResponse.TagIds, defaultResponse.FeatureId)}, &s.Suite),
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.service.EXPECT().ListBanners(gomock.Any(), test.expectedArgs).Return(test.mockedResponse, test.mockedErr)
			req := s.prepareReq(test.uri, method, nil, "")
			resp := s.doReq(req)
			s.compareResponse(resp, &StatusBodyPair{
				status: test.expectedStatus,
				body:   test.expectedBody,
			})
		})
	}
}

func (s *HandlersTest) TestSelectBannerVersionWrongRequest() {
	const (
		method = "PUT"
	)

	tests := []struct {
		name           string
		uri            string
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "wrong id",
			uri:            fmt.Sprintf("/banner/versions/%s/activate", "asd"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: strconv.ParseInt: parsing \"asd\": invalid syntax"), &s.Suite),
		},
		{
			name:           "no id",
			uri:            fmt.Sprintf("/banner/versions/%s/activate", ""),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: Key: 'IdParams.Id' Error:Field validation for 'Id' failed on the 'required' tag"), &s.Suite),
		},
		{
			name:           "no version",
			uri:            fmt.Sprintf("/banner/versions/%s/activate%s", "123", ""),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: Key: 'SelectBannersParams.Version' Error:Field validation for 'Version' failed on the 'required' tag"), &s.Suite),
		},
		{
			name:           "version < 1",
			uri:            fmt.Sprintf("/banner/versions/%s/activate%s", "123", "?version=-1"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: Key: 'SelectBannersParams.Version' Error:Field validation for 'Version' failed on the 'gte' tag"), &s.Suite),
		},
		{
			name:           "id < 0",
			uri:            fmt.Sprintf("/banner/versions/%s/activate%s", "-1", "?version=12"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: Key: 'IdParams.Id' Error:Field validation for 'Id' failed on the 'gt' tag"), &s.Suite),
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

func (s *HandlersTest) TestSelectBannerVersionService() {
	const (
		method  = "PUT"
		id      = 123
		version = 12
	)
	uri := fmt.Sprintf("/banner/versions/%d/activate?version=%d", id, version)

	tests := []struct {
		name           string
		mockedErr      error
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "ServiceError",
			mockedErr:      TestError("any"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   marshalBody(NewBannerErrorResponse("something went wrong"), &s.Suite),
		},
		{
			name:           "ServiceOK",
			expectedStatus: http.StatusOK,
			expectedBody:   []byte{},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.service.EXPECT().SelectBannerVersion(gomock.Any(), id, version).Return(test.mockedErr)
			req := s.prepareReq(uri, method, nil, "")
			resp := s.doReq(req)
			s.compareResponse(resp, &StatusBodyPair{
				status: test.expectedStatus,
				body:   test.expectedBody,
			})
		})
	}
}

func (s *HandlersTest) TestListBannerHistoryWrongRequest() {
	const (
		method     = "GET"
		defaultUri = "/banner/versions/"
	)

	tests := []struct {
		name           string
		uri            string
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "wrong id",
			uri:            defaultUri + "asd",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: strconv.ParseInt: parsing \"asd\": invalid syntax"), &s.Suite),
		},
		{
			name:           "no id",
			uri:            defaultUri,
			expectedStatus: http.StatusNotFound,
			expectedBody:   []byte("404 page not found"),
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

func (s *HandlersTest) TestListBannerHistoryService() {
	const (
		method = "GET"
		id     = 123
	)
	uri := fmt.Sprintf("/banner/versions/%d", id)
	defaultResponse := models.HistoryBanner{
		BaseBanner: models.BaseBanner{
			UserBanner: models.UserBanner{
				Content: map[string]any{"title": "some"},
			},
			FeatureId: 12,
			TagIds:    []int{1, 2, 3},
		},
		Version: 23,
	}
	tests := []struct {
		name           string
		mockedErr      error
		mockedResponse []models.HistoryBanner
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "ServiceError",
			mockedErr:      TestError("any"),
			mockedResponse: nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   marshalBody(NewBannerErrorResponse("something went wrong"), &s.Suite),
		},
		{
			name:           "ServiceOK",
			mockedErr:      nil,
			mockedResponse: []models.HistoryBanner{defaultResponse},
			expectedStatus: http.StatusOK,
			expectedBody:   marshalBody([]api.BannerVersionResponse{NewHistoryResponse(defaultResponse.Content, defaultResponse.TagIds, defaultResponse.FeatureId, defaultResponse.Version)}, &s.Suite),
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.service.EXPECT().ListBannerHistory(gomock.Any(), id).Return(test.mockedResponse, test.mockedErr)
			req := s.prepareReq(uri, method, nil, "")
			resp := s.doReq(req)
			s.compareResponse(resp, &StatusBodyPair{
				status: test.expectedStatus,
				body:   test.expectedBody,
			})
		})
	}
}

func (s *HandlersTest) TestDeleteByTagOrFeatureWrongRequest() {
	const (
		method     = "DELETE"
		defaultUri = "/banner/del"
	)

	tests := []struct {
		name           string
		uri            string
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "both empty",
			uri:            defaultUri,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: tag id, feature id both exist or not exist"), &s.Suite),
		},
		{
			name:           "both not empty",
			uri:            defaultUri + "?tag_id=1&feature_id=1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: tag id, feature id both exist or not exist"), &s.Suite),
		},
		{
			name:           "tag < 0",
			uri:            defaultUri + "?tag_id=-1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: Key: 'DeleteBannerParams.TagId' Error:Field validation for 'TagId' failed on the 'gte' tag"), &s.Suite),
		},
		{
			name:           "feature < 0",
			uri:            defaultUri + "?feature_id=-1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   marshalBody(NewBannerErrorResponse("request is invalid: error in param: Key: 'DeleteBannerParams.FeatureId' Error:Field validation for 'FeatureId' failed on the 'gte' tag"), &s.Suite),
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

func (s *HandlersTest) TestDeleteByTagOrFeatureService() {
	const (
		method     = "DELETE"
		tagId      = 123
		featureId  = 321
		defaultUri = "/banner/del"
	)

	tests := []struct {
		name           string
		uri            string
		expectedArgs   *models.BannerIdentOptions
		mockedErr      error
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name: "ServiceError",
			uri:  defaultUri + "?feature_id=" + strconv.Itoa(featureId),
			expectedArgs: &models.BannerIdentOptions{
				FeatureId: featureId,
				TagId:     models.ZeroValue,
			},
			mockedErr:      TestError("any"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   marshalBody(NewBannerErrorResponse("something went wrong"), &s.Suite),
		},
		{
			name: "ServiceOK",
			uri:  defaultUri + "?tag_id=" + strconv.Itoa(tagId),
			expectedArgs: &models.BannerIdentOptions{
				FeatureId: models.ZeroValue,
				TagId:     tagId,
			},
			mockedErr:      nil,
			expectedStatus: http.StatusAccepted,
			expectedBody:   []byte{},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.service.EXPECT().DeleteBannersByTagOrFeature(gomock.Any(), test.expectedArgs).Return(test.mockedErr)
			req := s.prepareReq(test.uri, method, nil, "")
			resp := s.doReq(req)
			s.compareResponse(resp, &StatusBodyPair{
				status: test.expectedStatus,
				body:   test.expectedBody,
			})
		})
	}
}
