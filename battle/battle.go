package battle

import (
	"errors"
	"fmt"

	"github.com/rozag/cabasi/dice"
)

// Log is an interface for logging the battle events.
type Log interface {
	// Roll logs a roll of a dice for a creature, a save of some sort.
	Roll(c Creature, roll uint8)

	// Attack logs an attack of an attacker on a defender with a specific attack
	// and damage dealt.
	Attack(attacker, defender Creature, attack Attack, damage uint8)
}

type Battle struct {
	rng         dice.RNG
	log         Log
	pickTargets PickTargets
}

// New creates a new Battle with the provided RNG and Log.
//
// The RNG is used for all the rolls.
//
// The Log is used for logging the battle events.
//
// The PickTargets is a function that picks the targets for the attackers.
//
// New returns an error if input is invalid in any way. The error has an
// `Unwrap() []error` method to get all the errors or `nil` if the inputs are
// valid.
func New(rng dice.RNG, log Log, pickTargets PickTargets) (*Battle, error) {
	var errs []error

	if rng == nil {
		errs = append(errs, errors.New("RNG must be provided"))
	}

	if log == nil {
		errs = append(errs, errors.New("Log must be provided"))
	}

	if pickTargets == nil {
		errs = append(errs, errors.New("PickTargets must be provided"))
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &Battle{rng, log, pickTargets}, nil
}

// Run simulates a battle between 2 groups of Creatures. It returns true if the
// players won, false otherwise.
//
// Run returns an error if input is invalid in any way. The error has an
// `Unwrap() []error` method to get all the errors or `nil` if the inputs are
// valid.
//
// Run doesn't modify the input creatures.
func (b *Battle) Run(players, monsters []Creature) (bool, error) {
	var errs []error

	if len(players) == 0 {
		errs = append(errs, errors.New("at least one player must be provided"))
	}
	for idx, player := range players {
		if err := player.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("invalid player at idx %d: %w", idx, err))
		}
	}

	if len(monsters) == 0 {
		errs = append(errs, errors.New("at least one monster must be provided"))
	}
	for idx, monster := range monsters {
		if err := monster.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("invalid monster at idx %d: %w", idx, err))
		}
	}

	ids := make(map[CreatureID]struct{})
	for idx, player := range players {
		if _, ok := ids[player.ID]; ok {
			errs = append(
				errs,
				fmt.Errorf("player at idx %d has non-unique ID %q", idx, player.ID),
			)
		} else {
			ids[player.ID] = struct{}{}
		}
	}
	for idx, monster := range monsters {
		if _, ok := ids[monster.ID]; ok {
			errs = append(
				errs,
				fmt.Errorf("monster at idx %d has non-unique ID %q", idx, monster.ID),
			)
		} else {
			ids[monster.ID] = struct{}{}
		}
	}

	if len(errs) > 0 {
		return false, errors.Join(errs...)
	}

	idToPlayers := make(map[CreatureID]Creature, len(players))
	for i := range players {
		idToPlayers[players[i].ID] = players[i].DeepCopy()
	}

	idToMonsters := make(map[CreatureID]Creature, len(monsters))
	for i := range monsters {
		idToMonsters[monsters[i].ID] = monsters[i].DeepCopy()
	}

	havePlayersWon := b.run(idToPlayers, idToMonsters)
	return havePlayersWon, nil
}

func (b *Battle) run(idToPlayers, idToMonsters map[CreatureID]Creature) bool {
	for {
		if true { // TODO: remove
			return true
		}

		// NOTE:
		// • Enemies must pass a WIL save to avoid fleeing when they take their
		//   first casualty and again when they lose half their number.
		// • Some groups may use their leader’s WIL in place of their own. Lone foes
		//   must save when they’re reduced to 0 HP.

		// TODO:
		// attackTargets is several steps:
		// 1. map targets to all attackers and their attacks
		// 2. resolve attacks
		// 3. handle fleeing as reducing STR to 0

		// playerTargets := b.pickTargets(players, monsters)
		//   true - can flee
		// attackTargets(rng, log, players, monsters, playerTargets, true)
		// if isAnyoneAlive(monsters) {
		// 	return true
		// }
		//
		// monsterTargets := b.pickTargets(monsters, players)
		//   false - cannot flee
		// attackTargets(rng, log, monsters, players, monsterTargets, false)
		// if isAnyoneAlive(players) {
		// 	return false
		// }

		// TODO:
		// each round while there are players and monsters alive (but monsters flee)
		// - all players pick targets
		// - all players attack, same target -> the highest hit
		// - all monsters pick targets
		// - all monsters attack, same target -> the highest hit

		// TODO: dead, incapacitated, or fleeing creatures have either their STR set
		// to 0 on attack resolve or have DEX or WIL as 0 because of some effect
	}
}
