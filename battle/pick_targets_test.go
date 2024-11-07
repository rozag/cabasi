package battle

import (
	"slices"
	"testing"

	"github.com/rozag/cabasi/dice"
)

func TestPickedTargetsEquals(t *testing.T) {
	tests := []struct {
		name string
		pt1  PickedTargets
		pt2  PickedTargets
		want bool
	}{
		{
			name: "EqualPickedTargets",
			pt1: PickedTargets{
				AttackerID:  "player-0",
				AttackIdx:   1,
				DefenderIDs: []CreatureID{"monster-0", "monster-1"},
			},
			pt2: PickedTargets{
				AttackerID:  "player-0",
				AttackIdx:   1,
				DefenderIDs: []CreatureID{"monster-0", "monster-1"},
			},
			want: true,
		},
		{
			name: "DifferentAttackerID",
			pt1: PickedTargets{
				AttackerID:  "player-0",
				AttackIdx:   1,
				DefenderIDs: []CreatureID{"monster-0", "monster-1"},
			},
			pt2: PickedTargets{
				AttackerID:  "player-1",
				AttackIdx:   1,
				DefenderIDs: []CreatureID{"monster-0", "monster-1"},
			},
			want: false,
		},
		{
			name: "DifferentAttackIdx",
			pt1: PickedTargets{
				AttackerID:  "player-0",
				AttackIdx:   1,
				DefenderIDs: []CreatureID{"monster-0", "monster-1"},
			},
			pt2: PickedTargets{
				AttackerID:  "player-0",
				AttackIdx:   2,
				DefenderIDs: []CreatureID{"monster-0", "monster-1"},
			},
			want: false,
		},
		{
			name: "DifferentDefenderIDs",
			pt1: PickedTargets{
				AttackerID:  "player-0",
				AttackIdx:   1,
				DefenderIDs: []CreatureID{"monster-0", "monster-1"},
			},
			pt2: PickedTargets{
				AttackerID:  "player-0",
				AttackIdx:   1,
				DefenderIDs: []CreatureID{"monster-2", "monster-3"},
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.pt1.Equals(test.pt2); got != test.want {
				t.Errorf("Equals() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestPickTargetsFirstAlive(t *testing.T) {
	spear := Attack{
		Name: "Spear", TargetCharacteristic: STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	tests := []struct {
		name                 string
		attackers, defenders []Creature
		want                 []PickedTargets
	}{
		{
			name:      "NilAttackers",
			attackers: nil,
			defenders: []Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: nil,
		},
		{
			name:      "EmptyAttackers",
			attackers: []Creature{},
			defenders: []Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: nil,
		},
		{
			name: "NilDefenders",
			attackers: []Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: nil,
			want:      nil,
		},
		{
			name: "EmptyDefenders",
			attackers: []Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []Creature{},
			want:      nil,
		},
		{
			name:      "NilAttackersNilDefenders",
			attackers: nil,
			defenders: nil,
			want:      nil,
		},
		{
			name:      "EmptyAttackersEmptyDefenders",
			attackers: []Creature{},
			defenders: []Creature{},
			want:      nil,
		},
		{
			name: "AllAttackersOut",
			attackers: []Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Doe", Attacks: []Attack{spear},
					STR: 8, DEX: 0, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-2", Name: "Joe Schmoe", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 0, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "monster-1", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: nil,
		},
		{
			name: "AllDefendersOut",
			attackers: []Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "monster-1", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 0, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "monster-2", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 0, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: nil,
		},
		{
			name: "AllAttackersOutAllDefendersOut",
			attackers: []Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Doe", Attacks: []Attack{spear},
					STR: 8, DEX: 0, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-2", Name: "Joe Schmoe", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 0, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "monster-1", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 0, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "monster-2", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 0, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: nil,
		},
		{
			name: "AttackersHaveNoAttacksLeft",
			attackers: []Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []Attack{
						{
							Name: "Sword", TargetCharacteristic: STR,
							Dice: dice.D6, DiceCnt: 1, Charges: 0,
							IsBlast: false,
						},
						{
							Name: "Paralyze", TargetCharacteristic: DEX,
							Dice: dice.D4, DiceCnt: 1, Charges: 0,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Doe",
					Attacks: []Attack{
						{
							Name: "Axe", TargetCharacteristic: STR,
							Dice: dice.D6, DiceCnt: 1, Charges: 0,
							IsBlast: false,
						},
						{
							Name: "Delirium", TargetCharacteristic: WIL,
							Dice: dice.D4, DiceCnt: 1, Charges: 0,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: nil,
		},
		{
			name: "PickFirstDefender",
			attackers: []Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Doe", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "monster-1", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: []PickedTargets{
				{
					AttackerID: "player-0", AttackIdx: 0,
					DefenderIDs: []CreatureID{"monster-0"},
				},
				{
					AttackerID: "player-1", AttackIdx: 0,
					DefenderIDs: []CreatureID{"monster-0"},
				},
			},
		},
		{
			name: "SkipFirstOutOfBattleDefenderAndPickSecond",
			attackers: []Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Doe", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "monster-1", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: []PickedTargets{
				{
					AttackerID: "player-0", AttackIdx: 0,
					DefenderIDs: []CreatureID{"monster-1"},
				},
				{
					AttackerID: "player-1", AttackIdx: 0,
					DefenderIDs: []CreatureID{"monster-1"},
				},
			},
		},
		{
			name: "PickSeveralDefendersForBlastAttack",
			attackers: []Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []Attack{
						{
							Name: "Fireball", TargetCharacteristic: STR,
							Dice: dice.D8, DiceCnt: 1, Charges: 1,
							IsBlast: true,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "monster-1", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: []PickedTargets{
				{
					AttackerID: "player-0", AttackIdx: 0,
					DefenderIDs: []CreatureID{"monster-0", "monster-1"},
				},
			},
		},
		{
			name: "PickAttackWithCharges",
			attackers: []Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []Attack{
						{
							Name: "Fireball", TargetCharacteristic: STR,
							Dice: dice.D8, DiceCnt: 1, Charges: 0,
							IsBlast: true,
						},
						spear,
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: []PickedTargets{
				{
					AttackerID: "player-0", AttackIdx: 1,
					DefenderIDs: []CreatureID{"monster-0"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := PickTargetsFirstAlive(test.attackers, test.defenders)
			slices.SortFunc(got, func(a, b PickedTargets) int {
				return a.AttackerID.CompareTo(b.AttackerID)
			})
			if !slices.EqualFunc(got, test.want, PickedTargets.Equals) {
				t.Fatalf("PickTargetsFirstAlive() = %v, want %v", got, test.want)
			}
		})
	}
}
