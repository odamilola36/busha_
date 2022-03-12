package dto

import "busha_/models"

type SwappiMoviesResponse struct {
	Count    int                  `json:"count"`
	Next     interface{}          `json:"next"`
	Previous interface{}          `json:"previous"`
	Results  []models.SwappiMovie `json:"results"`
}
