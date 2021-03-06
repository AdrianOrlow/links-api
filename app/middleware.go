package app

import (
	"github.com/AdrianOrlow/links-api/app/utils"
	"net/http"
)

func (a *App) adminOnly(h RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		status, err := utils.VerifyLoginJWT(authToken)
		if err != nil {
			w.WriteHeader(status)
			return
		}
		h(a.DB, w, r)
	}
}
