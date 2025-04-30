package pokemonRepository

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/extizout/pokemon-api-server/modules/pokemon"
	"github.com/extizout/pokemon-api-server/pkg/cache"
	"github.com/extizout/pokemon-api-server/pkg/httpclient"
)

type PokemonRepositoryService interface {
	GetPokemonByName(pctx context.Context, name string) (*pokemon.Pokemon, error)
	GetPokemons(pctx context.Context) ([]pokemon.PokemonExternalResponseDto, error)
}

type pokemonRepository struct {
	apiURL string
	client *httpclient.Client
	cache  cache.CacheService
}

func NewPokemonRepository(apiURL string, client *httpclient.Client, cache cache.CacheService) PokemonRepositoryService {
	return &pokemonRepository{
		apiURL: apiURL,
		client: client,
		cache:  cache,
	}
}

func (r *pokemonRepository) GetPokemonByName(pctx context.Context, name string) (*pokemon.Pokemon, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	cacheKey := "pokemon:" + name

	if cached, found := r.cache.Get(cacheKey); found {
		if p, ok := cached.(*pokemon.Pokemon); ok {
			return p, nil
		}
	}

	resp, err := r.client.Get(ctx, r.apiURL+"/pokemon/"+name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, pokemon.ErrPokemonNotFound
	}

	var p pokemon.Pokemon
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}
	r.cache.SetWithDefaultExpiration(cacheKey, &p)

	return &p, nil
}

func (r *pokemonRepository) GetPokemons(pctx context.Context) ([]pokemon.PokemonExternalResponseDto, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	cacheKey := "pokemonList"
	if cached, found := r.cache.Get(cacheKey); found {
		if result, ok := cached.([]pokemon.PokemonExternalResponseDto); ok {
			return result, nil
		}
	}

	u, err := url.Parse(r.apiURL + "/pokemon")
	if err != nil {
		return nil, errors.New("failed to parse URL")
	}
	q := u.Query()

	//WARN: hardcoded external service pagination to get all pokemons from their API
	q.Set("offset", strconv.Itoa(0))
	q.Set("limit", strconv.Itoa(10000))
	u.RawQuery = q.Encode()

	resp, err := r.client.Get(ctx, u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	external := new(pokemon.PokemonListExternalResponseDto)
	if err := json.NewDecoder(resp.Body).Decode(external); err != nil {
		return nil, errors.New("failed to decode response")
	}

	results := MapPokemonListExternalToInternal(*external)

	r.cache.SetWithDefaultExpiration(cacheKey, results)

	return results, nil
}
