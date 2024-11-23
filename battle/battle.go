package battle

import (
	"errors"
	"fmt"
	"math"
	"slices"

	"github.com/rozag/cabasi/atk"
	"github.com/rozag/cabasi/creat"
	"github.com/rozag/cabasi/dice"
)

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
	pickAttack  PickAttack
	pickTargets PickTargets
}

// New creates a new Battle with the provided RNG and strategies.
// The RNG is used for all the rolls.
// The PickAttack is a function that picks which attack the attacker will use.
// The PickTargets is a function that picks targets for an attack.
// New returns an error if input is invalid in any way. The error has an
// `Unwrap() []error` method to get all the errors or `nil` if the inputs are
// valid.
func New(
	rng dice.RNG,
	pickAttack PickAttack,
	pickTargets PickTargets,
) (*Battle, error) {
	var errs []error

	if rng == nil {
		errs = append(errs, errors.New("RNG must be provided"))
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

	battle := Battle{rng, pickAttack, pickTargets}
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

// run simulates a battle between 2 groups of Creatures. It returns true if the
// players win, false otherwise.
func (b *Battle) run(players, monsters []creat.Creature) bool {
	playerAtkIdxs := make([]int, len(players))
	playerTargets := make([][]uint, len(players))
	playerAttackers := make([][]attacker, len(players))
	playerUsedAttackIdxs := make([]int, len(players))
	damageToPlayers := make([]damage, len(players))

	monsterAtkIdxs := make([]int, len(monsters))
	monsterTargets := make([][]uint, len(monsters))
	monsterAttackers := make([][]attacker, len(monsters))
	monsterUsedAttackIdxs := make([]int, len(monsters))
	damageToMonsters := make([]damage, len(monsters))

	for {
		b.pickAttacksAndTargets(players, monsters, playerAtkIdxs, playerTargets)

		assignAttackers(monsterAttackers, playerTargets, playerAtkIdxs)
		if noAttackersAssigned(monsterAttackers) {
			// players cannot attack anyone, hence they lose
			return false
		}

		resolveAttacks(
			damageToMonsters, players, monsters, monsterAttackers,
			playerUsedAttackIdxs, b.rng,
		)
		if noDamageDone(damageToMonsters) {
			// players cannot deal any damage, hence they lose
			return false
		}

		applyDamageToMonsters(monsters, damageToMonsters, b.rng)
		if allOut(monsters) {
			// monsters are all out, hence players win
			return true
		}

		b.pickAttacksAndTargets(monsters, players, monsterAtkIdxs, monsterTargets)

		assignAttackers(playerAttackers, monsterTargets, monsterAtkIdxs)
		if noAttackersAssigned(monsterAttackers) {
			// monsters cannot attack anyone, hence players win
			return true
		}

		resolveAttacks(
			damageToPlayers, monsters, players, playerAttackers,
			monsterUsedAttackIdxs, b.rng,
		)
		if noDamageDone(damageToPlayers) {
			// monsters cannot deal any damage, hence players win
			return true
		}

		applyDamageToPlayers(players, damageToPlayers, b.rng)
		if allOut(players) {
			// players are all out, hence monsters win
			return false
		}
	}
}

func (b *Battle) pickAttacksAndTargets(
	attackers, defenders []creat.Creature,
	attackIndexes []int,
	targets [][]uint,
) {
	for i, attacker := range attackers {
		attackIdx := b.pickAttack(attacker, defenders)
		attackIndexes[i] = attackIdx
		if attackIdx < 0 {
			targets[i] = nil
		} else {
			targets[i] = b.pickTargets(attacker.Attacks[attackIdx], defenders)
		}
	}
}

type attacker struct{ attackerIdx, attackIdx uint }

// assignAttackers assigns attackers to targets.
// It receives attackers, targets, and attack indexes. It modifies attackers in
// place.
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
	if len(attackers) == 0 {
		return
	}

	// TODO:
	// - Attacks against detachments by individuals are Impaired (excluding Blast
	// damage).
	// - Attacks against individuals by detachments are Enhanced and deal Blast
	// damage.
	// NOTE: only this last part about detachments to individuals attacks is
	// relevant here

	for i := range attackers {
		attackers[i] = nil
	}

	if len(targets) == 0 ||
		len(attackIdxs) == 0 ||
		len(targets) != len(attackIdxs) {
		return
	}

	for attackerIdx, defenderIdxs := range targets {
		if len(defenderIdxs) == 0 {
			continue
		}

		attackIdx := attackIdxs[attackerIdx]
		if attackIdx < 0 {
			continue
		}

		for _, defenderIdx := range defenderIdxs {
			if defenderIdx >= uint(len(attackers)) {
				continue
			}

			attackers[defenderIdx] = append(
				attackers[defenderIdx],
				attacker{
					// Suppressing gosec "G115 integer overflow conversion int -> uint"
					// because int index will never overflow a uint variable.
					attackerIdx: uint(attackerIdx), //nolint:gosec
					attackIdx:   uint(attackIdx),
				},
			)
		}
	}

	for i := range attackers {
		attackers[i] = slices.Clip(attackers[i])
	}
}

// noAttackersAssigned returns true if no attackers are assigned to any target.
func noAttackersAssigned(attackers [][]attacker) bool {
	for _, assigned := range attackers {
		if len(assigned) > 0 {
			return false
		}
	}
	return true
}

type damage struct {
	characteristic atk.Characteristic
	value          uint8
}

// resolveAttacks computes the damage dealt to the defenders by the attackers
// (armor is taken into account) and decreases attacks' charges if they're not
// unlimited (-1).
// It receives damageToDefenders, attackers, defenders, assignedAttackers, and
// RNG. It modifies damageToDefenders in place.
// damageToDefenders is a slice of damage dealt to each defender.
// attackers is a slice of all attackers.
// defenders is a slice of all defenders.
// assignedAttackers is a slice of size of defenders, each element is a slice of
// attackers that target the defender with a particular attack.
// usedAttackIdxs is a slice of size of attackers, it doesn't matter what's
// inside because it's cleared before being used and acts as a reusable buffer.
// RNG is used for all the rolls.
func resolveAttacks(
	damageToDefenders []damage,
	attackers, defenders []creat.Creature,
	assignedAttackers [][]attacker,
	usedAttackIdxs []int,
	rng dice.RNG,
) {
	if len(damageToDefenders) == 0 {
		return
	}

	for i := range damageToDefenders {
		damageToDefenders[i].characteristic = atk.STR
		damageToDefenders[i].value = 0
	}

	if len(attackers) == 0 ||
		len(defenders) == 0 ||
		len(assignedAttackers) == 0 ||
		len(usedAttackIdxs) == 0 ||
		len(defenders) != len(damageToDefenders) ||
		len(defenders) != len(assignedAttackers) ||
		len(attackers) != len(usedAttackIdxs) ||
		rng == nil {
		return
	}

	for defenderIdx := range damageToDefenders {
		if defenders[defenderIdx].IsOut() ||
			len(assignedAttackers[defenderIdx]) == 0 {
			continue
		}

		maxDamageCharacteristic := atk.STR
		maxDamageValue := uint8(0)
		for _, assigned := range assignedAttackers[defenderIdx] {
			attackerIdx := assigned.attackerIdx
			if attackerIdx >= uint(len(attackers)) {
				continue
			}

			attacker := attackers[attackerIdx]
			if attacker.IsOut() {
				continue
			}

			attackIdx := assigned.attackIdx
			if attackIdx >= uint(len(attacker.Attacks)) {
				continue
			}

			attack := attacker.Attacks[attackIdx]
			if attack.Charges == 0 {
				continue
			}

			attackDice := attack.Dice
			isAttackerDetachment := attacker.IsDetachment
			isDefenderDetachment := defenders[defenderIdx].IsDetachment
			if isAttackerDetachment != isDefenderDetachment {
				if isAttackerDetachment && !isDefenderDetachment {
					attackDice = dice.D12
				} else if !attack.IsBlast {
					attackDice = dice.D4
				}
			}

			maxDmg := uint8(0)
			for range attack.DiceCnt {
				dmg := attackDice.Roll(rng)
				if dmg > maxDmg {
					maxDmg = dmg
				}
			}

			if maxDmg > maxDamageValue {
				maxDamageCharacteristic = attack.TargetCharacteristic
				maxDamageValue = maxDmg
			}
		}

		if maxDamageValue > 0 {
			if maxDamageCharacteristic == atk.STR &&
				defenders[defenderIdx].Armor > 0 {
				if maxDamageValue >= defenders[defenderIdx].Armor {
					maxDamageValue -= defenders[defenderIdx].Armor
				}
			}
			damageToDefenders[defenderIdx].characteristic = maxDamageCharacteristic
			damageToDefenders[defenderIdx].value = maxDamageValue
		}
	}

	for attackerIdx := range usedAttackIdxs {
		usedAttackIdxs[attackerIdx] = -1
	}

	for _, allAssigned := range assignedAttackers {
		for _, assigned := range allAssigned {
			attackerIdx := assigned.attackerIdx
			if attackerIdx >= uint(len(attackers)) {
				continue
			}

			attackIdx := assigned.attackIdx
			if attackIdx >= uint(len(attackers[attackerIdx].Attacks)) {
				continue
			}

			if attackIdx > uint(math.MaxInt) {
				continue
			}

			usedAttackIdxs[attackerIdx] = int(attackIdx)
		}
	}

	for attackerIdx, usedAttackIdx := range usedAttackIdxs {
		if usedAttackIdx < 0 {
			continue
		}

		if attackers[attackerIdx].Attacks[usedAttackIdx].Charges > 0 {
			attackers[attackerIdx].Attacks[usedAttackIdx].Charges--
		}
	}
}

// noDamageDone returns true if no damage is done after resolving the attacks.
func noDamageDone(damage []damage) bool {
	for _, dmg := range damage {
		if dmg.value > 0 {
			return false
		}
	}
	return true
}

// applyDamageToPlayers decreases player's characteristics according to damage
// received (armor is NOT taken into account) and handles critical damage (as
// reducing STR to 0).
// It receives players, damageToPlayers, and RNG. It modifies players in place.
// players is a slice of all players.
// damageToPlayers is a slice of damage dealt to each player.
// RNG is used for all the rolls.
func applyDamageToPlayers(
	players []creat.Creature,
	damageToPlayers []damage,
	rng dice.RNG,
) {
	if len(players) == 0 ||
		len(damageToPlayers) == 0 ||
		len(players) != len(damageToPlayers) ||
		rng == nil {
		return
	}

	for playerIdx := range players {
		if players[playerIdx].IsOut() {
			continue
		}

		value := damageToPlayers[playerIdx].value
		if value == 0 {
			continue
		}

		switch c := damageToPlayers[playerIdx].characteristic; c {
		case atk.STR:
			if value <= players[playerIdx].HP {
				players[playerIdx].HP -= value
				continue
			}

			value -= players[playerIdx].HP
			players[playerIdx].HP = 0

			if value >= players[playerIdx].STR {
				players[playerIdx].STR = 0
				continue
			}

			players[playerIdx].STR -= value
			if dice.D20.Roll(rng) > players[playerIdx].STR {
				players[playerIdx].STR = 0
			}

		case atk.DEX:
			players[playerIdx].DEX -= value

		case atk.WIL:
			players[playerIdx].WIL -= value

		default:
			panic(fmt.Errorf("unknown Characteristic: %d", c))
		}
	}
}

// applyDamageToMonsters decreases monster's characteristics according to damage
// received (armor is NOT taken into account) and handles fleeing (as reducing
// STR to 0).
// It receives monsters, damageToMonsters, and RNG. It modifies monsters in
// place.
// monsters is a slice of all monsters.
// damageToMonsters is a slice of damage dealt to each monster.
// RNG is used for all the rolls.
func applyDamageToMonsters(
	monsters []creat.Creature,
	damageToMonsters []damage,
	rng dice.RNG,
) {
	// TODO:
	// • Enemies must pass a WIL save to avoid fleeing when they take their
	//   first casualty and again when they lose half their number.
	// • Some groups may use their leader’s WIL in place of their own.
	// • Lone foes must save when they’re reduced to 0 HP.

	// TODO: handle fleeing (as reducing STR to 0)

	// TODO: dead, incapacitated, or fleeing creatures have either their STR set
	// to 0 on attack resolve or have DEX or WIL as 0 because of some effect
}

// allOut returns true if all creatures are out.
func allOut(creatures []creat.Creature) bool {
	for _, c := range creatures {
		if !c.IsOut() {
			return false
		}
	}
	return true
}
