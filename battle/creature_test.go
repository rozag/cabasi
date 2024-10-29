package battle

import (
	"testing"

	"github.com/rozag/cabasi/dice"
)

func TestCreatureValidate(t *testing.T) {
	spear := Attack{
		Name: "Spear", TargetCharacteristic: STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	tests := []struct {
		name       string
		creature   Creature
		wantErrCnt uint
	}{
		{
			name: "ValidCreature",
			creature: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			wantErrCnt: 0,
		},
		{
			name: "EmptyID",
			creature: Creature{
				ID: "", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "EmptyName",
			creature: Creature{
				ID: "monster-0", Name: "", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "NilAttacks",
			creature: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: nil,
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "EmptyAttacks",
			creature: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "InvalidAttack",
			creature: Creature{
				ID: "monster-0", Name: "Root Goblin",
				Attacks: []Attack{
					{
						Name: "", TargetCharacteristic: STR,
						Dice: dice.D6, DiceCnt: 1, Charges: -1,
						IsBlast: false,
					},
				},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "STRTooLow",
			creature: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 0, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "STRTooHigh",
			creature: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 21, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "DEXTooLow",
			creature: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 0, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "DEXTooHigh",
			creature: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 21, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "WILTooLow",
			creature: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 0, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "WILTooHigh",
			creature: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 21, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "HPTooLow",
			creature: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 0, Armor: 0,
				IsDetachment: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "ArmorTooHigh",
			creature: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 4,
				IsDetachment: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "MultipleErrors",
			creature: Creature{
				ID: "", Name: "", Attacks: []Attack{},
				STR: 21, DEX: 0, WIL: 21, HP: 0, Armor: 21,
				IsDetachment: true,
			},
			wantErrCnt: 8,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.creature.Validate()

			if test.wantErrCnt == 0 {
				if err != nil {
					t.Fatalf("Creature.Validate(): want nil, got %v", err)
				} else {
					return
				}
			}

			if err == nil {
				t.Fatalf("Creature.Validate(): want error, got nil")
			}

			jointErr, ok := err.(interface{ Unwrap() []error })
			if !ok {
				t.Fatalf(
					"Creature.Validate(): error must have `Unwrap() []error` method",
				)
			}

			errs := jointErr.Unwrap()
			if uint(len(errs)) != test.wantErrCnt {
				t.Fatalf(
					"Creature.Validate(): want %d errors, got %d",
					test.wantErrCnt, len(errs),
				)
			}
		})
	}
}
