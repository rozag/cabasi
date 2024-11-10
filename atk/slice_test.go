package atk

import (
	"testing"

	"github.com/rozag/cabasi/dice"
)

func TestAttackSliceEquals(t *testing.T) {
	knife := Attack{
		Name: "Knife", TargetCharacteristic: STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	spear := Attack{
		Name: "Spear", TargetCharacteristic: STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	tests := []struct {
		name        string
		this, other AttackSlice
		want        bool
	}{
		{
			name: "EqualNil",
			this: nil, other: nil,
			want: true,
		},
		{
			name: "EqualEmpty",
			this: AttackSlice{}, other: AttackSlice{},
			want: true,
		},
		{
			name: "NilNotEqualToEmpty",
			this: nil, other: AttackSlice{},
			want: false,
		},
		{
			name: "EmptyNotEqualToNil",
			this: AttackSlice{}, other: nil,
			want: false,
		},
		{
			name: "EqualNormal",
			this: AttackSlice{knife, spear}, other: AttackSlice{knife, spear},
			want: true,
		},
		{
			name: "DifferentLengths",
			this: AttackSlice{knife}, other: AttackSlice{knife, spear},
			want: false,
		},
		{
			name: "DifferentAttacks",
			this: AttackSlice{knife}, other: AttackSlice{spear},
			want: false,
		},
		{
			name: "OtherNil",
			this: AttackSlice{knife, spear}, other: nil,
			want: false,
		},
		{
			name: "OtherEmpty",
			this: AttackSlice{knife, spear}, other: AttackSlice{},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.this.Equals(test.other); got != test.want {
				t.Errorf("AttackSlice.Equals() = %v; want %v", got, test.want)
			}
		})
	}
}
