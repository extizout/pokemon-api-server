package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

type (
	Config struct {
		App             App
		Jwt             Jwt
		ExternalService ExternalService
		Cache           Cache
	}

	App struct {
		Name      string
		Url       string
		Stage     string
		BodyLimit string
		Version   string
	}

	Jwt struct {
		AccessSecretKey  string
		AccessDuration   int64
		RefreshSecretKey string
		RefreshDuration  int64
	}

	ExternalService struct {
		PokemonBaseUrl string
	}

	Cache struct {
		DefaultCacheDuration int64
		CleanupInterval      int64
	}
)

func LoadConfig(path string) Config {
	_ = godotenv.Load(path)

	return Config{
		App: App{
			Name:      getEnv("APP_NAME", "app"),
			Url:       getEnv("APP_URL", "0.0.0.0:3000"),
			Stage:     getEnv("APP_STAGE", "development"),
			BodyLimit: getEnv("APP_BODY_LIMIT", "2M"),
			Version:   getEnv("APP_VERSION", "1.0.0"),
		},
		Jwt: Jwt{
			AccessSecretKey: mustGetEnv("JWT_ACCESS_SECRET_KEY"),
			AccessDuration:  mustParseInt64("JWT_ACCESS_DURATION"),
			// RefreshSecretKey: mustGetEnv("JWT_REFRESH_SECRET_KEY"),
			// RefreshDuration:  mustParseInt64("JWT_REFRESH_DURATION"),
		},
		ExternalService: ExternalService{
			PokemonBaseUrl: mustGetEnv("POKEMON_BASE_URL"),
		},
		Cache: Cache{
			DefaultCacheDuration: mustParseInt64("DEFAULT_CACHE_DURATION"),
			CleanupInterval:      mustParseInt64("CLEANUP_CACHE_INTERVAL"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Required environment variable %s not set", key)
	}
	return value
}

func mustParseInt64(key string) int64 {
	valueStr := mustGetEnv(key)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid int64 value for %s", key)
	}
	return value
}
