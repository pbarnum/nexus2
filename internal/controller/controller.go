package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/msrevive/nexus2/internal/api/response"
	"github.com/msrevive/nexus2/internal/service"
	"github.com/msrevive/nexus2/internal/system"
)

type RouteHandler interface {
	ConfigureRoutes(router *mux.Router)
}

var _ RouteHandler = (*controller)(nil)

type Controller interface {
	RouteHandler
}

type controller struct {
	cfg     *system.ApiConfig
	service service.DebugService
}

func (c *controller) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/", c.TestRoot).Methods(http.MethodGet)
}

func (c *controller) TestRoot(w http.ResponseWriter, r *http.Request) {
	if c.cfg.Core.DebugMode {
		if err := c.service.Debug(r.Context()); err != nil {
			response.BadRequest(w, err)
			return
		}
	}

	response.Result(w, true)
}
