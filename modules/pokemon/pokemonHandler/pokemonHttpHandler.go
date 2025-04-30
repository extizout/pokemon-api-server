package pokemonHandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/extizout/pokemon-api-server/config"
	"github.com/extizout/pokemon-api-server/modules/pokemon"
	"github.com/extizout/pokemon-api-server/modules/pokemon/pokemonUsecase"
	"github.com/extizout/pokemon-api-server/pkg/httpclient"
	"github.com/extizout/pokemon-api-server/pkg/utils/request"
	"github.com/extizout/pokemon-api-server/pkg/utils/response"
	"github.com/labstack/echo/v4"
)

type (
	PokemonHttpHandlerService interface {
		GetPokemonByName(c echo.Context) error
		GetPokemonAbilityByName(c echo.Context) error
		GetRandomPokemon(c echo.Context) error
	}

	pokemonHttpHandler struct {
		cfg            *config.Config
		pokemonUsecase pokemonUsecase.PokemonUsecaseService
	}
)

func NewPokemonHttpHandler(cfg *config.Config, pokemonUsecase pokemonUsecase.PokemonUsecaseService) PokemonHttpHandlerService {
	return &pokemonHttpHandler{cfg: cfg, pokemonUsecase: pokemonUsecase}
}

func (h *pokemonHttpHandler) GetPokemonByName(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)
	requestDto := new(pokemon.PokemonNameRequestPathDto)
	if err := wrapper.Bind(requestDto); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	pokemonData, err := h.pokemonUsecase.GetPokemonByName(ctx, requestDto.Name)
	if err != nil {
		switch {
		case errors.Is(err, pokemon.ErrPokemonNotFound):
			return response.ErrorResponse(c, http.StatusNotFound, err.Error())
		case errors.Is(err, httpclient.ErrRequestTimeout) || errors.Is(err, httpclient.ErrConnectionError):
			return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		default:
			return response.ErrorResponse(c, http.StatusNotFound, err.Error())
		}
	}
	responseDto := pokemon.MapPokemonToResponse(pokemonData)
	return response.SuccessResponse(c, http.StatusOK, responseDto)
}

func (h *pokemonHttpHandler) GetPokemonAbilityByName(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)
	requestDto := new(pokemon.PokemonNameRequestPathDto)
	if err := wrapper.Bind(requestDto); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	pokemonData, err := h.pokemonUsecase.GetPokemonByName(ctx, requestDto.Name)
	if err != nil {
		switch {
		case errors.Is(err, pokemon.ErrPokemonNotFound):
			return response.ErrorResponse(c, http.StatusNotFound, err.Error())
		case errors.Is(err, httpclient.ErrRequestTimeout) || errors.Is(err, httpclient.ErrConnectionError):
			return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		default:
			return response.ErrorResponse(c, http.StatusNotFound, err.Error())
		}
	}

	responseDto := pokemon.MapPokemonToResponse(pokemonData)
	return response.SuccessResponse(c, http.StatusOK, responseDto.Abilities)
}

func (h *pokemonHttpHandler) GetRandomPokemon(c echo.Context) error {
	ctx := context.Background()

	pokemonData, err := h.pokemonUsecase.GetRandomPokemon(ctx)
	if err != nil {
		switch {
		case errors.Is(err, pokemon.ErrPokemonNotFound):
			return response.ErrorResponse(c, http.StatusNotFound, err.Error())
		case errors.Is(err, httpclient.ErrRequestTimeout) || errors.Is(err, httpclient.ErrConnectionError):
			return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		default:
			return response.ErrorResponse(c, http.StatusNotFound, err.Error())
		}
	}
	responseDto := pokemon.MapPokemonToResponse(pokemonData)
	return response.SuccessResponse(c, http.StatusOK, responseDto)
}
