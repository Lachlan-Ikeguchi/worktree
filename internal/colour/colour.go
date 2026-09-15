package colour

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

// Red returns the text wrapped in red color codes.
func Red(text string) string {
	return colorRed + text + colorReset
}

// Green returns the text wrapped in green color codes.
func Green(text string) string {
	return colorGreen + text + colorReset
}

// Yellow returns the text wrapped in yellow color codes.
func Yellow(text string) string {
	return colorYellow + text + colorReset
}

// Blue returns the text wrapped in blue color codes.
func Blue(text string) string {
	return colorBlue + text + colorReset
}

// Error returns a red "Error: " prefix string.
func Error() string {
	return colorRed + "Error: " + colorReset
}

// Warning returns a yellow "Warning: " prefix string.
func Warning() string {
	return colorYellow + "Warning: " + colorReset
}

// Info returns a blue "Info: " prefix string.
func Info() string {
	return colorBlue + "Info: " + colorReset
}

// Success returns a green "Success: " prefix string.
func Success() string {
	return colorGreen + "Success: " + colorReset
}
