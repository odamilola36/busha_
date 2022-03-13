package main

import (
	"busha_/config"
	"busha_/controllers"
	"busha_/service"
	"busha_/utils"
	"context"
	"github.com/gorilla/mux"
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

func main() {

	r := mux.NewRouter()

	installRoutes(r)

	r.Use(utils.RequestLogger)
	log.Printf("Starting API server on port %s and process id %d", Port, PID)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0" + Port,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}
	err := srv.ListenAndServe()
	log.Println("Server stopped")
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
