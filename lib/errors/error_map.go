package errors

import (
	"fmt"
	"net/http"
	"strings"
)

var svcError map[ServiceType]ErrorMessage

// ErrorMessage - Mapping Error Code as Human Message
type ErrorMessage map[Code]Message

type ServiceType int

// Message - Error Message Support Multi Language
type Message struct {
	StatusCode    int    `json:"status_code"`
	EN            string `json:"en"`
	ID            string `json:"id"`
	HasAnnotation bool
}

const (
	COMMON ServiceType = 1
)

func init() {
	svcError = map[ServiceType]ErrorMessage{
		COMMON: ErrorMessages,
	}
}

// TODO: We can change the status and error code as we see fit
var ErrorMessages = ErrorMessage{
	CodeValueInvalid: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeValueInvalid EN`,
		ID:         `Error CodeValueInvalid ID`,
	},
	CodeSQLBuilder: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeSQLBuilder EN`,
		ID:         `Error CodeSQLBuilder ID`,
	},
	CodeSQLTxBegin: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeSQLTxBegin EN`,
		ID:         `Error CodeSQLTxBegin ID`,
	},
	CodeSQLCreate: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeSQLCreate EN`,
		ID:         `Error CodeSQLCreate ID`,
	},
	CodeSQLTxCommit: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeSQLTxCommit EN`,
		ID:         `Error CodeSQLTxCommit ID`,
	},
	CodeSQLCannotRetrieveLastInsertID: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeSQLCannotRetrieveLastInsertID EN`,
		ID:         `Error CodeSQLCannotRetrieveLastInsertID ID`,
	},
	CodeHTTPBadRequest: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeHTTPBadRequest EN`,
		ID:         `Error CodeHTTPBadRequest ID`,
	},
	CodeHTTPNotFound: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeHTTPNotFound EN`,
		ID:         `Error CodeHTTPNotFound ID`,
	},
	CodeHTTPInternalServerError: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeHTTPInternalServerError EN`,
		ID:         `Error CodeHTTPInternalServerError ID`,
	},
	CodeHTTPUnmarshal: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeHTTPUnmarshal EN`,
		ID:         `Error CodeHTTPUnmarshal ID`,
	},
	CodeUsecase: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeUsecase EN`,
		ID:         `Error CodeUsecase ID`,
	},
	CodeDomain: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeDomain EN`,
		ID:         `Error CodeDomain ID`,
	},
	CodeCacheMarshal: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeCacheMarshal EN`,
		ID:         `Error CodeCacheMarshal ID`,
	},
	CodeCacheUnmarshal: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeCacheUnmarshal EN`,
		ID:         `Error CodeCacheUnmarshal ID`,
	},
	CodeCacheGetSimpleKey: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeCacheGetSimpleKey EN`,
		ID:         `Error CodeCacheGetSimpleKey ID`,
	},
	CodeCacheSetSimpleKey: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeCacheSetSimpleKey EN`,
		ID:         `Error CodeCacheSetSimpleKey ID`,
	},
	CodeCacheDeleteSimpleKey: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeCacheDeleteSimpleKey EN`,
		ID:         `Error CodeCacheDeleteSimpleKey ID`,
	},
	CodeCacheSetExpiration: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeCacheSetExpiration EN`,
		ID:         `Error CodeCacheSetExpiration ID`,
	},
	CodeCacheDecode: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeCacheDecode EN`,
		ID:         `Error CodeCacheDecode ID`,
	},
	CodeCacheLockNotAcquired: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeCacheLockNotAcquired EN`,
		ID:         `Error CodeCacheLockNotAcquired ID`,
	},
	CodeCacheLockFailed: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeCacheLockFailed EN`,
		ID:         `Error CodeCacheLockFailed ID`,
	},
	CodeCacheInvalidCastType: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeCacheInvalidCastType EN`,
		ID:         `Error CodeCacheInvalidCastType ID`,
	},
	CodeCacheNotFound: Message{
		StatusCode: http.StatusBadRequest,
		EN:         `Error CodeCacheNotFound EN`,
		ID:         `Error CodeCacheNotFound ID`,
	},
}

// AppError - Application Error Structure
type AppError struct {
	Code       Code    `json:"code"`
	Message    string  `json:"message"`
	DebugError *string `json:"debug,omitempty"`
	sys        error
}

func (e *AppError) Error() string {
	return e.sys.Error()
}

// Compile - Get Error Code and HTTP Status
// Common --> Service --> Default
func Compile(service ServiceType, err error, lang string) (int, AppError) {
	// Get Error Code
	code := ErrCode(err)

	// Get Common Error
	if errMessage, ok := svcError[COMMON][code]; ok {
		msg := errMessage.ID
		if lang == "EN" {
			msg = errMessage.EN
		}
		return errMessage.StatusCode, AppError{
			Code:    code,
			Message: msg,
		}
	}

	// Get Service Error
	if errMessages, ok := svcError[service]; ok {
		if errMessage, ok := errMessages[code]; ok {
			msg := errMessage.ID
			if lang == "EN" {
				msg = errMessage.EN
			}

			if errMessage.HasAnnotation {
				args := fmt.Sprintf("%q", err.Error())
				index := strings.Index(args, `\n`)
				if index > 0 {
					args = strings.TrimSpace(args[1:index])
				}
				msg = fmt.Sprintf(msg, args)
			}

			return errMessage.StatusCode, AppError{
				Code:    code,
				Message: msg,
			}
		}
		return http.StatusInternalServerError, AppError{
			Code:    code,
			Message: "error message not defined!",
		}
	}

	// Set Default Error
	return http.StatusInternalServerError, AppError{
		Code:    code,
		Message: "service error not defined!",
	}
}
