package entities

import "fmt"

type Category struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoxArtURL string `json:"box_art_url"`
}

func NewCategory() *Category {
	return &Category{}
}

func (c *Category) ToString() string {
	return fmt.Sprintf("ID: %s\tName: %s", c.ID, c.Name)
}
