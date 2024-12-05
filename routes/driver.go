package routes

import (
	"fleets/controllers"

	"github.com/gorilla/mux"
)

func DriverRoutes(server *mux.Router) {
	server.HandleFunc("/driver", controllers.GetAllDriverController()).Methods("GET")
	server.HandleFunc("/driver/{id}", controllers.GetDriverByIdController()).Methods("GET")
	server.HandleFunc("/driver/{id}", controllers.UpdateDriverByIdController()).Methods("PUT")
	server.HandleFunc("/driver", controllers.CreateNewDriverController()).Methods("POST")
	server.HandleFunc("/driver/{id}", controllers.DeleteDriverByIdController()).Methods("DELETE")
}
