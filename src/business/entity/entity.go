package entity

import x "github.com/mafazans/kumparan/lib/errors"

type Meta struct {
	Path       string      `json:"path"`
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Error      *x.AppError `json:"error,omitempty" swaggertype:"primitive,object"`
	Timestamp  string      `json:"timestamp"`
}

type AppError struct {
	Code    x.Code `json:"code"`
	Message string `json:"message"`
}

type CacheControl struct {
	MustRevalidate bool
}
