package user

import (
	"net/http"

	"github.com/moocss/chi-webserver/src/service"
)

func HandleFindUser(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
