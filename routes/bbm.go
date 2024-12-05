package routes

import (
	"fleets/controllers"

	"github.com/gorilla/mux"
)

func BbmRoutes(server *mux.Router) {
	server.HandleFunc("/bbm", controllers.GetAllBbmController()).Methods("GET")
	server.HandleFunc("/bbm/{id}", controllers.GetBbmByIdController()).Methods("GET")
	server.HandleFunc("/bbm/{id}", controllers.UpdateBbmByIdController()).Methods("PUT")
	server.HandleFunc("/bbm", controllers.CreateNewBbmController()).Methods("POST")
	server.HandleFunc("/bbm/{id}", controllers.DeleteBbmByIdController()).Methods("DELETE")
}
