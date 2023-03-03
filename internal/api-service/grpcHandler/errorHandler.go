package grpcHandler

import (
	"github.com/pkg/errors"

	"github.com/abhirajranjan/gochat/internal/api-service/serviceErrors"
	"github.com/abhirajranjan/gochat/pkg/logger"
	"github.com/go-playground/validator/v10"
)

type fieldErrMappingFunc func(field validator.FieldError) string

func handleBindingErr(err error) serviceErrors.IErr {
	return handleBindingErrWithFunc(err, DefaultFieldErrMappingFunc)
}

func handleBindingErrWithFunc(err error, f fieldErrMappingFunc) serviceErrors.IErr {
	var validationerr validator.ValidationErrors

	if errors.As(err, &validationerr) {
		errorArray := make(serviceErrors.ErrorArray, len(validationerr))
		for i, fielderr := range validationerr {
			errorArray[i] = serviceErrors.NewBindingErr(fielderr.Field(), f(fielderr))
		}
		return &errorArray
	}
	return nil
}

func DefaultFieldErrMappingFunc(field validator.FieldError) string {
	switch field.Tag() {
	case "required":
		return "this field is required"
	}
	return field.Error()
}

// handle any internal error that are returned by the grpc endpoint
//
// run when response.ErrCode == http.StatusInternalServerError
func handleInternalGrpcEndpointError(err error, logger logger.ILogger, wrap string) {
	logger.Error(errors.Wrap(err, wrap))
}

// handle grpc failover to communicate with grpc endpoint
func handleGrpcServerError(err error, logger logger.ILogger, wrap string) {
	logger.Error(errors.Wrap(err, wrap))
}

// logs validation error
func handleValidationErr(err error, logger logger.ILogger, wrap string) {
	var valerr serviceErrors.ValidationErr
	if errors.As(err, &valerr) {
		logger.Debug(valerr.To_json())
	} else {
		logger.Error(errors.Wrap(err, wrap))
	}
}
