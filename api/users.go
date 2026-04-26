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
	return getUser(login, nil)
}

func getUser(login string, e Executor) *UserResponse {
	request := newRequestWithExecutor(orShell(e))
	buf := request.Get("users", "-q login="+login)
	var response UserResponse
	if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return &response
}
