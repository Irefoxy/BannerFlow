package auth

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type AuthTest struct {
	suite.Suite
	auth       *Auth
	userToken  string
	adminToken string
}

func (s *AuthTest) SetupTest() {
	var err error
	s.auth = NewAuth()
	s.userToken, err = s.auth.GenerateToken(false)
	s.Require().NoError(err)
	s.adminToken, err = s.auth.GenerateToken(true)
	s.Require().NoError(err)
}

func (s *AuthTest) TestAuthenticate() {
	err := s.auth.Authenticate(s.userToken)
	s.Assert().NoError(err)
	err = s.auth.Authenticate(s.adminToken)
	s.Assert().NoError(err)
	err = s.auth.Authenticate("something")
	s.Assert().Error(err)
}

func (s *AuthTest) TestIsAdmin() {
	ok := s.auth.IsAdmin(s.adminToken)
	s.Assert().True(ok)
	ok = s.auth.IsAdmin(s.userToken)
	s.Assert().False(ok)
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthTest))
}
