package routes

import (
	"fleets/controllers"

	"github.com/gorilla/mux"
)

func VehicleTypeRoutes(server *mux.Router) {
	server.HandleFunc("/type", controllers.GetAllVehicleTypeController()).Methods("GET")
	server.HandleFunc("/type/{id}", controllers.GetVehicleTypeByIdController()).Methods("GET")
	server.HandleFunc("/type/{id}", controllers.UpdateVehicleTypeByIdController()).Methods("PUT")
	server.HandleFunc("/type", controllers.CreateNewVehicleTypeController()).Methods("POST")
	server.HandleFunc("/type/{id}", controllers.DeleteVehicleTypeByIdController()).Methods("DELETE")
}
