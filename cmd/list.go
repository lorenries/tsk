package cmd

import (
	"fmt"

	"github.com/lorenries/tsk/db"

	"github.com/mergestat/timediff"
	"github.com/muesli/coral"
	"github.com/muesli/termenv"
)

var listFlags struct {
	allFlag, doneFlag, tagsFlag bool
	tagFlag                     string
}

// listCmd represents the list command
var listCmd = &coral.Command{
	Use:   "list",
	Short: "Lists all of your tasks",
	Run: func(cmd *coral.Command, args []string) {
		if listFlags.tagsFlag {
			tags, err := db.AllTags()
			if err != nil {
				panic(err)
			}
			for _, tag := range tags {
				fmt.Println(tag)
			}
		} else {
			var tasks []db.Task
			var err error
			if listFlags.allFlag == true {
				tasks, err = db.AllTasks()
			} else if listFlags.doneFlag == true {
				tasks, err = db.CompletedTasks()
			} else {
				tasks, err = db.ActiveTasks()
			}

			if err != nil {
				panic(err)
			}

			for _, task := range tasks {
				added := termenv.String(timediff.TimeDiff(task.TimeAdded)).Faint()
				var checkbox string
				if task.Completed {
					checkbox = fmt.Sprint("[x]")
				} else {
					checkbox = fmt.Sprint("[ ]")
				}
				fmt.Println(checkbox, task.Key, task.Value, added)
			}
		}
	},
}

func init() {
	listCmd.Flags().BoolVarP(&listFlags.allFlag, "all", "a", false, "List active and completed tasks")
	listCmd.Flags().BoolVarP(&listFlags.doneFlag, "done", "d", false, "List completed tasks only")
	listCmd.Flags().BoolVar(&listFlags.tagsFlag, "tags", false, "List all tags")
	listCmd.Flags().StringVarP(&listFlags.tagFlag, "tag", "t", "", "List tasks for a given tag")
	rootCmd.AddCommand(listCmd)
}
