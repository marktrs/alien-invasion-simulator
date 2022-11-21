package model

type Alien struct {
	Name      string
	Location  *City
	NextCity  *City
	IsDead    bool
	IsTrapped bool
}
