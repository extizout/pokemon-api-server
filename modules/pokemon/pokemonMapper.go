package pokemon

func MapPokemonToResponse(p *Pokemon) *PokemonResponseDto {
	return &PokemonResponseDto{
		Id:        p.Id,
		Name:      p.Name,
		Weight:    p.Weight,
		Types:     mapTypesToDto(p.Types),
		Abilities: mapAbilitiesToDto(p.Abilities),
	}
}

func MapPokemonListToResponseList(pokemons []Pokemon) []PokemonResponseDto {
	responses := make([]PokemonResponseDto, len(pokemons))
	for i, p := range pokemons {
		responses[i] = *MapPokemonToResponse(&p)
	}
	return responses
}

func MapPokemonAbilityToResponse(p *Pokemon) PokemonAbilitiesResponseDto {
	return mapAbilitiesToDto(p.Abilities)
}

func mapTypesToDto(types []Types) []Type {
	result := make([]Type, len(types))
	for i, t := range types {
		result[i] = Type{Name: t.Type.Name}
	}
	return result
}

func mapAbilitiesToDto(abilities []Abilities) []Ability {
	result := make([]Ability, len(abilities))
	for i, a := range abilities {
		result[i] = Ability{Name: a.Ability.Name}
	}
	return result
}
