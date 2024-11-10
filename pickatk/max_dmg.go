package pickatk

import (
	"github.com/rozag/cabasi/battle"
	"github.com/rozag/cabasi/creat"
)

var _ battle.PickAttack = MaxDmg

// MaxDmg is a function that picks an attack that will deal the maximum damage
// to the defenders. STR damage > DEX damage > WIL damage.
// It receives an attacker and a slice of defenders.
// It returns index of the attack the attacker will use.
// It returns -1 if the attacker does not attack.
func MaxDmg(attacker creat.Creature, defenders []creat.Creature) int {
	return -1
}
