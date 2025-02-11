package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mikebionic/viewscount"
	"time"
)

var viewTracker *viewscount.ViewTracker

func InitializeViewTracker(db interface{}, minutesGap time.Duration) {
	viewTracker = viewscount.NewViewTracker(db, minutesGap)
}

func ViewCounterMiddleware(tableName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.Next()
			return
		}

		idInt := 0
		fmt.Sscanf(id, "%d", &idInt)

		if err := viewTracker.HandleView(c.Request, tableName, idInt); err != nil {
			fmt.Printf("Error tracking view: %v\n", err)
		}
		c.Next()
	}
}
