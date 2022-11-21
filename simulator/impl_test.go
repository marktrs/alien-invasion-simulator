package simulator

import (
	"testing"

	"github.com/marktrs/alien-invasion-assignment/model"
	"github.com/stretchr/testify/suite"
)

type SimulatorTestSuite struct {
	suite.Suite
	simulator *Simulator
}

func TestSimulatorTestSuite(t *testing.T) {
	suite.Run(t, new(SimulatorTestSuite))
}

func (suite *SimulatorTestSuite) SetupTest() {
	// Set default simulation config with mocked world data
	suite.simulator = NewSimulator(&model.World{
		Cities: map[string]*model.City{
			"Bar": {
				Name: "Bar",
				Connections: map[string]string{
					"south": "Foo",
					"west":  "Bee",
				},
			},
			"Baz": {
				Name: "Baz",
				Connections: map[string]string{
					"east": "Foo",
				},
			},
			"Bee": {
				Name: "Bee",
				Connections: map[string]string{
					"east": "Bar",
				},
			},
			"Foo": {
				Name: "Foo",
				Connections: map[string]string{
					"north": "Bar",
					"south": "Qu-ux",
					"west":  "Baz",
				},
			},
			"Qu-ux": {
				Name: "Qu-ux",
				Connections: map[string]string{
					"north": "Foo",
				},
			},
		},
	}, &Config{
		TotalAliveAliens: 4,
		TotalIterations:  10,
	})
}

func (suite *SimulatorTestSuite) TestStartSimulator() {
	suite.NoError(suite.simulator.Start())
}

func (suite *SimulatorTestSuite) TestStartSimulatorWithoutSurvivedCities() {
	suite.simulator.config = &Config{
		TotalAliveAliens: 1,
		TotalIterations:  2,
	}
	suite.simulator.world = &model.World{
		Cities: map[string]*model.City{
			"Bar": {
				Name:        "Bar",
				IsDestroyed: true,
			},
		},
	}

	suite.NoError(suite.simulator.Start())

	suite.Equal(suite.simulator.config.TotalIterations, uint64(1))
	suite.Equal(suite.simulator.config.TotalAliveAliens, uint64(1))
}

func (suite *SimulatorTestSuite) Test_spawnAlien() {
	suite.simulator.config = &Config{
		TotalAliveAliens: 1,
		TotalIterations:  1,
	}
	suite.simulator.world = &model.World{
		Cities: map[string]*model.City{
			"Bar": {
				Name: "Bar",
			},
		},
	}
	suite.simulator.spawnAliens()
	suite.Equal(len(suite.simulator.world.Cities["Bar"].Aliens), 1, "Bar city should have Alien1")
	suite.Equal(suite.simulator.aliens[AlienNamePrefix+"1"].Location.Name, "Bar", "Alien1 should spawned at Bar city")
}

func (suite *SimulatorTestSuite) Test_proceedNextMove() {
	suite.simulator.config = &Config{
		TotalAliveAliens: 1,
		TotalIterations:  1,
	}
	suite.simulator.world = &model.World{
		Cities: map[string]*model.City{
			"Bar": {
				Name: "Bar",
				Connections: map[string]string{
					"north": "Foo",
				},
			},
			"Foo": {
				Name: "Foo",
				Connections: map[string]string{
					"south": "Bar",
				},
			},
		},
	}

	suite.simulator.spawnAliens()
	alienCountInBar := len(suite.simulator.world.Cities["Bar"].Aliens)

	suite.simulator.proceedNextMove()

	if alienCountInBar > 0 {
		// if alien spawned at Bar city it should move to Foo city
		suite.Equal(len(suite.simulator.world.Cities["Foo"].Aliens), 1, "Foo city should have Alien1")
		suite.Equal(len(suite.simulator.world.Cities["Bar"].Aliens), 0, "Bar city should not have Alien1")
	} else {
		// else alien spawned at Foo city it should move to Bar city
		suite.Equal(len(suite.simulator.world.Cities["Bar"].Aliens), 1, "Bar city should have Alien1")
		suite.Equal(len(suite.simulator.world.Cities["Foo"].Aliens), 0, "Foo city should not have Alien1")
	}
}
