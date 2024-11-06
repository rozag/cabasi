package battle

import "slices"

// PickedTargets represents targets picked by an attacker for an attack with a
// specific index.
type PickedTargets struct {
	AttackerID  CreatureID
	AttackIdx   uint
	DefenderIDs []CreatureID
}

// Equals checks if the PickedTargets is equal to the other PickedTargets.
func (this PickedTargets) Equals(other PickedTargets) bool {
	return this.AttackerID == other.AttackerID &&
		this.AttackIdx == other.AttackIdx &&
		slices.Equal(this.DefenderIDs, other.DefenderIDs)
}

// PickTargets is a function that picks the targets for the attackers.
//
// It receives two CreatureID to Creature maps - the attackers and the
// defenders.
//
// It returns a slice of PickedTargets.
//
// It returns nil if there are no attackers or no defenders to attack.
type PickTargets func(a, d map[CreatureID]Creature) []PickedTargets

// PickTargetsFirstAlive is a PickTargets function that picks the first
// not-out-of-battle defender for each attacker.
func PickTargetsFirstAlive(a, d map[CreatureID]Creature) []PickedTargets {
	// TODO:
	return nil
}
