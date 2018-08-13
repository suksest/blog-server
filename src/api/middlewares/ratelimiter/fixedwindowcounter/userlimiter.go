package fixedwindowcounter

import (
	"encoding/json"
	"fmt"
	"redis"
	"strconv"
	"strings"
	"time"

	b64 "encoding/base64"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Payload struct {
	Hash string `json:"hash"`
	Exp  string `json:"exp"`
	Jti  string `json:"jti"`
}

func getJWTPayload(header, authScheme string, c echo.Context) (string, error) {
	auth := c.Request().Header.Get(echo.HeaderAuthorization)
	l := len(authScheme)

	if len(auth) > l+1 && auth[:l] == authScheme {
		token := auth[l+1:]
		tokens := strings.Split(token, ".") //prefix_id_time
		return tokens[1], nil
	}

	return "", middleware.ErrJWTMissing
}

func getDecodedPayload(payload string) (string, error) {
	payloadDecoded, err := b64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}
	return payloadDecoded, nil
}

func getPayloadMap(payload []byte) Payload {
	payloadObj := Payload{}
	err := json.Unmarshal(payload, &payloadObj)
	if err != nil {
		fmt.Printf(payloadObj.Hash)
	}
	fmt.Printf(payloadObj.Hash)
	return payloadObj
}

func UserLimiter(config *Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			payload, err := getJWTPayload(echo.HeaderAuthorization, "Bearer", c)
			if err != nil {
				fmt.Printf(fmt.Sprint(err))
			}
			fmt.Printf("\nTOKEN:" + payload + "\n")
			id := payload

			exp, err := getDecodedPayload(payload)

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
