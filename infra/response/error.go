package response

import (
	"errors"
	"net/http"
)


var (
	// general error
	ErrNotFound = errors.New("data not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbiddenAccess = errors.New("forbiddden access")
)

var (
	// auth error
	ErrEmailRequired = errors.New("email is required")
	ErrEmailInvalid = errors.New("email is invalid")
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordInvalid  = errors.New("password is must be at least 6 characters")
	ErrEmailAlredyExist = errors.New("email already exist")
	ErrPasswordNotMatch = errors.New("password not match")
)

type Error struct {
	Message string
	Code string
	HttpCode int
}

func NewError(message string, code string, httpCode int) Error {
	return Error{
		Message: message,
		Code: code,
		HttpCode: httpCode,
	}
}

func (e Error) Error() string {
	return e.Message
}

var (
	ErrorGeneral = NewError("general error", "99999", http.StatusInternalServerError)
	ErrorBadRequest = NewError("bad request", "40000", http.StatusBadRequest)
	ErrorNotFound = NewError(ErrNotFound.Error(), "40400", http.StatusNotFound)
	ErrorUnauthorized    = NewError(ErrUnauthorized.Error(), "40100", http.StatusUnauthorized)
	ErrorForbiddenAccess = NewError(ErrForbiddenAccess.Error(), "40100", http.StatusForbidden)
)

var (
	ErrorMapping = map[string]Error {
		ErrNotFound.Error(): ErrorNotFound,
		ErrUnauthorized.Error(): ErrorUnauthorized,
		ErrForbiddenAccess.Error(): ErrorForbiddenAccess,
	}
)