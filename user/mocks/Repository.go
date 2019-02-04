package mocks

import (
	"context"

	"github.com/budhilaw/blog/models"
	"github.com/stretchr/testify/mock"
)

// Repository mock
type Repository struct {
	mock.Mock
}

// Fetch provides a mock function with given fields: ctx, num
func (_m *Repository) Fetch(c context.Context, num int64) ([]*models.User, error) {
	ret := _m.Called(c, num)

	var r0 []*models.User
	if rf, ok := ret.Get(0).(func(context.Context, int64) []*models.User); ok {
		r0 = rf(c, num)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(c, num)
	} else {
		r1 = ret.Error(2)
	}

	return r0, r1

}

// Store provides a mock function with given fields: ctx, m
func (_m *Repository) Store(c context.Context, m *models.User) error {
	ret := _m.Called(c, m)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User) error); ok {
		r0 = rf(c, m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Login provides a mock function with given fields: context, string, string, models
func (_m *Repository) Login(c context.Context, email string, pass string, m *models.User) error {
	ret := _m.Called(c, m)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *models.User) error); ok {
		r0 = rf(c, email, pass, m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
