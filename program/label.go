package program

import "fmt"

type Label interface {
	Word
	isLabel()
}

type AbstractIndex int

func (a *AbstractIndex) String() string {
	return fmt.Sprintf("label(%d)", *a)
}

func (a *AbstractIndex) isLabel() {}
