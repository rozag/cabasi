package battle

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rozag/cabasi/dice"
)

type Characteristic uint8

const (
	STR Characteristic = iota
	DEX
	WIL
)

// String returns the string representation of the Characteristic.
func (c Characteristic) String() string {
	switch c {
	case STR:
		return "STR"
	case DEX:
		return "DEX"
	case WIL:
		return "WIL"
	default:
		panic(fmt.Errorf("unknown Characteristic: %d", c))
	}
}

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

// String returns the string representation of the Attack.
func (a *Attack) String() string {
	return fmt.Sprintf(
		"Attack{"+
			"Name: %q"+
			", TargetCharacteristic: %s"+
			", Dice: %s"+
			", DiceCnt: %d"+
			", Charges: %d"+
			", IsBlast: %t"+
			"}",
		a.Name,
		a.TargetCharacteristic,
		a.Dice,
		a.DiceCnt,
		a.Charges,
		a.IsBlast,
	)
}

// Validate checks if the freshly created attack is valid. It returns an error
// with `Unwrap() []error` method to get all the errors or `nil` if the attack
// is valid.
func (a *Attack) Validate() error {
	var errs []error

	if len(a.Name) == 0 {
		errs = append(errs, errors.New("attack must have a name"))
	}

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

// Equals checks if the Attack is equal to the other Attack.
func (this *Attack) Equals(other *Attack) bool {
	return this.Name == other.Name &&
		this.TargetCharacteristic == other.TargetCharacteristic &&
		this.Dice == other.Dice &&
		this.DiceCnt == other.DiceCnt &&
		this.Charges == other.Charges &&
		this.IsBlast == other.IsBlast
}

// DeepCopy creates a deep copy of the Attack.
func (a *Attack) DeepCopy() *Attack {
	return &Attack{
		Name:                 a.Name,
		TargetCharacteristic: a.TargetCharacteristic,
		Dice:                 a.Dice,
		DiceCnt:              a.DiceCnt,
		Charges:              a.Charges,
		IsBlast:              a.IsBlast,
	}
}

// AttackSlice is a `[]Attack` with helper methods.
type AttackSlice []Attack

// String returns the string representation of the AttackSlice.
func (a AttackSlice) String() string {
	var sb strings.Builder
	sb.WriteString("[]Attack{")
	for i, attack := range a {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(attack.String())
	}
	sb.WriteString("}")
	return sb.String()
}

// Equals checks if the AttackSlice is equal to the other AttackSlice.
func (this AttackSlice) Equals(other AttackSlice) bool {
	if len(this) != len(other) {
		return false
	}

	for i := range this {
		if !this[i].Equals(&other[i]) {
			return false
		}
	}

	return true
}
