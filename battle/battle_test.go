package battle

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/rozag/cabasi/atk"
	"github.com/rozag/cabasi/creat"
	"github.com/rozag/cabasi/dice"
	"github.com/rozag/cabasi/pickatk"
	"github.com/rozag/cabasi/picktargets"
)

type minRNG struct{}

func (minRNG) UintN(uint) uint { return 0 }

type maxRNG struct{}

func (maxRNG) UintN(n uint) uint { return n - 1 }

type sequenceRNG struct {
	seq []uint
	idx uint
}

func (s *sequenceRNG) UintN(n uint) uint {
	if n == 0 {
		panic(errors.New("sequenceRNG: n must be greater than 0"))
	}
	if len(s.seq) == 0 {
		panic(errors.New("sequenceRNG: no values"))
	}
	if s.idx >= uint(len(s.seq)) {
		panic(fmt.Errorf(
			"sequenceRNG: no more values: idx=%d, len(seq)=%d", s.idx, len(s.seq),
		))
	}
	value := s.seq[s.idx]
	if value >= n {
		panic(
			fmt.Errorf("sequenceRNG: value must be from [0, n=%d), got %d", n, value),
		)
	}
	s.idx++
	return value
}

func newDeterministicRNG() *rand.Rand {
	// Suppressing gosec "G404 Use of weak random number generator" in tests.
	return rand.New(rand.NewPCG(0x436169726E, 0x525047)) //nolint:gosec
}

func dummyPickAttack(creat.Creature, []creat.Creature) int { return -1 }

func dummyPickTargets(atk.Attack, []creat.Creature) []uint { return nil }

func TestNewValidation(t *testing.T) {
	rng := minRNG{}
	tests := []struct {
		rng         dice.RNG
		pickAttack  PickAttack
		pickTargets PickTargets
		name        string
		wantErrCnt  uint
	}{
		{
			name:       "ValidNew",
			rng:        rng,
			pickAttack: dummyPickAttack, pickTargets: dummyPickTargets,
			wantErrCnt: 0,
		},
		{
			name:       "NoRNG",
			rng:        nil,
			pickAttack: dummyPickAttack, pickTargets: dummyPickTargets,
			wantErrCnt: 1,
		},
		{
			name:       "NoPickAttack",
			rng:        rng,
			pickAttack: nil, pickTargets: dummyPickTargets,
			wantErrCnt: 1,
		},
		{
			name:       "NoPickTargets",
			rng:        rng,
			pickAttack: dummyPickAttack, pickTargets: nil,
			wantErrCnt: 1,
		},
		{
			name:       "MultipleErrors",
			rng:        nil,
			pickAttack: nil, pickTargets: nil,
			wantErrCnt: 3,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(test.rng, test.pickAttack, test.pickTargets)

			if test.wantErrCnt == 0 {
				if err != nil {
					t.Fatalf("New(): want nil error, got %v", err)
				} else {
					return
				}
			}

			if err == nil {
				t.Fatalf("New(): want error, got nil")
			}

			jointErr, ok := err.(interface{ Unwrap() []error })
			if !ok {
				t.Fatalf("New(): error must have `Unwrap() []error` method")
			}

			errs := jointErr.Unwrap()
			if uint(len(errs)) != test.wantErrCnt {
				t.Fatalf("New(): want %d errors, got %d", test.wantErrCnt, len(errs))
			}
		})
	}
}

func TestRunValidation(t *testing.T) {
	spear := atk.Attack{
		Name: "Spear", TargetCharacteristic: atk.STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	player := creat.Creature{
		ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
		STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
		IsDetachment: false,
	}
	monster := creat.Creature{
		ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
		STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
		IsDetachment: false,
	}
	tests := []struct {
		name              string
		players, monsters []creat.Creature
		wantErrCnt        uint
	}{
		{
			name:    "ValidRun",
			players: []creat.Creature{player}, monsters: []creat.Creature{monster},
			wantErrCnt: 0,
		},
		{
			name:    "NilPlayers",
			players: nil, monsters: []creat.Creature{monster},
			wantErrCnt: 1,
		},
		{
			name:    "EmptyPlayers",
			players: []creat.Creature{}, monsters: []creat.Creature{monster},
			wantErrCnt: 1,
		},
		{
			name:    "NilMonsters",
			players: []creat.Creature{player}, monsters: nil,
			wantErrCnt: 1,
		},
		{
			name:    "EmptyMonsters",
			players: []creat.Creature{player}, monsters: []creat.Creature{},
			wantErrCnt: 1,
		},
		{
			name: "InvalidPlayer",
			players: []creat.Creature{
				{
					ID: "", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			monsters:   []creat.Creature{monster},
			wantErrCnt: 1,
		},
		{
			name:    "InvalidMonster",
			players: []creat.Creature{player},
			monsters: []creat.Creature{
				{
					ID: "", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			wantErrCnt: 1,
		},
		{
			name: "NonUniqueCreatureID",
			players: []creat.Creature{
				{
					ID: "creature", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			monsters: []creat.Creature{
				{
					ID: "creature", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			wantErrCnt: 1,
		},
		{
			name:    "MultipleErrors",
			players: nil, monsters: nil,
			wantErrCnt: 2,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := New(minRNG{}, dummyPickAttack, dummyPickTargets)
			if err != nil {
				t.Fatalf("New(): want nil error, got %v", err)
			}

			_, err = b.Run(test.players, test.monsters)

			if test.wantErrCnt == 0 {
				if err != nil {
					t.Fatalf("Run(): want nil error, got %v", err)
				} else {
					return
				}
			}

			if err == nil {
				t.Fatalf("Run(): want error, got nil")
			}

			jointErr, ok := err.(interface{ Unwrap() []error })
			if !ok {
				t.Fatalf("Run(): error must have `Unwrap() []error` method")
			}

			errs := jointErr.Unwrap()
			if uint(len(errs)) != test.wantErrCnt {
				t.Fatalf("Run(): want %d errors, got %d", test.wantErrCnt, len(errs))
			}
		})
	}
}

func TestRunDoesNotMutateCreatures(t *testing.T) {
	b, err := New(newDeterministicRNG(), pickatk.MaxDmg, picktargets.FirstAlive)
	if err != nil {
		t.Fatalf("New(): want nil error, got %v", err)
	}

	spear := atk.Attack{
		Name: "Spear", TargetCharacteristic: atk.STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	originalPlayers := []creat.Creature{
		{
			ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
			STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
			IsDetachment: false,
		},
	}
	originalMonsters := []creat.Creature{
		{
			ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
			STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
			IsDetachment: false,
		},
	}

	players := make([]creat.Creature, len(originalPlayers))
	for i, player := range originalPlayers {
		copied := player.DeepCopy()
		players[i] = copied
	}

	monsters := make([]creat.Creature, len(originalMonsters))
	for i, monster := range originalMonsters {
		copied := monster.DeepCopy()
		monsters[i] = copied
	}

	_, err = b.Run(players, monsters)
	if err != nil {
		t.Fatalf("Run(): want nil error, got %v", err)
	}

	if !creat.CreatureSlice(originalPlayers).Equals(players) {
		t.Fatalf(
			"Run(): players mutated, want %v, got %v", originalPlayers, players,
		)
	}
	if !creat.CreatureSlice(originalMonsters).Equals(monsters) {
		t.Fatalf(
			"Run(): monsters mutated, want %v, got %v", originalMonsters, monsters,
		)
	}
}

func TestAssignAttackers(t *testing.T) {
	tests := []struct {
		name       string
		attackers  [][]attacker
		targets    [][]uint
		attackIdxs []int
		want       [][]attacker
	}{
		{
			name:       "EmptyAttackers",
			attackers:  [][]attacker{},
			targets:    [][]uint{{0, 1}, {1, 2}},
			attackIdxs: []int{0, 1},
			want:       [][]attacker{},
		},
		{
			name:       "NilAttackers",
			attackers:  nil,
			targets:    [][]uint{{0}, {1, 2}},
			attackIdxs: []int{0, 1},
			want:       nil,
		},
		{
			name:       "EmptyTargets",
			attackers:  [][]attacker{nil, nil, nil},
			targets:    [][]uint{},
			attackIdxs: []int{0, 1},
			want:       [][]attacker{nil, nil, nil},
		},
		{
			name:       "NilTargets",
			attackers:  [][]attacker{nil, nil, nil},
			targets:    nil,
			attackIdxs: []int{0, 1},
			want:       [][]attacker{nil, nil, nil},
		},
		{
			name:       "AllTargetsEmpty",
			attackers:  [][]attacker{nil, nil, nil},
			targets:    [][]uint{{}, {}},
			attackIdxs: []int{0, 1},
			want:       [][]attacker{nil, nil, nil},
		},
		{
			name:       "AllTargetsNil",
			attackers:  [][]attacker{nil, nil, nil},
			targets:    [][]uint{nil, nil},
			attackIdxs: []int{0, 1},
			want:       [][]attacker{nil, nil, nil},
		},
		{
			name:       "InvalidIdxsInTargets",
			attackers:  [][]attacker{nil, nil, nil},
			targets:    [][]uint{{2, 3}, {4}},
			attackIdxs: []int{1, 0},
			want:       [][]attacker{nil, nil, {{attackerIdx: 0, attackIdx: 1}}},
		},
		{
			name:       "EmptyAttackIdxs",
			attackers:  [][]attacker{nil, nil, nil},
			targets:    [][]uint{{0, 1}, {1, 2}},
			attackIdxs: []int{},
			want:       [][]attacker{nil, nil, nil},
		},
		{
			name:       "NilAttackIdxs",
			attackers:  [][]attacker{nil, nil, nil},
			targets:    [][]uint{{0, 1}, {1, 2}},
			attackIdxs: nil,
			want:       [][]attacker{nil, nil, nil},
		},
		{
			name:       "NegativeAttackIdxs",
			attackers:  [][]attacker{nil, nil, nil},
			targets:    [][]uint{{0, 1}, {1, 2}},
			attackIdxs: []int{1, -1},
			want: [][]attacker{
				{{attackerIdx: 0, attackIdx: 1}},
				{{attackerIdx: 0, attackIdx: 1}},
				nil,
			},
		},
		{
			name: "DirtyAttackersReset",
			attackers: [][]attacker{
				{{attackerIdx: 0, attackIdx: 0}},
				{{attackerIdx: 1, attackIdx: 0}},
				{{attackerIdx: 2, attackIdx: 0}},
			},
			targets:    [][]uint{nil, {1}, nil},
			attackIdxs: []int{-1, 0, -1},
			want:       [][]attacker{nil, {{attackerIdx: 1, attackIdx: 0}}, nil},
		},
		{
			name: "DirtyAttackersResetWithInvalidInputs",
			attackers: [][]attacker{
				{{attackerIdx: 0, attackIdx: 0}},
				{{attackerIdx: 1, attackIdx: 0}},
				{{attackerIdx: 2, attackIdx: 0}},
			},
			targets:    nil,
			attackIdxs: nil,
			want:       [][]attacker{nil, nil, nil},
		},
		{
			name:       "TargetsAndAttackIdxsOfDifferentLength",
			attackers:  [][]attacker{nil, nil, nil},
			targets:    [][]uint{{0, 1}, {1, 2}},
			attackIdxs: []int{0},
			want:       [][]attacker{nil, nil, nil},
		},
		{
			name:       "AllAttackersAttack",
			attackers:  [][]attacker{nil, nil, nil},
			targets:    [][]uint{{0}, {1}, {2}, {0}},
			attackIdxs: []int{0, 1, 2, 3},
			want: [][]attacker{
				{{attackerIdx: 0, attackIdx: 0}, {attackerIdx: 3, attackIdx: 3}},
				{{attackerIdx: 1, attackIdx: 1}},
				{{attackerIdx: 2, attackIdx: 2}},
			},
		},
		{
			name:       "SomeAttackersAttack",
			attackers:  [][]attacker{nil, nil, nil},
			targets:    [][]uint{{0}, nil, nil, {0}},
			attackIdxs: []int{0, -1, -1, 3},
			want: [][]attacker{
				{{attackerIdx: 0, attackIdx: 0}, {attackerIdx: 3, attackIdx: 3}},
				nil,
				nil,
			},
		},
		{
			name:       "SomeAttackersAttackMultipleTargets",
			attackers:  [][]attacker{nil, nil, nil},
			targets:    [][]uint{{0, 1}, {2}, {1, 2}, {0}},
			attackIdxs: []int{0, 1, 2, 3},
			want: [][]attacker{
				{{attackerIdx: 0, attackIdx: 0}, {attackerIdx: 3, attackIdx: 3}},
				{{attackerIdx: 0, attackIdx: 0}, {attackerIdx: 2, attackIdx: 2}},
				{{attackerIdx: 1, attackIdx: 1}, {attackerIdx: 2, attackIdx: 2}},
			},
		},
		{
			name:       "AllAttackersAttackMultipleTargets",
			attackers:  [][]attacker{nil, nil, nil},
			targets:    [][]uint{{0, 1, 2}, {0, 1, 2}, {0, 1, 2}, {0, 1, 2}},
			attackIdxs: []int{0, 1, 2, 3},
			want: [][]attacker{
				{
					{attackerIdx: 0, attackIdx: 0},
					{attackerIdx: 1, attackIdx: 1},
					{attackerIdx: 2, attackIdx: 2},
					{attackerIdx: 3, attackIdx: 3},
				},
				{
					{attackerIdx: 0, attackIdx: 0},
					{attackerIdx: 1, attackIdx: 1},
					{attackerIdx: 2, attackIdx: 2},
					{attackerIdx: 3, attackIdx: 3},
				},
				{
					{attackerIdx: 0, attackIdx: 0},
					{attackerIdx: 1, attackIdx: 1},
					{attackerIdx: 2, attackIdx: 2},
					{attackerIdx: 3, attackIdx: 3},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assignAttackers(test.attackers, test.targets, test.attackIdxs)
			if !slices.EqualFunc(test.attackers, test.want, slices.Equal) {
				t.Fatalf(
					"assignAttackers(): want %v, got %v", test.want, test.attackers,
				)
			}
		})
	}
}

func TestNoAttackersAssigned(t *testing.T) {
	tests := []struct {
		name      string
		attackers [][]attacker
		want      bool
	}{
		{
			name:      "EmptyAttackers",
			attackers: [][]attacker{},
			want:      true,
		},
		{
			name:      "NilAttackers",
			attackers: nil,
			want:      true,
		},
		{
			name:      "AllAttackersEmpty",
			attackers: [][]attacker{{}, {}, {}},
			want:      true,
		},
		{
			name:      "AllAttackersNil",
			attackers: [][]attacker{nil, nil, nil},
			want:      true,
		},
		{
			name:      "SomeAttackersAssigned",
			attackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}, nil, {}},
			want:      false,
		},
		{
			name: "AllAttackersAssigned",
			attackers: [][]attacker{
				{{attackerIdx: 0, attackIdx: 0}, {attackerIdx: 1, attackIdx: 1}},
				{{attackerIdx: 2, attackIdx: 2}},
				{{attackerIdx: 3, attackIdx: 3}},
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			initial := make([][]attacker, len(test.attackers))
			for i, assigned := range test.attackers {
				initial[i] = make([]attacker, len(assigned))
				copy(initial[i], assigned)
			}

			if got := noAttackersAssigned(test.attackers); got != test.want {
				t.Fatalf("noAttackersAssigned(): want %t, got %t", test.want, got)
			}

			if !slices.EqualFunc(test.attackers, initial, slices.Equal) {
				t.Fatalf(
					"noAttackersAssigned(): attackers mutated, want %v, got %v",
					initial, test.attackers,
				)
			}
		})
	}
}

func TestResolveAttacks(t *testing.T) {
	spear := atk.Attack{
		Name: "Spear", TargetCharacteristic: atk.STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	player0 := creat.Creature{
		ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
		STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
		IsDetachment: false,
	}
	monster0 := creat.Creature{
		ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
		STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
		IsDetachment: false,
	}
	tests := []struct {
		name                 string
		damageToDefenders    []damage
		attackers, defenders []creat.Creature
		assignedAttackers    [][]attacker
		usedAttackIdxs       []int
		rng                  dice.RNG
		wantDamage           []damage
		wantAttackers        []creat.Creature
	}{
		{
			name:              "EmptyDamageToDefenders",
			damageToDefenders: []damage{},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        []damage{},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "NilDamageToDefenders",
			damageToDefenders: nil,
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        nil,
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "EmptyAttackers",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         []creat.Creature{},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     []creat.Creature{},
		},
		{
			name:              "NilAttackers",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         nil,
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     nil,
		},
		{
			name:              "EmptyDefenders",
			damageToDefenders: []damage{},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        []damage{},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "NilDefenders",
			damageToDefenders: []damage{},
			attackers:         []creat.Creature{player0},
			defenders:         nil,
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        []damage{},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "EmptyAssignedAttackers",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "NilAssignedAttackers",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: nil,
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "AllAssignedAttackersEmpty",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{}},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "AllAssignedAttackersNil",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{nil},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "EmptyUsedAttackIndexes",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    []int{},
			rng:               maxRNG{},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "NilUsedAttackIndexes",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    nil,
			rng:               maxRNG{},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "NilRNG",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    []int{42},
			rng:               nil,
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name: "DamageToDefendersAndDefendersOfDifferentLength",
			damageToDefenders: []damage{
				{characteristic: atk.STR, value: 0},
				{characteristic: atk.STR, value: 0},
			},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage: []damage{
				{characteristic: atk.STR, value: 0},
				{characteristic: atk.STR, value: 0},
			},
			wantAttackers: []creat.Creature{player0},
		},
		{
			name:              "DefendersAndAssignedAttackersOfDifferentLength",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{
				{{attackerIdx: 0, attackIdx: 0}},
				{{attackerIdx: 1, attackIdx: 0}},
			},
			usedAttackIdxs: []int{42},
			rng:            maxRNG{},
			wantDamage:     []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:  []creat.Creature{player0},
		},
		{
			name:              "AttackersAndUsedAttackIndexesOfDifferentLength",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    []int{-123456, 123455},
			rng:               maxRNG{},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "AssignedAttackersOut",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 0, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 0, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{
				{attackerIdx: 0, attackIdx: 0},
				{attackerIdx: 1, attackIdx: 0},
			}},
			usedAttackIdxs: []int{42, -10},
			rng:            maxRNG{},
			wantDamage:     []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 0, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 0, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
		},
		{
			name:              "TargetedDefendersOut",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         []creat.Creature{player0},
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 0, Armor: 0,
					IsDetachment: false,
				},
			},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "InvalidAttackerIndexes",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{{attackerIdx: 1, attackIdx: 0}}},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "InvalidAttackIndexes",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 1}}},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "AttackWithNoCharges",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Fireball", TargetCharacteristic: atk.STR,
							Dice: dice.D6, DiceCnt: 1, Charges: 0,
							IsBlast: true,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Fireball", TargetCharacteristic: atk.STR,
							Dice: dice.D6, DiceCnt: 1, Charges: 0,
							IsBlast: true,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
		},
		{
			name:              "DirtyDamageToDefendersReset",
			damageToDefenders: []damage{{characteristic: atk.DEX, value: 6}},
			attackers:         []creat.Creature{player0},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{},
			usedAttackIdxs:    []int{42},
			rng:               maxRNG{},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "AttackCannotPenetrateArmor",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers:         []creat.Creature{player0},
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 3,
					IsDetachment: false,
				},
			},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    []int{42},
			rng:               &sequenceRNG{seq: []uint{2}, idx: 0},
			wantDamage:        []damage{{characteristic: atk.STR, value: 0}},
			wantAttackers:     []creat.Creature{player0},
		},
		{
			name:              "SeveralAttacksToDifferentCharacteristics",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Spear", TargetCharacteristic: atk.STR,
							Dice: dice.D6, DiceCnt: 1, Charges: -1,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Delirium", TargetCharacteristic: atk.WIL,
							Dice: dice.D8, DiceCnt: 1, Charges: 1,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{
				{attackerIdx: 0, attackIdx: 0},
				{attackerIdx: 1, attackIdx: 0},
			}},
			usedAttackIdxs: []int{42, -10},
			rng:            maxRNG{},
			wantDamage:     []damage{{characteristic: atk.WIL, value: 8}},
			wantAttackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Spear", TargetCharacteristic: atk.STR,
							Dice: dice.D6, DiceCnt: 1, Charges: -1,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Delirium", TargetCharacteristic: atk.WIL,
							Dice: dice.D8, DiceCnt: 1, Charges: 0,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
		},
		{
			name:              "SingleAttackWithMultipleDice",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Fireball", TargetCharacteristic: atk.STR,
							Dice: dice.D8, DiceCnt: 2, Charges: 2,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders:         []creat.Creature{monster0},
			assignedAttackers: [][]attacker{{{attackerIdx: 0, attackIdx: 0}}},
			usedAttackIdxs:    []int{-10},
			rng:               &sequenceRNG{seq: []uint{3, 6}, idx: 0},
			wantDamage:        []damage{{characteristic: atk.STR, value: 7}},
			wantAttackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Fireball", TargetCharacteristic: atk.STR,
							Dice: dice.D8, DiceCnt: 2, Charges: 1,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
		},
		{
			name:              "MultipleAttacksWithMultipleDice",
			damageToDefenders: []damage{{characteristic: atk.STR, value: 0}},
			attackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Fireball", TargetCharacteristic: atk.STR,
							Dice: dice.D8, DiceCnt: 2, Charges: 1,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Sword", TargetCharacteristic: atk.STR,
							Dice: dice.D6, DiceCnt: 2, Charges: -1,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 1,
					IsDetachment: false,
				},
			},
			assignedAttackers: [][]attacker{{
				{attackerIdx: 0, attackIdx: 0},
				{attackerIdx: 1, attackIdx: 0},
			}},
			usedAttackIdxs: []int{-10, 42},
			rng:            &sequenceRNG{seq: []uint{0, 4, 3, 5}, idx: 0},
			wantDamage:     []damage{{characteristic: atk.STR, value: 5}},
			wantAttackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Fireball", TargetCharacteristic: atk.STR,
							Dice: dice.D8, DiceCnt: 2, Charges: 0,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Sword", TargetCharacteristic: atk.STR,
							Dice: dice.D6, DiceCnt: 2, Charges: -1,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
		},
		{
			name: "AllAttackersAttackDifferentTargets",
			damageToDefenders: []damage{
				{characteristic: atk.STR, value: 0},
				{characteristic: atk.STR, value: 0},
				{characteristic: atk.STR, value: 0},
			},
			attackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{
						spear,
						{
							Name: "Delirium", TargetCharacteristic: atk.WIL,
							Dice: dice.D8, DiceCnt: 1, Charges: 1,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed",
					Attacks: []atk.Attack{spear},
					STR:     8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-2", Name: "John Doe",
					Attacks: []atk.Attack{
						{
							Name: "Fireball", TargetCharacteristic: atk.STR,
							Dice: dice.D6, DiceCnt: 1, Charges: 2,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 1,
					IsDetachment: false,
				},
				{
					ID: "monster-1", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 2,
					IsDetachment: false,
				},
				{
					ID: "monster-2", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 3,
					IsDetachment: false,
				},
			},
			assignedAttackers: [][]attacker{
				{{attackerIdx: 0, attackIdx: 1}},
				{{attackerIdx: 1, attackIdx: 0}},
				{{attackerIdx: 2, attackIdx: 0}},
			},
			usedAttackIdxs: []int{-10, 42, 123456},
			rng:            maxRNG{},
			wantDamage: []damage{
				{characteristic: atk.WIL, value: 8},
				{characteristic: atk.STR, value: 4},
				{characteristic: atk.STR, value: 3},
			},
			wantAttackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{
						spear,
						{
							Name: "Delirium", TargetCharacteristic: atk.WIL,
							Dice: dice.D8, DiceCnt: 1, Charges: 0,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed",
					Attacks: []atk.Attack{spear},
					STR:     8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-2", Name: "John Doe",
					Attacks: []atk.Attack{
						{
							Name: "Fireball", TargetCharacteristic: atk.STR,
							Dice: dice.D6, DiceCnt: 1, Charges: 1,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
		},
		{
			name: "SomeAttackersAttackDifferentTargets",
			damageToDefenders: []damage{
				{characteristic: atk.STR, value: 0},
				{characteristic: atk.STR, value: 0},
				{characteristic: atk.STR, value: 0},
			},
			attackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{spear},
					STR:     8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed",
					Attacks: []atk.Attack{spear},
					STR:     8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-2", Name: "John Doe",
					Attacks: []atk.Attack{
						{
							Name: "Fireball", TargetCharacteristic: atk.STR,
							Dice: dice.D6, DiceCnt: 1, Charges: 1,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 1,
					IsDetachment: false,
				},
				{
					ID: "monster-1", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 2,
					IsDetachment: false,
				},
				{
					ID: "monster-2", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 3,
					IsDetachment: false,
				},
			},
			assignedAttackers: [][]attacker{
				{{attackerIdx: 2, attackIdx: 0}},
				nil,
				{{attackerIdx: 0, attackIdx: 0}},
			},
			usedAttackIdxs: []int{-10, 42, 123456},
			rng:            maxRNG{},
			wantDamage: []damage{
				{characteristic: atk.STR, value: 5},
				{characteristic: atk.STR, value: 0},
				{characteristic: atk.STR, value: 3},
			},
			wantAttackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{spear},
					STR:     8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed",
					Attacks: []atk.Attack{spear},
					STR:     8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-2", Name: "John Doe",
					Attacks: []atk.Attack{
						{
							Name: "Fireball", TargetCharacteristic: atk.STR,
							Dice: dice.D6, DiceCnt: 1, Charges: 0,
							IsBlast: false,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
		},
		{
			name: "SomeAttackersAttackMultipleTargets",
			damageToDefenders: []damage{
				{characteristic: atk.STR, value: 0},
				{characteristic: atk.STR, value: 0},
				{characteristic: atk.STR, value: 0},
			},
			attackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{spear},
					STR:     8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed",
					Attacks: []atk.Attack{spear},
					STR:     8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-2", Name: "John Doe",
					Attacks: []atk.Attack{
						{
							Name: "Paralyze", TargetCharacteristic: atk.DEX,
							Dice: dice.D6, DiceCnt: 1, Charges: 1,
							IsBlast: true,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 1,
					IsDetachment: false,
				},
				{
					ID: "monster-1", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 2,
					IsDetachment: false,
				},
				{
					ID: "monster-2", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 3,
					IsDetachment: false,
				},
			},
			assignedAttackers: [][]attacker{
				{{attackerIdx: 2, attackIdx: 0}},
				{{attackerIdx: 2, attackIdx: 0}},
				{{attackerIdx: 2, attackIdx: 0}},
			},
			usedAttackIdxs: []int{-10, 42, 123456},
			rng:            maxRNG{},
			wantDamage: []damage{
				{characteristic: atk.DEX, value: 6},
				{characteristic: atk.DEX, value: 6},
				{characteristic: atk.DEX, value: 6},
			},
			wantAttackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{spear},
					STR:     8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed",
					Attacks: []atk.Attack{spear},
					STR:     8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-2", Name: "John Doe",
					Attacks: []atk.Attack{
						{
							Name: "Paralyze", TargetCharacteristic: atk.DEX,
							Dice: dice.D6, DiceCnt: 1, Charges: 0,
							IsBlast: true,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
		},
		{
			name: "AllAttackersAttackMultipleTargets",
			damageToDefenders: []damage{
				{characteristic: atk.STR, value: 0},
				{characteristic: atk.STR, value: 0},
				{characteristic: atk.STR, value: 0},
			},
			attackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Delirium", TargetCharacteristic: atk.WIL,
							Dice: dice.D4, DiceCnt: 1, Charges: 1,
							IsBlast: true,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Fireball", TargetCharacteristic: atk.STR,
							Dice: dice.D8, DiceCnt: 1, Charges: 1,
							IsBlast: true,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-2", Name: "John Doe",
					Attacks: []atk.Attack{
						{
							Name: "Paralyze", TargetCharacteristic: atk.DEX,
							Dice: dice.D6, DiceCnt: 1, Charges: 1,
							IsBlast: true,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			defenders: []creat.Creature{
				{
					ID: "monster-0", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 1,
					IsDetachment: false,
				},
				{
					ID: "monster-1", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 2,
					IsDetachment: false,
				},
				{
					ID: "monster-2", Name: "Root Goblin", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 3,
					IsDetachment: false,
				},
			},
			assignedAttackers: [][]attacker{
				{
					{attackerIdx: 0, attackIdx: 0},
					{attackerIdx: 1, attackIdx: 0},
					{attackerIdx: 2, attackIdx: 0},
				},
				{
					{attackerIdx: 0, attackIdx: 0},
					{attackerIdx: 1, attackIdx: 0},
					{attackerIdx: 2, attackIdx: 0},
				},
				{
					{attackerIdx: 0, attackIdx: 0},
					{attackerIdx: 1, attackIdx: 0},
					{attackerIdx: 2, attackIdx: 0},
				},
			},
			usedAttackIdxs: []int{-10, 42, 123456},
			rng:            maxRNG{},
			wantDamage: []damage{
				{characteristic: atk.STR, value: 7},
				{characteristic: atk.STR, value: 6},
				{characteristic: atk.STR, value: 5},
			},
			wantAttackers: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Delirium", TargetCharacteristic: atk.WIL,
							Dice: dice.D4, DiceCnt: 1, Charges: 0,
							IsBlast: true,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed",
					Attacks: []atk.Attack{
						{
							Name: "Fireball", TargetCharacteristic: atk.STR,
							Dice: dice.D8, DiceCnt: 1, Charges: 0,
							IsBlast: true,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-2", Name: "John Doe",
					Attacks: []atk.Attack{
						{
							Name: "Paralyze", TargetCharacteristic: atk.DEX,
							Dice: dice.D6, DiceCnt: 1, Charges: 0,
							IsBlast: true,
						},
					},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resolveAttacks(
				test.damageToDefenders,
				test.attackers,
				test.defenders,
				test.assignedAttackers,
				test.usedAttackIdxs,
				test.rng,
			)
			if !slices.Equal(test.damageToDefenders, test.wantDamage) {
				t.Errorf(
					"resolveAttacks(): damage: want %v, got %v",
					test.wantDamage, test.damageToDefenders,
				)
			}
			if !creat.CreatureSlice(test.attackers).Equals(test.wantAttackers) {
				t.Errorf(
					"resolveAttacks(): attackers' attacks mismatch: want %v, got %v",
					test.wantAttackers, test.attackers,
				)
			}
		})
	}
}

func TestNoDamageDone(t *testing.T) {
	tests := []struct {
		name   string
		damage []damage
		want   bool
	}{
		{
			name:   "EmptyDamage",
			damage: []damage{},
			want:   true,
		},
		{
			name:   "NilDamage",
			damage: nil,
			want:   true,
		},
		{
			name: "AllDamageZero",
			damage: []damage{
				{characteristic: atk.STR, value: 0},
				{characteristic: atk.DEX, value: 0},
				{characteristic: atk.WIL, value: 0},
			},
			want: true,
		},
		{
			name: "SomeDamageSet",
			damage: []damage{
				{characteristic: atk.STR, value: 0},
				{characteristic: atk.DEX, value: 1},
				{characteristic: atk.WIL, value: 0},
			},
			want: false,
		},
		{
			name: "AllDamageSet",
			damage: []damage{
				{characteristic: atk.STR, value: 1},
				{characteristic: atk.DEX, value: 2},
				{characteristic: atk.WIL, value: 3},
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			initial := make([]damage, len(test.damage))
			copy(initial, test.damage)

			if got := noDamageDone(test.damage); got != test.want {
				t.Fatalf("noDamageDone(): want %t, got %t", test.want, got)
			}

			if !slices.Equal(test.damage, initial) {
				t.Fatalf(
					"noDamageDone(): damage mutated, want %v, got %v",
					initial, test.damage,
				)
			}
		})
	}
}

func TestApplyDamageToPlayers(t *testing.T) {
	spear := atk.Attack{
		Name: "Spear", TargetCharacteristic: atk.STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	player := creat.Creature{
		ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
		STR: 0, DEX: 14, WIL: 8, HP: 0, Armor: 0,
		IsDetachment: false,
	}
	tests := []struct {
		name            string
		players         []creat.Creature
		damageToPlayers []damage
		rng             dice.RNG
		want            []creat.Creature
	}{
		{
			name:            "EmptyPlayers",
			players:         []creat.Creature{},
			damageToPlayers: []damage{{characteristic: atk.STR, value: 4}},
			rng:             maxRNG{},
			want:            []creat.Creature{},
		},
		{
			name:            "NilPlayers",
			players:         nil,
			damageToPlayers: []damage{{characteristic: atk.STR, value: 4}},
			rng:             maxRNG{},
			want:            nil,
		},
		{
			name:            "EmptyDamage",
			players:         []creat.Creature{player},
			damageToPlayers: []damage{},
			rng:             maxRNG{},
			want:            []creat.Creature{player},
		},
		{
			name:            "NilDamage",
			players:         []creat.Creature{player},
			damageToPlayers: nil,
			rng:             maxRNG{},
			want:            []creat.Creature{player},
		},
		{
			name:            "NilRNG",
			players:         []creat.Creature{player},
			damageToPlayers: []damage{{characteristic: atk.STR, value: 4}},
			rng:             nil,
			want:            []creat.Creature{player},
		},
		{
			name:    "PlayersShorterThanDamage",
			players: []creat.Creature{player},
			damageToPlayers: []damage{
				{characteristic: atk.STR, value: 4},
				{characteristic: atk.DEX, value: 2},
			},
			rng:  maxRNG{},
			want: []creat.Creature{player},
		},
		{
			name: "PlayersLongerThanDamage",
			players: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			damageToPlayers: []damage{{characteristic: atk.STR, value: 4}},
			rng:             maxRNG{},
			want: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
		},
		{
			name: "AllPlayersOut",
			players: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 0, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 0, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			damageToPlayers: []damage{
				{characteristic: atk.STR, value: 4},
				{characteristic: atk.STR, value: 4},
			},
			rng: maxRNG{},
			want: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 0, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 0, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
		},
		{
			name:            "AllDamageZero",
			players:         []creat.Creature{player},
			damageToPlayers: []damage{{characteristic: atk.STR, value: 0}},
			rng:             maxRNG{},
			want:            []creat.Creature{player},
		},
		{
			name: "DamageToSTRReducesHP",
			players: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 1,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 1,
					IsDetachment: false,
				},
			},
			damageToPlayers: []damage{
				{characteristic: atk.STR, value: 3},
				{characteristic: atk.STR, value: 4},
			},
			rng: maxRNG{},
			want: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 1, Armor: 1,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 0, Armor: 1,
					IsDetachment: false,
				},
			},
		},
		{
			name: "DamageToSTRSuccessfulSave",
			players: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 1,
					IsDetachment: false,
				},
			},
			damageToPlayers: []damage{{characteristic: atk.STR, value: 7}},
			rng:             minRNG{},
			want: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 5, DEX: 14, WIL: 8, HP: 0, Armor: 1,
					IsDetachment: false,
				},
			},
		},
		{
			name: "DamageToSTRFailedSave",
			players: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 1,
					IsDetachment: false,
				},
			},
			damageToPlayers: []damage{{characteristic: atk.STR, value: 7}},
			rng:             maxRNG{},
			want: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 0, Armor: 1,
					IsDetachment: false,
				},
			},
		},
		{
			name: "DamageToSTRKills",
			players: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 1,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 1,
					IsDetachment: false,
				},
			},
			damageToPlayers: []damage{
				{characteristic: atk.STR, value: 12},
				{characteristic: atk.STR, value: 15},
			},
			rng: maxRNG{},
			want: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 0, Armor: 1,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed", Attacks: []atk.Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 0, Armor: 1,
					IsDetachment: false,
				},
			},
		},
		{
			name: "DamageToDEX",
			players: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 1,
					IsDetachment: false,
				},
			},
			damageToPlayers: []damage{{characteristic: atk.DEX, value: 8}},
			rng:             maxRNG{},
			want: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 6, WIL: 8, HP: 4, Armor: 1,
					IsDetachment: false,
				},
			},
		},
		{
			name: "DamageToWIL",
			players: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 1,
					IsDetachment: false,
				},
			},
			damageToPlayers: []damage{{characteristic: atk.WIL, value: 7}},
			rng:             maxRNG{},
			want: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 1, HP: 4, Armor: 1,
					IsDetachment: false,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			applyDamageToPlayers(test.players, test.damageToPlayers, test.rng)
			if !creat.CreatureSlice(test.players).Equals(test.want) {
				t.Fatalf(
					"applyDamageToPlayers(): players mismatch: want %v, got %v",
					test.want, test.players,
				)
			}
		})
	}
}

func TestApplyDamageToMonsters(t *testing.T) {
	t.Fatalf("not implemented") // TODO:
	tests := []struct {
		name             string
		monsters         []creat.Creature
		damageToMonsters []damage
		rng              dice.RNG
		want             []creat.Creature
	}{
		// TODO:
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			applyDamageToMonsters(test.monsters, test.damageToMonsters, test.rng)
			if !creat.CreatureSlice(test.monsters).Equals(test.want) {
				t.Fatalf(
					"applyDamageToMonsters(): monsters mismatch: want %v, got %v",
					test.want, test.monsters,
				)
			}
		})
	}
}

func TestAllOut(t *testing.T) {
	spear := atk.Attack{
		Name: "Spear", TargetCharacteristic: atk.STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	tests := []struct {
		name      string
		creatures []creat.Creature
		want      bool
	}{
		{
			name:      "EmptyCreatures",
			creatures: []creat.Creature{},
			want:      true,
		},
		{
			name:      "NilCreatures",
			creatures: nil,
			want:      true,
		},
		{
			name: "AllCreaturesOut",
			creatures: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 0, DEX: 14, WIL: 8, HP: 0, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 0, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-2", Name: "John Doe", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 0, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: true,
		},
		{
			name: "SomeCreaturesOut",
			creatures: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 0, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-2", Name: "John Doe", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 0, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: false,
		},
		{
			name: "NoCreaturesOut",
			creatures: []creat.Creature{
				{
					ID: "player-0", Name: "John Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-1", Name: "Jane Appleseed", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
				{
					ID: "player-2", Name: "John Doe", Attacks: []atk.Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := allOut(test.creatures); got != test.want {
				t.Fatalf("allOut(): want %t, got %t", test.want, got)
			}
		})
	}
}
