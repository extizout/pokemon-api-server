package pokemon

import "errors"

type (
	Pokemon struct {
		Id        int         `json:"id"`
		Name      string      `json:"name"`
		Weight    int         `json:"weight"`
		Types     []Types     `json:"types"`
		Abilities []Abilities `json:"abilities"`
	}

	Abilities struct {
		Ability Ability
	}

	Ability struct {
		Name string `json:"name"`
	}

	Types struct {
		Slot int  `json:"slot"`
		Type Type `json:"type"`
	}

	Type struct {
		Name string `json:"name"`
	}
)

var (
	ErrPokemonNotFound = errors.New("pokemon not found")
	ErrExternalService = errors.New("external service error")
)
