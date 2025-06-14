package router

import (
	"errors"
	"net/http"

	apiError "github.com/dsaime/goods-and-projects/internal/controller/http2/api_error"
	"github.com/dsaime/goods-and-projects/internal/domain"
)

func errHttpResponse(err error) (data ResponseError, httpStatus int) {
	apiErr := apiErrorFrom(err)
	responseError := ResponseError{
		Code:    apiErr.Code(),
		Message: apiErr.Message(),
		Details: nil,
	}
	if withDetails, ok := (err).(interface{ Details() map[string]any }); ok {
		responseError.Details = withDetails.Details()
	}

	return responseError, errHttpStatus(apiErr)
}

func errHttpStatus(err apiError.Error) int {
	switch code := err.Code(); code {
	case ErrCommonNotFound.Code():
		return http.StatusNotFound
	case ErrInternalJsonEncode.Code(),
		ErrInternalWriteResponse.Code():
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

func apiErrorFrom(err error) apiError.Error {
	var apiErr apiError.Error
	if errors.As(err, &apiErr) {
		return apiErr
	}

	switch {
	case errors.Is(err, ErrJsonMarshalResponseData):
		return ErrInternalJsonEncode
	case errors.Is(err, ErrWriteResponseBytes):
		return ErrInternalWriteResponse
	case errors.Is(err, domain.ErrGoodNotFound):
		return ErrCommonNotFound
	}

	return apiError.New(0, err.Error())
}

var (
	ErrInternalJsonEncode    = apiError.New(1, "errors.internal.jsonEncode")
	ErrInternalWriteResponse = apiError.New(1, "errors.internal.writeResponse")
	ErrCommonNotFound        = apiError.New(3, "errors.common.notFound")
)
