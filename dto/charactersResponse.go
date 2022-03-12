package dto

import "busha_/models"

type CharactersResponse struct {
	Count           int                `json:"total"`
	TotalHeightCm   string             `json:"total_height_cm"`
	TotalHeightFeet string             `json:"total_height_feet"`
	Characters      []models.Character `json:"characters"`
}
