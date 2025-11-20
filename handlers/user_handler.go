package handlers

import (
 "database/sql"
 "net/http"
 "github.com/gin-gonic/gin"
 "myapp/models"
)

func GetUsers(db *sql.DB) gin.HandlerFunc {
 return func(c *gin.Context) {
  rows, err := db.Query("SELECT id, name, email FROM users")
  if err != nil {
   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
   return
  }
  defer rows.Close()

  var users []models.User
  for rows.Next() {
   var user models.User
   if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
   }
   users = append(users, user)
  }
  c.JSON(http.StatusOK, users)
 }
}
