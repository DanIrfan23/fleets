package routes

import (
	"fleets/controllers"

	"github.com/gorilla/mux"
)

func UserRoutes(server *mux.Router) {
	server.HandleFunc("/user/login/data", controllers.GetUserLoginDataController()).Methods("GET")
}
