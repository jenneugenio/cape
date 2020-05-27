// package ui contains functionality for prompting for input, colorizing
// output, and creating animations.
package ui

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/juju/ansiterm"
	"github.com/leekchan/gtf"
	"os"
	"text/template"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"

	errors "github.com/capeprivacy/cape/partyerrors"
)

// Stdout is an implementation of UI that prints to stdout
// while taking into account the state of a users terminal.
type Stdout struct {
	Attached   bool
	Animations bool
	Colors     bool
}

// NewStdout returns a configured Stdout struct
func NewStdout(colors bool, animations bool) *Stdout {
	return &Stdout{
		Colors:     colors,
		Animations: animations,
		Attached:   Attached(),
	}
}

func (u *Stdout) prompt(p *promptui.Prompt) (string, error) {
	if !u.Attached {
		return "", errors.New(NotAttachedCause, "Can't prompt for question, a terminal is not attached to stdout.")
	}

	// We mutate the promptui errors so we can display them nicely inside our
	// system!
	result, err := p.Run()
	if err != nil && err != promptui.ErrAbort {
		return "", err
	}
	if err == promptui.ErrAbort {
		return "", ErrAborted
	}

	return result, nil
}

// Notify is used to notify the user about some piece of information including
// error messages or warning.
func (u *Stdout) Notify(t NotificationType, msg string, args ...interface{}) error {
	shorthand := fmt.Sprintf("%s %s", t.Symbol, t.Text)
	if u.CanColorize() {
		shorthand = colorize(t.Color, shorthand)
	}

	out := fmt.Sprintf(shorthand+": "+msg, args...)
	fmt.Fprintf(os.Stdout, out+"\n")
	return nil
}

// Secret prompts the user to answer a terminal question that will be masked
func (u *Stdout) Secret(question string, validator promptui.ValidateFunc) (string, error) {
	p := &promptui.Prompt{
		Label:    question,
		Validate: validator,
		Mask:     '*',
	}

	return u.prompt(p)
}

// Question promps the user to answer a terminal question
func (u *Stdout) Question(question string, validator promptui.ValidateFunc) (string, error) {
	p := &promptui.Prompt{
		Label:    question,
		Validate: validator,
	}

	return u.prompt(p)
}

// Confirm prompts the user with a confirmation dialog
//
// Confirmation dialogs are usually used to ask the user if they really want to
// perform an action. If stdout is not attached to a terminal then an error
// is returned.
func (u *Stdout) Confirm(question string) error {
	// TODO: Come back and configure the prompt template for coloring and
	// everything else that is fun!
	p := &promptui.Prompt{
		Label:     question,
		IsConfirm: true,
	}

	r, err := u.prompt(p)
	if err != nil {
		return err
	}

	if r != "y" {
		return ErrAborted
	}

	return nil
}

// Table prints the provided header and body to the Stdout.
func (u *Stdout) Table(header TableHeader, body TableBody) error {
	w := ansiterm.NewTabWriter(os.Stdout, 2, 0, 4, ' ', 0)

	if u.CanColorize() {
		w.SetStyle(ansiterm.Bold)
	}

	for _, h := range header {
		fmt.Fprintf(w, "%s\t", h)
	}

	w.Reset()
	fmt.Fprintln(w)

	for _, row := range body {
		for _, itm := range row {
			fmt.Fprintf(w, "%s\t", itm)
		}

		fmt.Fprintln(w)
	}

	return w.Flush()
}

// Details prints the provided details to the Stdout. Details is a Stdout Component
// with labelled key and value pairs.
func (u *Stdout) Details(details Details) error {
	w := ansiterm.NewTabWriter(os.Stdout, 2, 0, 4, ' ', 0)
	for label, value := range details {
		if u.CanColorize() {
			label = faded(label)
		}

		out := ""
		switch v := value.(type) {
		case fmt.Stringer:
			out = v.String()
		case string:
			out = v
		default:
			return ErrCantDisplay
		}

		fmt.Fprintf(w, "%s:\t%s\n", label, out)
	}

	// Always make whitespace after a table
	fmt.Fprintln(w)
	return w.Flush()
}

func (u *Stdout) funcMap() template.FuncMap {
	return template.FuncMap{
		"faded": func(t string) string {
			if !u.CanColorize() {
				return t
			}
			return faded(t)
		},
		"bold": func(t string) string {
			if !u.CanColorize() {
				return t
			}
			return bold(t)
		},
		"italic": func(t string) string {
			if !u.CanColorize() {
				return t
			}
			return italic(t)
		},
		"color": func(c ansiterm.Color, t string) string {
			if !u.CanColorize() {
				return t
			}
			return colorize(c, t)
		},
	}
}

// Template takes in a text/template style template and renders it
func (u *Stdout) Template(t string, args interface{}) error {
	tmpl, err := template.New("template").
		Funcs(gtf.GtfTextFuncMap).
		Funcs(sprig.TxtFuncMap()).
		Funcs(u.funcMap()).
		Parse(t)
	if err != nil {
		return err
	}

	w := ansiterm.NewTabWriter(os.Stdout, 2, 0, 4, ' ', 0)
	err = tmpl.Execute(w, args)
	if err != nil {
		return err
	}

	return w.Flush()
}

// CanColorized returns whether or not colorization can be supported
func (u *Stdout) CanColorize() bool {
	return u.Colors && u.Attached
}

// Attached return a boolean representing whether or not the current session is
// attached to a terminal or not.
func Attached() bool {
	return readline.IsTerminal(int(os.Stdout.Fd()))
}
