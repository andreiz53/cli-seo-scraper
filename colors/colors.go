package colors

import "github.com/fatih/color"

var (
	Error   = color.New(color.FgRed).SprintFunc()
	Success = color.New(color.FgGreen).SprintFunc()
	Warning = color.New(color.FgYellow).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()
	Info    = color.New(color.Bold, color.FgCyan).SprintFunc()
)
