package pickatk

import (
	"github.com/rozag/cabasi/battle"
	"github.com/rozag/cabasi/creat"
)

var _ battle.PickAttack = MaxDmg

// MaxDmg is a function that picks an attack that will deal the maximum damage
// to the defenders.
// It receives an attacker and a slice of defenders.
// It returns index of the attack the attacker will use.
// It returns -1 if the attacker does not attack.
func MaxDmg(attacker creat.Creature, defenders []creat.Creature) int {
	if attacker.IsOut() {
		return -1
	}

	attackableCnt := 0
	for _, defender := range defenders {
		if !defender.IsOut() {
			attackableCnt++
		}
	}
	if attackableCnt == 0 {
		return -1
	}

	maxDmg := uint(0)
	maxDmgIdx := -1
	for idx, attack := range attacker.Attacks {
		if attack.Charges == 0 {
			continue
		}

		dmg := uint(attack.Dice)
		if attack.IsBlast {
			dmg *= uint(len(defenders))
		}

		if dmg > maxDmg {
			maxDmg = dmg
			maxDmgIdx = idx
		}
	}

	return maxDmgIdx
}
