package battle

import (
	"testing"

	"github.com/rozag/cabasi/dice"
)

type dummyRNG struct{}

func (dummyRNG) UintN(n uint) uint { return 0 }

type dummyLog struct{}

func (dummyLog) Roll(c Creature, roll uint8) {}
func (dummyLog) Attack(
	attacker, defender Creature,
	attack Attack,
	damage uint8,
) {
}

func dummyPickTargets(a, d map[CreatureID]Creature) []PickedTargets {
	return nil
}

func TestNewValidation(t *testing.T) {
	rng := dummyRNG{}
	log := dummyLog{}
	tests := []struct {
		rng         dice.RNG
		log         Log
		pickTargets PickTargets
		name        string
		wantErrCnt  uint
	}{
		{
			name: "ValidNew",
			rng:  rng, log: log, pickTargets: dummyPickTargets,
			wantErrCnt: 0,
		},
		{
			name: "NoRNG",
			rng:  nil, log: log, pickTargets: dummyPickTargets,
			wantErrCnt: 1,
		},
		{
			name: "NoLog",
			rng:  rng, log: nil, pickTargets: dummyPickTargets,
			wantErrCnt: 1,
		},
		{
			name: "NoPickTargets",
			rng:  rng, log: log, pickTargets: nil,
			wantErrCnt: 1,
		},
		{
			name: "MultipleErrors",
			rng:  nil, log: nil, pickTargets: nil,
			wantErrCnt: 3,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(test.rng, test.log, test.pickTargets)

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
	spear := Attack{
		Name: "Spear", TargetCharacteristic: STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	player := Creature{
		ID: "player-0", Name: "John Appleseed", Attacks: []Attack{spear},
		STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
		IsDetachment: false,
	}
	monster := Creature{
		ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
		STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
		IsDetachment: false,
	}
	tests := []struct {
		name              string
		players, monsters []Creature
		wantErrCnt        uint
	}{
		{
			name:    "ValidRun",
			players: []Creature{player}, monsters: []Creature{monster},
			wantErrCnt: 0,
		},
		{
			name:    "NilPlayers",
			players: nil, monsters: []Creature{monster},
			wantErrCnt: 1,
		},
		{
			name:    "EmptyPlayers",
			players: []Creature{}, monsters: []Creature{monster},
			wantErrCnt: 1,
		},
		{
			name:    "NilMonsters",
			players: []Creature{player}, monsters: nil,
			wantErrCnt: 1,
		},
		{
			name:    "EmptyMonsters",
			players: []Creature{player}, monsters: []Creature{},
			wantErrCnt: 1,
		},
		{
			name: "InvalidPlayer",
			players: []Creature{
				{
					ID: "", Name: "John Appleseed", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			monsters:   []Creature{monster},
			wantErrCnt: 1,
		},
		{
			name:    "InvalidMonster",
			players: []Creature{player},
			monsters: []Creature{
				{
					ID: "", Name: "Root Goblin", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			wantErrCnt: 1,
		},
		{
			name: "NonUniqueCreatureID",
			players: []Creature{
				{
					ID: "creature", Name: "John Appleseed", Attacks: []Attack{spear},
					STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
					IsDetachment: false,
				},
			},
			monsters: []Creature{
				{
					ID: "creature", Name: "Root Goblin", Attacks: []Attack{spear},
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
			b, err := New(dummyRNG{}, dummyLog{}, dummyPickTargets)
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
	b, err := New(dummyRNG{}, dummyLog{}, dummyPickTargets)
	if err != nil {
		t.Fatalf("New(): want nil error, got %v", err)
	}

	spear := Attack{
		Name: "Spear", TargetCharacteristic: STR,
		Dice: dice.D6, DiceCnt: 1, Charges: -1,
		IsBlast: false,
	}
	originalPlayers := []Creature{
		{
			ID: "player-0", Name: "John Appleseed", Attacks: []Attack{spear},
			STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
			IsDetachment: false,
		},
	}
	originalMonsters := []Creature{
		{
			ID: "monster-0", Name: "Root Goblin", Attacks: []Attack{spear},
			STR: 8, DEX: 14, WIL: 8, HP: 4, Armor: 0,
			IsDetachment: false,
		},
	}

	players := make([]Creature, len(originalPlayers))
	for i, player := range originalPlayers {
		copied := player.DeepCopy()
		players[i] = copied
	}

	monsters := make([]Creature, len(originalMonsters))
	for i, monster := range originalMonsters {
		copied := monster.DeepCopy()
		monsters[i] = copied
	}

	_, err = b.Run(players, monsters)
	if err != nil {
		t.Fatalf("Run(): want nil error, got %v", err)
	}

	if !CreatureSlice(originalPlayers).Equals(players) {
		t.Fatalf(
			"Run(): players mutated, want %v, got %v", originalPlayers, players,
		)
	}
	if !CreatureSlice(originalMonsters).Equals(monsters) {
		t.Fatalf(
			"Run(): monsters mutated, want %v, got %v", originalMonsters, monsters,
		)
	}
}
