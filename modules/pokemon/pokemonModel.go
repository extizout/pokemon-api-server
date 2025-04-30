package pokemon

type (
	PokemonResponseDto struct {
		Id        int       `json:"id"`
		Name      string    `json:"name"`
		Weight    int       `json:"weight"`
		Types     []Type    `json:"types"`
		Abilities []Ability `json:"abilities"`
	}

	PokemonAbilitiesResponseDto []Ability

	PokemonListExternalResponseDto struct {
		Count    int                          `json:"count"`
		Next     *string                      `json:"next"`
		Previous *string                      `json:"previous"`
		Results  []PokemonExternalResponseDto `json:"results"`
	}

	PokemonExternalResponseDto struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	PokemonNameRequestPathDto struct {
		Name string `param:"name" validate:"required"`
	}
)
