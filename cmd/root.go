package cmd

import (
	"fmt"
	"os"
	"tsk/app"
	"tsk/db"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/coral"
)

var rootCmd = &coral.Command{
	Use:   "task",
	Short: "Manage your to-do list",
	Run: func(cmd *coral.Command, args []string) {
		tasks, err := db.ActiveTasks()
		if err != nil {
			panic(err)
		}
		m := app.NewModel(tasks)
		program := tea.NewProgram(m)
		if err := program.Start(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}
