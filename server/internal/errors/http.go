package errors

import (
	stderrors "errors"
	"net/http"

	"gorm.io/gorm"
)

var (
	errBadRequest = stderrors.New("bad request")
	errNotFound   = stderrors.New("not found")
	errInternal   = stderrors.New("internal server error")
)

type categorizedError struct {
	category error
	cause    error
}

func (err categorizedError) Error() string {
	if err.cause != nil {
		return err.cause.Error()
	}
	if err.category != nil {
		return err.category.Error()
	}
	return ""
}

func (err categorizedError) Unwrap() error {
	return err.cause
}

func (err categorizedError) Is(target error) bool {
	return target != nil && err.category == target
}

func BadRequest(cause error) error {
	return categorizedError{category: errBadRequest, cause: cause}
}

func Internal(cause error) error {
	return categorizedError{category: errInternal, cause: cause}
}

func MapStoreError(err error) error {
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		return categorizedError{category: errNotFound, cause: err}
	}

	return categorizedError{category: errInternal, cause: err}
}

func Write(w http.ResponseWriter, err error) {
	switch {
	case stderrors.Is(err, errBadRequest):
		http.Error(w, errBadRequest.Error(), http.StatusUnprocessableEntity)
	case stderrors.Is(err, errNotFound):
		http.Error(w, errNotFound.Error(), http.StatusNotFound)
	default:
		http.Error(w, errInternal.Error(), http.StatusInternalServerError)
	}
}
