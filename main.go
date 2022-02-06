package main

import (
	"fmt"
	"os"
	"path/filepath"
	"tsk/cmd"
	"tsk/db"

	"github.com/mitchellh/go-homedir"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".config", "tsk")
	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}

	dbPath := filepath.Join(configDir, "tsk.db")
	err = db.Open(dbPath)
	if err != nil {
		return err
	}

	err = cmd.Execute()
	if err != nil {
		return err
	}

	return nil
}
