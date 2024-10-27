package battle

import (
	"errors"
	"fmt"

	"github.com/rozag/cabasi/dice"
)

type Characteristic uint8

const (
	STR Characteristic = iota
	DEX
	WIL
)

// Attack represents a single attack on a characteristic of a creature - a blunt
// attack, a special ability, a spell, etc.
type Attack struct {
	Name                 string
	TargetCharacteristic Characteristic
	Dice                 dice.Dice
	DiceCnt              uint8
	Charges              int8 // <0 means infinite
	IsBlast              bool
}

// Validate checks if the freshly created attack is valid. It returns an error
// with `Unwrap() []error` method to get all the errors or `nil` if the attack
// is valid.
func (a *Attack) Validate() error {
	var errs []error

	switch a.TargetCharacteristic {
	case STR, DEX, WIL:
		// OK
	default:
		errs = append(
			errs,
			fmt.Errorf("invalid target characteristic: %d", a.TargetCharacteristic),
		)
	}

	switch a.Dice {
	case dice.D4, dice.D6, dice.D8, dice.D10, dice.D12, dice.D20:
		// OK
	default:
		errs = append(errs, fmt.Errorf("invalid dice: %d", a.Dice))
	}

	if a.DiceCnt == 0 {
		errs = append(errs, errors.New("dice count must be at least 1"))
	}

	return errors.Join(errs...)
}