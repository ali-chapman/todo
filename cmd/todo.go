package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func listTodos(done bool) {
	var path string = todoPath()
	if done {
		path = donePath()
	}
	readFile, err := os.Open(path)

	if err != nil {
		log.Fatalf("Error opening todo file: %v", err)
		return
	}
	defer readFile.Close()

	todos := getTodos(readFile)
	printTodos(todos)
}

func printTodos(todos []string) {
	if len(todos) == 0 {
		fmt.Println("No todos found.")
		return
	}
	for idx, todo := range todos {
		fmt.Printf("[%d]: %s\n", idx, todo)
	}
}

func addTodo(todo string) {
	file, err := os.OpenFile(todoPath(), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error opening todo file: %v", err)
		return
	}
	defer file.Close()

	todos := append(getTodos(file), todo)
	if err := saveTodos(todos, file); err != nil {
		log.Fatalf("Error saving todos: %v", err)
		return
	}
	listTodos(false)
}

func removeTodo(idx int) {
	file, err := os.OpenFile(todoPath(), os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Error opening todo file: %v", err)
		return
	}
	defer file.Close()

	todos := getTodos(file)
	if idx < 0 || idx >= len(todos) {
		log.Fatalf("Index out of range: %d", idx)
		return
	}

	completedTodo := todos[idx]
	todos = append(todos[:idx], todos[idx+1:]...)
	if err := saveTodos(todos, file); err != nil {
		log.Fatalf("Error saving todos: %v", err)
		return
	}

	// Append to done file with timestamp
	doneFile, err := os.OpenFile(donePath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening done file: %v", err)
		return
	}
	defer doneFile.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if _, err := fmt.Fprintf(doneFile, "%s - %s\n", timestamp, completedTodo); err != nil {
		log.Fatalf("Error writing to done file: %v", err)
		return
	}

	fmt.Printf("Completed: %s\n", completedTodo)
	printTodos(todos)
}

func editTodo(idx int, newTodo string) {
	file, err := os.OpenFile(todoPath(), os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Error opening todo file: %v", err)
		return
	}
	defer file.Close()

	todos := getTodos(file)
	if idx < 0 || idx >= len(todos) {
		log.Fatalf("Index out of range: %d", idx)
		return
	}
	todos[idx] = newTodo
	if err := saveTodos(todos, file); err != nil {
		log.Fatalf("Error saving todos: %v", err)
		return
	}
	printTodos(todos)
}

func saveTodos(todos []string, file *os.File) error {
	file.Truncate(0)
	file.Seek(0, 0)
	for _, todo := range todos {
		if _, err := file.WriteString(todo + "\n"); err != nil {
			return err
		}
	}
	return nil
}

func getTodos(file *os.File) []string {
	var todos []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		todos = append(todos, scanner.Text())
	}
	return todos
}

func todoPath() string {
	return os.Getenv("HOME") + "/.config/.todo"
}

func donePath() string {
	return os.Getenv("HOME") + "/.config/.done"
}
