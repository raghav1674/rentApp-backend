package middlewares

import (
	"errors"
	"net/http"
	customerr "sample-web/errors"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        errs := c.Errors
        if len(errs) == 0 {
            return
        }
        lastErr := errs.Last().Err
        var appErr *customerr.AppError
        if errors.As(lastErr, &appErr) {
            c.JSON(appErr.Code, gin.H{"error": appErr.Message})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
    }
}
