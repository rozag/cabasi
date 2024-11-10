package creat

import "strings"

// CreatureSlice is a `[]Creature` with helper methods.
type CreatureSlice []Creature

// String returns the string representation of the CreatureSlice.
func (cs CreatureSlice) String() string {
	var sb strings.Builder
	sb.WriteString("[]Creature{")
	for i, c := range cs {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(c.String())
	}
	sb.WriteString("}")
	return sb.String()
}

// Equals checks if the CreatureSlice is equal to the other CreatureSlice.
func (cs CreatureSlice) Equals(other CreatureSlice) bool {
	if cs == nil && other == nil {
		return true
	}

	if cs == nil || other == nil {
		return false
	}

	if len(cs) != len(other) {
		return false
	}

	for i := range cs {
		if !cs[i].Equals(&other[i]) {
			return false
		}
	}

	return true
}
