package ui

import (
	"bytes"

	"github.com/juju/ansiterm"
)

func faded(text string) string {
	ctx := ansiterm.Context{
		Foreground: ansiterm.DarkGray,
	}

	return colorize(ctx, text)
}

func colorize(ctx ansiterm.Context, text string) string {
	buf := bytes.Buffer{}
	w := ansiterm.NewWriter(&buf)
	w.SetColorCapable(true)
	ctx.Fprintf(w, text)
	return buf.String()
}
