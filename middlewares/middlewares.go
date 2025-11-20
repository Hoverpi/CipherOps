package middlewares

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

func ValidateSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionCookie, err := ctx.Cookie("session")
		if err != nil {
			// ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			// 	"error": "missing session cookie",
			// })
			// return
			ctx.Redirect(http.StatusSeeOther, "/login")
            ctx.Abort()
		}
		log.Println("DEBUG: SessionCookie=",sessionCookie)
		
		ctx.Set("session", sessionCookie)

		ctx.Next()
	}
}