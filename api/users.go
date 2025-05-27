package api

import (
	"encoding/json"
	"log"

	"github.com/RootControl/twitch/entities"
)

type UserResponse struct {
	Data []entities.User `json:"data"`
}

func GetUser(login string) *UserResponse {
	request := NewApiRequest()

	responseBuffer := request.Get("users", "-q login="+login)

	var response UserResponse
	err := json.Unmarshal(responseBuffer.Bytes(), &response)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return &response
}
