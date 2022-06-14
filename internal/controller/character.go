package controller

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/msrevive/nexus2/internal/api/response"
	"github.com/msrevive/nexus2/internal/ent"
	"github.com/msrevive/nexus2/internal/helper"
	"github.com/msrevive/nexus2/internal/log"
	"github.com/msrevive/nexus2/internal/service"
	"github.com/msrevive/nexus2/internal/system"
)

var _ Controller = (*charController)(nil)
var _ RouteHandler = (*controller)(nil)

type charController struct {
	controller
	service service.CharacterService
}

func NewCharacterController(cfg *system.ApiConfig, svc service.CharacterService) Controller {
	return &charController{
		controller: controller{
			cfg: cfg,
		},
		service: svc,
	}
}

func (c *charController) ConfigureRoutes(router *mux.Router) {
	r := router.PathPrefix("/character").Subrouter()
	r.HandleFunc("/", c.GetAllCharacters).Methods(http.MethodGet)
	r.HandleFunc("/id/{uid}", c.GetCharacterByID).Methods(http.MethodGet)
	r.HandleFunc("/{steamid:[0-9]+}", c.GetCharacters).Methods(http.MethodGet)
	r.HandleFunc("/{steamid:[0-9]+}/{slot:[0-9]}", c.GetCharacter).Methods(http.MethodGet)
	r.HandleFunc("/export/{steamid:[0-9]+}/{slot:[0-9]}", c.ExportCharacter).Methods(http.MethodGet)
	r.HandleFunc("/", c.PostCharacter).Methods(http.MethodPost)
	r.HandleFunc("/{uid}", c.PutCharacter).Methods(http.MethodPut)
	r.HandleFunc("/{uid}", c.DeleteCharacter).Methods(http.MethodDelete)
}

//GET /character/
func (c *charController) GetAllCharacters(w http.ResponseWriter, r *http.Request) {
	chars, err := c.service.CharactersGetAll(r.Context())
	if err != nil {
		log.Log.Errorln(err)
		response.BadRequest(w, err)
		return
	}

	response.OK(w, chars)
}

//GET /character/{steamid}
func (c *charController) GetCharacters(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	steamid := vars["steamid"]

	chars, err := c.service.CharactersGetBySteamid(r.Context(), steamid)
	if err != nil {
		log.Log.Errorln(err)
		response.BadRequest(w, err)
		return
	}

	isBanned := c.cfg.EnforceAndVerifyBanned(steamid)
	isAdmin := c.cfg.IsSteamIdAdmin(steamid)
	response.OKChar(w, isBanned, isAdmin, chars)
}

//GET /character/{steamid}/{slot}
func (c *charController) GetCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	steamid := vars["steamid"]
	slot, err := strconv.Atoi(vars["slot"])
	if err != nil {
		log.Log.Errorln(err)
		response.BadRequest(w, err)
		return
	}

	char, err := c.service.CharacterGetBySteamidSlot(r.Context(), steamid, slot)
	if err != nil {
		log.Log.Errorln(err)
		response.BadRequest(w, err)
		return
	}

	isBanned := c.cfg.EnforceAndVerifyBanned(steamid)
	isAdmin := c.cfg.IsSteamIdAdmin(steamid)
	response.OKChar(w, isBanned, isAdmin, char)
}

//GET /character/export/{steamid}/{slot}
func (c *charController) ExportCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	steamid := vars["steamid"]
	slot, err := strconv.Atoi(vars["slot"])
	if err != nil {
		log.Log.Errorln(err)
		response.BadRequest(w, err)
		return
	}

	char, err := c.service.CharacterGetBySteamidSlot(r.Context(), steamid, slot)
	if err != nil {
		log.Log.Errorln(err)
		response.BadRequest(w, err)
		return
	}

	file, path, err := helper.GenerateCharFile(steamid, slot, char.Data)
	if err != nil {
		log.Log.Errorln(err)
		response.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", path))
	io.Copy(w, file)
}

//GET /character/id/{uid}
func (c *charController) GetCharacterByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := uuid.Parse(vars["uid"])
	if err != nil {
		log.Log.Errorln(err)
		response.BadRequest(w, err)
		return
	}

	char, err := c.service.CharacterGetByID(r.Context(), uid)
	if err != nil {
		log.Log.Errorln(err)
		response.BadRequest(w, err)
		return
	}

	isBanned := c.cfg.EnforceAndVerifyBanned(char.Edges.Player.Steamid)
	isAdmin := c.cfg.IsSteamIdAdmin(char.Edges.Player.Steamid)
	response.OKChar(w, isBanned, isAdmin, char)
}

//POST /character/
func (c *charController) PostCharacter(w http.ResponseWriter, r *http.Request) {
	var newChar ent.Character
	err := json.NewDecoder(r.Body).Decode(&newChar)
	if err != nil {
		log.Log.Errorln(err)
		response.BadRequest(w, err)
		return
	}

	char, err := c.service.CharacterCreate(r.Context(), newChar)
	if err != nil {
		log.Log.Errorln(err)
		response.Error(w, err)
		return
	}

	response.OK(w, char)
}

//PUT /character/{uid}
func (c *charController) PutCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := uuid.Parse(vars["uid"])
	if err != nil {
		log.Log.Errorln(err)
		response.BadRequest(w, err)
		return
	}

	var updateChar ent.Character
	err = json.NewDecoder(r.Body).Decode(&updateChar)
	if err != nil {
		log.Log.Errorln(err)
		response.BadRequest(w, err)
		return
	}

	char, err := c.service.CharacterUpdate(r.Context(), uid, updateChar)
	if err != nil {
		log.Log.Errorln(err)
		response.Error(w, err)
		return
	}

	response.OK(w, char)
}

//DELETE /character/{uid}
func (c *charController) DeleteCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := uuid.Parse(vars["uid"])
	if err != nil {
		log.Log.Errorln(err)
		response.BadRequest(w, err)
		return
	}

	err = c.service.CharacterDelete(r.Context(), uid)
	if err != nil {
		log.Log.Errorln(err)
		response.Error(w, err)
		return
	}

	response.OK(w, uid)
}
