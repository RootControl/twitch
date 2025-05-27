package entities

import "fmt"

type Category struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoxArtURL string `json:"box_art_url"`
}

func NewCategory(id, name, boxArtURL string) *Category {
	return &Category{
		ID:        id,
		Name:      name,
		BoxArtURL: boxArtURL,
	}
}

func (c *Category) ToString() string {
	return fmt.Sprintf("ID: %s\tName: %s", c.ID, c.Name)
}
