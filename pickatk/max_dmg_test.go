package pickatk

import (
	"testing"

	"github.com/rozag/cabasi/atk"
	"github.com/rozag/cabasi/creat"
	"github.com/rozag/cabasi/dice"
)

func TestMaxDmg(t *testing.T) {
	spear := atk.Attack{
		Name: "Spear", TargetCharacteristic: atk.STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	player := creat.Creature{
		ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
		STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
		IsDetachment: false,
	}
	tests := []struct {
		name      string
		defenders []creat.Creature
		attacker  creat.Creature
		want      int
	}{
		{
			name:      "NilDefenders",
			attacker:  player,
			defenders: nil,
			want:      -1,
		},
		{
			name:      "EmptyDefenders",
			attacker:  player,
			defenders: []creat.Creature{},
			want:      -1,
		},
		{
			name: "AttackerIsOut",
			attacker: creat.Creature{
				ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
				STR: 0, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: -1,
		},
		{
			name: "AttackerHasNoAttacks",
			attacker: creat.Creature{
				ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: -1,
		},
		{
			name:     "AllDefendersOut",
			attacker: player,
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
			},
			want: -1,
		},
		{
			name: "NoAttacksLeft",
			attacker: creat.Creature{
				ID: "player-0", Name: "John Appleseed",
				Attacks: []atk.Attack{
					{
						Name: "Fire Bolt", TargetCharacteristic: atk.STR,
						Dice: dice.D6, DiceCnt: 1, Charges: 0,
						IsBlast: false,
					},
				},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: -1,
		},
		{
			name: "PickSingleAttack",
			attacker: creat.Creature{
				ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
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
			want: 0,
		},
		{
			name: "PickTheMostPowerfulAttack",
			attacker: creat.Creature{
				ID: "player-0", Name: "John Appleseed",
				Attacks: []atk.Attack{
					{
						Name: "Spear", TargetCharacteristic: atk.STR,
						Dice: dice.D6, DiceCnt: 1, Charges: -1,
						IsBlast: false,
					},
					{
						Name: "Longsword", TargetCharacteristic: atk.STR,
						Dice: dice.D8, DiceCnt: 1, Charges: -1,
						IsBlast: false,
					},
				},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
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
			want: 1,
		},
		{
			name: "PickBlastAttackForMultipleDefenders",
			attacker: creat.Creature{
				ID: "player-0", Name: "John Appleseed",
				Attacks: []atk.Attack{
					{
						Name: "Spear", TargetCharacteristic: atk.STR,
						Dice: dice.D10, DiceCnt: 1, Charges: -1,
						IsBlast: false,
					},
					{
						Name: "Fireball", TargetCharacteristic: atk.STR,
						Dice: dice.D6, DiceCnt: 1, Charges: 1,
						IsBlast: true,
					},
					{
						Name: "Longsword", TargetCharacteristic: atk.STR,
						Dice: dice.D12, DiceCnt: 1, Charges: -1,
						IsBlast: false,
					},
				},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
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
				{
					ID: "monster-2", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: 1,
		},
		{
			name: "PickMorePowerfulAttackInsteadOfBlastForSingleDefender",
			attacker: creat.Creature{
				ID: "player-0", Name: "John Appleseed",
				Attacks: []atk.Attack{
					{
						Name: "Spear", TargetCharacteristic: atk.STR,
						Dice: dice.D10, DiceCnt: 1, Charges: -1,
						IsBlast: false,
					},
					{
						Name: "Fireball", TargetCharacteristic: atk.STR,
						Dice: dice.D6, DiceCnt: 1, Charges: 1,
						IsBlast: true,
					},
					{
						Name: "Longsword", TargetCharacteristic: atk.STR,
						Dice: dice.D12, DiceCnt: 1, Charges: -1,
						IsBlast: false,
					},
				},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: 2,
		},
		{
			name: "DoNotPickPowerfulAttackWithNoCharges",
			attacker: creat.Creature{
				ID: "player-0", Name: "John Appleseed",
				Attacks: []atk.Attack{
					{
						Name: "Magic Spear", TargetCharacteristic: atk.STR,
						Dice: dice.D12, DiceCnt: 1, Charges: 0,
						IsBlast: false,
					},
					{
						Name: "Knife", TargetCharacteristic: atk.STR,
						Dice: dice.D6, DiceCnt: 1, Charges: -1,
						IsBlast: false,
					},
				},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: 1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := MaxDmg(test.attacker, test.defenders); got != test.want {
				t.Errorf("MaxDmg() = %v, want %v", got, test.want)
			}
		})
	}
}
