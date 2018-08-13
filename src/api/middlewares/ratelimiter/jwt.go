package ratelimiter

import (
	"encoding/json"
	"strings"

	b64 "encoding/base64"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Payload struct {
	Hash string `json:"hash"`
	Exp  int64  `json:"exp"`
	Jti  string `json:"jti"`
}

func GetJWTPayload(header, authScheme string, c echo.Context) (string, error) {
	auth := c.Request().Header.Get(echo.HeaderAuthorization)
	l := len(authScheme)

	if len(auth) > l+1 && auth[:l] == authScheme {
		token := auth[l+1:]
		tokens := strings.Split(token, ".") //prefix_id_time
		return tokens[1], nil
	}

	return "", middleware.ErrJWTMissing
}

func GetDecodedPayload(payload string) string {
	payloadDecoded, _ := b64.StdEncoding.DecodeString(payload)
	payloadStr := string(payloadDecoded) + "}"
	// fmt.Printf(payloadStr + "\n")

	return payloadStr
}

func GetPayloadMap(payload []byte) (Payload, error) {
	payloadObj := Payload{}
	err := json.Unmarshal(payload, &payloadObj)
	if err != nil {
		return payloadObj, err
	}
	return payloadObj, nil
}
