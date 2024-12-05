package routes

import (
	"fleets/controllers"

	"github.com/gorilla/mux"
)

func CarRoutes(server *mux.Router) {
	server.HandleFunc("/car", controllers.GetAllCarController()).Methods("GET")
	server.HandleFunc("/car/{id}", controllers.GetCarByIdController()).Methods("GET")
	server.HandleFunc("/car", controllers.CreateNewCarController()).Methods("CREATE")
	server.HandleFunc("/car/{id}", controllers.UpdateCarByIdController()).Methods("PUT")
	server.HandleFunc("/car/{id}", controllers.DeleteCarByIdController()).Methods("DELETE")
}
