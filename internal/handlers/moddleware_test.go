package handlers

import (
	e "BannerFlow/internal/domain/errors"
	"BannerFlow/internal/handlers/mocks"
	"BannerFlow/pkg/api"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"net/http"
	"strconv"
	"testing"
)

func TestMiddleware(t *testing.T) {
	suite.Run(t, new(AuthTest))
	suite.Run(t, new(ErrorHandlerTest))
	suite.Run(t, new(GetTokenTest))
}

var errs = map[int]error{
	http.StatusInternalServerError: e.ErrorInternal,
	http.StatusBadRequest:          e.ErrorBadRequest,
	http.StatusUnauthorized:        e.ErrorAuthenticationFailed,
	http.StatusForbidden:           e.ErrorNoPermission,
	http.StatusNotFound:            e.ErrorNotFound,
}

var nums = map[error]int{
	e.ErrorInternal:             http.StatusInternalServerError,
	e.ErrorBadRequest:           http.StatusBadRequest,
	e.ErrorAuthenticationFailed: http.StatusUnauthorized,
	e.ErrorNoPermission:         http.StatusForbidden,
	e.ErrorNotFound:             http.StatusNotFound,
}

type ErrorHandlerTest struct {
	TestSuite
}

func (s *ErrorHandlerTest) SetupTest() {
	s.InitRouter()
	s.router.GET("/:id", callError)
	s.StartSrv()
}

func (s *ErrorHandlerTest) TestErrors() {
	const method = "GET"
	internalResponse := "something went wrong"
	badRequestResponse := e.ErrorBadRequest.Error()

	tests := []struct {
		name string
		err  error
		body []byte
	}{
		{name: "Internal", err: e.ErrorInternal, body: marshalBody(api.BannerErrorResponse{Error: &internalResponse}, &s.Suite)},
		{name: "BadRequest", err: e.ErrorBadRequest, body: marshalBody(api.BannerErrorResponse{Error: &badRequestResponse}, &s.Suite)},
		{name: "Unathorized", err: e.ErrorAuthenticationFailed, body: []byte{}},
		{name: "Forbidden", err: e.ErrorNoPermission, body: []byte{}},
		{name: "NotFound", err: e.ErrorNotFound, body: []byte{}},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			req := s.prepareReq("/"+strconv.Itoa(nums[test.err]), method, nil, "")
			r := s.doReq(req)
			s.compareResponse(r, &StatusBodyPair{
				status: nums[test.err],
				body:   test.body,
			})
		})
	}
}

type GetTokenTest struct {
	TestSuite
	generator *mocks.MockTokenGenerator
}

func (s *GetTokenTest) SetupTest() {
	ctrl := s.CTRL()
	s.generator = mocks.NewMockTokenGenerator(ctrl)
	s.handlers.generator = s.generator

	s.InitRouter()
	s.router.GET("/generate/*admin", s.handlers.handleTokenGeneration)
	s.StartSrv()
}

func (s *GetTokenTest) TestGetUserTokenOK() {
	const (
		method    = "GET"
		userToken = "asdf1234"
		uri       = "/generate/"
	)

	s.generator.EXPECT().GenerateToken(false).Return(userToken, nil)

	req := s.prepareReq(uri, method, nil, "")
	r := s.doReq(req)
	s.compareResponse(r, &StatusBodyPair{
		status: http.StatusOK,
		body:   marshalBody(api.TokenResponse{Token: userToken}, &s.Suite),
	})
}

func (s *GetTokenTest) TestGetUserTokenError() {
	const (
		method = "GET"
		uri    = "/generate/"
	)
	errorResponse := "something went wrong"
	s.generator.EXPECT().GenerateToken(false).Return("", TestError("any"))

	req := s.prepareReq(uri, method, nil, "")
	r := s.doReq(req)
	s.compareResponse(r, &StatusBodyPair{
		status: http.StatusInternalServerError,
		body:   marshalBody(api.BannerErrorResponse{Error: &errorResponse}, &s.Suite),
	})
}

func (s *GetTokenTest) TestGetAdminTokenOK() {
	const (
		method     = "GET"
		adminToken = "asdf5674"
	)
	const uri = "/generate/admin"
	s.generator.EXPECT().GenerateToken(true).Return(adminToken, nil)

	req := s.prepareReq(uri, method, nil, "")
	r := s.doReq(req)
	s.compareResponse(r, &StatusBodyPair{
		status: http.StatusOK,
		body:   marshalBody(api.TokenResponse{Token: adminToken}, &s.Suite),
	})
}

func (s *GetTokenTest) TestGetAdminTokenError() {
	const (
		method = "GET"
		uri    = "/generate/"
	)
	errorResponse := "something went wrong"
	s.generator.EXPECT().GenerateToken(false).Return("", TestError("any"))

	req := s.prepareReq(uri, method, nil, "")
	r := s.doReq(req)
	s.compareResponse(r, &StatusBodyPair{
		status: http.StatusInternalServerError,
		body:   marshalBody(api.BannerErrorResponse{Error: &errorResponse}, &s.Suite),
	})
}

type AuthTest struct {
	TestSuite
	authorizer    *mocks.MockAuthorizer
	authenticator *mocks.MockAuthenticator
}

func (s *AuthTest) SetupTest() {
	ctrl := s.CTRL()

	s.authorizer = mocks.NewMockAuthorizer(ctrl)
	s.authenticator = mocks.NewMockAuthenticator(ctrl)

	s.handlers.authenticator = s.authenticator
	s.handlers.authorizer = s.authorizer

	s.InitRouter()
	s.router.GET("/authenticate", s.handlers.authenticate, callOK)
	s.router.GET("/authorize", s.handlers.authenticate, s.handlers.authorize, callOK)
	s.StartSrv()
}

func (s *AuthTest) TestAuthenticateOK() {
	const (
		method = "GET"
		uri    = "/authenticate"
		token  = "asdf1234"
	)

	s.authenticator.EXPECT().Authenticate(token).Return(nil)

	req := s.prepareReq(uri, method, nil, token)
	r := s.doReq(req)
	s.compareResponse(r, &StatusBodyPair{
		status: http.StatusOK,
		body:   []byte{},
	})
}

func (s *AuthTest) TestAuthenticateError() {
	const (
		method = "GET"
		uri    = "/authenticate"
		token  = "asdf1234"
	)

	s.authenticator.EXPECT().Authenticate(gomock.Any()).Return(TestError("any"))

	req := s.prepareReq(uri, method, nil, token)
	r := s.doReq(req)
	s.compareResponse(r, &StatusBodyPair{
		status: http.StatusUnauthorized,
		body:   []byte{},
	})
}

func (s *AuthTest) TestAuthorizeOK() {
	const (
		method  = "GET"
		uri     = "/authorize"
		token   = "asdf5674"
		isAdmin = true
	)
	s.authenticator.EXPECT().Authenticate(gomock.Any()).Return(nil)
	s.authorizer.EXPECT().IsAdmin(token).Return(isAdmin)

	req := s.prepareReq(uri, method, nil, token)
	r := s.doReq(req)
	s.compareResponse(r, &StatusBodyPair{
		status: http.StatusOK,
		body:   []byte{},
	})
}

func (s *AuthTest) TestAuthorizeError() {
	const (
		method  = "GET"
		uri     = "/authorize"
		token   = "asdf5674"
		isAdmin = false
	)
	s.authenticator.EXPECT().Authenticate(gomock.Any()).Return(nil)
	s.authorizer.EXPECT().IsAdmin(gomock.Any()).Return(isAdmin)

	req := s.prepareReq(uri, method, nil, token)
	r := s.doReq(req)
	s.compareResponse(r, &StatusBodyPair{
		status: http.StatusForbidden,
		body:   []byte{},
	})
}
