package picktargets

import (
	"slices"

	"github.com/rozag/cabasi/creat"
)

// FirstAlive is a function that picks the first available not-out-of-battle
// defenders for the selected attack.
// It receives an attacker, a picked attack index, and a slice of defenders.
// It returns a slice of picked defender indexes.
// It returns nil if there are no defenders to attack.
func FirstAlive(
	attacker creat.Creature,
	pickedAttackIdx uint,
	defenders []creat.Creature,
) []uint {
	if len(defenders) == 0 {
		return nil
	}

	// TODO: test if attacker is out

	// TODO: test if pickedAttackIdx >= len(attacker.Attacks)

	// TODO:
	// - Attacks against detachments by individuals are Impaired (excluding Blast
	// damage).
	// - Attacks against individuals by detachments are Enhanced and deal Blast
	// damage.
	// NOTE: only this last part about detachments to individuals attacks is
	// relevant here

	attack := attacker.Attacks[pickedAttackIdx]
	if attack.Charges == 0 {
		return nil
	}

	targetsCnt := 1
	if attack.IsBlast {
		targetsCnt = len(defenders)
	}

	indexes := make([]uint, 0, targetsCnt)
	for idx, defender := range defenders {
		if len(indexes) >= targetsCnt {
			break
		}

		if !defender.IsOut() {
			// Suppressing gosec "G115 integer overflow conversion int -> uint"
			// because int index will never overflow a uint variable.
			indexes = append(indexes, uint(idx)) //nolint:gosec
		}
	}

	if len(indexes) == 0 {
		return nil
	}

	return slices.Clip(indexes)
}
