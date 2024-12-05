package routes

import (
	"fleets/controllers"

	"github.com/gorilla/mux"
)

func RelationshipRoutes(server *mux.Router) {
	server.HandleFunc("/relationship/user/{id}", controllers.GetRelationshipByUserIdController()).Methods("GET")
	server.HandleFunc("/relationship", controllers.CreateRelationshipController()).Methods("POST")
	server.HandleFunc("/relationship/{id}", controllers.GetRelationshipByIdController()).Methods("GET")
	server.HandleFunc("/relationship/{id}", controllers.UpdateStatusRelationshipController()).Methods("PUT")
}
