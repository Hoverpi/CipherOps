package db

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/jackc/pgx/v5/stdlib"
    "CipherOps/config"
)

func InitDB(cfg config.Config) *sql.DB {
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
                        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
    log.Println("BD ENV data: ", dsn)
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        log.Fatal(err)
    }
    
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }
    log.Println("Successfully connected to the database")
    return db
}
