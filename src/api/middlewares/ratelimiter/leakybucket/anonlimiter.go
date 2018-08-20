package leakybucket

import (
	"net/http"

	"github.com/labstack/echo"
)

//AnonLimiter handle request for anonymous user
func AnonLimiter(config *Bucket) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.RealIP()

			if Limiter(config, c, id) {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusTooManyRequests)
		}
	}
}
