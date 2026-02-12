package controller

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"postificus/internal/domain"
	"postificus/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repo for Service
type MockActivityRepo struct {
	mock.Mock
}

func (m *MockActivityRepo) GetCredentials(ctx context.Context, userID int, platform string) (*domain.UserCredential, error) {
	args := m.Called(ctx, userID, platform)
	if args.Get(0) == nil {
		fmt.Printf("Mock GetCredentials returning nil, error: %v\n", args.Error(1))
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserCredential), args.Error(1)
}

func (m *MockActivityRepo) SaveCredentials(ctx context.Context, userID int, platform string, credentials map[string]string) error {
	return nil
}
func (m *MockActivityRepo) GetAllCredentials(ctx context.Context, userID int) ([]domain.UserCredential, error) {
	return nil, nil
}
func (m *MockActivityRepo) UpsertUnifiedPost(ctx context.Context, userID int, post domain.UnifiedPost) error {
	return nil
}
func (m *MockActivityRepo) GetUnifiedPosts(ctx context.Context, userID int, limit int) ([]domain.UnifiedPost, error) {
	return nil, nil
}

func TestActivityController_GetMediumActivity_NoCreds(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/medium/activity", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock Chain: Controller -> Service -> Repo
	mockRepo := new(MockActivityRepo)
	svc := service.NewActivityService(mockRepo) // Uses real service
	ctrl := NewActivityController(svc)

	// User ID assumption (middleware usually sets this, but for test we might need to modify controller or assume default)
	// The controller reads UserID from... wait, let's check code.
	// It extracts from middleware "user_id". We need to set it.
	c.Set("user_id", 1)

	// Expectation: Service calls GetCredentials
	mockRepo.On("GetCredentials", mock.Anything, 1, "medium").Return(nil, nil) // No creds found

	// Execute
	err := ctrl.GetMediumActivity(c)
	t.Logf("GetMediumActivity returned: %v", err)

	// Assert
	if assert.Error(t, err) {
		he, ok := err.(*echo.HTTPError)
		if assert.True(t, ok) {
			assert.Equal(t, http.StatusInternalServerError, he.Code)
			assert.Contains(t, he.Message, "credentials missing")
		}
	}

	mockRepo.AssertExpectations(t)
}
