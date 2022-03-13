package service

import (
	"busha_/dto"
	"busha_/models"
	"busha_/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type CharacterService interface {
	GetCharacters(req map[string]string) (*dto.CharactersResponse, error)
	genderFilter(gender string, ch []models.Character) ([]models.Character, error)
	getCharactersFromService(movieId string) ([]models.Character, error)
	sortCharacters(ch []models.Character, sort string, order string) ([]models.Character, error)
	getHeight(ch []models.Character) (float64, float64, float64)
}

type characterService struct {
	redisClient *redis.Client
}

func NewCharacterService(client *redis.Client) *characterService {
	return &characterService{
		redisClient: client,
	}
}

func (c *characterService) GetCharacters(req map[string]string) (*dto.CharactersResponse, error) {

	var charactersResponse *dto.CharactersResponse
	var characters, filteredCharacters []models.Character

	var err error
	movieId := req["movieId"]

	result, err := c.redisClient.Get(context.TODO(), movieId).Result()
	if err == redis.Nil {
		// If the key does not exist, call the api to get the
		characters, err = c.getCharactersFromService(movieId)

		if err != nil {
			utils.ErrorLineLogger(err)
			return charactersResponse, err
		}
		// Save the characters in redis
		v, err := json.Marshal(characters)
		if err != nil {
			utils.ErrorLineLogger(err)
			return charactersResponse, err
		}
		cmd := c.redisClient.Set(context.TODO(), movieId, v, 30*time.Second)
		if err := cmd.Err(); err != nil {
			utils.ErrorLineLogger(err)
			return charactersResponse, err
		}
	} else {
		// If the key exists, get the characters from redis
		err = json.Unmarshal([]byte(result), &characters)
		if err != nil {
			utils.ErrorLineLogger(err)
			return charactersResponse, err
		}
	}

	sortBy := strings.ToLower(strings.TrimSpace(req["sort_by"]))
	order := strings.ToLower(strings.TrimSpace(req["order"]))
	filterBy := strings.ToLower(strings.TrimSpace(req["filter"]))

	filteredCharacters, err = c.sortCharacters(characters, sortBy, order)
	if err != nil {
		utils.ErrorLineLogger(err)
		return charactersResponse, err
	}
	filteredCharacters, err = c.genderFilter(filterBy, filteredCharacters)
	if err != nil {
		utils.ErrorLineLogger(err)
		return charactersResponse, err
	}
	cmHeight, ftHeight, inHeight := c.getHeight(filteredCharacters)
	charactersResponse = &dto.CharactersResponse{
		Characters:      filteredCharacters,
		Count:           len(filteredCharacters),
		TotalHeightCm:   fmt.Sprintf("%.2f", cmHeight),
		TotalHeightFeet: fmt.Sprintf("%.0f%s and %.2f%s", ftHeight, "ft", inHeight, "inches"),
	}

	return charactersResponse, nil
}

func (c *characterService) sortCharacters(ch []models.Character, sortBy string, order string) ([]models.Character, error) {

	if "name" == strings.ToLower(strings.TrimSpace(sortBy)) {
		sort.Slice(ch, func(i, j int) bool {
			if order == "asc" {
				return ch[i].Name < ch[j].Name
			}
			return ch[i].Name > ch[j].Name
		})
	}

	if "gender" == strings.ToLower(strings.TrimSpace(sortBy)) {
		sort.Slice(ch, func(i, j int) bool {
			if order == "asc" {
				return ch[i].Gender < ch[j].Gender
			}
			return ch[i].Gender > ch[j].Gender
		})
	}

	if "height" == strings.ToLower(strings.TrimSpace(sortBy)) {
		sort.Slice(ch, func(i, j int) bool {
			atoi, _ := strconv.Atoi(ch[i].Height)
			i2, _ := strconv.Atoi(ch[j].Height)
			if order == "asc" {
				return atoi < i2
			}
			return atoi > i2
		})
	}

	return ch, nil
}

func (c *characterService) getHeight(ch []models.Character) (float64, float64, float64) {
	var sum = 0.0
	for _, character := range ch {
		if v, err := strconv.ParseFloat(character.Height, 64); err == nil {
			sum += v
		}
	}

	feet, dec := math.Modf(sum / 30.48)

	return sum, feet, math.Round(dec*12*100) / 100
}

func (c *characterService) genderFilter(gender string, ch []models.Character) ([]models.Character, error) {
	genderStr := strings.TrimSpace(strings.ToLower(gender))
	if genderStr == "male" || genderStr == "female" || genderStr == "n\\a" {
		var characters []models.Character
		for _, character := range ch {
			if strings.ToLower(character.Gender) == genderStr {
				characters = append(characters, character)
			}
		}
		return characters, nil
	}
	return ch, nil
}

func (c *characterService) getCharactersFromService(movieId string) ([]models.Character, error) {

	var movies models.SwappiMovie

	res, err := http.Get("https://swapi.dev/api/films/" + movieId)

	if err != nil {
		utils.ErrorLineLogger(err)
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("error in getting characters from service, status code: " + res.Status)
	}
	err = json.NewDecoder(res.Body).Decode(&movies)

	if err != nil {
		utils.ErrorLineLogger(err)
		return nil, err
	}

	var characters []models.Character

	for _, characterUrl := range movies.Characters {
		var character models.Character
		res, err := http.Get(characterUrl)
		if err != nil {
			utils.ErrorLineLogger(err)
			return nil, err
		}
		err = json.NewDecoder(res.Body).Decode(&character)
		if err != nil {
			utils.ErrorLineLogger(err)
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, nil
}
