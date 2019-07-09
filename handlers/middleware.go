package handlers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

// NewTokenAuthMiddleware returns a middleware that checks token in header.
func NewTokenAuthMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			token := ctx.Request().Header.Get("Token")

			// Validate token
			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
				}

				return hmacSecret, nil
			})

			if _, ok := parsedToken.Claims.(jwt.MapClaims); !ok || !parsedToken.Valid {
				return err
			}

			if user, err := findUserByToken(db, token); err == nil {
				log.Debugf("%s from authenticated user %s", ctx.Request().Method, user.Email)
				ctx.Set("current_user", user)
				return next(ctx)
			}

			log.Errorf("%s failed with invalid token", ctx.Request().Method)
			return echo.NewHTTPError(http.StatusUnauthorized, "valid token is not presented in header")
		}
	}
}