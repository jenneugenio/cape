package models

type Label string

func (l Label) String() string {
	return string(l)
}
