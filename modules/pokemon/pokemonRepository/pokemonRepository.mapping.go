package pokemonRepository

import "github.com/extizout/pokemon-api-server/modules/pokemon"

func MapPokemonListExternalToInternal(external pokemon.PokemonListExternalResponseDto) []pokemon.PokemonExternalResponseDto {
	internalList := make([]pokemon.PokemonExternalResponseDto, len(external.Results))

	for i, ext := range external.Results {
		internalList[i] = pokemon.PokemonExternalResponseDto{
			Name: ext.Name,
			URL:  ext.URL,
		}
	}

	return internalList
}
