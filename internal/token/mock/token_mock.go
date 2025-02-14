package mock

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/mock"
)


type MockTokenImpl struct {
	mock.Mock
}

func (m *MockTokenImpl) GenerateJWT(userName string) (string, error) {
	args := m.Called(userName)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockTokenImpl) ParseJWT(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Token), args.Error(1) 
}