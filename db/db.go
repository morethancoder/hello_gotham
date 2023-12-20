package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func Init(dotenvPath string) (*sql.DB, error) {
    err := godotenv.Load(dotenvPath)
    if err != nil {return nil, err}

    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASS")
    dbName := os.Getenv("DB_NAME")
    dbURL := fmt.Sprintf("%s:%s@/", dbUser, dbPass)

    db, err := sql.Open("mysql", dbURL)
    if err != nil {return nil, err}

    // Check if the database exists
    rows, err := db.Query("SHOW DATABASES LIKE " + fmt.Sprintf(`"%s"`, dbName))
    if err != nil {return nil, err}
    defer rows.Close()

    if !rows.Next() {
        // Create the database if it doesn't exist
        _, err := db.Exec("CREATE DATABASE " + dbName)
        if err != nil {return nil, err}
    }

    // Close the current connection and reconnect to the specific database
    db.Close()

    dbURLWithDbName:= fmt.Sprintf("%s:%s@/%s", dbUser, dbPass, dbName)
    db, err = sql.Open("mysql", dbURLWithDbName)
    if err != nil {return nil, err}

    return db, nil
}

