package picktargets

import (
	"slices"
	"testing"

	"github.com/rozag/cabasi/atk"
	"github.com/rozag/cabasi/creat"
	"github.com/rozag/cabasi/dice"
)

func TestFirstAlive(t *testing.T) {
	spear := atk.Attack{
		Name: "Spear", TargetCharacteristic: atk.STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	tests := []struct {
		name      string
		attack    atk.Attack
		defenders []creat.Creature
		want      []uint
	}{
		{
			name:      "NilDefenders",
			attack:    spear,
			defenders: nil,
			want:      nil,
		},
		{
			name:      "EmptyDefenders",
			attack:    spear,
			defenders: []creat.Creature{},
			want:      nil,
		},
		{
			name:   "AllDefendersOut",
			attack: spear,
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "monster-1", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 0, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "monster-2", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 0, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: nil,
		},
		{
			name: "AttackHasNoCharges",
			attack: atk.Attack{
				Name: "Sword", TargetCharacteristic: atk.STR,
				Dice: dice.D6, DiceCnt: 1, Charges: 0,
				IsBlast: false,
			},
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: nil,
		},
		{
			name:   "PickFirstDefender",
			attack: spear,
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "monster-1", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: []uint{0},
		},
		{
			name:   "SkipFirstOutOfBattleDefenderAndPickSecond",
			attack: spear,
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "monster-1", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: []uint{1},
		},
		{
			name: "PickSeveralDefendersForBlastAttack",
			attack: atk.Attack{
				Name: "Fireball", TargetCharacteristic: atk.STR,
				Dice: dice.D8, DiceCnt: 1, Charges: 1,
				IsBlast: true,
			},
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "monster-1", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: []uint{0, 1},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := FirstAlive(test.attack, test.defenders)
			if !slices.Equal(got, test.want) {
				t.Fatalf("PickTargetsFirstAlive() = %v, want %v", got, test.want)
			}
		})
	}
}
