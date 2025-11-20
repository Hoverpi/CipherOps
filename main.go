package main

import (
	"fmt"
	"log"
	// local imports
	"CipherOps/utils"
	"CipherOps/config"
	"CipherOps/db"
	"CipherOps/routes"
)

func main() {
	cfg := config.LoadConfig()
	dbConnection := db.InitDB(cfg)
	router := routes.SetupRouter(dbConnection)

	fmt.Println(firewall.Pr())

	log.Println("Server: http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}