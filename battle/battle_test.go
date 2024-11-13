package battle

import (
	"slices"
	"testing"

	"github.com/rozag/cabasi/atk"
	"github.com/rozag/cabasi/creat"
	"github.com/rozag/cabasi/dice"
)

type dummyRNG struct{}

func (dummyRNG) UintN(uint) uint { return 0 }

type dummyLog struct{}

func (dummyLog) Roll(creat.Creature, uint8) {}

func (dummyLog) Attack(creat.Creature, []creat.Creature, atk.Attack, uint) {}

func dummyPickAttack(creat.Creature, []creat.Creature) int { return -1 }

func dummyPickTargets(atk.Attack, []creat.Creature) []uint { return nil }

func TestNewValidation(t *testing.T) {
	rng := dummyRNG{}
	log := dummyLog{}
	tests := []struct {
		rng         dice.RNG
		log         Log
		pickAttack  PickAttack
		pickTargets PickTargets
		name        string
		wantErrCnt  uint
	}{
		{
			name: "ValidNew",
			rng:  rng, log: log,
			pickAttack: dummyPickAttack, pickTargets: dummyPickTargets,
			wantErrCnt: 0,
		},
		{
			name: "NoRNG",
			rng:  nil, log: log,
			pickAttack: dummyPickAttack, pickTargets: dummyPickTargets,
			wantErrCnt: 1,
		},
		{
			name: "NoLog",
			rng:  rng, log: nil,
			pickAttack: dummyPickAttack, pickTargets: dummyPickTargets,
			wantErrCnt: 1,
		},
		{
			name: "NoPickAttack",
			rng:  rng, log: log,
			pickAttack: nil, pickTargets: dummyPickTargets,
			wantErrCnt: 1,
		},
		{
			name: "NoPickTargets",
			rng:  rng, log: log,
			pickAttack: dummyPickAttack, pickTargets: nil,
			wantErrCnt: 1,
		},
		{
			name: "MultipleErrors",
			rng:  nil, log: nil,
			pickAttack: nil, pickTargets: nil,
			wantErrCnt: 4,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(test.rng, test.log, test.pickAttack, test.pickTargets)

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
			b, err := New(dummyRNG{}, dummyLog{}, dummyPickAttack, dummyPickTargets)
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
	b, err := New(dummyRNG{}, dummyLog{}, dummyPickAttack, dummyPickTargets)
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
