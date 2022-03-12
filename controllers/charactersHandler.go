package controllers

import (
	"busha_/dto"
	"busha_/service"
	"github.com/gorilla/mux"
	"net/http"
)

type CharactersHandler interface {
	GetCharacters(w http.ResponseWriter, r *http.Request)
}

type charactersHandler struct {
	charactersService service.CharacterService
}

func NewCharactersHandler(c service.CharacterService) CharactersHandler {
	return &charactersHandler{
		charactersService: c,
	}
}

func (c *charactersHandler) GetCharacters(w http.ResponseWriter, r *http.Request) {
	//r.Context()

	query := r.URL.Query()
	req := map[string]string{
		"sort_by": query.Get("sort_by"),
		"order":   query.Get("order"),
		"filter":  query.Get("filter"),
		"movieId": mux.Vars(r)["id"],
	}
	charactersResp, err := c.charactersService.GetCharacters(req)

	if err != nil {
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Error:   err,
			Data:    nil,
		})
		return
	}

	dto.BuildResponse(w, &dto.Response{
		Status:  http.StatusOK,
		Message: "Successfully retrieved characters",
		Error:   nil,
		Data:    charactersResp,
	})
}
