package grpchandler

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/sixsat/user-manager-be-challenge/domain"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func invalidRequestError(err error) error {
	var violations []*errdetails.BadRequest_FieldViolation
	if validationErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
		for _, fieldErr := range validationErrs {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       fieldErr.Field(),
				Description: fieldErr.Error(),
			})
		}
	}

	st, _ := status.
		New(codes.InvalidArgument, "invalid request").
		WithDetails(&errdetails.BadRequest{FieldViolations: violations})

	return st.Err()
}

func serviceError(err error) error {
	if bizErr, ok := errors.AsType[domain.BizErr](err); ok {
		code := codes.FailedPrecondition
		if bizErr == domain.BizErrUserAlreadyExists {
			code = codes.AlreadyExists
		}
		return status.Error(code, bizErr.Error())
	}

	if errors.Is(err, domain.ErrUserNotFound) {
		return status.Error(codes.NotFound, domain.ErrUserNotFound.Error())
	}

	return status.Error(codes.Internal, "internal server error")
}
