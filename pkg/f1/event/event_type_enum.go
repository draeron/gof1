// Code generated by go-enum
// DO NOT EDIT!

package event

import (
	"fmt"
)

const (
	// Pressed is a Type of type Pressed.
	Pressed Type = iota
	// Released is a Type of type Released.
	Released
	// Changed is a Type of type Changed.
	Changed
	// Increment is a Type of type Increment.
	Increment
	// Decrement is a Type of type Decrement.
	Decrement
)

const _TypeName = "PressedReleasedChangedIncrementDecrement"

var _TypeMap = map[Type]string{
	0: _TypeName[0:7],
	1: _TypeName[7:15],
	2: _TypeName[15:22],
	3: _TypeName[22:31],
	4: _TypeName[31:40],
}

// String implements the Stringer interface.
func (x Type) String() string {
	if str, ok := _TypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Type(%d)", x)
}

var _TypeValue = map[string]Type{
	_TypeName[0:7]:   0,
	_TypeName[7:15]:  1,
	_TypeName[15:22]: 2,
	_TypeName[22:31]: 3,
	_TypeName[31:40]: 4,
}

// ParseType attempts to convert a string to a Type
func ParseType(name string) (Type, error) {
	if x, ok := _TypeValue[name]; ok {
		return x, nil
	}
	return Type(0), fmt.Errorf("%s is not a valid Type", name)
}
