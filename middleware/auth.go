package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userDetail struct {
	ID    primitive.ObjectID `bson:"_id"`
	Email string             `json:"email" bson:"email"`
	Role  string             `json:"role" bson:"role"`
}

// IsLoggedIn is middleware for checking logged in user
func IsLoggedIn(db mongo.Database) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// get identity
			header := c.Request().Header.Get("Authorization")
			bearer := strings.Split(header, " ")
			if len(bearer) != 2 {
				return echo.ErrUnauthorized
			}

			if bearer[0] != "Bearer" {
				return echo.ErrUnauthorized
			}

			token, err := jwt.Parse(bearer[1], func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(viper.GetString(`app.secret`)), nil
			})

			if err != nil {
				return echo.ErrUnauthorized
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID, err := primitive.ObjectIDFromHex(claims["id"].(string))
				if err != nil {
					return echo.ErrUnauthorized
				}

				var user userDetail
				err = db.Collection("users").FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)

				if err != nil {
					return echo.ErrUnauthorized
				}

				c.Set("me", user)
			} else {
				return echo.ErrUnauthorized
			}

			return next(c)
		}
	}
}

// IsAdmin is middleware to check whether logged in user is admin or not
func IsAdmin(db mongo.Database) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// get from context
			me := c.Get("me")
			// convert back to model type from interface
			user := me.(userDetail)
			if user.Role != "administrator" {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}
