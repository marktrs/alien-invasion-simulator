package simulator

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/marktrs/alien-invasion-assignment/alien"
	"github.com/marktrs/alien-invasion-assignment/model"
)

var AlienNamePrefix = "Alien"

type Config struct {
	TotalAliveAliens uint64
	TotalIterations  uint64
	IsLoggerEnabled  bool
}

type Simulator struct {
	world  *model.World
	aliens map[string]*model.Alien
	config *Config
}

func NewSimulator(
	world *model.World,
	config *Config,
) *Simulator {
	return &Simulator{
		world:  world,
		aliens: map[string]*model.Alien{},
		config: config,
	}
}

func (s *Simulator) Start() error {
	s.spawnAliens()

	if s.config.IsLoggerEnabled {
		s.initialReport()
	}

	// iterate through the loop to move alien and count the destroyed cities
	for {
		if s.config.IsLoggerEnabled {
			fmt.Println("\nIteration left", s.config.TotalIterations)
			fmt.Println("---------------------------")
			s.status()
		}

		// end iteration loop if reached max iteration count
		if s.config.TotalIterations == 0 {
			s.summaryReport("\nSimulation end : Max iteration is reached")
			return nil
		}
		s.config.TotalIterations--

		// If there are survived alien then they will make a new move
		s.proceedNextMove()

		if s.config.IsLoggerEnabled {
			fmt.Println("-----------Moved-----------")
			s.status()
		}

		// get survived and destroyed cities and calculate number of survived alien
		// in separated loop to avoid high number of big(O)
		survivedCities, destroyedCities := s.getSurvivedAndDestroyedCities()
		if len(survivedCities) == 0 {
			s.summaryReport("\nSimulation end : All cities was destroyed or isolated")
			return nil
		}

		diedAliens := s.getDiedAliens(destroyedCities)

		// decrease total alien amount in this simulator with died number
		s.config.TotalAliveAliens -= uint64(len(diedAliens))

		trappedCities := s.getTrappedCities(survivedCities)
		trappedAliens := s.getTrappedAliens(trappedCities)

		// decrease total alien amount in this simulator with died number
		s.config.TotalAliveAliens -= uint64(len(trappedAliens))

		if s.config.IsLoggerEnabled {
			fmt.Println("-----------Result-----------")
			s.status()
		}

		if s.config.TotalAliveAliens == 0 {
			s.summaryReport("\nSimulation end : All aliens is dead or trapped")
			return nil
		}

	}
}

func (s *Simulator) getSurvivedAndDestroyedCities() (survivedCities, destroyedCities map[string]*model.City) {
	survivedCities = make(map[string]*model.City)
	destroyedCities = make(map[string]*model.City)

	for _, city := range s.world.Cities {
		if city.IsDestroyed {
			continue
		}

		// if there are cities with more than 2 aliens they should fight and marked as died
		if len(city.Aliens) > 1 {
			city.IsDestroyed = true
			destroyedCities[city.Name] = city
		} else {
			survivedCities[city.Name] = city
		}
	}

	return survivedCities, destroyedCities

}

func (s *Simulator) getDiedAliens(destroyedCities map[string]*model.City) []string {
	diedAlien := make([]string, 0)

	// find died alien from destroyed cities
	for _, city := range destroyedCities {
		fmt.Printf("%s city has been destroyed by", city.Name)
		for _, alien := range city.Aliens {
			fmt.Printf(" %s", alien.Name)
			alien.IsDead = true
			diedAlien = append(diedAlien, alien.Name)
		}
		fmt.Print("\n")
	}

	return diedAlien
}

func (s *Simulator) getTrappedCities(survivedCities map[string]*model.City) map[string]*model.City {
	trappedCities := make(map[string]*model.City)
	// after battle zone city was destroyed count the trapped neighbor city
	for _, city := range survivedCities {
		for _, neighbor := range city.Connections {
			// if neighbor city is not destroyed, those alien not trapped
			if !s.world.Cities[neighbor].IsDestroyed {
				break
			}

			// if all connection was destroyed it's a trapped city
			city.IsIsolated = true
			trappedCities[city.Name] = city
		}
	}

	return trappedCities
}

func (s *Simulator) getTrappedAliens(trappedCities map[string]*model.City) []string {
	// trappedAlien collect and count total trapped alien from trapped city
	trappedAlien := make([]string, 0)
	for _, city := range trappedCities {
		for k, alien := range city.Aliens {
			alien.IsTrapped = true
			trappedAlien = append(trappedAlien, k)
		}
	}

	return trappedAlien
}

func (s *Simulator) proceedNextMove() {
	for _, alien := range s.aliens {
		// ignore died, trapped alien since they can't move
		if alien.IsDead || alien.IsTrapped {
			continue
		}

		availableCities := make([]string, 0, len(alien.Location.Connections))
		for _, city := range alien.Location.Connections {
			if s.world.Cities[city].IsDestroyed {
				continue
			}
			availableCities = append(availableCities, city)
		}

		// random next move from available city
		if len(availableCities) == 0 {
			break
		}

		// random next move from available city
		nextCity := availableCities[rand.Intn(len(availableCities))]

		// remove them from current city to prevent existing at two location
		delete(s.world.Cities[alien.Location.Name].Aliens, alien.Name)

		// moving to the next city
		alien.Location = s.world.Cities[nextCity]
		s.world.Cities[nextCity].Aliens[alien.Name] = alien

	}
}

// spawnAliens pick a random city and assigned to alien
func (s *Simulator) spawnAliens() {
	var spawnCityName string
	cityNames := make([]string, 0)

	// fetch city to invade
	for name, city := range s.world.Cities {
		city.Aliens = make(map[string]*model.Alien)
		cityNames = append(cityNames, name)
	}

	// random city name and assigned to alien
	for i := uint64(0); i < s.config.TotalAliveAliens; i++ {
		spawnCityName = getRandomCityName(cityNames)
		alienName := AlienNamePrefix + strconv.Itoa(int(i+1))
		newAlien := alien.NewAlien(alienName, s.world.Cities[spawnCityName])
		s.aliens[alienName] = newAlien
		s.world.Cities[spawnCityName].Aliens[alienName] = newAlien
	}
}

func (s *Simulator) status() {
	fmt.Println("Aliens status:")
	for _, alien := range s.aliens {
		fmt.Println(alien.Name,
			"location:", alien.Location.Name,
			"isDead:", alien.IsDead,
			"isTrapped:", alien.IsTrapped)
	}

	fmt.Println("\nCity status:")
	for _, city := range s.world.Cities {
		fmt.Println(city.Name,
			"isDestroyed:", city.IsDestroyed,
			"IsIsolated:", city.IsIsolated,
			"aliensCount:", len(city.Aliens),
		)
	}
}

func (s *Simulator) initialReport() {
	fmt.Println("Initial report")
	fmt.Println("----- City Report -----")
	for _, city := range s.world.Cities {
		fmt.Printf("%s ", city.Name)
		for direction, name := range city.Connections {
			neighbor := s.world.Cities[name]
			if !neighbor.IsDestroyed {
				fmt.Printf("%s=%s ", direction, neighbor.Name)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Println("---- Aliens Report ----")
	for _, alien := range s.aliens {
		fmt.Printf("%s is spawned at %s", alien.Name, alien.Location.Name)
	}
}

func (s *Simulator) summaryReport(message string) {
	fmt.Printf("\n %s \n", message)
	fmt.Println("----- City Report -----")
	for _, city := range s.world.Cities {
		fmt.Printf("%s ", city.Name)
		if city.IsDestroyed {
			fmt.Println("has been destroyed")
			continue
		}

		if city.IsIsolated {
			fmt.Println("is isolated")
			continue
		}

		for direction, name := range city.Connections {
			neighbor := s.world.Cities[name]
			if !neighbor.IsDestroyed {
				fmt.Printf("%s=%s ", direction, neighbor.Name)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Println("---- Aliens Report ----")
	for _, alien := range s.aliens {
		fmt.Printf("%s ", alien.Name)
		if alien.IsDead {
			fmt.Println("is dead")
		} else if alien.IsTrapped {
			fmt.Println("is trapped")
		} else {
			fmt.Println("is alive")
		}
	}
}

func getRandomCityName(cities []string) string {
	index := rand.Intn(len(cities))
	return cities[index]
}
