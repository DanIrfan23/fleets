package routes

import (
	"fleets/controllers"

	"github.com/gorilla/mux"
)

func SpbuRoutes(server *mux.Router) {
	server.HandleFunc("/spbu", controllers.GetAllSpbuController()).Methods("GET")
	server.HandleFunc("/spbu/{id}", controllers.GetSpbuByIdController()).Methods("GET")
	server.HandleFunc("/spbu/{id}", controllers.UpdateSpbuByIdController()).Methods("PUT")
	server.HandleFunc("/spbu", controllers.CreateNewSpbuController()).Methods("POST")
	server.HandleFunc("/spbu/{id}", controllers.DeleteSpbuByIdController()).Methods("DELETE")
}
