package apperr

import "net/http"

var (
	ErrBadRequest = AppError{
		Code:   "ERR-001",
		Title:  http.StatusText(http.StatusBadRequest),
		Detail: "bad request",
	}
	ErrUnauthorized = AppError{
		Code:   "ERR-002",
		Title:  http.StatusText(http.StatusUnauthorized),
		Detail: "Invalid credentials",
	}
	ErrForbidden = AppError{
		Code:   "ERR-003",
		Title:  http.StatusText(http.StatusForbidden),
		Detail: "User is not authorized to access this resource",
	}
	ErrNotFound = AppError{
		Code:   "ERR-004",
		Title:  http.StatusText(http.StatusNotFound),
		Detail: "not found",
	}
	ErrInternalServerError = AppError{
		Code:   "ERR-005",
		Title:  http.StatusText(http.StatusInternalServerError),
		Detail: "internal server error",
	}
	ErrUnsupportedMediaType = AppError{
		Code:  "ERR-006",
		Title: http.StatusText(http.StatusUnsupportedMediaType),
	}
	ErrValidationRequest = AppError{
		Code:  "ERR-007",
		Title: "problem validation request",
		Err:   ErrBadRequest,
	}
	ErrDecodingRequest = AppError{
		Code:  "ERR-008",
		Title: "problem decoding request",
		Err:   ErrBadRequest,
	}
	ErrAffectedMoreThanOneRow = AppError{
		Code:   "ERR-009",
		Title:  http.StatusText(http.StatusInternalServerError),
		Detail: "affected more than one row",
		Err:    ErrInternalServerError,
	}
	ErrNoFoundForeignKey = AppError{
		Code:   "ERR-010",
		Title:  http.StatusText(http.StatusInternalServerError),
		Detail: "not found foreign key",
		Err:    ErrInternalServerError,
	}
	ErrNoBearerToken = AppError{
		Code:   "ERR-011",
		Detail: "no bearer token found",
		Err:    ErrUnauthorized,
	}
	ErrInvalidToken = AppError{
		Code:   "ERR-012",
		Detail: "invalid token",
		Err:    ErrUnauthorized,
	}
	ErrSendMail = AppError{
		Code:   "ERR-013",
		Detail: "error send mail",
		Err:    ErrInternalServerError,
	}
	ErrExecuteTemplate = AppError{
		Code:   "ERR-014",
		Detail: "error execute template",
		Err:    ErrInternalServerError,
	}
	ErrUserNotActivated = AppError{
		Code:   "ERR-015",
		Detail: "user not activated",
		Err:    ErrForbidden,
	}
	ErrInvalidOTP = AppError{
		Code:   "ERR-016",
		Detail: "invalid otp",
		Err:    ErrForbidden,
	}
	ErrAlreadyExists = AppError{
		Code:   "ERR-017",
		Detail: "already exists",
		Err:    ErrBadRequest,
	}
)
