package routes

import (
	"fleets/controllers"

	"github.com/gorilla/mux"
)

func SiteRoutes(server *mux.Router) {
	server.HandleFunc("/site", controllers.GetAllSiteController()).Methods("GET")
	server.HandleFunc("/site/{id}", controllers.GetSiteByIdController()).Methods("GET")
	server.HandleFunc("/site/{id}", controllers.UpdateSiteByIdController()).Methods("PUT")
	server.HandleFunc("/site", controllers.CreateNewSiteController()).Methods("POST")
	server.HandleFunc("/site/{id}", controllers.DeleteSiteByIdController()).Methods("DELETE")
}
