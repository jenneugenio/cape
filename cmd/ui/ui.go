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

// Attached return a boolean representing whether or not the current session is
// attached to a terminal or not.
func Attached() bool {
	return readline.IsTerminal(int(os.Stdout.Fd()))
}
