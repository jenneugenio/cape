package ui

import (
	"github.com/juju/ansiterm"
)

// NotificationType represents a type of notification such as a warning or error
type NotificationType struct {
	Symbol string
	Text   string
	Color  ansiterm.Color
}

var (
	// Warn represents a warning notification to the user
	Warn = NotificationType{
		Symbol: "⚠",
		Text:   "Warning",
		Color:  ansiterm.Yellow,
	}

	// Error represents an error notification to the user
	Error = NotificationType{
		Symbol: "✗",
		Text:   "Error",
		Color:  ansiterm.Red,
	}

	// Remember represents a reminder notification to the user
	Remember = NotificationType{
		Symbol: "‼",
		Text:   "Remember",
		Color:  ansiterm.Green,
	}
)
