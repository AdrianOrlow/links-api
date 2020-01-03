package app

import (
	"github.com/AdrianOrlow/links-api/app/handler"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

func (a *App) setRouters() {
	r := a.Router
	a.Get(r, "/{hashIdOrLink}", a.handleRequest(handler.HandleRedirect))
	a.Get(r, "/oauth/google/login", a.handleRequest(handler.HandleGoogleLogin))
	a.Get(r, "/oauth/google/callback", a.handleRequest(handler.HandleGoogleCallback))

	v1 := r.PathPrefix("/v1").Subrouter()
	a.Get(v1, "/links", a.adminOnly(handler.GetAllLinks))
	a.Post(v1, "/links", a.adminOnly(handler.CreateLink))
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(h RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(a.DB, w, r)
	}
}

func (a *App) Get(r *mux.Router, path string, f func(w http.ResponseWriter, r *http.Request)) {
	r.HandleFunc(path, f).Methods("GET")
}

func (a *App) Post(r *mux.Router, path string, f func(w http.ResponseWriter, r *http.Request)) {
	r.HandleFunc(path, f).Methods("POST")
}

func (a *App) Put(r *mux.Router, path string, f func(w http.ResponseWriter, r *http.Request)) {
	r.HandleFunc(path, f).Methods("PUT")
}

func (a *App) Delete(r *mux.Router, path string, f func(w http.ResponseWriter, r *http.Request)) {
	r.HandleFunc(path, f).Methods("DELETE")
}