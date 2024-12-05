package main

import (
	"fleets/configs"
	"fleets/controllers"
	"fleets/middlewares"
	"fleets/responses"
	"fleets/routes"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	configs.InitDatabase()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		responses.SuccessResponse(w, http.StatusOK, "Welcome to GO-FLEETS REST API", "")
	}).Methods("GET")

	r.HandleFunc("/login", controllers.LoginController()).Methods("POST")

	protectedRouter := r.PathPrefix("/api").Subrouter()
	protectedRouter.Use(middlewares.JwtMiddleware)

	routes.MapRoutes(protectedRouter)

	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")), handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Cookie"}),
		handlers.AllowCredentials(),
	)(r))
}
