package users

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/moocss/chi-webserver/src/model"
	"github.com/moocss/chi-webserver/src/pkg/errno"
	"github.com/moocss/chi-webserver/src/service"
)

func HandleDelete(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(chi.URLParam(r, "id"))
		user := &model.User{}
		user.ID = uint64(userId)

		if err := userService.DeleteUser(user); err != nil {
			model.SendResult(w, nil, errno.ErrDatabase, )
			return
		}

		model.SendResult(w, nil, nil)
	}
}
