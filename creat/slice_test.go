package creat

import (
	"testing"

	"github.com/rozag/cabasi/atk"
	"github.com/rozag/cabasi/dice"
)

func TestCreatureSliceEquals(t *testing.T) {
	spear := atk.Attack{
		Name: "Spear", TargetCharacteristic: atk.STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	hero := Creature{
		ID: "player-0", Name: "Jane Appleseed", Attacks: []atk.Attack{spear},
		STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
		IsDetachment: false,
	}
	monster := Creature{
		ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
		STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
		IsDetachment: false,
	}
	tests := []struct {
		name        string
		this, other CreatureSlice
		want        bool
	}{
		{
			name: "EqualNil",
			this: nil, other: nil,
			want: true,
		},
		{
			name: "EqualEmpty",
			this: CreatureSlice{}, other: CreatureSlice{},
			want: true,
		},
		{
			name: "NilNotEqualToEmpty",
			this: nil, other: CreatureSlice{},
			want: false,
		},
		{
			name: "EmptyNotEqualToNil",
			this: CreatureSlice{}, other: nil,
			want: false,
		},
		{
			name: "EqualNormal",
			this: CreatureSlice{hero, monster}, other: CreatureSlice{hero, monster},
			want: true,
		},
		{
			name: "DifferentLengths",
			this: CreatureSlice{hero}, other: CreatureSlice{hero, monster},
			want: false,
		},
		{
			name: "DifferentCreatures",
			this: CreatureSlice{hero}, other: CreatureSlice{monster},
			want: false,
		},
		{
			name: "OtherNil",
			this: CreatureSlice{hero, monster}, other: nil,
			want: false,
		},
		{
			name: "OtherEmpty",
			this: CreatureSlice{hero, monster}, other: CreatureSlice{},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.this.Equals(test.other); got != test.want {
				t.Errorf("CreatureSlice.Equals() = %v; want %v", got, test.want)
			}
		})
	}
}
