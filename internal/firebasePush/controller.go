package firebasePush

import (
	"texApi/config"
	"texApi/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func Controller(router *gin.Engine) {
	group := router.Group(config.ENV.API_PREFIX + "/auth/")

	group.POST("/save-notification-token/", middlewares.Guard, SaveNotificationToken)
}
