package cmd

import (
	"fmt"
	"strconv"

	"github.com/lorenries/tsk/db"

	"github.com/muesli/coral"
)

// doCmd represents the do command
var doCmd = &coral.Command{
	Use:   "do",
	Short: "Marks a task as done",
	Run: func(cmd *coral.Command, args []string) {
		var complete []int

		for _, arg := range args {
			task, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Error parsing argument:", arg)
			} else {
				complete = append(complete, task)
			}
		}

		for _, key := range complete {
			_, err := db.MarkDone(key)
			if err != nil {
				panic(err)
			}
		}

		fmt.Println(complete)
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
