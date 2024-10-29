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

func TestRunValidation(t *testing.T) {
	rng := dummyRNG{}
	log := dummyLog{}
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
		rng               dice.RNG
		log               Log
		players, monsters []Creature
		wantErrCnt        uint
	}{
		{
			name: "ValidRun", rng: rng, log: log,
			players: []Creature{player}, monsters: []Creature{monster},
			wantErrCnt: 0,
		},
		{
			name: "NoRNG", rng: nil, log: log,
			players: []Creature{player}, monsters: []Creature{monster},
			wantErrCnt: 1,
		},
		{
			name: "NoLog", rng: rng, log: nil,
			players: []Creature{player}, monsters: []Creature{monster},
			wantErrCnt: 1,
		},
		{
			name: "NilPlayers", rng: rng, log: log,
			players: nil, monsters: []Creature{monster},
			wantErrCnt: 1,
		},
		{
			name: "EmptyPlayers", rng: rng, log: log,
			players: []Creature{}, monsters: []Creature{monster},
			wantErrCnt: 1,
		},
		{
			name: "NilMonsters", rng: rng, log: log,
			players: []Creature{player}, monsters: nil,
			wantErrCnt: 1,
		},
		{
			name: "EmptyMonsters", rng: rng, log: log,
			players: []Creature{player}, monsters: []Creature{},
			wantErrCnt: 1,
		},
		{
			name: "InvalidPlayer", rng: rng, log: log,
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
			name: "InvalidMonster", rng: rng, log: log,
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
			name: "NonUniqueCreatureID", rng: rng, log: log,
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
			name: "MultipleErrors", rng: nil, log: nil,
			players: nil, monsters: nil,
			wantErrCnt: 4,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Run(test.rng, test.log, test.players, test.monsters)

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
