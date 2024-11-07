package battle

import "slices"

// PickedTargets represents targets picked by an attacker for an attack with a
// specific index.
type PickedTargets struct {
	AttackerID  CreatureID
	DefenderIDs []CreatureID
	AttackIdx   uint
}

// Equals checks if the PickedTargets is equal to the other PickedTargets.
func (this PickedTargets) Equals(other PickedTargets) bool {
	return this.AttackerID == other.AttackerID &&
		this.AttackIdx == other.AttackIdx &&
		slices.Equal(this.DefenderIDs, other.DefenderIDs)
}

// PickTargets is a function that picks the targets for the attackers.
//
// It receives two Creature slices - the attackers and the defenders.
//
// It returns a slice of PickedTargets.
//
// It returns nil if there are no attackers or no defenders to attack.
type PickTargets func(attackers, defenders []Creature) []PickedTargets

// PickTargetsFirstAlive is a PickTargets function that picks the first
// available attack and the first available not-out-of-battle defender for each
// attacker.
func PickTargetsFirstAlive(attackers, defenders []Creature) []PickedTargets {
	if len(attackers) == 0 || len(defenders) == 0 {
		return nil
	}

	picked := make([]PickedTargets, 0, len(attackers))
	for _, attacker := range attackers {
		if attacker.IsOut() {
			continue
		}

		attackIdx := pickFirstAvailableAttack(attacker)
		if attackIdx < 0 {
			continue
		}

		var targetsCnt uint
		if attacker.Attacks[attackIdx].IsBlast {
			targetsCnt = uint(len(defenders))
		} else {
			targetsCnt = 1
		}

		defenderIDs := pickFirstAvailableDefenders(targetsCnt, defenders)
		if len(defenderIDs) == 0 {
			continue
		}

		picked = append(picked, PickedTargets{
			AttackerID: attacker.ID, DefenderIDs: defenderIDs,
			AttackIdx: uint(attackIdx),
		})
	}

	if len(picked) == 0 {
		return nil
	}

	return slices.Clip(picked)
}

func pickFirstAvailableAttack(attacker Creature) int {
	for idx, attack := range attacker.Attacks {
		if attack.Charges != 0 {
			return idx
		}
	}
	return -1
}

func pickFirstAvailableDefenders(
	targetsCnt uint,
	defenders []Creature,
) []CreatureID {
	targets := make([]CreatureID, 0, targetsCnt)
	for _, defender := range defenders {
		if uint(len(targets)) >= targetsCnt {
			break
		}
		if !defender.IsOut() {
			targets = append(targets, defender.ID)
		}
	}
	return slices.Clip(targets)
}
