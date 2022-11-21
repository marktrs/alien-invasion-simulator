package world

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/marktrs/alien-invasion-assignment/model"
)

var (
	ErrInvalidDirectionName   = errors.New("invalid direction name")
	ErrDuplicateDirection     = errors.New("duplicated direction name is not allowed")
	ErrMaxConnectionExceed    = errors.New("max connection exceed")
	ErrSelfConnection         = errors.New("self connection is not allowed")
	ErrDuplicatedConnection   = errors.New("duplicated connection is not allowed")
	ErrInconsistentConnection = errors.New("connection is not consistent")
)

var AllowedDirectionNames = []string{"north", "south", "east", "west"}
var OppositeDirections = map[string]string{
	"north": "south",
	"south": "north",
	"east":  "west",
	"west":  "east",
}

func NewFromFile(path string) (*model.World, error) {
	world := &model.World{
		Cities: make(map[string]*model.City),
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		city := &model.City{
			Connections: make(map[string]string),
		}
		// trim zero-width whitespace
		trimmedInput := strings.ReplaceAll(scanner.Text(), "\u200b", "")

		// extract city name and connections from current line
		inputs := strings.Split(trimmedInput, " ")
		city.Name = inputs[0]

		// check if city is registered in the world and load existing city
		// if not registered the new city to the world
		if world.Cities[city.Name] != nil {
			city = world.Cities[city.Name]
		} else {
			world.Cities[city.Name] = city
		}

		for i, chunk := range inputs[1:] {
			// Only max 4 connections (0-3)
			if i > 3 {
				return nil, ErrMaxConnectionExceed
			}

			// c[0] is the direction and c[1] is the name of neighbor city on the same direction
			c := strings.Split(chunk, "=")
			direction, err := validateAllowedDirectionNames(strings.Trim(c[0], " "))
			if err != nil {
				return nil, err
			}

			neighborCityName := c[1]

			// If neighbor city is already registered
			if world.Cities[neighborCityName] == nil {
				// register their neighbor city
				world.Cities[neighborCityName] = &model.City{
					Name: neighborCityName,
					Connections: map[string]string{
						OppositeDirections[direction]: city.Name,
					},
				}
			}

			// Current city should connect to neighbor city
			if err = connect(world, city, direction, neighborCityName); err != nil {
				return nil, err
			}

			// then opposite direction of neighbor should be current city name
			if err = connect(world, world.Cities[neighborCityName], OppositeDirections[direction], city.Name); err != nil {
				return nil, err
			}

			// update current city in city mapper
			world.Cities[city.Name] = city
		}
	}

	return world, nil
}

func validateAllowedDirectionNames(name string) (string, error) {
	name = strings.ToLower(name)
	for _, v := range AllowedDirectionNames {
		if v == name {
			return name, nil
		}
	}

	return "", ErrInvalidDirectionName
}

func connect(world *model.World, city *model.City, direction string, neighbor string) error {
	// self connection is not allowed
	if city.Name == neighbor {
		return ErrSelfConnection
	}

	if city.Connections[direction] != "" {
		if city.Connections[direction] != neighbor {
			return ErrInconsistentConnection
		}
	}

	city.Connections[direction] = neighbor

	return nil
}
