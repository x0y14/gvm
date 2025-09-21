package gvm

import "fmt"

type Directive interface {
	Word
	isDirective()
}

type DefineLabel int

func (d DefineLabel) String() string {
	return fmt.Sprintf("label(%d):", d)
}

func (d DefineLabel) isDirective() {}
