package models

import "time"

type SwappiMovie struct {
	Title        string    `json:"title"`
	EpisodeId    int       `json:"episode_id"`
	OpeningCrawl string    `json:"opening_crawl"`
	Director     string    `json:"director"`
	Producer     string    `json:"producer"`
	ReleaseDate  string    `json:"release_date"`
	Characters   []string  `json:"characters"`
	Planets      []string  `json:"planets"`
	Starships    []string  `json:"starships"`
	Vehicles     []string  `json:"vehicles"`
	Species      []string  `json:"species"`
	Created      time.Time `json:"created"`
	Edited       time.Time `json:"edited"`
	Url          string    `json:"url"`
}

type Movie struct {
	CommentCount int       `json:"comment_count"`
	Title        string    `json:"title"`
	OpeningCrawl string    `json:"opening_crawl"`
	ReleaseDate  time.Time `json:"-"`
}

func (s SwappiMovie) ToMovie() Movie {
	v, _ := time.Parse("2006-01-02", s.ReleaseDate)
	return Movie{
		Title:        s.Title,
		OpeningCrawl: s.OpeningCrawl,
		ReleaseDate:  v,
	}
}
