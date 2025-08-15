package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
)

var DB_PATH = databasePath()

const (
	Pending = "pending"
	Done    = "done"
	All     = "all"
)

type todo struct {
	ID          int
	Title       string
	Done        bool
	CreatedAt   sql.NullString
	CompletedAt sql.NullString
	Tags        []string
}

type todoFormat struct {
	status          string
	showCreatedAt   bool
	showCompletedAt bool
	showTags        bool
}

var defaultFormat = todoFormat{
	status:          Pending,
	showCreatedAt:   false,
	showCompletedAt: false,
	showTags:        false,
}

func listTodos(format todoFormat, tagFilter string) {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare the SQL query
	query := "SELECT id, title, done, created_at, completed_at, tags FROM todos"
	var args []interface{}

	conditions := []string{}
	if format.status != All {
		conditions = append(conditions, "done = ?")
		args = append(args, format.status == Done)
	}
	if tagFilter != "" {
		conditions = append(conditions, "JSON_EXTRACT(tags, '$') LIKE ?")
		args = append(args, "%\""+tagFilter+"\"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var todos []todo
	for rows.Next() {
		var t todo
		var tagsJSON string
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt, &t.CompletedAt, &tagsJSON); err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal([]byte(tagsJSON), &t.Tags); err != nil {
			log.Fatal(err)
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	for _, t := range todos {
		printTodo(t, format)
	}
}

func printTodo(t todo, format todoFormat) {
	status := " "
	if t.Done {
		status = "x"
	}
	var displayMessage string = fmt.Sprintf("[%s] %d: %s", status, t.ID, t.Title)
	if format.showCreatedAt && t.CreatedAt.Valid {
		displayMessage += fmt.Sprintf(" (Created: %s)", t.CreatedAt.String)
	}
	if format.showCompletedAt && t.CompletedAt.Valid {
		displayMessage += fmt.Sprintf(" (Completed: %s)", t.CompletedAt.String)
	}
	if format.showTags {
		displayMessage += fmt.Sprintf(" (Tags: %v)", t.Tags)
	}
	fmt.Println(displayMessage)
}

func addTodo(todo string) {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare the SQL statement, insert todo as title, done as false, and created_at as current timestamp
	stmt, err := db.Prepare("INSERT INTO todos (title, done, tags) VALUES (?, 0, ?)")

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	title, tags := extractTags(todo)
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		log.Fatal(err)
	}

	// Execute the statement
	_, err = stmt.Exec(title, string(tagsJSON))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Added todo: %s\n", todo)
}

func completeTodo(idx int) {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare the SQL statement to mark todo as done
	stmt, err := db.Prepare("UPDATE todos SET done = 1, completed_at = CURRENT_TIMESTAMP WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the statement with the provided index
	res, err := stmt.Exec(idx)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rowsAffected == 0 {
		fmt.Printf("No todo found with index %d\n", idx)
	} else {
		fmt.Printf("Marked todo %d as completed.\n", idx)
	}
}

func deleteTodo(idx int) {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare the SQL statement to delete the todo
	stmt, err := db.Prepare("DELETE FROM todos WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the statement with the provided index
	res, err := stmt.Exec(idx)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rowsAffected == 0 {
		fmt.Printf("No todo found with index %d\n", idx)
	} else {
		fmt.Printf("Deleted todo %d.\n", idx)
	}
}

func editTodo(idx int, newTodo string) {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	title, tags := extractTags(newTodo)
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the SQL statement to update the todo title and tags
	stmt, err := db.Prepare("UPDATE todos SET title = ?, tags = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the statement with the new title, tags, and index
	res, err := stmt.Exec(title, string(tagsJSON), idx)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rowsAffected == 0 {
		fmt.Printf("No todo found with index %d\n", idx)
	} else {
		fmt.Printf("Updated todo %d to: %s\n", idx, newTodo)
	}
}

func extractTags(todo string) (string, []string) {
	var tags []string
	var title string

	// Split the todo string by spaces
	words := strings.FieldsSeq(todo)
	for word := range words {
		after, hasPrefix := strings.CutPrefix(word, "@")
		if hasPrefix {
			tags = append(tags, after)
		} else {
			title += word + " "
		}
	}

	// Trim any trailing space from the title
	title = strings.TrimSpace(title)

	return title, tags
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return nil, err
	}
	// create todos table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME,
		tags TEXT DEFAULT '[]'
	);`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func databasePath() string {
	envPath := os.Getenv("TODO_PATH")
	if envPath != "" {
		return envPath
	} else {
		return os.Getenv("HOME") + "/.config/.todo.db"
	}
}
