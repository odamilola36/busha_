package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	MovieId     int    `json:"movie_id"`
	Body        string `valid:"MaxSize(500)" json:"body"`
	CommenterIP string `json:"commenter_ip"`
}
