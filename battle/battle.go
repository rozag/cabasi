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

// Run simulates a battle between 2 groups of Creatures. It returns true if the
// players won, false otherwise.
//
// The RNG is used for all the rolls.
//
// The Log is used for logging the battle events.
//
// Run returns an error if input is invalid in any way. The error has an
// `Unwrap() []error` method to get all the errors or `nil` if the inputs are
// valid.
//
// Run doesn't modify the input creatures.
func Run(rng dice.RNG, log Log, players, monsters []Creature) (bool, error) {
	var errs []error

	if rng == nil {
		errs = append(errs, errors.New("RNG must be provided"))
	}

	if log == nil {
		errs = append(errs, errors.New("Log must be provided"))
	}

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

	ids := make(map[string]struct{})
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

	playersCopy := make([]Creature, len(players))
	for i, player := range players {
		copied := player.DeepCopy()
		playersCopy[i] = copied
	}

	monstersCopy := make([]Creature, len(monsters))
	for i, monster := range monsters {
		copied := monster.DeepCopy()
		monstersCopy[i] = copied
	}

	havePlayersWon := run(rng, log, playersCopy, monstersCopy)
	return havePlayersWon, nil
}

// TODO: battle.New(rng, log)? More deps/settings will come later

func run(rng dice.RNG, log Log, players, monsters []Creature) bool {
	for {
		if true { // TODO: remove
			return true
		}

		// playerTargets := pickTargets(players, monsters)
		//   true - can flee
		// monsters := attackTargets(rng, log, players, monsters, playerTargets, true)
		// if len(monsters) == 0 {
		// 	return true
		// }
		//
		// monsterTargets := pickTargets(monsters, players)
		//   false - cannot flee
		// players := attackTargets(rng, log, monsters, players, monsterTargets, false)
		// if len(players) == 0 {
		// 	return false
		// }

		// TODO:
		// each round while there are players and monsters alive (but monsters flee)
		// - all players pick targets
		// - all players attack, same target -> the highest hit
		// - all monsters pick targets
		// - all monsters attack, same target -> the highest hit

		// TODO: remove dead or fleeing creatures from the slices on attack resolve
	}
}
