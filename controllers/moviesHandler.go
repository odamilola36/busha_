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
