package echo

import (
	"strconv"

	"github.com/labstack/echo"

	"github.com/ulule/limiter"
)

// Middleware is the middleware for basic http.Handler.
type Middleware struct {
	Limiter        *limiter.Limiter
	OnError        ErrorHandler
	OnLimitReached LimitReachedHandler
	KeyGetter      KeyGetter
}

// NewMiddleware return a new instance of a basic HTTP middleware.
func NewMiddleware(limiter *limiter.Limiter, options ...Option) echo.MiddlewareFunc {
	middleware := &Middleware{
		Limiter:        limiter,
		OnError:        DefaultErrorHandler,
		OnLimitReached: DefaultLimitReachedHandler,
		KeyGetter:      DefaultKeyGetter,
	}

	for _, option := range options {
		option.apply(middleware)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := middleware.KeyGetter(c)
			context, err := middleware.Limiter.Get(c, key)
			if err != nil {
				middleware.OnError(c, err)
				return echo.ErrForbidden
			}

			c.Response().Header().Set("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
			c.Response().Header().Set("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
			c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

			if context.Reached {
				middleware.OnLimitReached(c)
				return echo.ErrForbidden
			}

			return next(c)
		}
	}
}
