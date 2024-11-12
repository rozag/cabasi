package battle

import (
	"errors"
	"fmt"

	"github.com/rozag/cabasi/atk"
	"github.com/rozag/cabasi/creat"
	"github.com/rozag/cabasi/dice"
)

// Log is an interface for logging the battle events.
type Log interface {
	// Roll logs a roll of a dice for a creature, a save of some sort.
	Roll(c creat.Creature, roll uint8)

	// Attack logs an attack of an attacker on a defender with a specific attack
	// and damage dealt.
	Attack(
		attacker creat.Creature,
		defenders []creat.Creature,
		attack atk.Attack,
		damage uint,
	)
}

// PickAttack is a function that picks which attack the attacker will use.
// It receives an attacker and a slice of defenders.
// It returns index of the attack the attacker will use.
// It returns -1 if the attacker does not attack.
type PickAttack func(attacker creat.Creature, defenders []creat.Creature) int

// PickTargets is a function that picks targets for the selected attack.
// It receives a picked attack and a slice of defenders.
// It returns a slice of picked defender indexes.
// It returns nil if there are no defenders to attack.
type PickTargets func(attack atk.Attack, defenders []creat.Creature) []uint

// Battle represents a battle between 2 parties.
type Battle struct {
	rng         dice.RNG
	log         Log
	pickAttack  PickAttack
	pickTargets PickTargets
}

// New creates a new Battle with the provided RNG and Log.
// The RNG is used for all the rolls.
// The Log is used for logging the battle events.
// The PickAttack is a function that picks which attack the attacker will use.
// The PickTargets is a function that picks targets for an attack.
// New returns an error if input is invalid in any way. The error has an
// `Unwrap() []error` method to get all the errors or `nil` if the inputs are
// valid.
func New(
	rng dice.RNG,
	log Log,
	pickAttack PickAttack,
	pickTargets PickTargets,
) (*Battle, error) {
	var errs []error

	if rng == nil {
		errs = append(errs, errors.New("RNG must be provided"))
	}

	if log == nil {
		errs = append(errs, errors.New("Log must be provided"))
	}

	if pickAttack == nil {
		errs = append(errs, errors.New("PickAttack must be provided"))
	}

	if pickTargets == nil {
		errs = append(errs, errors.New("PickTargets must be provided"))
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	battle := Battle{rng, log, pickAttack, pickTargets}
	return &battle, nil
}

// Run simulates a battle between 2 groups of Creatures. It returns true if the
// players won, false otherwise.
//
// Run returns an error if input is invalid in any way. The error has an
// `Unwrap() []error` method to get all the errors or `nil` if the inputs are
// valid.
//
// Run doesn't modify the input creatures.
func (b *Battle) Run(players, monsters []creat.Creature) (bool, error) {
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

	ids := make(map[creat.ID]struct{})
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

	playersCopy := make([]creat.Creature, len(players))
	for i, player := range players {
		copied := player.DeepCopy()
		playersCopy[i] = copied
	}

	monstersCopy := make([]creat.Creature, len(monsters))
	for i, monster := range monsters {
		copied := monster.DeepCopy()
		monstersCopy[i] = copied
	}

	havePlayersWon := b.run(playersCopy, monstersCopy)
	return havePlayersWon, nil
}

func (b *Battle) run(players, monsters []creat.Creature) bool {
	playerAtkIdxs := make([]int, len(players))
	playerTargets := make([][]uint, len(players))
	playerAttackers := make([][]attacker, len(players))

	monsterAtkIdxs := make([]int, len(monsters))
	monsterTargets := make([][]uint, len(monsters))
	monsterAttackers := make([][]attacker, len(monsters))

	for {
		_ = players  // TODO: remove
		_ = monsters // TODO: remove
		if true {    // TODO: remove
			return true
		}

		// TODO:
		// • Enemies must pass a WIL save to avoid fleeing when they take their
		//   first casualty and again when they lose half their number.
		// • Some groups may use their leader’s WIL in place of their own. Lone foes
		//   must save when they’re reduced to 0 HP.

		// TODO:
		// attackTargets is several steps:
		// 1. map targets to all attackers and their attacks
		// 2. resolve attacks
		// 3. handle critical damage and fleeing (as reducing STR to 0)

		// TODO:
		// each round while there are players and monsters alive (but monsters flee)
		// - all players pick targets
		// - all players attack, same target -> the highest hit
		// - all monsters pick targets
		// - all monsters attack, same target -> the highest hit

		// TODO: dead, incapacitated, or fleeing creatures have either their STR set
		// to 0 on attack resolve or have DEX or WIL as 0 because of some effect

		for i, player := range players {
			atkIdx := b.pickAttack(player, monsters)
			playerAtkIdxs[i] = atkIdx
			if atkIdx < 0 {
				playerTargets[i] = nil
			} else {
				playerTargets[i] = b.pickTargets(player.Attacks[atkIdx], monsters)
			}
		}

		assignAttackers(monsterAttackers, playerTargets, playerAtkIdxs)

		// TODO:

		for i, monster := range monsters {
			atkIdx := b.pickAttack(monster, players)
			monsterAtkIdxs[i] = atkIdx
			if atkIdx < 0 {
				monsterTargets[i] = nil
			} else {
				monsterTargets[i] = b.pickTargets(monster.Attacks[atkIdx], players)
			}
		}

		assignAttackers(playerAttackers, monsterTargets, monsterAtkIdxs)
	}
}

type attacker struct {
	attackerIdx uint
	attackIdx   uint
}

// assignAttackers assigns attackers to targets.
// It receives attackers, targets, and attack indexes.
// It modifies attackers in place and doesn't return anything.
// attackers is a slice of size of defenders, each element is a slice of
// attackers that target the defender with a particular attack.
// targets is a slice of size of attackers, each element is a slice of defender
// indexes that are targeted by the attacker.
// attackIdxs is a slice of size of attackers, each element is an index of the
// attack the attacker will use.
func assignAttackers(
	attackers [][]attacker,
	targets [][]uint,
	attackIdxs []int,
) {
	// TODO:
}
