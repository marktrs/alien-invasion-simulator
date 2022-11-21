package alien

import (
	"github.com/marktrs/alien-invasion-assignment/model"
)

func NewAlien(name string, location *model.City) *model.Alien {
	return &model.Alien{
		Name:     name,
		Location: location,
	}
}
