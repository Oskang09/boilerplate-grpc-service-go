package api

import (
	"context"
	"net/http"
	"service/app/errcode"
	v1 "service/protobuf/go/v1"
)

func (e *example) ExampleHandler(ctx context.Context, req *v1.ExampleRequest) (*v1.ExampleResponse, error) {
	res := new(v1.ExampleResponse)

	var i struct {
		Data string `validate:"required"`
	}

	if err := e.bind(req, i); err != nil {
		return nil, responseError(ctx, http.StatusBadRequest, errcode.InvalidRequest, err.Error())
	}

	if err := e.validate(i); err != nil {
		return nil, responseError(ctx, http.StatusUnprocessableEntity, errcode.ValidationErros, err.Error())
	}

	return res, nil
}
