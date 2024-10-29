package battle

import (
	"errors"
	"fmt"
	"slices"
)

const (
	CharacteristicMin = 1
	CharacteristicMax = 20
)

const HPMin = 1

const ArmorMax = 3

// Creature represents a creature in a battle - a player or a monster.
type Creature struct {
	ID           string
	Name         string
	Attacks      []Attack
	STR          uint8
	DEX          uint8
	WIL          uint8
	HP           uint8
	Armor        uint8
	IsDetachment bool
}

// Validate checks if the freshly created creature is valid. It returns an error
// with `Unwrap() []error` method to get all the errors or `nil` if the creature
// is valid. Validate is not meant to be used on a Creature in the middle of
// a battle (with decreased characteristics), but rather on a freshly created
// one.
func (c Creature) Validate() error {
	var errs []error

	if len(c.ID) == 0 {
		errs = append(errs, errors.New("creature must have an ID"))
	}

	if len(c.Name) == 0 {
		errs = append(errs, errors.New("creature must have a name"))
	}

	if len(c.Attacks) == 0 {
		errs = append(errs, errors.New("creature must have at least one attack"))
	}
	for idx, attack := range c.Attacks {
		if err := attack.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("invalid attack at idx %d: %w", idx, err))
		}
	}

	if c.STR < CharacteristicMin || c.STR > CharacteristicMax {
		errs = append(errs, fmt.Errorf(
			"STR must be between %d and %d, got %d",
			CharacteristicMin, CharacteristicMax, c.STR,
		))
	}

	if c.DEX < CharacteristicMin || c.DEX > CharacteristicMax {
		errs = append(errs, fmt.Errorf(
			"DEX must be between %d and %d, got %d",
			CharacteristicMin, CharacteristicMax, c.DEX,
		))
	}

	if c.WIL < CharacteristicMin || c.WIL > CharacteristicMax {
		errs = append(errs, fmt.Errorf(
			"WIL must be between %d and %d, got %d",
			CharacteristicMin, CharacteristicMax, c.WIL,
		))
	}

	if c.HP < HPMin {
		errs = append(errs, fmt.Errorf(
			"HP must be at least %d, got %d", HPMin, c.HP,
		))
	}

	if c.Armor > ArmorMax {
		errs = append(errs, fmt.Errorf(
			"Armor must be at most %d, got %d", ArmorMax, c.Armor,
		))
	}

	return errors.Join(errs...)
}

// Equals checks if the Creature is equal to the other Creature.
func (this Creature) Equals(other Creature) bool {
	return this.ID == other.ID &&
		this.Name == other.Name &&
		this.STR == other.STR &&
		this.DEX == other.DEX &&
		this.WIL == other.WIL &&
		this.HP == other.HP &&
		this.Armor == other.Armor &&
		this.IsDetachment == other.IsDetachment &&
		slices.EqualFunc(this.Attacks, other.Attacks, Attack.Equals)
}

// DeepCopy creates a deep copy of the Creature.
func (c Creature) DeepCopy() Creature {
	// TODO:
	return c
}
