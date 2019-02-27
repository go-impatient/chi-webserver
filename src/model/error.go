package model

import (
	"net/http"

	"github.com/moocss/chi-webserver/src/pkg/render"
)

// The Error contains error relevant information.
type Error struct {
	// The general error message
	Error string `json:"error"`

	// The http error code.
	ErrorCode int `json:"error_code"`

	// The http error code.
	ErrorDescription string `json:"error_description"`
}

func SendError(w http.ResponseWriter) {
	render.JSON(w, &Error{
		Error:            http.StatusText(http.StatusNotFound),
		ErrorCode:        http.StatusNotFound,
		ErrorDescription: "page not found",
	}, http.StatusNotFound)
}

