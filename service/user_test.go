package service

import (
	"IAMbyGo/repo"
	"testing"

	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestService(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

type serviceTestSuite struct {
	suite.Suite
	service      User
	mockUserRepo *repo.MockUser
	mockTools    *mockTools
}

func (s *serviceTestSuite) SetupTest() {
	s.mockUserRepo = &repo.MockUser{}
	s.mockTools = &mockTools{}
	s.service = &user{
		userRepo:       s.mockUserRepo,
		privateMethods: s.mockTools,
	}
}

func (s *serviceTestSuite) TestHashUserpasswordWhenCreateUser() {

	s.mockUserRepo.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.service.CreateUser("", "", "")
	s.mockUserRepo.AssertCalled(s.T(), "CreateUser", mock.Anything, mock.Anything, mock.Anything)
}

func (s *serviceTestSuite) Test2() {

	s.mockUserRepo.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.service.CreateUser("", "", "")
	s.mockUserRepo.AssertCalled(s.T(), "CreateUser", mock.Anything, mock.Anything, mock.Anything)
}

func (s *serviceTestSuite) TestCreateUserUseHashedPassword() {
	username := "abc"
	password := "ABC"
	email := "abc email"
	// userID := "123"
	hashedPassword := "AABBCC"

	s.mockUserRepo.On("CreateUser", mock.Anything, username, password).Return(nil)
	s.mockTools.On("hashPassword", password).Return(hashedPassword, nil)
	s.service.CreateUser(username, email, password)
	s.mockUserRepo.AssertCalled(s.T(), "CreateUser", mock.Anything, username, password)
	s.mockTools.AssertCalled(s.T(), "hashPassword", password)

}
