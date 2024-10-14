package middlewares

import "github.com/gin-gonic/gin"

func Cors(ctx *gin.Context) {

	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
	ctx.Writer.Header().Set("AllowFiles", "*")
	ctx.Writer.Header().Set("AllowCredentials", "*")

	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(204)
		return
	}

	ctx.Next()
}
