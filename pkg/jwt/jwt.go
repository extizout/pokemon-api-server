package jwt

import (
	"errors"
	"math"
	"time"

	"github.com/extizout/pokemon-api-server/pkg/authentication"
	"github.com/golang-jwt/jwt/v5"
)

type (
	jwtAuthFactory interface {
		SignToken() string
	}

	AuthMapClaims struct {
		*authentication.Claims
		jwt.RegisteredClaims
	}

	authConcreate struct {
		Secret []byte
		Claims *AuthMapClaims `json:"claims"`
	}

	accessToken struct {
		*authConcreate
	}
	refreshToken struct {
		*authConcreate
	}
)

var (
	ErrTokenMalformed               = errors.New("token format is invalid")
	ErrTokenExpired                 = errors.New("token is expired")
	ErrTokenInvalid                 = errors.New("token is invalid")
	ErrTokenUnexpectedSigningMethod = errors.New("unexpected signing method")
)

func (a *authConcreate) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.Claims)
	signedString, _ := token.SignedString(a.Secret)

	return signedString
}

func jwtTimeDurationCal(sec int64) *jwt.NumericDate {
	durationTimeAdded := time.Now().Add(time.Duration(sec * int64(math.Pow(10, 9))))
	return jwt.NewNumericDate(durationTimeAdded)
}

func jwtTimeRepeatAdaptor(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func NewAccessToken(secret string, expiredAt int64, claims *authentication.Claims) jwtAuthFactory {
	return &accessToken{
		authConcreate: &authConcreate{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					//TODO: make this dynamic based on the environment
					Issuer:    "Pokemon-API",
					Subject:   "access-token",
					Audience:  []string{"localhost"},
					ExpiresAt: jwtTimeDurationCal(expiredAt),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}

func NewRefreshToken(secret string, expiredAt int64, claims *authentication.Claims) jwtAuthFactory {
	return &refreshToken{
		authConcreate: &authConcreate{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					//TODO: make this dynamic based on the environment
					Issuer:    "Pokemon-API",
					Subject:   "refresh-token",
					Audience:  []string{"localhost"},
					ExpiresAt: jwtTimeDurationCal(expiredAt),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}

func ReloadRefreshToken(secret string, expiredAt int64, claims *authentication.Claims) string {
	payload := &refreshToken{
		authConcreate: &authConcreate{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					//TODO: make this dynamic based on the environment
					Issuer:    "Pokemon-API",
					Subject:   "refresh-token",
					Audience:  []string{"localhost"},
					ExpiresAt: jwtTimeRepeatAdaptor(expiredAt),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
	return payload.SignToken()
}

func ParseToken(secret string, tokenString string) (*AuthMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthMapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenUnexpectedSigningMethod
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else {
			return nil, ErrTokenInvalid
		}
	}

	if claims, ok := token.Claims.(*AuthMapClaims); ok {
		return claims, nil
	}

	return nil, nil
}
