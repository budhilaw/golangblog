package middleware

import (
	"github.com/labstack/echo"
)

type defaultMiddleware struct{}

// ResponseError ...
type ResponseError struct {
	Message string `json:"message"`
}

func (m *defaultMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

// InitMiddleware ...
func InitMiddleware() *defaultMiddleware {
	return &defaultMiddleware{}
}
