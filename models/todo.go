package models

import (
	"database/sql"
	"fmt"
	"time"
)

type TodoInstance struct {
	ID    int       `json:"id"`
	Title string    `json:"title"`
	Text  string    `json:"text"`
	Done  bool      `json:"done"`
	Date  time.Time `json:"date"`
    SessionToken string `json:"session_token"` 
}

const TodosTableKey string = "todos"

func InitTodosTable(dbClient *sql.DB) error {
    query := fmt.Sprintf("SHOW TABLES LIKE '%s';", TodosTableKey)
    rows, err := dbClient.Query(query); if err != nil {return err}
    defer rows.Close()

    if !rows.Next() {
        query := fmt.Sprintf(`
            CREATE TABLE  %s (
                id INT AUTO_INCREMENT PRIMARY KEY,
                title VARCHAR(255) NOT NULL,
                text TEXT,
                done BOOLEAN,
                date DATETIME,
                session_token CHAR(43) NOT NULL
            );
        `, TodosTableKey)

        _, err := dbClient.Exec(query)
        return err
    }
	return nil
}


func (i *TodoInstance) Create(dbClient *sql.DB) error {
	query := fmt.Sprintf(`
    INSERT INTO %s (title, text, done, date, session_token) 
    VALUES (?, ?, ?, ?, ?);
    `, TodosTableKey)

	_, err := dbClient.Exec(query, i.Title, i.Text, i.Done, i.Date, i.SessionToken)
	return err
}

func (i *TodoInstance) Read(dbClient *sql.DB) error {
	query := fmt.Sprintf(`
    SELECT title, text, done, date, session_token FROM %s WHERE id=?;
    `, TodosTableKey)

	var dateStr string
	err := dbClient.QueryRow(query, i.ID).Scan(&i.Title, &i.Text, &i.Done, &dateStr, &i.SessionToken)
	if err != nil {
		return err
	}

	i.Date, err = time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return err
	}
	return nil
}

func (i *TodoInstance) Update(dbClient *sql.DB) error {
	query := fmt.Sprintf(`
        UPDATE %s 
        SET title=?, text=?, done=?, date=?, session_token=?
        WHERE id=?;
    `, TodosTableKey)

	_, err := dbClient.Exec(query, i.Title, i.Text, i.Done, i.Date, i.SessionToken, i.ID)
	if err != nil {
		return err
	}
	return nil
}

func (i *TodoInstance) Delete(dbClient *sql.DB) error {
	query := fmt.Sprintf(`
        DELETE FROM %s
        WHERE id=?;
    `, TodosTableKey)

	_, err := dbClient.Exec(query, i.ID)
	if err != nil {
		return err
	}
	return nil
}

func ReadAllSessionTodos(dbClient *sql.DB, sessionID string) ([]TodoInstance, error) {
	query := fmt.Sprintf(`
        SELECT id, title, text, done, date FROM %s WHERE session_token=?;
    `, TodosTableKey)

	rows, err := dbClient.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []TodoInstance
	for rows.Next() {
		var todo TodoInstance
		var dateStr string

		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Text, &todo.Done, &dateStr); err != nil {
			return nil, err
		}

		todo.Date, err = time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

