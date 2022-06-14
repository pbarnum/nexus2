package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/msrevive/nexus2/internal/api/response"
	"github.com/msrevive/nexus2/internal/service"
	"github.com/msrevive/nexus2/internal/system"
)

var _ Controller = (*apiController)(nil)
var _ RouteHandler = (*apiController)(nil)

type apiController struct {
	controller
}

func NewApiController(cfg *system.ApiConfig, svc service.DebugService) Controller {
	return &apiController{
		controller: controller{
			cfg:     cfg,
			service: svc,
		},
	}
}

func (c *apiController) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/", c.TestRoot).Methods(http.MethodGet)
	router.HandleFunc("/ping", c.GetPing).Methods(http.MethodGet)
	router.HandleFunc("/map/{name}/{hash}", c.GetMapVerify).Methods(http.MethodGet)
	router.HandleFunc("/ban/{steamid:[0-9]+}", c.GetBanVerify).Methods(http.MethodGet)
	router.HandleFunc("/sc/{hash}", c.GetSCVerify).Methods(http.MethodGet)
}

//GET map/{name}/{hash}
func (c *apiController) GetMapVerify(w http.ResponseWriter, r *http.Request) {
	if !c.cfg.Verify.EnforceMap {
		response.Result(w, true)
		return
	}

	vars := mux.Vars(r)
	name := vars["name"]
	hash, err := strconv.ParseUint(vars["hash"], 10, 32)
	if err != nil {
		response.BadRequest(w, err)
		return
	}

	if c.cfg.VerifyMapName(name, uint32(hash)) {
		response.Result(w, true)
		return
	}

	response.Result(w, false)
}

//GET ban/{steamid}
//in this case false means player isn't banned
func (c *apiController) GetBanVerify(w http.ResponseWriter, r *http.Request) {
	if !c.cfg.Verify.EnforceBan {
		response.Result(w, false)
		return
	}

	vars := mux.Vars(r)
	steamid := vars["steamid"]

	if c.cfg.IsSteamIdBanned(steamid) {
		response.Result(w, true)
		return
	}

	response.Result(w, false)
}

//GET sc/{hash}
func (c *apiController) GetSCVerify(w http.ResponseWriter, r *http.Request) {
	if !c.cfg.Verify.EnforceSC {
		response.Result(w, true)
		return
	}

	vars := mux.Vars(r)
	hash, err := strconv.ParseUint(vars["hash"], 10, 32)
	if err != nil {
		response.BadRequest(w, err)
		return
	}

	if c.cfg.Verify.SCHash == uint32(hash) {
		response.Result(w, true)
		return
	}

	response.Result(w, false)
}

//GET ping
func (c *apiController) GetPing(w http.ResponseWriter, r *http.Request) {
	response.Result(w, true)
}
