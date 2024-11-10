package atk

import "strings"

// AttackSlice is a `[]Attack` with helper methods.
type AttackSlice []Attack

// String returns the string representation of the AttackSlice.
func (as AttackSlice) String() string {
	var sb strings.Builder
	sb.WriteString("[]Attack{")
	for i, attack := range as {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(attack.String())
	}
	sb.WriteString("}")
	return sb.String()
}

// Equals checks if the AttackSlice is equal to the other AttackSlice.
func (as AttackSlice) Equals(other AttackSlice) bool {
	if as == nil && other == nil {
		return true
	}

	if as == nil || other == nil {
		return false
	}

	if len(as) != len(other) {
		return false
	}

	for i := range as {
		if !as[i].Equals(&other[i]) {
			return false
		}
	}

	return true
}
