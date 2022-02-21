package shieldbuilder

import "strings"

// shield creates a shield struct
 type shield struct {
	front bool
	back bool
	right bool
	left bool
}

// shBuilder creates a shield(sh?) builder struct
type shBuilder struct {
	code string
}

// NewShieldBuilder constructor
func NewShieldBuilder() *shBuilder {
	return new(shBuilder)
}

// shieldBuilder methods
func (sh *shBuilder) RaiseFront() *shBuilder {
	sh.code += "F" 
	return sh
}

func (sh *shBuilder) RaiseBack() *shBuilder {
	sh.code += "B" 
	return sh
}

func (sh *shBuilder) RaiseRight() *shBuilder {
	sh.code += "R" 
	return sh
}

func (sh *shBuilder) RaiseLeft() *shBuilder {
	sh.code += "L" 
	return sh
}

// Build builds a shield
func (sh *shBuilder) Build() *shield {
	code := sh.code
	return &shield{
		front: strings.Contains(code, "F"),
		back: strings.Contains(code, "B"),
		right: strings.Contains(code, "R"),
		left: strings.Contains(code, "L"),
	}
}