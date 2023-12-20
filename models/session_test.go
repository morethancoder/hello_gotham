package models_test

import (
	"morethancoder/hello_gotham/configs"
	"morethancoder/hello_gotham/db"
	"morethancoder/hello_gotham/models"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestSession(t *testing.T) { 
    err := godotenv.Load(configs.DotenvPath); if err != nil { t.Fatalf(err.Error()) }
    dbName := os.Getenv("DB_NAME")
    
    dbClient, err := db.Init(configs.DotenvPath); if err != nil {t.Fatalf(err.Error())}
    err = models.InitSessionsTable(dbClient); if err != nil {t.Fatalf(err.Error())}
    // other tests go here
    _, err = dbClient.Exec("DROP TABLE "+ models.SessionsTableKey); if err != nil {t.Fatalf(err.Error())}
    _, err = dbClient.Exec("DROP DATABASE IF EXISTS "+ dbName); if err != nil {t.Fatalf(err.Error())}

}

