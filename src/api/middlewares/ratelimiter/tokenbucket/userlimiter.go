package tokenbucket

import (
	"api/middlewares/ratelimiter"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

//UserLimiter handle registered user
func UserLimiter(config *Bucket) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			payload, err := ratelimiter.GetJWTPayload(echo.HeaderAuthorization, "Bearer", c)
			if err != nil {
				fmt.Printf(fmt.Sprint(err))
			}
			// fmt.Printf("\nTOKEN:" + payload + "\n")
			exp := ratelimiter.GetDecodedPayload(payload)
			payloadObj, err := ratelimiter.GetPayloadMap([]byte(exp))
			// fmt.Printf(payloadObj.Hash + "\n")
			if err != nil {
				fmt.Printf(fmt.Sprint(err))
			}

			id := payloadObj.Hash

			if Limiter(config, c, id) {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusTooManyRequests)
		}
	}
}
