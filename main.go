package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	// local imports
	"CipherOps/utils/firewall"
)

// Middleware
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

func getRegister(ctx *gin.Context) {
	ctx.File("./static/register.html")
}

func panelHandler(ctx *gin.Context) {
	sessionCookie := ctx.GetString("session")
	log.Println("DEBUG: Handler session =", sessionCookie)

	// Return something so the client doesn't hang
	ctx.JSON(http.StatusOK, gin.H{
		"message": "welcome to the protected panel",
		"session": sessionCookie,
	})
}

func main() {
	router := gin.Default()
	router.Static("/static", "./static")

	router.GET("/", func (ctx *gin.Context) {
		ctx.File("./static/index.html")
	})
	router.GET("/login", func (ctx *gin.Context) {
		ctx.File("./static/login.html")
	})
	router.GET("/register", func (ctx *gin.Context) {
		ctx.File("./static/register.html")
	})

	protected := router.Group("/")
	protected.Use(ValidateSession()) 
	{
		protected.GET("/panel", panelHandler)
	}

	fmt.Println(firewall.Pr())

	log.Println("Server: http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}