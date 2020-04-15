package ui

import (
	"strings"

	"github.com/juju/ansiterm"
)

func colorize(color ansiterm.Color, text string) string {
	ctx := ansiterm.Context{
		Foreground: color,
	}

	return stylize(ctx, text)
}

func faded(text string) string {
	ctx := ansiterm.Context{
		Foreground: ansiterm.DarkGray,
	}

	return stylize(ctx, text)
}

func bold(text string) string {
	ctx := ansiterm.Context{
		Styles: []ansiterm.Style{ansiterm.Bold},
	}

	return stylize(ctx, text)
}

func stylize(ctx ansiterm.Context, text string) string {
	builder := &strings.Builder{}
	w := ansiterm.NewWriter(builder)
	w.SetColorCapable(true)
	ctx.Fprintf(w, text)
	return builder.String()
}
