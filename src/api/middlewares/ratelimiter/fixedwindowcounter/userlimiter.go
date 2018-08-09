package fixedwindowcounter

import (
	"encoding/json"
	"fmt"
	"log"
	"redis"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func getJWT(header, authScheme string, c echo.Context) (string, error) {
	auth := c.Request().Header.Get(echo.HeaderAuthorization)
	l := len(authScheme)
	if len(auth) > l+1 && auth[:l] == authScheme {
		return auth[l+1:], nil
	}
	return "", middleware.ErrJWTMissing
}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := "" // Value
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

func UserLimiter(config *Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := getJWT(echo.HeaderAuthorization, "Bearer", c)
			if err != nil {
				fmt.Printf(fmt.Sprint(err))
			}
			fmt.Printf("\nTOKEN:" + token + "\n")
			id := token

			r := redis.RedisConnect()
			defer r.Close()

			keys, err := r.Do("KEYS", config.Prefix+":"+id+":*")
			if err != nil {
				fmt.Printf(fmt.Sprint(err))
			}
			if fmt.Sprint(keys) == "[]" {
				Init(config, id)
			} else {
				for _, k := range keys.([]interface{}) {
					key := string(k.([]byte))
					params := strings.Split(key, ":") //prefix_id_time

					savedTime, err := strconv.ParseInt(params[2], 10, 64)
					if err != nil {
						fmt.Printf(fmt.Sprint(err))
					}
					reqTime := time.Now().Unix()
					duration := reqTime - savedTime
					if duration > GetPeriodInt(config.Period) {
						Init(config, id)
						_, err := r.Do("DEL", k.([]byte))
						if err != nil {
							fmt.Printf(fmt.Sprint(err))
						}
					} else {
						var counter uint
						reply, err := r.Do("GET", k.([]byte))
						if err != nil {
							fmt.Printf(fmt.Sprint(err))
						}
						if err := json.Unmarshal(reply.([]byte), &counter); err != nil {
							fmt.Printf(fmt.Sprint(err))
						}
						if counter < config.MaxRequest {
							_, err := r.Do("INCR", k.([]byte))
							if err != nil {
								fmt.Printf(fmt.Sprint(err))
							}
							SetHeader(c, config.MaxRequest, config.MaxRequest-counter, GetPeriodInt(config.Period)-duration)
							return next(c)
						} else {
							fmt.Printf("Limit Exceeded for :" + id + "\n")
							SetHeader(c, config.MaxRequest, 0, GetPeriodInt(config.Period)-duration)
							return echo.ErrForbidden
						}
					}
				}
			}
			SetHeader(c, config.MaxRequest, config.MaxRequest, GetPeriodInt(config.Period))
			return next(c)
		}
	}
}
