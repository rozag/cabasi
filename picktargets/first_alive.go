package picktargets

import (
	"slices"

	"github.com/rozag/cabasi/atk"
	"github.com/rozag/cabasi/battle"
	"github.com/rozag/cabasi/creat"
)

var _ battle.PickTargets = FirstAlive

// FirstAlive is a function that picks the first available not-out-of-battle
// defenders for the selected attack.
// It receives a picked attack and a slice of defenders.
// It returns a slice of picked defender indexes.
// It returns nil if there are no defenders to attack.
func FirstAlive(attack atk.Attack, defenders []creat.Creature) []uint {
	if len(defenders) == 0 {
		return nil
	}

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
			// Suppressing gosec G115 "integer overflow conversion int -> uint"
			// because int index will never overflow a uint variable.
			indexes = append(indexes, uint(idx)) //nolint:gosec
		}
	}

	if len(indexes) == 0 {
		return nil
	}

	return slices.Clip(indexes)
}
