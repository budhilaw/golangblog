package mocks

import (
	"context"

	"github.com/budhilaw/blog/models"
	mock "github.com/stretchr/testify/mock"
)

// Usecase mock
type Usecase struct {
	mock.Mock
}

// Store provides a mock function with given fields ctx, m
func (_m *Usecase) Store(c context.Context, m *models.User) error {
	ret := _m.Called(c, m)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User) error); ok {
		r0 = rf(c, m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
