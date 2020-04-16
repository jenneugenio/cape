package ui

import (
	"github.com/manifoldco/promptui"
)

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
// colors in the Stdout
type Details map[string]interface{}

// UI makes it easy to present prompts, animation, and other ui enhancements
type UI interface {
	Notify(t NotificationType, msg string, args ...interface{}) error
	Secret(question string, validator promptui.ValidateFunc) (string, error)
	Question(question string, validator promptui.ValidateFunc) (string, error)
	Confirm(question string) error
	Table(header TableHeader, body TableBody) error
	Template(t string, args interface{}) error
	Details(details Details) error
}
