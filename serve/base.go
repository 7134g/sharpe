package serve

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type BaseContext struct {
	*gin.Context
}

func (c *BaseContext) QueryInt(key string) (int, error) {
	value := c.Context.Query(key)
	return strconv.Atoi(value)
}

func (c *BaseContext) QueryIntDefault(key string) int {
	value := c.Context.Query(key)
	if value == "" {
		return 0
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return v
}
