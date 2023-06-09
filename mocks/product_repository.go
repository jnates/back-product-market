// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import (
	model "backend_crudgo/domain/products/domain/model"
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "backend_crudgo/types"
)

// ProductRepository is an autogenerated mock type for the ProductRepository type
type ProductRepository struct {
	mock.Mock
}

// CreateProduct provides a mock function with given fields: ctx, product
func (_m *ProductRepository) CreateProduct(ctx context.Context, product *model.Product) (*types.CreateResponse, error) {
	ret := _m.Called(ctx, product)

	var r0 *types.CreateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Product) (*types.CreateResponse, error)); ok {
		return rf(ctx, product)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Product) *types.CreateResponse); ok {
		r0 = rf(ctx, product)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.CreateResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Product) error); ok {
		r1 = rf(ctx, product)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteProduct provides a mock function with given fields: ctx, id
func (_m *ProductRepository) DeleteProduct(ctx context.Context, id string) (*types.GenericResponse, error) {
	ret := _m.Called(ctx, id)

	var r0 *types.GenericResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*types.GenericResponse, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *types.GenericResponse); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GenericResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProduct provides a mock function with given fields: ctx, id
func (_m *ProductRepository) GetProduct(ctx context.Context, id string) (*types.GenericResponse, error) {
	ret := _m.Called(ctx, id)

	var r0 *types.GenericResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*types.GenericResponse, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *types.GenericResponse); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GenericResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProducts provides a mock function with given fields: ctx
func (_m *ProductRepository) GetProducts(ctx context.Context) (*types.GenericResponse, error) {
	ret := _m.Called(ctx)

	var r0 *types.GenericResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*types.GenericResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *types.GenericResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GenericResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProduct provides a mock function with given fields: ctx, id, product
func (_m *ProductRepository) UpdateProduct(ctx context.Context, id string, product *model.Product) (*types.GenericResponse, error) {
	ret := _m.Called(ctx, id, product)

	var r0 *types.GenericResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.Product) (*types.GenericResponse, error)); ok {
		return rf(ctx, id, product)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.Product) *types.GenericResponse); ok {
		r0 = rf(ctx, id, product)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GenericResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *model.Product) error); ok {
		r1 = rf(ctx, id, product)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewProductRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewProductRepository creates a new instance of ProductRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProductRepository(t mockConstructorTestingTNewProductRepository) *ProductRepository {
	mock := &ProductRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
