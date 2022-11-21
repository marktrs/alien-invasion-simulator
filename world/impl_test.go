package world

import (
	"testing"

	"github.com/marktrs/alien-invasion-assignment/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WorldTestSuite struct {
	suite.Suite
}

func TestWorldTestSuite(t *testing.T) {
	suite.Run(t, new(WorldTestSuite))
}

func (suite *WorldTestSuite) TestNewFromFile() {
	path := "./fixtures/map1.txt"

	world, err := NewFromFile(path)
	if err != nil {
		suite.Error(err)
	}

	assert.Equal(suite.T(), &model.World{
		Cities: map[string]*model.City{
			"Bar": {
				Name: "Bar",
				Connections: map[string]string{
					"south": "Foo",
					"west":  "Bee",
				},
				IsDestroyed: false,
			},
			"Baz": {
				Name: "Baz",
				Connections: map[string]string{
					"east": "Foo",
				},
				IsDestroyed: false,
			},
			"Bee": {
				Name: "Bee",
				Connections: map[string]string{
					"east": "Bar",
				},
				IsDestroyed: false,
			},
			"Foo": {
				Name: "Foo",
				Connections: map[string]string{
					"north": "Bar",
					"south": "Qu-ux",
					"west":  "Baz",
				},
				IsDestroyed: false,
			},
			"Qu-ux": {
				Name: "Qu-ux",
				Connections: map[string]string{
					"north": "Foo",
				},
				IsDestroyed: false,
			},
		},
	}, world)
}

func (suite *WorldTestSuite) TestNewFromFileWithInvalidDirectionName() {
	path := "./fixtures/map2.txt"
	_, err := NewFromFile(path)
	suite.EqualError(err, ErrInvalidDirectionName.Error())
}

func (suite *WorldTestSuite) TestNewFromFileWithDuplicateDirection() {
	path := "./fixtures/map3.txt"
	_, err := NewFromFile(path)
	suite.EqualError(err, ErrInconsistentConnection.Error())
}

func (suite *WorldTestSuite) TestNewFromFileWithMaxConnectionExceed() {
	path := "./fixtures/map4.txt"
	_, err := NewFromFile(path)
	suite.EqualError(err, ErrMaxConnectionExceed.Error())
}

func (suite *WorldTestSuite) TestNewFromFileWithDuplicatedCityNameInput() {
	path := "./fixtures/map5.txt"
	_, err := NewFromFile(path)
	suite.EqualError(err, ErrInconsistentConnection.Error())
}

func (suite *WorldTestSuite) TestNewFromFileWithInconsistentConnection() {
	path := "./fixtures/map6.txt"
	_, err := NewFromFile(path)
	suite.EqualError(err, ErrInconsistentConnection.Error())
}

func (suite *WorldTestSuite) TestNewFromFileWithSelfConnection() {
	path := "./fixtures/map7.txt"
	_, err := NewFromFile(path)
	suite.EqualError(err, ErrSelfConnection.Error())
}
