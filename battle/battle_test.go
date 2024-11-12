package battle

import (
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
	// TODO:
	t.Fatalf("TestAssignAttackers() not implemented")
}
