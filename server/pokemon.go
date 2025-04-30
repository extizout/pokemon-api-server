package server

import (
	"math"
	"time"

	"github.com/extizout/pokemon-api-server/modules/pokemon/pokemonHandler"
	"github.com/extizout/pokemon-api-server/modules/pokemon/pokemonRepository"
	"github.com/extizout/pokemon-api-server/modules/pokemon/pokemonUsecase"
	"github.com/extizout/pokemon-api-server/modules/user"
	"github.com/extizout/pokemon-api-server/pkg/authentication/authenticationRepository"
	"github.com/extizout/pokemon-api-server/pkg/authentication/authenticationUsecase"
	"github.com/extizout/pokemon-api-server/pkg/cache"
	"github.com/extizout/pokemon-api-server/pkg/httpclient"
	"github.com/extizout/pokemon-api-server/pkg/middleware"
)

func (s *server) pokemonService(usersInMemory *[]user.User) {
	client := httpclient.NewClient(0*time.Second, httpclient.RetryConfig{
		MaxRetries: 3,
		Backoff:    1 * time.Second,
	})

	defaultCacheDuration := time.Duration(s.cfg.Cache.DefaultCacheDuration * int64(math.Pow(10, 9)))
	cleanupInterval := time.Duration(s.cfg.Cache.CleanupInterval * int64(math.Pow(10, 9)))

	cache := cache.NewMemoryCache(defaultCacheDuration, cleanupInterval)
	repository := pokemonRepository.NewPokemonRepository(s.cfg.ExternalService.PokemonBaseUrl, client, cache)
	usecase := pokemonUsecase.NewPokemonUsecase(repository)
	handler := pokemonHandler.NewPokemonHttpHandler(s.cfg, usecase)

	jwtMiddleware := middleware.NewJWTMiddleware(s.cfg.Jwt.AccessSecretKey)

	//TODO: remove coupling, this is just for poc
	authRepository := authenticationRepository.NewAuthenticationRepository(usersInMemory)
	authUsecase := authenticationUsecase.NewAuthenticationUsecase(s.cfg, authRepository)
	userMiddleware := middleware.NewUserVerificationMiddleware(authUsecase)

	pokemonRouter := s.app.Group("/api/v1/pokemon")

	pokemonRouter.Use(jwtMiddleware)
	pokemonRouter.Use(userMiddleware)
	pokemonRouter.GET("/:name", handler.GetPokemonByName)
	pokemonRouter.GET("/:name/ability", handler.GetPokemonAbilityByName)
	pokemonRouter.GET("/random", handler.GetRandomPokemon)
}
