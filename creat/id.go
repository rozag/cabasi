package creat

// ID is a unique identifier of a creature.
type ID string

// CompareTo returns an integer comparing two ids lexicographically. The result
// will be 0 if this == other, -1 if this < other, and +1 if this > other.
func (id ID) CompareTo(other ID) int {
	if id < other {
		return -1
	}
	if id > other {
		return 1
	}
	return 0
}
