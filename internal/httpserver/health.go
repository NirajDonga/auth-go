package httpserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"time": time.Now().UTC(),
	})
}
