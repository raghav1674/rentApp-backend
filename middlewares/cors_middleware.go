package middlewares

import (
	"net/http"
	"regexp"
	"sample-web/configs"

	"github.com/gin-gonic/gin"
)


func CORSMiddleware(cfg configs.CORSConfig ) gin.HandlerFunc {
	return func(ctx *gin.Context){
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", cfg.AllowOrigin)
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", cfg.AllowHeaders)
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", cfg.AllowMethods)
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", cfg.AllowCredentials)

		// Handle preflight requests
		if ctx.Request.Method == "OPTIONS" {
			ctx.Writer.Header().Set("Access-Control-Max-Age", cfg.MaxAge)
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		re, err := regexp.Compile(cfg.AllowOrigin)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !re.MatchString(ctx.Request.Header.Get("Origin")) {
			ctx.AbortWithStatus(http.StatusForbidden) 
		}


		ctx.Next()
	}
}