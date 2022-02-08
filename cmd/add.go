package cmd

import (
	"fmt"
	"strings"

	"github.com/lorenries/tsk/db"

	"github.com/muesli/coral"
)

var tagFlag bool

// addCmd represents the add command
var addCmd = &coral.Command{
	Use:   "add",
	Short: "Adds a task to your to-do list",
	Run: func(cmd *coral.Command, args []string) {
		input := strings.Join(args, " ")
		if tagFlag {
			tag, err := db.CreateTag(input)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Created tag \"%s\"!", tag.Value)
		} else {
			task, err := db.CreateTask(input)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Created task \"%s\"!", task.Value)
		}
	},
}

func init() {
	addCmd.Flags().BoolVarP(&tagFlag, "tag", "t", false, "Adds a new tag")
	rootCmd.AddCommand(addCmd)
}
