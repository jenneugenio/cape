// package ui contains functionality for prompting for input, colorizing
// output, and creating animations.
package ui

import (
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

// NewUI returns a configured UI struct
func NewUI(cfg *config.Config) (*UI, error) {
	return &UI{
		Config:   cfg,
		Attached: Attached(),
	}, nil
}

// Confirm prompts the user with a confirmation dialog
//
// Confirmation dialogs are usually used to ask the user if they really want to
// perform an action. If stdout is not attached to a terminal then an error
// is returned.
func (u *UI) Confirm(question string) error {
	if !u.Attached {
		return errors.New(NotAttachedCause, "Can't prompt for confirmation, a terminal is not attached to stdout.")
	}

	// TODO: Come back and configure the prompt template for coloring and
	// everything else that is fun!
	prompt := &promptui.Prompt{
		Label:     question,
		IsConfirm: true,
	}

	// We mutate the promptui errors so we can display them nicely inside our
	// system!
	result, err := prompt.Run()
	if err != nil && err != promptui.ErrAbort {
		return err
	}
	if err == promptui.ErrAbort {
		return ErrAborted
	}
	if result != "y" {
		return ErrAborted
	}

	return nil
}

// Attached return a boolean representing whether or not the current session is
// attached to a terminal or not.
func Attached() bool {
	return readline.IsTerminal(int(os.Stdout.Fd()))
}
