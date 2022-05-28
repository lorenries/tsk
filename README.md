# tsk

`tsk` is a terminal user interface (TUI) to-do list manager.

**Installation**

```
brew tap lorenries/tap
brew install tsk
```

**Usage**

`tsk` can be used in CLI mode or TUI mode. In TUI mode, the available keybindings are:

- ↑/k - up
- ↓/j - down
- →/l/pgdn - next page
- ←/h/pgup - previous page
- g/home - go to start 
- G/end - go to end
- enter - toggle complete
- x - delete
- / - filter
- a - add item
- P - toggle pagination
- q - quit
- ? - toggle help

In CLI mode:

```
➜  tsk git:(main) ✗ tsk --help
Manage your to-do list

Usage:
  task [flags]
  task [command]

Available Commands:
  add         Adds a task to your to-do list
  completion  Generate the autocompletion script for the specified shell
  do          Marks a task as done
  help        Help about any command
  list        Lists all of your tasks

Flags:
  -h, --help   help for task

Use "task [command] --help" for more information about a command.
```

**Todo**

- Can't mark items complete in filter mode
- `i` to edit an existing item
- tab for completed items
- categories / tags
