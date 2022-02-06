package cmd

import (
	"fmt"
	"strings"
	"tsk/db"

	"github.com/muesli/coral"
)

// addCmd represents the add command
var addCmd = &coral.Command{
	Use:   "add",
	Short: "Adds a task to your to-do list",
	Run: func(cmd *coral.Command, args []string) {
		input := strings.Join(args, " ")
		task, err := db.CreateTask(input)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created task \"%s\"!", task.Value)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
