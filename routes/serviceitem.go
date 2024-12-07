package routes

import (
	"fleets/controllers"

	"github.com/gorilla/mux"
)

func ServiceItemRoutes(server *mux.Router) {
	server.HandleFunc("/service-item/lastid", controllers.GetServiceItemLastIdController()).Methods("GET")
	server.HandleFunc("/service-item", controllers.GetAllServiceItemController()).Methods("GET")
	server.HandleFunc("/service-item/{id}", controllers.GetServiceItemByIdController()).Methods("GET")
	server.HandleFunc("/service-item/{id}", controllers.UpdateServiceItemByIdController()).Methods("PUT")
	server.HandleFunc("/service-item", controllers.CreateNewServiceItemController()).Methods("POST")
	server.HandleFunc("/service-item/{id}", controllers.DeleteServiceItemByIdController()).Methods("DELETE")
}
