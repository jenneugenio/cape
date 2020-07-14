package models

type Name string

func (n Name) String() string {
	return string(n)
}
