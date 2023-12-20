package models

import (
	"database/sql"
	"fmt"
	"time"
)

type SessionInstance struct {
    Token string `json:"token"`
    Data []byte `json:"data"`
    Expiry time.Time `json:"expiry"`
}

const SessionsTableKey string = "sessions"

func InitSessionsTable(dbClient *sql.DB) error {
    query := fmt.Sprintf("SHOW TABLES LIKE '%s';", SessionsTableKey)
    rows, err := dbClient.Query(query); if err != nil {return err}
    defer rows.Close()

    if !rows.Next() {
        query := fmt.Sprintf(`
            CREATE TABLE %s (
                token CHAR(43) PRIMARY KEY,
                data BLOB NOT NULL,
                expiry TIMESTAMP(6) NOT NULL
            );
         `, SessionsTableKey)

        _, err := dbClient.Exec(query); if err != nil {return err}
        query = fmt.Sprintf("CREATE INDEX sessions_expiry_idx ON %s (expiry);", SessionsTableKey)
        _, err = dbClient.Exec(query); if err != nil {return err}

    }
	return nil
}

