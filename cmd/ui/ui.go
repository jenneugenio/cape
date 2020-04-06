// package ui contains functionality for prompting for input, colorizing
// output, and creating animations.
package ui

import (
	"fmt"
	"github.com/juju/ansiterm"
	"os"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"

	"github.com/dropoutlabs/cape/cmd/config"
	errors "github.com/dropoutlabs/cape/partyerrors"
)

// UI makes it easy to present prompts, animation, and other ui enhancements
// while taking into account the state of a users terminal.
type UI struct {
	Config   *config.Config
	Attached bool
}

// TableHeader is used with the `Table` function. It will print the header of your table in bold.
// E.g.
// [ "Name", "Number", "Gibberish" ]
type TableHeader []string

// TableBody is used with the `Table` function. It is a slice of string slices that contain the body of your table.
// E.g.
// [
//   [ "Ben", "100", "jdkslajdklas" ],
//   [ "Ian", "100", "jdklsajdklsa" ],
//   [ "Justin", "100", "kljdaslkdaskjl" ],
// ]
type TableBody [][]string

// Details is used to colorize a series of key and value pairs when listing
// colors in the UI
type Details map[string]interface{}

// NewUI returns a configured UI struct
func NewUI(cfg *config.Config) (*UI, error) {
	return &UI{
		Config:   cfg,
		Attached: Attached(),
	}, nil
}

func (u *UI) prompt(p *promptui.Prompt) (string, error) {
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
func (u *UI) Notify(t NotificationType, msg string, args ...interface{}) error {
	shorthand := fmt.Sprintf("%s %s", t.Symbol, t.Text)
	if u.CanColorize() {
		shorthand = colorize(t.Color, shorthand)
	}

	out := fmt.Sprintf(shorthand+": "+msg, args...)
	fmt.Fprintf(os.Stdout, out+"\n")
	return nil
}

// Secret prompts the user to answer a terminal question that will be masked
func (u *UI) Secret(question string, validator promptui.ValidateFunc) (string, error) {
	p := &promptui.Prompt{
		Label:    question,
		Validate: validator,
		Mask:     '*',
	}

	return u.prompt(p)
}

// Question promps the user to answer a terminal question
func (u *UI) Question(question string, validator promptui.ValidateFunc) (string, error) {
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
func (u *UI) Confirm(question string) error {
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

// Table prints the provided header and body to the UI.
func (u *UI) Table(header TableHeader, body TableBody) error {
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

// Details prints the provided details to the UI. Details is a UI Component
// with labelled key and value pairs.
func (u *UI) Details(details Details) error {
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

	return w.Flush()
}

// CanColorized returns whether or not colorization can be supported
func (u *UI) CanColorize() bool {
	return u.Config.UI.Colors && u.Attached
}

// Attached return a boolean representing whether or not the current session is
// attached to a terminal or not.
func Attached() bool {
	return readline.IsTerminal(int(os.Stdout.Fd()))
}
