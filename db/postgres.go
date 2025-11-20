package db

import (
 "database/sql"
 "fmt"
 _ "github.com/lib/pq"
 "myapp/config"
)

func InitDB(cfg config.Config) *sql.DB {
 dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
  cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

 db, err := sql.Open("postgres", dsn)
 if err != nil {
  panic(err)
 }
 if err := db.Ping(); err != nil {
  panic(err)
 }
 return db
}
