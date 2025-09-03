package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

var DB_PATH = databasePath()
var dbInstance *todoDb

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

type displayConfig struct {
	status          string
	showCreatedAt   bool
	showCompletedAt bool
	showTags        bool
	format          string
}

type todoDb struct {
	db *sql.DB
}

func (todoDb *todoDb) listTodos(config displayConfig, tagFilter string) {
	// Prepare the SQL query
	query := "SELECT id, title, done, created_at, completed_at, tags FROM todos"
	var args []any

	conditions := []string{}
	if config.status != All {
		conditions = append(conditions, "done = ?")
		args = append(args, config.status == Done)
	}
	if tagFilter != "" {
		conditions = append(conditions, "tags LIKE ?")
		args = append(args, "%\"%"+tagFilter+"%\"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := todoDb.db.Query(query, args...)
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

	printTodos(todos, config)
}

func printTodos(todos []todo, config displayConfig) {
	if len(todos) == 0 {
		return
	}

	switch config.format {
	case "json":
		printJson(todos)
	case "csv":
		printCsv(todos)
	case "txt":
		printTxt(todos)
	case "table":
		printTable(todos, config)
	default:
		log.Fatalf("Unknown format: %s", config.format)
	}
}

func printTable(todos []todo, config displayConfig) {
	// Create a tabwriter for aligned output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print column headers
	printHeaders(w, config)

	for _, t := range todos {
		printTodo(w, t, config)
	}
}

func printJson(todos []todo) {
	type jsonTodo struct {
		ID          int      `json:"id"`
		Title       string   `json:"title"`
		Done        bool     `json:"done"`
		CreatedAt   string   `json:"created_at"`
		CompletedAt string   `json:"completed_at"`
		Tags        []string `json:"tags"`
	}

	var jsonTodos []jsonTodo
	for _, t := range todos {
		jt := jsonTodo{
			ID:    t.ID,
			Title: t.Title,
			Done:  t.Done,
			Tags:  t.Tags,
		}

		jt.CreatedAt = unwrapNullString(t.CreatedAt)
		jt.CompletedAt = unwrapNullString(t.CompletedAt)

		jsonTodos = append(jsonTodos, jt)
	}

	jsonData, err := json.MarshalIndent(jsonTodos, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
}

func printTxt(todos []todo) {
	if len(todos) == 0 {
		return
	}

	for _, t := range todos {
		status := " "
		if t.Done {
			status = "x"
		}
		tagsStr := ""
		if len(t.Tags) > 0 {
			tagsStr = " " + strings.Join(t.Tags, " ")
		}
		fmt.Printf("- [%s] %s%s\n", status, t.Title, tagsStr)
	}
}

func printCsv(todos []todo) {
	if len(todos) == 0 {
		return
	}

	type csvTodo struct {
		ID          int
		Title       string
		Done        bool
		CreatedAt   string
		CompletedAt string
		Tags        string
	}

	var csvTodos []csvTodo
	for _, t := range todos {
		ct := csvTodo{
			ID:    t.ID,
			Title: t.Title,
			Done:  t.Done,
			Tags:  strings.Join(t.Tags, ","),
		}

		ct.CreatedAt = unwrapNullString(t.CreatedAt)
		ct.CompletedAt = unwrapNullString(t.CompletedAt)

		csvTodos = append(csvTodos, ct)
	}

	fmt.Println("ID,Title,Done,CreatedAt,CompletedAt,Tags")
	for _, ct := range csvTodos {
		fmt.Printf("%d,%q,%t,%s,%s,%q\n", ct.ID, ct.Title, ct.Done, ct.CreatedAt, ct.CompletedAt, ct.Tags)
	}
}

func unwrapNullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func printHeaders(w *tabwriter.Writer, config displayConfig) {
	var headers []string

	if config.status == All {
		headers = append(headers, "Status")
	}

	headers = append(headers, "ID")
	headers = append(headers, "Title")

	if config.showTags {
		headers = append(headers, "Tags")
	}

	if config.showCreatedAt {
		headers = append(headers, "Created")
	}

	if config.showCompletedAt {
		headers = append(headers, "Completed")
	}

	fmt.Fprintln(w, strings.Join(headers, "\t"))
}

func printTodo(w *tabwriter.Writer, t todo, config displayConfig) {
	var columns []string

	if config.status == All {
		status := " "
		if t.Done {
			status = "x"
		}
		columns = append(columns, fmt.Sprintf("[%s]", status))
	}

	columns = append(columns, fmt.Sprintf("%d", t.ID))
	columns = append(columns, t.Title)

	if config.showTags {
		tagsStr := fmt.Sprintf("%v", t.Tags)
		columns = append(columns, tagsStr)
	}

	if config.showCreatedAt {
		var createdAt string = ""
		if t.CreatedAt.Valid {
			createdAt = formatRelativeTime(t.CreatedAt.String)
		}
		columns = append(columns, createdAt)
	}

	if config.showCompletedAt {
		var completedAt string = ""
		if t.CompletedAt.Valid {
			completedAt = formatRelativeTime(t.CompletedAt.String)
		}
		columns = append(columns, completedAt)
	}

	fmt.Fprintln(w, strings.Join(columns, "\t"))
}

func (todoDb *todoDb) addTodo(todo string) {
	// Prepare the SQL statement, insert todo as title, done as false, and created_at as current timestamp
	stmt, err := todoDb.db.Prepare("INSERT INTO todos (title, done, tags) VALUES (?, 0, ?)")

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

func (todoDb *todoDb) completeTodo(idx int) {
	// Prepare the SQL statement to mark todo as done
	stmt, err := todoDb.db.Prepare("UPDATE todos SET done = 1, completed_at = CURRENT_TIMESTAMP WHERE id = ?")
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

func (todoDb *todoDb) deleteTodo(idx int) {
	// Prepare the SQL statement to delete the todo
	stmt, err := todoDb.db.Prepare("DELETE FROM todos WHERE id = ?")
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

func (todoDb *todoDb) editTodo(idx int, newTodo string) {
	title, tags := extractTags(newTodo)
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the SQL statement to update the todo title and tags
	stmt, err := todoDb.db.Prepare("UPDATE todos SET title = ?, tags = ? WHERE id = ?")
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

func formatRelativeTime(timestamp string) string {
	t, err := time.Parse("2006-01-02T15:04:05Z", timestamp)
	if err != nil {
		return timestamp
	}

	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	}
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	return t.Format("Jan 2")
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

func connect() (*todoDb, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

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

	// Create indexes for better performance
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_todos_done ON todos(done);`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at);`)
	if err != nil {
		return nil, err
	}

	dbInstance = &todoDb{db}
	return dbInstance, nil
}

func databasePath() string {
	envPath := os.Getenv("TODO_PATH")
	if envPath != "" {
		return envPath
	} else {
		return os.Getenv("HOME") + "/.config/.todo.db"
	}
}
