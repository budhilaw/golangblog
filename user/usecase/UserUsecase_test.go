package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"

	"github.com/budhilaw/blog/models"

	mocks "github.com/budhilaw/blog/user/mocks"
	ucase "github.com/budhilaw/blog/user/usecase"
)

func TestFetch(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	now := time.Now()
	mockUser := &models.User{
		Name:      "Kai",
		Email:     "kai@gmail.com",
		Password:  "1550",
		CreadedAt: now,
		UpdatedAt: now,
	}

	mockListUser := make([]*models.User, 0)
	mockListUser = append(mockListUser, mockUser)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Fetch", mock.Anything, mock.AnythingOfType("int")).Return(mockListUser, nil).Once()
		u := ucase.New(mockUserRepo, time.Second*2)
		num := int64(1)
		list, err := u.Fetch(context.TODO(), num)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListUser))

		mockUserRepo.AssertExpectations(t)
	})
}
func TestStore(t *testing.T) {
	mockUserRepo := new(mocks.Repository)

	now := time.Now()
	mockUser := models.User{
		Name:      "Kai",
		Email:     "kai@gmail.com",
		Password:  "1550",
		CreadedAt: now,
		UpdatedAt: now,
	}

	t.Run("success", func(t *testing.T) {
		tempMockUser := mockUser
		tempMockUser.ID = 0
		mockUserRepo.On("Store", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()

		u := ucase.New(mockUserRepo, time.Second*2)
		err := u.Store(context.TODO(), &tempMockUser)

		assert.NoError(t, err)
		assert.Equal(t, mockUser.Email, tempMockUser.Email)
		mockUserRepo.AssertExpectations(t)
	})

}
