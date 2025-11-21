package main

import (
	//"log"
	// local imports
	"CipherOps/automation"
	//"CipherOps/config"
	//"CipherOps/db"
	//"CipherOps/routes"
)

func main() {
	// Automate
	automate.SetupNecessaryPkgs()

	//cfg := config.LoadConfig()
	//dbConnection := db.InitDB(cfg)
	//router := routes.SetupRouter(dbConnection)

	//log.Println("Server: http://localhost:8080")
	//if err := router.Run(":8080"); err != nil {
	//	log.Fatalf("server failed: %v", err)
	//}
}