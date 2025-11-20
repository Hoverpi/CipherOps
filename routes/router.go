package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"CipherOps/handlers"
	"CipherOps/middlewares"
)

func SetupRouter(db *sql.DB) *gin.Engine {
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
	protected.Use(middlewares.ValidateSession()) 
	{
		protected.GET("/panel", handlers.PanelHandler)
	}

	return router
}
