package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var removeFlag bool
var editFlag bool
var doneFlag bool

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A simple CLI todo application",
	Long:  "Todo is a command line application that allows you to manage your tasks efficiently.",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag && len(args) > 0 {
			idx, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid index: %v\n", err)
				os.Exit(1)
			}
			if len(args) < 2 {
				fmt.Fprintln(os.Stderr, "Please provide a new todo item to edit.")
				os.Exit(1)
			}
			newTodo := args[1]
			editTodo(idx, newTodo)
		} else if removeFlag && len(args) > 0 {
			idx, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid index: %v\n", err)
				os.Exit(1)
			}
			removeTodo(idx)
		} else if len(args) > 0 {
			fmt.Println("Adding todo:", args[0])
			addTodo(args[0])
		} else {
			listTodos(doneFlag)
		}
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&removeFlag, "remove", "r", false, "Remove todo by index")
	rootCmd.Flags().BoolVarP(&editFlag, "edit", "e", false, "Edit todo by index")
	rootCmd.Flags().BoolVarP(&doneFlag, "done", "d", false, "List completed todos")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		os.Exit(1)
	}
}
