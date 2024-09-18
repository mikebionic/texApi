package services

import (
	"fmt"
	"net/http"
	"texApi/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

func HandleConnections(ctx *gin.Context) {
	token := ctx.Param("token")

	if token == "" {
		ctx.AbortWithStatus(401)
		return
	}

	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(
		token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.ENV.ACCESS_KEY), nil
		},
	)

	if err != nil {
		ctx.AbortWithStatus(403)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	ws, _ := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	user := claims["role"].(string) + fmt.Sprintf("%v", claims["id"].(float64))

	config.SocketClients[user] = ws
}
