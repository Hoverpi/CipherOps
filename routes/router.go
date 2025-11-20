package routes

import (
 "database/sql"
 "github.com/gin-gonic/gin"
 "myapp/handlers"
)

func SetupRouter(db *sql.DB) *gin.Engine {
 r := gin.Default()
 r.GET("/users", handlers.GetUsers(db))
 return r
}
