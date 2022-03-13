package controllers

import (
	"busha_/dto"
	"busha_/models"
	"busha_/service"
	"net/http"
)

type MovieController interface {
	MovieListController(w http.ResponseWriter, r *http.Request)
}

type movieController struct {
	movieService service.MoviesService
}

func NewMovieController(m service.MoviesService) *movieController {
	return &movieController{
		movieService: m,
	}
}

// @Summary      Get all movies
// @Description  Get all movies in order of their release date from earliest to newest in the cache or from swapi if the cache is empty
// @Produce  json
// @Success 200 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /api/v1/movies [get]
func (m *movieController) MovieListController(w http.ResponseWriter, r *http.Request) {
	var movies []models.Movie
	var err error

	movies, err = m.movieService.GetMovies()

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
		Message: "Successfully retrieved movies",
		Error:   nil,
		Data:    movies,
	})
}
