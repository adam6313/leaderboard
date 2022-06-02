package v1

import (
	"leaderboard/pkg/response"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// C -
type C struct {
	iris.Context
}

// HandleFunc custom iris context
func HandleFunc(handler func(*C)) func(iris.Context) {
	return func(c iris.Context) {
		customerContext := &C{
			c,
		}
		handler(customerContext)
	}
}

// Cros for iris cros middleware
func Cros() context.Handler {
	return func(ctx context.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.Header("Content-Type", "application/json")
		ctx.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		ctx.Next()
	}
}

// R this fn for success response
func (c *C) R(data interface{}) {
	c.StatusCode(iris.StatusOK)
	c.JSON(response.Success(c.Request().Context(), data))
}

// E this fn for error response
func (c *C) E(err error) {
	c.StatusCode(iris.StatusOK)
	c.JSON(response.Error(c.Request().Context(), err))
}
