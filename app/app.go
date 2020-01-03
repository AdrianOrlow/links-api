package app

import (
	"github.com/AdrianOrlow/links-api/app/handler"
	"github.com/AdrianOrlow/links-api/app/utils"
	"github.com/AdrianOrlow/links-api/config"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize(config *config.Config) {
	err := a.InitializeDatabase(config)
	if err != nil {
		log.Fatal(err)
	}

	err = utils.Initialize(config)
	if err != nil {
		log.Fatal(err)
	}

	handler.InitializeAuth(config)

	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) Run(host string) {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	log.Print("Listening on " + host)
	log.Fatal(http.ListenAndServe(host, handlers.CORS(headersOk, originsOk, methodsOk)(a.Router)))
}
