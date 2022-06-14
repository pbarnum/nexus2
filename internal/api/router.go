package api

import (
	"github.com/gorilla/mux"
	mw "github.com/msrevive/nexus2/internal/api/middleware"
	"github.com/msrevive/nexus2/internal/controller"
	"github.com/msrevive/nexus2/internal/ent"
	"github.com/msrevive/nexus2/internal/service"
	"github.com/msrevive/nexus2/internal/system"
)

func NewRouter(cfg *system.ApiConfig, client *ent.Client) mux.Router {
	router := mux.NewRouter()

	//middleware
	router.Use(mw.Log, mw.PanicRecovery)
	if cfg.RateLimit.Enable {
		router.Use(mw.RateLimit(cfg))
	}
	router.Use(mw.Auth(cfg))

	// API routes
	apiRouter := router.PathPrefix(cfg.Core.RootPath).Subrouter()
	for _, c := range []controller.RouteHandler{
		controller.NewApiController(cfg, service.NewDebugService(client)),
		controller.NewCharacterController(cfg, service.NewCharacterService(client)),
	} {
		c.ConfigureRoutes(apiRouter)
	}

	return router
}
