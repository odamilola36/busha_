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

// @Summary get all characters for a single movie
// @Description  gets all characters for a single movie
// @Produce  json
// @Param movie_id path int true "Movie ID"
// @Param sort_by query string false "Sort by height or name or gender"
// @Param order query string false "asc or desc"
// @Param filter query string false "male, female, n/a"
// @Success 200 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /api/v1/characters/{movie_id} [get]
func (c *charactersHandler) GetCharacters(w http.ResponseWriter, r *http.Request) {

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
