package battle

import (
	"testing"

	"github.com/rozag/cabasi/dice"
)

func TestAttackValidate(t *testing.T) {
	tests := []struct {
		name       string
		attack     Attack
		wantErrCnt uint
	}{
		{
			name: "ValidAttack",
			attack: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			wantErrCnt: 0,
		},
		{
			name: "EmptyName",
			attack: Attack{
				Name: "", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "UnknownTargetCharacteristic",
			attack: Attack{
				Name: "Knife", TargetCharacteristic: Characteristic(42),
				Dice: dice.D6, DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "UnknownDice",
			attack: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.Dice(42), DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "InvalidDiceCnt",
			attack: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 0, Charges: -1,
				IsBlast: false,
			},
			wantErrCnt: 1,
		},
		{
			name: "MultipleErrors",
			attack: Attack{
				Name: "", TargetCharacteristic: Characteristic(42),
				Dice: dice.Dice(42), DiceCnt: 0, Charges: -1,
				IsBlast: true,
			},
			wantErrCnt: 4,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.attack.Validate()

			if test.wantErrCnt == 0 {
				if err != nil {
					t.Fatalf("Attack.Validate(): want nil, got %v", err)
				} else {
					return
				}
			}

			if err == nil {
				t.Fatalf("Attack.Validate(): want error, got nil")
			}

			jointErr, ok := err.(interface{ Unwrap() []error })
			if !ok {
				t.Fatalf("Attack.Validate(): error must have `Unwrap() []error` method")
			}

			errs := jointErr.Unwrap()
			if uint(len(errs)) != test.wantErrCnt {
				t.Fatalf(
					"Attack.Validate(): want %d errors, got %d",
					test.wantErrCnt, len(errs),
				)
			}
		})
	}
}

func TestAttackEquals(t *testing.T) {
	tests := []struct {
		name        string
		this, other Attack
		want        bool
	}{
		{
			name: "EqualAttacks",
			this: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			other: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			want: true,
		},
		{
			name: "DifferentName",
			this: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			other: Attack{
				Name: "Sword", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			want: false,
		},
		{
			name: "DifferentTargetCharacteristic",
			this: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			other: Attack{
				Name: "Knife", TargetCharacteristic: DEX,
				Dice: dice.D6, DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			want: false,
		},
		{
			name: "DifferentDice",
			this: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			other: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.D8, DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			want: false,
		},
		{
			name: "DifferentDiceCnt",
			this: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			other: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 2, Charges: -1,
				IsBlast: false,
			},
			want: false,
		},
		{
			name: "DifferentCharges",
			this: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			other: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 1, Charges: 1,
				IsBlast: false,
			},
			want: false,
		},
		{
			name: "DifferentIsBlast",
			this: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 1, Charges: -1,
				IsBlast: false,
			},
			other: Attack{
				Name: "Knife", TargetCharacteristic: STR,
				Dice: dice.D6, DiceCnt: 1, Charges: -1,
				IsBlast: true,
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.this.Equals(&test.other); got != test.want {
				t.Fatalf("Attack.Equals() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestAttackDeepCopy(t *testing.T) {
	original := Attack{
		Name: "Knife", TargetCharacteristic: STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	copied := original.DeepCopy()

	if !original.Equals(&copied) {
		t.Fatalf("Attack.DeepCopy() = %v, want %v", copied, original)
	}

	copied.Name = "Sword"
	copied.TargetCharacteristic = DEX
	copied.Dice = dice.D8
	copied.DiceCnt = 2
	copied.Charges = 1
	copied.IsBlast = true

	if original.Equals(&copied) {
		t.Errorf("modifying the copy affected the original: %v", original)
	}

	if original.Name == copied.Name {
		t.Errorf("original.Name == copied.Name")
	}
	if original.TargetCharacteristic == copied.TargetCharacteristic {
		t.Errorf("original.TargetCharacteristic == copied.TargetCharacteristic")
	}
	if original.Dice == copied.Dice {
		t.Errorf("original.Dice == copied.Dice")
	}
	if original.DiceCnt == copied.DiceCnt {
		t.Errorf("original.DiceCnt == copied.DiceCnt")
	}
	if original.Charges == copied.Charges {
		t.Errorf("original.Charges == copied.Charges")
	}
	if original.IsBlast == copied.IsBlast {
		t.Errorf("original.IsBlast == copied.IsBlast")
	}
}
