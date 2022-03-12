package service

import (
	"busha_/dto"
	"busha_/models"
	"busha_/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"net/http"
	"sort"
	"time"
)

const baseUrl = "https://swapi.dev/api/films"

type MoviesService interface {
	GetMovies() ([]models.Movie, error)
	getMoviesFromService() ([]models.SwappiMovie, error)
	getTransformedMovies() ([]models.Movie, error)
}

type moviesService struct {
	redisClient    *redis.Client
	commentService CommentService
}

func NewMoviesService(redisClient *redis.Client, commentService CommentService) *moviesService {
	return &moviesService{
		redisClient:    redisClient,
		commentService: commentService,
	}
}

func (s *moviesService) GetMovies() ([]models.Movie, error) {

	var movies []models.Movie
	var err error

	cache, err := s.redisClient.Get(context.TODO(), "movies").Result()

	if err == redis.Nil {
		movies, err = s.getTransformedMovies()
		if err != nil {
			utils.ErrorLineLogger(err)
			return movies, err
		}
		marshal, err := json.Marshal(movies)
		if err != nil {
			utils.ErrorLineLogger(err)
			return movies, err
		}
		s.redisClient.Set(context.TODO(), "movies", marshal, 30*time.Second)
	} else {
		err = json.Unmarshal([]byte(cache), &movies)
		if err != nil {
			utils.ErrorLineLogger(err)
			return movies, err
		}
	}
	return movies, nil
}

func (s *moviesService) getTransformedMovies() ([]models.Movie, error) {
	var movies []models.Movie

	moviesS, err := s.getMoviesFromService()

	if err != nil {
		utils.ErrorLineLogger(err)
		return movies, err
	}

	for _, movie := range moviesS {
		tr := movie.ToMovie()
		comments, err := s.commentService.GetCommentsList(movie.EpisodeId)
		if err != nil {
			utils.ErrorLineLogger(err)
			return movies, err
		}
		tr.CommentCount = len(comments)
		movies = append(movies, tr)
	}

	sort.Slice(movies, func(i, j int) bool {
		return movies[i].ReleaseDate.After(movies[j].ReleaseDate)
	})

	return movies, nil

}

func (s *moviesService) getMoviesFromService() ([]models.SwappiMovie, error) {

	var result dto.SwappiMoviesResponse

	res, err := http.Get(baseUrl)

	if err != nil {
		utils.ErrorLineLogger(err)
		return result.Results, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("error in getting movies from service, status code: " + res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(&result)

	if err != nil {
		utils.ErrorLineLogger(err)
		return result.Results, err
	}

	return result.Results, nil
}
