package main

import (
	"fmt"
	"net/http"
)

type InternalError struct {
	message string
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("parse %v: internal error", e.message)
}

type appError struct {
	Error   InternalError
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		fmt.Println("error: ", err)
		http.Error(w, err.Error.Error(), err.Code)
	}
}

func testAppHandler(w http.ResponseWriter, r *http.Request) *appError {
	// some validation
	// err := fmt.Errorf("test new error")
	return &appError{
		Message: "test error",
		Code:    505,
		Error:   InternalError{message: "some error"},
	}
}
