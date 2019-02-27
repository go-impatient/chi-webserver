package users

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/moocss/chi-webserver/src/model"
	"github.com/moocss/chi-webserver/src/pkg/errno"
	"github.com/moocss/chi-webserver/src/pkg/log"
	"github.com/moocss/chi-webserver/src/service"
)

// HandleFind 控制器， 按照用户名查询
func HandleFind(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// cxt := r.Context()
		username := chi.URLParam(r, "username")

		user, err := userService.FindUser(username)

		if err != nil {
			model.SendResult(w,nil, errno.ErrUserNotFound)
			return
		}

		model.SendResult(w, user, nil)
	}
}

// HandleFindById 控制器， 按照用户ID查询
func HandleFindById(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(chi.URLParam(r, "id"))

		log.Debugf("用户ID: %s", userId)

		user, err := userService.FindUserById(uint64(userId))

		if err != nil {
			model.SendResult(w, nil, errno.ErrUserNotFound)
			return
		}

		model.SendResult(w, user, nil)
	}
}
