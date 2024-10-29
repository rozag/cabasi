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

func TestCreatureEquals(t *testing.T) {
	spear := Attack{
		Name: "Spear", TargetCharacteristic: STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	tests := []struct {
		name        string
		this, other Creature
		want        bool
	}{
		{
			name: "EqualCreatures",
			this: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			other: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			want: true,
		},
		{
			name: "DifferentID",
			this: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			other: Creature{
				ID: "monster-1", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			want: false,
		},
		{
			name: "DifferentName",
			this: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			other: Creature{
				ID: "monster-0", Name: "Boot Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			want: false,
		},
		{
			name: "DifferentAttacks",
			this: Creature{
				ID: "monster-0", Name: "Root Goblin",
				Attacks: []Attack{
					{
						Name: "Spear", TargetCharacteristic: STR,
						Dice: dice.D6, DiceCnt: 1, Charges: -1,
						IsBlast: false,
					},
				},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			other: Creature{
				ID: "monster-0", Name: "Root Goblin",
				Attacks: []Attack{
					{
						Name: "Sword", TargetCharacteristic: STR,
						Dice: dice.D6, DiceCnt: 1, Charges: -1,
						IsBlast: false,
					},
				},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			want: false,
		},
		{
			name: "DifferentSTR",
			this: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			other: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 9, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			want: false,
		},
		{
			name: "DifferentDEX",
			this: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			other: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 15, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			want: false,
		},
		{
			name: "DifferentWIL",
			this: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			other: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 9, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			want: false,
		},
		{
			name: "DifferentHP",
			this: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			other: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 5, Armor: 0,
				IsDetachment: false,
			},
			want: false,
		},
		{
			name: "DifferentArmor",
			this: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			other: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 1,
				IsDetachment: false,
			},
			want: false,
		},
		{
			name: "DifferentIsDetachment",
			this: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: false,
			},
			other: Creature{
				ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
				STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
				IsDetachment: true,
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.this.Equals(&test.other); got != test.want {
				t.Fatalf("Creature.Equals() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestCreatureDeepCopy(t *testing.T) {
	spear := Attack{
		Name: "Spear", TargetCharacteristic: STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	original := Creature{
		ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
		STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
		IsDetachment: false,
	}
	copied := original.DeepCopy()

	if !original.Equals(copied) {
		t.Fatalf("Creature.DeepCopy() = %v, want %v", copied, original)
	}

	copied.ID = "monster-1"
	copied.Name = "Boot Goblin"
	copied.Attacks[0].Name = "Sword"
	copied.STR = 9
	copied.DEX = 15
	copied.WIL = 9
	copied.HP = 5
	copied.Armor = 1
	copied.IsDetachment = true

	if original.Equals(copied) {
		t.Errorf("modifying the copy affected the original: %v", original)
	}

	if original.ID == copied.ID {
		t.Errorf("original.ID == copied.ID")
	}
	if original.Name == copied.Name {
		t.Errorf("original.Name == copied.Name")
	}
	if attackSlice(original.Attacks).equal(attackSlice(copied.Attacks)) {
		t.Errorf("original.Attacks == copied.Attacks")
	}
	if original.STR == copied.STR {
		t.Errorf("original.STR == copied.STR")
	}
	if original.DEX == copied.DEX {
		t.Errorf("original.DEX == copied.DEX")
	}
	if original.WIL == copied.WIL {
		t.Errorf("original.WIL == copied.WIL")
	}
	if original.HP == copied.HP {
		t.Errorf("original.HP == copied.HP")
	}
	if original.Armor == copied.Armor {
		t.Errorf("original.Armor == copied.Armor")
	}
	if original.IsDetachment == copied.IsDetachment {
		t.Errorf("original.IsDetachment == copied.IsDetachment")
	}
}
