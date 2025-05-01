package middlewares

import (
	"errors"
	"net/http"
	customerr "sample-web/errors"
	"sample-web/utils"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        
        ctx.Next()

        _, span := utils.Tracer().Start(ctx.Request.Context(), "middlewares.ErrorHandler")
		defer span.End()

        errs := ctx.Errors
        if len(errs) == 0 {
            span.AddEvent("no errors")
            return
        }
        lastErr := errs.Last().Err
        
        span.AddEvent("some error occurred")
        
        var appErr *customerr.AppError
        if errors.As(lastErr, &appErr) {
            span.AddEvent("returning AppError")
            ctx.JSON(appErr.Code, gin.H{"error": appErr.Message})
            return
        }
        span.AddEvent("returning generic error")
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
    }
}
