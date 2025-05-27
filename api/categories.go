package api

import (
	"encoding/json"
	"log"

	"github.com/RootControl/twitch/entities"
)

const SEARCH_CATEGORIES = "search/categories"

type CategoriesResponse struct {
	Data       []entities.Category `json:"data"`
	Pagination entities.Pagination `json:"pagination"`
}

func NewCategoriesResponse() *CategoriesResponse {
	return &CategoriesResponse{
		Data:       make([]entities.Category, 0),
		Pagination: entities.Pagination{},
	}
}

func GetCategories(categoryName string) CategoriesResponse {
	request := NewApiRequest()

	responseBuffer := request.Get(SEARCH_CATEGORIES, "-q query="+categoryName, "-q first=50")

	var response CategoriesResponse
	err := json.Unmarshal(responseBuffer.Bytes(), &response)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return response
}
