package service

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"postificus/internal/domain"
	_ "postificus/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCredentialsRepository
type MockCredentialsRepository struct {
	mock.Mock
}

func (m *MockCredentialsRepository) SaveCredentials(ctx context.Context, userID string, platform string, credentials map[string]string) error {
	args := m.Called(ctx, userID, platform, credentials)
	return args.Error(0)
}

func (m *MockCredentialsRepository) GetCredentials(ctx context.Context, userID string, platform string) (*domain.UserCredential, error) {
	args := m.Called(ctx, userID, platform)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserCredential), args.Error(1)
}

func (m *MockCredentialsRepository) GetAllCredentials(ctx context.Context, userID string) ([]domain.UserCredential, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.UserCredential), args.Error(1)
}

// Mock Automation (Swap function)
// In a real scenario, we'd use an interface. For this test, we assume the browser calls will fail if not mocked or we skip the actual browser part if possible.
// However, since we can't easily swap package-level functions in Go without them being variables,
// and we didn't refactor that yet, we will test the 'GetCredentials' part and ensure Breaker is invoked.
//
// TODO: Refactor browser calls to an interface for better testing.
// For now, we will test that credentials failure returns error BEFORE browser is called.

func TestPublishService_HandlePublishTask_CredsError(t *testing.T) {
	// Setup
	os.Unsetenv("MEDIUM_UID")
	os.Unsetenv("MEDIUM_SID")
	os.Unsetenv("MEDIUM_XSRF")

	mockRepo := new(MockCredentialsRepository)
	svc := NewPublishService(mockRepo)

	payload := PublishPayload{
		UserID:   DefaultUserID(),
		Platform: "medium",
		Title:    "Test Title",
		Content:  "Content",
	}
	payloadBytes, _ := json.Marshal(payload)

	// Expectation: GetCredentials fails
	mockRepo.On("GetCredentials", mock.Anything, DefaultUserID(), "medium").Return(nil, errors.New("db error"))

	// Execute
	err := svc.HandlePublishTask(payloadBytes)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "medium credentials missing")

	mockRepo.AssertExpectations(t)
}

func TestPublishService_HandlePublishTask_NoCreds(t *testing.T) {
	// Setup
	t.Setenv("MEDIUM_UID", "")
	t.Setenv("MEDIUM_SID", "")
	t.Setenv("MEDIUM_XSRF", "")

	mockRepo := new(MockCredentialsRepository)
	svc := NewPublishService(mockRepo)

	payload := PublishPayload{
		UserID:   DefaultUserID(),
		Platform: "medium",
		Title:    "Test Title",
		Content:  "Content",
	}
	payloadBytes, _ := json.Marshal(payload)

	// Expectation: GetCredentials returns nil
	mockRepo.On("GetCredentials", mock.Anything, DefaultUserID(), "medium").Return(nil, nil)

	// Execute
	err := svc.HandlePublishTask(payloadBytes)
	t.Logf("HandlePublishTask returned: %v", err)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "medium credentials missing")

	mockRepo.AssertExpectations(t)
}
