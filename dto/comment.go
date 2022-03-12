package dto

type CreateComment struct {
	MovieId int    `json:"movie_id"`
	Body    string `valid:"MaxSize(500)" json:"body"`
}
