package models_test

import (
	"morethancoder/hello_gotham/configs"
	"morethancoder/hello_gotham/db"
	"morethancoder/hello_gotham/models"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)


func TestTodo(t *testing.T) { 
    err := godotenv.Load(configs.DotenvPath); if err != nil { t.Fatalf(err.Error()) }
    dbName := os.Getenv("DB_NAME")
    
    dbClient, err := db.Init(configs.DotenvPath); if err != nil {t.Fatalf(err.Error())}
    err = models.InitTodosTable(dbClient); if err != nil {t.Fatalf(err.Error())}
    sessionToken := "1234567890123456789012345678901234567890123"
    i := models.TodoInstance{
        Title: "Test todo",
        Text: "",
        Date: time.Now(), 
        Done: false,
        SessionToken: sessionToken,
    }
    err = i.Create(dbClient); if err != nil {t.Fatalf(err.Error())}
    var foo models.TodoInstance = models.TodoInstance{ID: 1}

    err = foo.Read(dbClient); if err !=  nil {t.Fatalf(err.Error())}
    t.Logf(`
    Title: %s,
    Text: %s,
    Done: %v,
    Date: %v
    `, i.Title, i.Text, i.Done, i.Date)
    
    t.Log("Updating todo Text")
    foo.Text = "New Text"
    err = foo.Update(dbClient); if err != nil {t.Fatalf(err.Error())}
    todos, err := models.ReadAllSessionTodos(dbClient, sessionToken); if err != nil {t.Fatalf(err.Error())}
    for index, item := range todos {
        t.Logf("\nTODO:%d\nID:%d\nTitle:%s\nText:%s\nDone:%v\nDate:%s\n", 
        index,
        item.ID,
        item.Title,
        item.Text,
        item.Done,
        item.Date)
    }

    err = foo.Delete(dbClient); if err != nil {t.Fatalf(err.Error())}
    _, err = dbClient.Exec("DROP TABLE "+ models.TodosTableKey); if err != nil {t.Fatalf(err.Error())}
    _, err = dbClient.Exec("DROP DATABASE IF EXISTS "+ dbName); if err != nil {t.Fatalf(err.Error())}

}
