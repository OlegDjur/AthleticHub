package controller

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type UserClaims struct {
	UID   string   `json:"uid"`
	Roles []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

func JWTMiddleware(secret []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing bearer token")
			}
			tokenStr := strings.TrimSpace(h[len("Bearer "):])

			claims := new(UserClaims)
			token, err := jwt.ParseWithClaims(
				tokenStr, claims,
				func(t *jwt.Token) (any, error) {
					if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
						return nil, echo.NewHTTPError(http.StatusUnauthorized, "alg mismatch")
					}
					return secret, nil
				},
				jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
				jwt.WithExpirationRequired(), // требуем exp
			)
			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}
			if !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			c.Set("claims", claims)
			c.Set("sub", claims.Subject)
			return next(c)
		}
	}
}
