package creat

import (
	"errors"
	"fmt"

	"github.com/rozag/cabasi/atk"
)

const (
	// CharacteristicMin is the minimum value of a newly created creature's
	// characteristic.
	CharacteristicMin = 1
	// CharacteristicMax is the maximum value of a newly created creature's
	// characteristic.
	CharacteristicMax = 20
)

// HPMin is the minimum value of a newly created creature's HP.
const HPMin = 1

// ArmorMax is the maximum value of a creature's armor.
const ArmorMax = 3

// Creature represents a creature in a battle - a player or a monster.
type Creature struct {
	ID           ID
	Name         string
	Attacks      []atk.Attack
	STR          uint8
	DEX          uint8
	WIL          uint8
	HP           uint8
	Armor        uint8
	IsDetachment bool
}

// IsOut checks if the Creature is out of the battle - if any of its core
// characteristics is zero.
func (c *Creature) IsOut() bool {
	return c.STR == 0 || c.DEX == 0 || c.WIL == 0
}

// String returns the string representation of the Creature.
func (c *Creature) String() string {
	return fmt.Sprintf(
		"Creature{"+
			"ID: %q"+
			", Name: %q"+
			", Attacks: %s"+
			", STR: %d"+
			", DEX: %d"+
			", WIL: %d"+
			", HP: %d"+
			", Armor: %d"+
			", IsDetachment: %t"+
			"}",
		c.ID,
		c.Name,
		atk.AttackSlice(c.Attacks),
		c.STR,
		c.DEX,
		c.WIL,
		c.HP,
		c.Armor,
		c.IsDetachment,
	)
}

// Validate checks if the freshly created creature is valid. It returns an error
// with `Unwrap() []error` method to get all the errors or `nil` if the creature
// is valid. Validate is not meant to be used on a Creature in the middle of
// a battle (with decreased characteristics), but rather on a freshly created
// one.
func (c *Creature) Validate() error {
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
func (c *Creature) Equals(other *Creature) bool {
	return c.ID == other.ID &&
		c.Name == other.Name &&
		c.STR == other.STR &&
		c.DEX == other.DEX &&
		c.WIL == other.WIL &&
		c.HP == other.HP &&
		c.Armor == other.Armor &&
		c.IsDetachment == other.IsDetachment &&
		atk.AttackSlice(c.Attacks).Equals(atk.AttackSlice(other.Attacks))
}

// DeepCopy creates a deep copy of the Creature.
func (c *Creature) DeepCopy() Creature {
	attacks := make([]atk.Attack, len(c.Attacks))
	for i, attack := range c.Attacks {
		copied := attack.DeepCopy()
		attacks[i] = copied
	}
	return Creature{
		ID:           c.ID,
		Name:         c.Name,
		Attacks:      attacks,
		STR:          c.STR,
		DEX:          c.DEX,
		WIL:          c.WIL,
		HP:           c.HP,
		Armor:        c.Armor,
		IsDetachment: c.IsDetachment,
	}
}
