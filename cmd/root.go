package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
)

var completeFlag bool
var deleteFlag bool
var editFlag bool
var statusFlag string
var showCreatedAtFlag bool
var showCompletedAtFlag bool
var hideTagsFlag bool
var tagFilterFlag string

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A simple CLI todo application",
	Long: `Todo is a command line application that allows you to manage your tasks efficiently.

Todos are stored in a SQLite database, and you can add, edit, complete, delete, and list them with various filters.
By default, the database is stored in ~/.config/.todo.db, but you can change this by setting the TODO_PATH environment variable.`,
	Run: func(cmd *cobra.Command, args []string) {
		statusMap := map[string]string{
			"a":       "all",
			"all":     "all",
			"d":       "done",
			"done":    "done",
			"p":       "pending",
			"pending": "pending",
		}
		format := todoFormat{
			status:          statusMap[statusFlag],
			showCreatedAt:   showCreatedAtFlag,
			showCompletedAt: showCompletedAtFlag,
			showTags:        !hideTagsFlag,
		}
		db, err := connect()
		if err != nil {
			log.Fatal(err)
		}

		if editFlag && len(args) > 0 {
			idx := parseIndex(args)

			if len(args) < 2 {
				fmt.Fprintln(os.Stderr, "Please provide a new todo item to edit.")
				os.Exit(1)
			}

			newTodo := strings.Join(args[1:], " ")
			db.editTodo(idx, newTodo)
		} else if completeFlag {
			idx := parseIndex(args)
			db.completeTodo(idx)
		} else if deleteFlag {
			idx := parseIndex(args)
			db.deleteTodo(idx)
		} else if len(args) > 0 {
			todo := strings.Join(args, " ")
			db.addTodo(todo)
		}
		db.listTodos(format, tagFilterFlag)
	},
}

func parseIndex(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please provide an index.")
		os.Exit(1)
	}

	idx, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid index: %v\n", err)
		os.Exit(1)
	}

	return idx
}

func init() {
	rootCmd.Flags().BoolVarP(&completeFlag, "done", "d", false, "Complete todo by index")
	rootCmd.Flags().BoolVarP(&editFlag, "edit", "e", false, "Edit todo by index")
	rootCmd.Flags().BoolVarP(&showCreatedAtFlag, "created", "c", false, "Show creation date of todos")
	rootCmd.Flags().BoolVarP(&showCompletedAtFlag, "completed", "C", false, "Show completion date of todos")
	rootCmd.Flags().BoolVarP(&hideTagsFlag, "hide-tags", "n", false, "Hide tags of todos")
	rootCmd.Flags().StringVarP(&statusFlag, "status", "s", "pending", "Filter todos by status (all|a, done|d, pending|p)")
	rootCmd.Flags().StringVarP(&tagFilterFlag, "tag", "T", "", "Filter todos by tag")
	rootCmd.Flags().BoolVarP(&deleteFlag, "delete", "x", false, "Delete todo by index")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		os.Exit(1)
	}
}
