package ginutils

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetClientIPFromXFF returns the client ID using x-forwarded-for headers before relying on c.ClientIP().
// This assumes use of a trusted proxy.
func GetClientIPFromXFF(c *gin.Context) string {
	forwardHeader := c.Request.Header.Get("x-forwarded-for")
	firstAddress := strings.Split(forwardHeader, ",")[0]
	if net.ParseIP(strings.TrimSpace(firstAddress)) != nil {
		return firstAddress
	}

	return c.ClientIP()
}
