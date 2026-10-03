package utils

import "github.com/pterm/pterm"

// PrintSuccess prints a green success message to the terminal.
func PrintSuccess(msg string) {
	pterm.Success.Println(msg)
}

// PrintError prints a red error message to the terminal.
func PrintError(msg string) {
	pterm.Error.Println(msg)
}

// PrintInfo prints a cyan informational message to the terminal.
func PrintInfo(msg string) {
	pterm.Info.Println(msg)
}

// PrintWarning prints a yellow warning message to the terminal.
func PrintWarning(msg string) {
	pterm.Warning.Println(msg)
}

// PrintDebug prints a gray debug message, but only when verbose mode is enabled.
func PrintDebug(msg string) {
	pterm.Debug.Println(msg)
}
