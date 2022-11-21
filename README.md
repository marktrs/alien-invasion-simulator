# Alien Invasion Simulator
![Coverage](https://img.shields.io/badge/Coverage-76.7%25-brightgreen)

- [Background](#background)
- [Usage](#usage)
- [Notes](#notes)

## Background

Mad aliens are about to invade the earth and you are tasked with simulating the invasion. You are given a map containing the names of cities in the non-existent world of X. The map is in a file, with one city per line. The city name is first,
followed by 1-4 directions (north, south, east, or west). Each one represents a road to another city that lies in that direction.

For example:

Foo north=Bar west=Baz south=Qu-ux

Bar south=Foo west=Bee

The city and each of the pairs are separated by a single space, and the
directions are separated from their respective cities with an equals (=) sign. 

You should create N aliens, where N is specified as a command-line argument.
These aliens start out at random places on the map, and wander around randomly, following links. Each iteration, the aliens can travel in any of the directions
leading out of a city. In our example above, an alien that starts at Foo can go
north to Bar, west to Baz, or south to Qu-ux.
When two aliens end up in the same place, they fight, and in the process kill each other and destroy the city. When a city is destroyed, it is removed from the map, and so are any roads that lead into or out of it. In our example above, if Bar were destroyed the map would now be something like:

Foo west=Baz south=Qu-ux

Once a city is destroyed, aliens can no longer travel to or through it. This
may lead to aliens getting "trapped".
You should create a program that reads in the world map, creates N aliens, and unleashes them. The program should run until all the aliens have been destroyed, or each alien has moved at least 10,000 times. When two aliens fight, print out a message like:

Bar has been destroyed by alien 10 and alien 34!

## Assumptions

*World / City*

* City names do not contain special characters
* Total cities in the world should be equal to or greater than 1
* Self-connection is not allowed
* Duplicate directions name is not allowed
* Duplicated connections to the same city are not allowed
* Only four cardinal directions north, south, east, and west are allowed
* The city can be isolate without any connections
* City connection is bidirectional where A north=B then B south=A

*Alien*

* Total aliens should be equal to or greater than 2
* Spawned aliens at the same location will not fight
* Aliens could not be present in 2 locations at the same time
* Alien fight and the destroyed city will proceed after all aliens are moved to the destination city

*Iteration*

* Total iterations should be equal to or greater than 1


## Usage

### Prerequisites

- installed [Golang 1.17](https://golang.org/)

### Start the application

Application will accept following flags :

```sh
$ go run . --help
-total-aliens uint
      Total number of aliens for simulation (default 10)
-total-iterations uint
      Total numbers of iteration for the simulation (default 10000)
-logging
      Logging each alien move on iteration
-world-file-path string
      World map data file path contains cities location data (default "data/world.txt")
```

Example 

```sh
$ go run . -total_aliens 4 -total_iterations 20 -logging true -output-file-path data/output.txt
```

Testing

```sh
$ go test ./...
```
