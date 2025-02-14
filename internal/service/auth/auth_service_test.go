package auth_test

import (
	"avito/internal/logger"
	repo_mock "avito/internal/repository/mock"
	"avito/internal/service/auth"
	token_mock "avito/internal/token/mock"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticate(t *testing.T) {
	mockRepo := &repo_mock.MockRepositoryImpl{}
	mockLogger := logger.NewLogger("debug")
	mockToken := &token_mock.MockTokenImpl{}

	authService := auth.NewAuthService(mockRepo, mockLogger, mockToken)

	mockRepo.On("FindUserByName", mock.Anything, "testuser").Return(nil, nil) // говорим бд его не находить 
	mockToken.On("GenerateJWT", "testuser").Return("mocked_token", nil) // говорим ей создать такой токен для нового юзера
	mockRepo.On("CreateUser", mock.Anything, "testuser", mock.Anything, "mocked_token").Return(nil) // говорим что успешно его создали

	token, err := authService.Authenticate(context.Background(), "testuser", "password")	

	assert.NoError(t, err)

	assert.Equal(t, "mocked_token", token)
}