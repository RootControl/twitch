package api

import (
	"github.com/RootControl/twitch/entities"
)

const (
	SEARCH_CATEGORIES = "search/categories"
	TOP_GAMES         = "games/top"
)

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

func GetCategories(categoryName string, limit int) (CategoriesResponse, error) {
	return getCategories(categoryName, limit, nil)
}

func getCategories(categoryName string, limit int, e Executor) (CategoriesResponse, error) {
	return fetch[CategoriesResponse](e, SEARCH_CATEGORIES,
		Q("query", categoryName),
		Q("first", clampLimit(limit)),
	)
}

// GetTopGames returns the categories with the most current viewers.
func GetTopGames(limit int) (CategoriesResponse, error) {
	return getTopGames(limit, nil)
}

func getTopGames(limit int, e Executor) (CategoriesResponse, error) {
	return fetch[CategoriesResponse](e, TOP_GAMES, Q("first", clampLimit(limit)))
}
