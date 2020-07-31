package ui

import "github.com/manifoldco/promptui"

// Mock is a fake UI that will track what it was called with
type Mock struct {
	Calls []*Call
}

// Call is a helper struct to record types of calls
type Call struct {
	Name string
	Args []interface{}
}

func (m *Mock) appendCall(name string, args ...interface{}) {
	call := &Call{
		Name: name,
		Args: args,
	}

	m.Calls = append(m.Calls, call)
}

// Template from the UI interface
func (m *Mock) Template(t string, args interface{}) error {
	m.appendCall("template", t, args)
	return nil
}

// Details from the UI interface
func (m *Mock) Details(details Details) error {
	m.appendCall("details", details)
	return nil
}

// Notify from the UI interface
func (m *Mock) Notify(t NotificationType, msg string, args ...interface{}) error {
	m.appendCall("notify", t, args)
	return nil
}

// Secret from the UI interface
func (m *Mock) Secret(question string, validator promptui.ValidateFunc) (string, error) {
	m.appendCall("secret", question, validator)
	return "abc", nil
}

// Question from the UI interface
func (m *Mock) Question(question string, validator promptui.ValidateFunc) (string, error) {
	m.appendCall("question", question, validator)
	return "yes", nil
}

// Confirm from the UI interface
func (m *Mock) Confirm(question string) error {
	m.appendCall("confirm", question)
	return nil
}

// Table from the UI interface
func (m *Mock) Table(header TableHeader, body TableBody) error {
	m.appendCall("table", header, body)
	return nil
}
