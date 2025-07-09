package main

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
)

// Example: [!] This is an error.
func PrintError(format string, a ...interface{}) {
	fmt.Printf("[%s] %s\n", red("!"), fmt.Sprintf(format, a...))
}

// Example: [*] This is some info.
func PrintInfo(format string, a ...interface{}) {
	fmt.Printf("[%s] %s\n", yellow("*"), fmt.Sprintf(format, a...))
}

// Example: [+] Scan complete.
func PrintSuccess(format string, a ...interface{}) {
	fmt.Printf("[%s] %s\n", green("+"), fmt.Sprintf(format, a...))
}
