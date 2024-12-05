package routes

import (
	"github.com/gorilla/mux"
)

func MapRoutes(server *mux.Router) {
	UserRoutes(server)
	RelationshipRoutes(server)

	// Master
	CarRoutes(server)
	VehicleTypeRoutes(server)
	DriverRoutes(server)
	ServiceItemRoutes(server)
	SpbuRoutes(server)
	SiteRoutes(server)
	EstimationRoutes(server)
	BbmRoutes(server)
}
