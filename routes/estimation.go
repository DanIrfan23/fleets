package routes

import (
	"fleets/controllers"

	"github.com/gorilla/mux"
)

func EstimationRoutes(server *mux.Router) {
	server.HandleFunc("/estimation", controllers.GetAllEstimationController()).Methods("GET")
	server.HandleFunc("/estimation/{id}", controllers.GetEstimationByIdController()).Methods("GET")
	server.HandleFunc("/estimation/{id}", controllers.UpdateEstimationByIdController()).Methods("PUT")
	server.HandleFunc("/estimation", controllers.CreateNewEstimationController()).Methods("POST")
	server.HandleFunc("/estimation/{id}", controllers.DeleteEstimationByIdController()).Methods("DELETE")
}
