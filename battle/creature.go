package battle

import (
	"errors"
	"fmt"
	"strings"
)

const (
	CharacteristicMin = 1
	CharacteristicMax = 20
)

const HPMin = 1

const ArmorMax = 3

// CreatureID is a unique identifier of a creature.
type CreatureID string

// CompareTo returns an integer comparing two ids lexicographically. The result
// will be 0 if this == other, -1 if this < other, and +1 if this > other.
func (this CreatureID) CompareTo(other CreatureID) int {
	if this < other {
		return -1
	}
	if this > other {
		return 1
	}
	return 0
}

// Creature represents a creature in a battle - a player or a monster.
type Creature struct {
	ID           CreatureID
	Name         string
	Attacks      []Attack
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
		AttackSlice(c.Attacks),
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
func (this *Creature) Equals(other *Creature) bool {
	return this.ID == other.ID &&
		this.Name == other.Name &&
		this.STR == other.STR &&
		this.DEX == other.DEX &&
		this.WIL == other.WIL &&
		this.HP == other.HP &&
		this.Armor == other.Armor &&
		this.IsDetachment == other.IsDetachment &&
		AttackSlice(this.Attacks).Equals(AttackSlice(other.Attacks))
}

// DeepCopy creates a deep copy of the Creature.
func (c *Creature) DeepCopy() Creature {
	attacks := make([]Attack, len(c.Attacks))
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

// CreatureSlice is a `[]Creature` with helper methods.
type CreatureSlice []Creature

// String returns the string representation of the CreatureSlice.
func (cs CreatureSlice) String() string {
	var sb strings.Builder
	sb.WriteString("[]Creature{")
	for i, c := range cs {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(c.String())
	}
	sb.WriteString("}")
	return sb.String()
}

// Equals checks if the CreatureSlice is equal to the other CreatureSlice.
func (this CreatureSlice) Equals(other CreatureSlice) bool {
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
