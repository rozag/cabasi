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
