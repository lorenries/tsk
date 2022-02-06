package cmd

import (
	"fmt"
	"tsk/db"

	"github.com/mergestat/timediff"
	"github.com/muesli/coral"
	"github.com/muesli/termenv"
)

var all bool
var done bool

// listCmd represents the list command
var listCmd = &coral.Command{
	Use:   "list",
	Short: "Lists all of your tasks",
	Run: func(cmd *coral.Command, args []string) {
		var tasks []db.Task
		var err error
		if all == true {
			tasks, err = db.AllTasks()
		} else if done == true {
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
	},
}

func init() {
	listCmd.Flags().BoolVarP(&all, "all", "a", false, "List active and completed tasks")
	listCmd.Flags().BoolVarP(&done, "done", "d", false, "List completed tasks only")
	rootCmd.AddCommand(listCmd)
}
