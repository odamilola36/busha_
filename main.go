package main

import (
	"busha_/config"
	"busha_/controllers"
	_ "busha_/docs"
	"busha_/service"
	"busha_/utils"
	"context"
	"github.com/gorilla/mux"
	_ "github.com/swaggo/files"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	moviesController    controllers.MovieController
	commentsController  controllers.CommentHandler
	characterController controllers.CharactersHandler
	PID                 = os.Getppid()
	Port                = ":8080"
)

func init() {
	dbConf := config.SetupDb()
	redisConf := config.SetupRedis()
	redisClient := redisConf.RedisClient()
	redisClient.FlushAll(context.Background())
	Db := dbConf.Database()
	err := dbConf.Migrate()
	if err != nil {
		utils.ErrorLineLogger(err)
	}
	commentService := service.NewCommentService(Db)
	movieService := service.NewMoviesService(redisClient, commentService)
	characterService := service.NewCharacterService(redisClient)

	moviesController = controllers.NewMovieController(movieService)
	characterController = controllers.NewCharactersHandler(characterService)
	commentsController = controllers.NewCommentHandler(commentService)

}

// @title        Movie API
// @version      1.0

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  odamilola36

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host		  52.207.233.252:8080
// @BasePath  	  /
func main() {

	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://52.207.233.252:8080/swagger/doc.json"), //The url pointing to API definition
	))

	subrouter := r.PathPrefix("/api/v1").Subrouter()
	installRoutes(subrouter)

	subrouter.Use(utils.RequestLogger)
	log.Printf("Starting API server on port %s and process id %d", Port, PID)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0" + Port,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		utils.ErrorLineLogger(err)
		log.Fatal(err)
	}
}

func installRoutes(r *mux.Router) {
	r.HandleFunc("/movies", moviesController.MovieListController).Methods("GET")
	r.HandleFunc("/characters/{id:[1-9]+}", characterController.GetCharacters).Methods("GET")
	r.HandleFunc("/comments/{id:[1-9]+}", commentsController.GetAllComments).Methods("GET")
	r.HandleFunc("/comment", commentsController.CreateComment).Methods("POST")
	r.HandleFunc("/comment/{commentId:[1-9]+}", commentsController.DeleteComment).Methods("DELETE")
	r.HandleFunc("/comment/{commentId:[1-9]+}", commentsController.UpdateComment).Methods("PATCH")
	r.HandleFunc("/comment/{commentId:[1-9]+}/{movieId:[1-9]+}", commentsController.GetComment).Methods("GET")
}
