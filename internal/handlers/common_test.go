package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
)

type StatusBodyPair struct {
	status int
	body   []byte
}

type TestError string

func (e TestError) Error() string {
	return string(e)
}

type TestSuite struct {
	suite.Suite
	router   *gin.Engine
	handlers *HandlerBuilder
	srv      *httptest.Server
	client   *http.Client
}

func (st *TestSuite) getCtrl() func() *gomock.Controller {
	ctrl := gomock.NewController(st.T())
	return func() *gomock.Controller {
		return ctrl
	}
}

func (st *TestSuite) CTRL() *gomock.Controller {
	getter := st.getCtrl()
	return getter()
}

func (st *TestSuite) SetupSuite() {
	st.handlers = &HandlerBuilder{
		logger: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	}
}

func (st *TestSuite) StartSrv() {
	st.srv = httptest.NewServer(st.router)
	st.client = st.srv.Client()
}

func (st *TestSuite) InitRouter() {
	st.router = gin.Default()
	st.router.Use(st.handlers.errorMiddleware)
}

func (st *TestSuite) TearDownTest() {
	st.srv.Close()
}

func (st *TestSuite) compareResponse(r, expected *StatusBodyPair) {
	st.Assert().Equal(expected.status, r.status)
	if st.Equal(len(expected.body), len(r.body)) {
		st.Assert().Equal(expected.body, r.body)
	}
}

func (st *TestSuite) prepareReq(path, token string) *http.Request {
	uri, err := url.JoinPath(st.srv.URL, path)
	st.Require().NoError(err)
	req, err := http.NewRequest("GET", uri, nil)
	st.Require().NoError(err)
	if token != "" {
		req.Header.Set("token", token)
	}
	return req
}

func (st *TestSuite) doReq(req *http.Request) *StatusBodyPair {
	response, err := st.client.Do(req)
	st.Require().NoError(err)
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	st.Require().NoError(err)
	return &StatusBodyPair{
		status: response.StatusCode,
		body:   b,
	}
}

func callOK(c *gin.Context) {
	c.Status(http.StatusOK)
}

func callError(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)
	collectErrors(c, errs[id])
}

func marshalBody(expected any, suite *suite.Suite) []byte {
	b, err := json.Marshal(&expected)
	suite.NoError(err)
	return b
}
