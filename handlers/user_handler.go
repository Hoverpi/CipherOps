package handlers

import (
    "log"
    // "database/sql"
    "net/http"
    "github.com/gin-gonic/gin"
    // "CipherOps/models"
)

// func GetUsers(db *sql.DB) gin.HandlerFunc {
//  return func(c *gin.Context) {
//   rows, err := db.Query("SELECT id, name, email FROM users")
//   if err != nil {
//    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//    return
//   }
//   defer rows.Close()

//   var users []models.User
//   for rows.Next() {
//    var user models.User
//    if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
//     c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//     return
//    }
//    users = append(users, user)
//   }
//   c.JSON(http.StatusOK, users)
//  }
// }

func PanelHandler(ctx *gin.Context) {
	sessionCookie := ctx.GetString("session")
	log.Println("DEBUG: Handler session =", sessionCookie)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "welcome to the protected panel",
		"session": sessionCookie,
	})
}