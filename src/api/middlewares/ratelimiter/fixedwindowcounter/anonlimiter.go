package fixedwindowcounter

import (
	"github.com/labstack/echo"
)

func AnonLimiter(config *Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.RealIP()

			if Limiter(config, c, id) {
				return next(c)
			} else {
				return echo.ErrForbidden
			}
		}
	}
}
