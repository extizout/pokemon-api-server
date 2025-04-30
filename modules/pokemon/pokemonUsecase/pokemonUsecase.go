package pokemonUsecase

import (
	"context"
	"math/rand"
	"time"

	"github.com/extizout/pokemon-api-server/modules/pokemon"
	"github.com/extizout/pokemon-api-server/modules/pokemon/pokemonRepository"
)

type (
	PokemonUsecaseService interface {
		GetPokemonByName(pctx context.Context, name string) (*pokemon.Pokemon, error)
		GetRandomPokemon(pctx context.Context) (*pokemon.Pokemon, error)
	}
	pokemonUsecase struct {
		pokemonRepository pokemonRepository.PokemonRepositoryService
	}
)

func NewPokemonUsecase(pokemonRepository pokemonRepository.PokemonRepositoryService) PokemonUsecaseService {
	return &pokemonUsecase{pokemonRepository: pokemonRepository}
}

func (u *pokemonUsecase) GetPokemonByName(pctx context.Context, name string) (*pokemon.Pokemon, error) {
	data, err := u.pokemonRepository.GetPokemonByName(pctx, name)
	if err != nil {
		return nil, err
	}

	return &pokemon.Pokemon{
		Id:        data.Id,
		Name:      data.Name,
		Weight:    data.Weight,
		Types:     data.Types,
		Abilities: data.Abilities,
	}, nil
}

func (u *pokemonUsecase) GetRandomPokemon(pctx context.Context) (*pokemon.Pokemon, error) {
	pokemons, err := u.pokemonRepository.GetPokemons(pctx)
	if err != nil {
		return nil, err
	}
	if len(pokemons) == 0 {
		return nil, pokemon.ErrPokemonNotFound
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(pokemons))

	randomPokemonName := pokemons[randomIndex].Name

	randomizedPokemon, err := u.pokemonRepository.GetPokemonByName(pctx, randomPokemonName)
	if err != nil {
		return nil, err
	}

	return randomizedPokemon, nil
}
