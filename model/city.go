package model

type City struct {
	Name        string
	Connections map[string]string
	Aliens      map[string]*Alien
	IsDestroyed bool
	IsIsolated   bool
}
