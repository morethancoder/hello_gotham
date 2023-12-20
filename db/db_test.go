package db_test

import (
	"fmt"
	"morethancoder/hello_gotham/configs"
	"morethancoder/hello_gotham/db"
	"os"
	"testing"

	"github.com/joho/godotenv"
)


func TestInit(t *testing.T) {
    err := godotenv.Load(configs.DotenvPath); if err != nil { t.Fatalf(err.Error()) }
    dbName := os.Getenv("DB_NAME")
    
    dbClient, err := db.Init(configs.DotenvPath)
    if err != nil { t.Fatalf("Unable to create dbclient: %v", err) }
    defer dbClient.Close()
    rows, err := dbClient.Query("SHOW DATABASES LIKE " + fmt.Sprintf(`"%s"`,dbName))
    if err != nil { t.Fatalf("Database check exsistance failed: %v", err) }
    defer rows.Close()

    if !rows.Next() { 
        t.Fatalf("Database doesnt exist after creation: %v", err) 
    } else {
        _, err := dbClient.Exec("DROP DATABASE IF EXISTS " + dbName)
        if err != nil { t.Fatalf("Unable to delete database: %v", err) }
        t.Logf("Database (%s) Created and Deleted seccussfully!", dbName)
    }
}

